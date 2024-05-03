// Copyright 2020 Self Group Ltd. All Rights Reserved.

package selfsdk

import (
	"crypto/ed25519"
	"encoding/json"
	"time"

	"github.com/joinself/self-go-sdk/authentication"
	"github.com/joinself/self-go-sdk/chat"
	"github.com/joinself/self-go-sdk/documents"
	"github.com/joinself/self-go-sdk/fact"
	"github.com/joinself/self-go-sdk/identity"
	"github.com/joinself/self-go-sdk/messaging"
	"github.com/joinself/self-go-sdk/pkg/object"
	"github.com/joinself/self-go-sdk/voice"
)

// RestTransport defines the interface required for the sdk to perform
// operations against self's rest api
type RestTransport interface {
	Get(path string) ([]byte, error)
	Post(path string, ctype string, data []byte) ([]byte, error)
	Put(path string, ctype string, data []byte) ([]byte, error)
	Delete(path string) ([]byte, error)
	BuildURL(path string) string
}

// WebsocketTransport defines the interface required for the sdk to perform
// operations against self's websocket services
type WebsocketTransport interface {
	Send(recipients []string, mtype string, priority int, data []byte) error
	SendAsync(recipients []string, mtype string, priority int, data []byte, callback func(error))
	Receive() ([]byte, string, int64, []byte, error)
	Connect() error
	Close() error
}

// MessagingClient defines the interface required for the sdk to perform
// operations against self's messaging service
type MessagingClient interface {
	Start() bool
	Send(recipients []string, mtype string, plaintext []byte) error
	SendAsync(recipients []string, mtype string, plaintext []byte, callback func(error))
	Request(recipients []string, cid string, mtype string, data []byte, timeout time.Duration) (string, []byte, error)
	Register(cid string)
	Wait(cid string, timeout time.Duration) (string, []byte, error)
	Subscribe(msgType string, sub func(sender string, payload []byte))
	Close() error
}

// PKIClient defines the interface required for the sdk to perform
// retrieving identity and device public keys from self
type PKIClient interface {
	GetHistory(selfID string) ([]json.RawMessage, error)
	GetDeviceKey(selfID, deviceID string) ([]byte, error)
	SetDeviceKeys(selfID, deviceID string, pkb []byte) error
	ListDeviceKeys(selfID, deviceID string) ([]byte, error)
}

// Storage the storage interface that is used to handle persistence across
type Storage interface {
	AccountCreate(inboxID string, secretKey ed25519.PrivateKey) error
	AccountOffset(inboxID string) (int64, error)
	Encrypt(from string, to []string, plaintext []byte) ([]byte, error)
	Decrypt(from, to string, offset int64, ciphertext []byte) ([]byte, error)
	Close() error
}

type remoteFile interface {
	SetObject(data []byte) (*object.EncryptedObject, error)
	GetObject(link, key string) ([]byte, error)
}

// Client handles all interactions with self services
type Client struct {
	config     Config
	connectors *Connectors
}

// New creates a new self client
func New(cfg Config) (*Client, error) {
	err := cfg.validate()
	if err != nil {
		return nil, err
	}

	err = cfg.load()
	if err != nil {
		return nil, err
	}

	client := &Client{
		config:     cfg,
		connectors: cfg.Connectors,
	}

	var utcZone = time.FixedZone("UTC", 0)
	time.Local = utcZone

	return client, nil
}

func (c *Client) Start() error {
	c.MessagingService().Start()
	return c.connectors.Websocket.Connect()
}

// FactService returns a client for working with facts
func (c *Client) FactService() *fact.Service {
	cfg := fact.Config{
		SelfID:      c.config.SelfAppID,
		DeviceID:    c.config.DeviceID,
		KeyID:       c.config.kid,
		Environment: c.config.Environment,
		PrivateKey:  c.config.sk,
		Rest:        c.connectors.Rest,
		PKI:         c.connectors.PKI,
		Messaging:   c.connectors.Messaging,
	}
	return fact.NewService(cfg)
}

// IdentityService returns a client for working with identities
func (c *Client) IdentityService() *identity.Service {
	cfg := identity.Config{
		Rest: c.connectors.Rest,
		PKI:  c.connectors.PKI,
	}
	return identity.NewService(cfg)
}

// AuthenticationService returns a client for working with authentication
func (c *Client) AuthenticationService() *authentication.Service {
	return authentication.NewService(authentication.Config{
		Requester: c.FactService(),
	})
}

// MessagingService returns a client for working with messages
func (c *Client) MessagingService() *messaging.Service {
	cfg := messaging.Config{
		SelfID:     c.config.SelfAppID,
		DeviceID:   c.config.DeviceID,
		PrivateKey: c.config.sk,
		KeyID:      c.config.kid,
		Rest:       c.connectors.Rest,
		PKI:        c.connectors.PKI,
		Messaging:  c.connectors.Messaging,
	}

	return messaging.NewService(cfg)
}

// ChatService returns a client for interacting with chat
func (c *Client) ChatService() *chat.Service {
	cfg := chat.Config{
		SelfID:           c.config.SelfAppID,
		PrivateKey:       c.config.sk,
		KeyID:            c.config.kid,
		Rest:             c.connectors.Rest,
		MessagingService: c.MessagingService(),
		MessagingClient:  c.connectors.Messaging,
		FileInteractor:   c.connectors.FileInteractor,
		Environment:      c.config.Environment,
	}

	return chat.NewService(cfg)
}

// DocsService returns a client for interacting with document signatures.
func (c *Client) DocsService() *documents.Service {
	cfg := documents.Config{
		SelfID:         c.config.SelfAppID,
		DeviceID:       c.config.DeviceID,
		PrivateKey:     c.config.sk,
		KeyID:          c.config.kid,
		Messaging:      c.connectors.Messaging,
		Rest:           c.connectors.Rest,
		PKI:            c.connectors.PKI,
		FileInteractor: c.connectors.FileInteractor,
	}

	return documents.NewService(cfg)
}

// VoiceService returns a client for managing voice call negotiation.
func (c *Client) VoiceService() *voice.Service {
	cfg := voice.Config{
		SelfID:           c.config.SelfAppID,
		DeviceID:         c.config.DeviceID,
		PrivateKey:       c.config.sk,
		KeyID:            c.config.kid,
		Rest:             c.connectors.Rest,
		MessagingService: c.MessagingService(),
		MessagingClient:  c.connectors.Messaging,
	}
	return voice.NewService(cfg)
}

// Rest provides access to the rest client to interact used by the sdk.
func (c *Client) Rest() RestTransport {
	return c.config.Connectors.Rest
}

// Close gracefully closes the self client
func (c *Client) Close() error {
	err := c.connectors.Websocket.Close()
	if err != nil {
		return err
	}

	err = c.connectors.Messaging.Close()
	if err != nil {
		return err
	}

	return c.connectors.Storage.Close()
}

// SelfAppID returns the current SelfAppID for this app.
func (c *Client) SelfAppID() string {
	return c.config.SelfAppID
}
