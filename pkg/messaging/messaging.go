// Copyright 2020 Self Group Ltd. All Rights Reserved.

package messaging

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/pkg/helpers"
	"github.com/joinself/self-go-sdk/pkg/ntp"
	"github.com/joinself/self-go-sdk/pkg/storage"
	"github.com/joinself/self-go-sdk/pkg/transport"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/ed25519"
)

type priority int

const (
	priorityInvisible priority = iota
	priorityVisible
)

var (
	decoder = base64.RawURLEncoding

	priorities = map[string]priority{
		"chat.invite":                 priorityVisible,
		"chat.join":                   priorityInvisible,
		"chat.message":                priorityVisible,
		"chat.message.delete":         priorityInvisible,
		"chat.message.delivered":      priorityInvisible,
		"chat.message.edit":           priorityInvisible,
		"chat.message.read":           priorityInvisible,
		"chat.remove":                 priorityInvisible,
		"document.sign.req":           priorityVisible,
		"identities.authenticate.req": priorityVisible,
		"identities.connections.req":  priorityVisible,
		"identities.facts.query.req":  priorityVisible,
		"identities.facts.issue":      priorityVisible,
		"identities.notify":           priorityVisible,
		"sessions.recover":            priorityInvisible,
	}
)

type response struct {
	sender  string
	payload []byte
}

// Transport the stateful connection used to send and receive messages
type Transport interface {
	Send(recipients []string, mtype string, priority int, data []byte) error
	SendAsync(recipients []string, mtype string, priority int, data []byte, callback func(err error))
	Receive() ([]byte, string, int64, []byte, error)
	Close() error
}

// Storage the storage provider used to encrypt and decrypt messages
type Storage interface {
	AccountOffset(inboxID string) (int64, error)
	Encrypt(from string, to []string, plaintext []byte) ([]byte, error)
	Decrypt(from, to string, offset int64, ciphertext []byte) ([]byte, error)
}

// Config messaging configuration for connecting to self messaging
type Config struct {
	SelfID       string
	DeviceID     string
	KeyID        string
	PrivateKey   ed25519.PrivateKey
	MessagingURL string
	APIURL       string
	Transport    Transport
	Storage      Storage
}

// Client default implementation of a messaging client
type Client struct {
	config        Config
	storage       Storage
	transport     Transport
	responses     sync.Map
	subscriptions sync.Map
	inboxID       string
	closing       chan struct{}
	closed        chan struct{}
	started       bool
}

// New create a new messaging client
func New(config Config) (*Client, error) {
	if config.Storage == nil {
		return nil, errors.New("no storage implementation provided")
	}

	inboxID := fmt.Sprintf("%s:%s", config.SelfID, config.DeviceID)

	if config.Transport == nil {
		offset, err := config.Storage.AccountOffset(inboxID)
		if err != nil {
			return nil, err
		}

		cfg := transport.WebsocketConfig{
			SelfID:       config.SelfID,
			DeviceID:     config.DeviceID,
			PrivateKey:   config.PrivateKey,
			MessagingURL: config.MessagingURL,
			Offset:       offset,
		}

		ws, err := transport.NewWebsocket(cfg)
		if err != nil {
			return nil, err
		}

		config.Transport = ws
	}

	c := Client{
		config:    config,
		responses: sync.Map{},
		transport: config.Transport,
		storage:   config.Storage,
		inboxID:   inboxID,
		closing:   make(chan struct{}, 1),
		closed:    make(chan struct{}, 1),
		started:   false,
	}

	return &c, nil
}

// Start starts the connection with self network.
func (c *Client) Start() bool {
	if c.started {
		return false
	}

	c.started = true
	go c.reader()

	return true
}

// Send sends an encypted message to recipients
func (c *Client) Send(recipients []string, mtype string, plaintext []byte) error {
	ciphertext, err := c.storage.Encrypt(c.inboxID, recipients, plaintext)
	if err != nil {
		return err
	}

	return c.transport.Send(recipients, mtype, int(selectPriority(mtype)), ciphertext)
}

// SendAsync sends an encypted message to recipients asynchromously, returning the servers response via the provided callback
func (c *Client) SendAsync(recipients []string, mtype string, plaintext []byte, callback func(error)) {
	ciphertext, err := c.storage.Encrypt(c.inboxID, recipients, plaintext)
	if err != nil {
		callback(err)
		return
	}

	c.transport.SendAsync(recipients, mtype, int(selectPriority(mtype)), ciphertext, callback)
}

// Request sends a request to a specified identity and blocks until response is received
func (c *Client) Request(recipients []string, cid string, mtype string, data []byte, timeout time.Duration) (string, []byte, error) {
	err := c.Send(recipients, mtype, data)
	if err != nil {
		return "", nil, err
	}

	return c.Wait(cid, timeout)
}

// Register registers a conversation
func (c *Client) Register(cid string) {
	c.responses.LoadOrStore(cid, make(chan response, 1))
}

// Wait waits for a response from a given conversation
func (c *Client) Wait(cid string, timeout time.Duration) (string, []byte, error) {
	r, _ := c.responses.LoadOrStore(cid, make(chan response, 1))

	if timeout == 0 {
		resp := <-r.(chan response)
		return resp.sender, resp.payload, nil
	}

	select {
	case resp := <-r.(chan response):
		return resp.sender, resp.payload, nil
	case <-time.After(timeout):
		return "", nil, errors.New("request timed out")
	}
}

// Subscribe subscribes to a given message type
// @param {String} message type to subscribe to [authentication.RequestAuthentication|
// authentication.ResponseAuthentication|fact.RequestInformation|fact.ResponseInformation]
func (c *Client) Subscribe(msgType string, sub func(sender string, payload []byte)) {
	c.subscriptions.Store(msgType, sub)
}

// Close gracefully closes down the messaging cient
func (c *Client) Close() error {
	c.closing <- struct{}{}
	<-c.closed
	return nil
}

func (c *Client) reader() {
	for {
		// check if reader has been closed
		select {
		case <-c.closing:
			c.closed <- struct{}{}
			return
		default:
		}

		id, sender, offset, ciphertext, err := c.transport.Receive()
		if err != nil {
			if !errors.Is(err, transport.ErrChannelClosed) {
				log.Println("messaging:", err)
			}
			continue
		}

		plaintext, err := c.storage.Decrypt(sender, c.inboxID, offset, ciphertext)
		if err != nil {
			log.Println("messaging:", err)
			if errors.Is(err, storage.ErrDecryptionFailed) {
				log.Println("messaging: decryption failed, re-establishing session with sender")

				// we've entered a failure state with the current
				// session, so let the sender know their message
				// failed to decrypt using a new session
				resp, err := helpers.PrepareJWS(
					map[string]interface{}{
						"typ":           "sessions.recover",
						"jti":           uuid.New().String(),
						"iss":           c.config.SelfID,
						"sub":           strings.Split(sender, ":")[0],
						"aud":           strings.Split(sender, ":")[0],
						"iat":           ntp.TimeFunc().Format(time.RFC3339),
						"exp":           ntp.TimeFunc().Add(time.Hour * 10240).Format(time.RFC3339),
						"from_event_id": string(id),
					},
					c.config.KeyID,
					c.config.PrivateKey,
				)

				if err != nil {
					log.Println("messaging:", err)
					continue
				}

				ciphertext, err = c.storage.Encrypt(c.inboxID, []string{sender}, resp)
				if err != nil {
					log.Println("messaging:", err)
					continue
				}

				mtype := "sessions.recover"

				c.transport.SendAsync([]string{sender}, mtype, int(selectPriority(mtype)), ciphertext, func(err error) {
					if err != nil {
						log.Println("messaging:", err)
					}
				})
			}
			continue
		}

		encPayload := gjson.GetBytes(plaintext, "payload").String()
		if encPayload == "" {
			log.Println("messaging: invalid jws message")
			continue
		}

		payload, err := decoder.DecodeString(encPayload)
		if err != nil {
			log.Println("messaging:", err)
			continue
		}

		cid := gjson.GetBytes(payload, "cid").String()

		ch, ok := c.responses.LoadAndDelete(cid)
		if ok {
			ch.(chan response) <- response{sender, plaintext}
			continue
		}

		typ := gjson.GetBytes(payload, "typ").String()
		fn, ok := c.subscriptions.Load(typ)
		if ok {
			go fn.(func(sender string, plaintext []byte))(sender, plaintext)
			continue
		}

		fn, ok = c.subscriptions.Load("*")
		if ok {
			go fn.(func(sender string, plaintext []byte))(sender, plaintext)
			continue
		}
	}
}

func selectPriority(mtype string) priority {
	p, ok := priorities[mtype]
	if ok {
		return p
	}

	return priorityVisible
}
