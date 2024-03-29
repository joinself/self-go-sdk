// Copyright 2020 Self Group Ltd. All Rights Reserved.

package transport

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joinself/self-go-sdk/pkg/pqueue"
	"github.com/joinself/self-go-sdk/pkg/protos/msgprotov2"
	"golang.org/x/crypto/ed25519"
)

const (
	priorityClose = iota
	priorityPong
	priorityNotification
	priorityMessage
)

// ErrChannelClosed returned when the websocket connection is shut down manually
var ErrChannelClosed = errors.New("channel closed")

type (
	sigclose bool
	sigpong  bool
)

type msg struct {
	msg    Message
	offset int64
}

// WebsocketConfig configuration for connecting to a websocket
type WebsocketConfig struct {
	MessagingURL string
	SelfID       string
	DeviceID     string
	KeyID        string
	Offset       int64
	PrivateKey   ed25519.PrivateKey
	TCPDeadline  time.Duration
	InboxSize    int
	OnConnect    func()
	OnDisconnect func(err error)
	OnPing       func()
	messagingID  string
}

func (c *WebsocketConfig) load() {
	c.messagingID = fmt.Sprintf(
		"%s:%s",
		c.SelfID,
		c.DeviceID,
	)
}

// Websocket websocket client for self messaging
type Websocket struct {
	config    WebsocketConfig
	ws        *websocket.Conn
	queue     *pqueue.Queue
	inbox     chan msg
	responses sync.Map
	offset    int64
	closed    int32
	shutdown  int32
	enc       Encoder
}

type event struct {
	id   string
	data []byte
	err  chan error
	cb   func(err error)
}

// NewWebsocket creates a new websocket connection
func NewWebsocket(config WebsocketConfig) (*Websocket, error) {
	config.load()

	var enc Encoder

	if strings.Contains(config.MessagingURL, "/v1/messaging") {
		enc = newEncoderV1()
	} else {
		enc = newEncoderV2()
	}

	c := Websocket{
		config:    config,
		queue:     pqueue.New(4, config.InboxSize),
		inbox:     make(chan msg, config.InboxSize),
		responses: sync.Map{},
		offset:    config.Offset,
		closed:    1,
		enc:       enc,
	}

	return &c, nil
}

// Send send a message to given recipients. recipient is a combination of "selfID:deviceID"
func (c *Websocket) Send(recipients []string, mtype string, priority int, data []byte) error {
	for _, r := range recipients {
		id := uuid.New().String()

		msg, err := c.enc.MarshalMessage(id, c.config.messagingID, r, mtype, priority, data)
		if err != nil {
			return err
		}

		e := event{
			id:   id,
			data: msg,
			err:  make(chan error, 1),
		}

		c.queue.Push(priorityMessage, &e)

		err = <-e.err
		if err != nil {
			return err
		}
	}

	return nil
}

// SendAsync send a message to given recipients with a callback to handle the server response
func (c *Websocket) SendAsync(recipients []string, mtype string, priority int, data []byte, callback func(err error)) {
	for _, r := range recipients {
		id := uuid.New().String()

		msg, err := c.enc.MarshalMessage(id, c.config.messagingID, r, mtype, priority, data)
		if err != nil {
			callback(err)
			return
		}

		e := event{
			id:   id,
			data: msg,
			cb:   callback,
		}

		c.queue.Push(priorityMessage, &e)
	}
}

// SendAsync send a message with a given id to a single recipient, with a callback to handle the server response
func (c *Websocket) SendAsyncWithID(id, recipient string, mtype string, priority int, data []byte, callback func(err error)) {
	msg, err := c.enc.MarshalMessage(id, c.config.messagingID, recipient, mtype, priority, data)
	if err != nil {
		callback(err)
		return
	}

	e := event{
		id:   id,
		data: msg,
		cb:   callback,
	}

	c.queue.Push(priorityMessage, &e)
}

// Receive receive a message
func (c *Websocket) Receive() ([]byte, string, int64, []byte, error) {
	m, ok := <-c.inbox
	if !ok {
		return nil, "", -1, nil, ErrChannelClosed
	}

	return m.msg.Id(), string(m.msg.Sender()), m.offset, m.msg.CiphertextBytes(), nil
}

// Close closes the messaging clients persistent connection
func (c *Websocket) Close() error {
	atomic.StoreInt32(&c.shutdown, 1)

	c.close(nil)

	// wait for subscribers to drain
	for len(c.inbox) > 0 {
		time.Sleep(time.Millisecond)
	}

	return nil
}

func (c *Websocket) pingHandler(string) error {
	if c.config.OnPing != nil {
		c.config.OnPing()
	}

	deadline := time.Now().Add(c.config.TCPDeadline)

	c.queue.Push(priorityPong, sigpong(true))

	return c.ws.SetReadDeadline(deadline)
}

func (c *Websocket) Connect() error {
	if !atomic.CompareAndSwapInt32(&c.closed, 1, 0) {
		return errors.New("could not connect")
	}

	var ws *websocket.Conn
	var connected bool

	defer func(success *bool) {
		if !(*success) {
			// if it failed to reconnect, set the connection status to closed
			atomic.CompareAndSwapInt32(&c.closed, 0, 1)

			// close the connection
			if ws == nil {
				return
			}

			log.Println("[websocket] closing errored connection")

			err := ws.Close()
			if err != nil {
				log.Println("[websocket]", err)
			}
		}
	}(&connected)

	token, err := GenerateToken(c.config.SelfID, c.config.KeyID, c.config.PrivateKey)
	if err != nil {
		return err
	}

	ws, _, err = websocket.DefaultDialer.Dial(c.config.MessagingURL, nil)
	if err != nil {
		return err
	}

	ws.SetPingHandler(c.pingHandler)

	c.ws = ws

	auth, err := c.enc.MarshalAuth(c.config.DeviceID, token, c.offset)
	if err != nil {
		return err
	}

	err = c.ws.WriteMessage(websocket.BinaryMessage, auth)
	if err != nil {
		return err
	}

	c.ws.SetReadDeadline(time.Now().Add(c.config.TCPDeadline))
	_, data, err := c.ws.ReadMessage()
	if err != nil {
		log.Println("[websocket] authentication timeout:", err.Error())
		return err
	}

	resp, err := c.enc.UnmarshalNotification(data)
	if err != nil {
		return err
	}

	switch resp.Msgtype() {
	case msgprotov2.MsgTypeACK:
	case msgprotov2.MsgTypeERR:
		return errors.New(string(resp.Error()))
	default:
		return errors.New("unknown authentication response")
	}

	connected = true

	go c.reader()
	go c.writer()

	if c.config.OnConnect != nil {
		c.config.OnConnect()
	}

	return nil
}

func (c *Websocket) reader() {
	for {
		if c.isShutdown() {
			close(c.inbox)
			return
		}

		if c.isClosed() {
			return
		}

		_, data, err := c.ws.ReadMessage()
		if err != nil {
			if c.isShutdown() {
				close(c.inbox)
			} else {
				c.reconnect(err)
			}
			return
		}

		hdr, err := c.enc.UnmarshalHeader(data)
		if err != nil {
			if c.isShutdown() {
				close(c.inbox)
			} else {
				c.reconnect(err)
			}
			return
		}

		switch hdr.Msgtype() {
		case msgprotov2.MsgTypeACK, msgprotov2.MsgTypeERR:
			n, err := c.enc.UnmarshalNotification(data)
			if err != nil {
				log.Printf("[websocket] failed to unmarshal notification: %s", err.Error())
				continue
			}

			pch, ok := c.responses.Load(string(n.Id()))
			if !ok {
				continue
			}

			c.responses.Delete(string(n.Id()))

			var rerr error

			if n.Msgtype() == msgprotov2.MsgTypeERR {
				rerr = errors.New(string(n.Error()))
			}

			rev := pch.(*event)

			if rev.cb != nil {
				rev.cb(rerr)
			} else {
				rev.err <- rerr
			}

		case msgprotov2.MsgTypeMSG:
			m, _, offset, err := c.enc.UnmarshalMessage(data)
			if err != nil {
				log.Printf("[websocket] failed to unmarshal notification: %s", err.Error())
				continue
			}

			c.offset = offset

			c.inbox <- msg{msg: m, offset: offset}
		}
	}
}

func (c *Websocket) writer() {
	var err error

	for {
		p, e := c.queue.PopWithPrioriry()

		switch p {
		case priorityClose:
			for i := priorityClose; i <= priorityMessage; i++ {
				c.queue.Flush(i)
			}
			return
		case priorityPong:
			deadline := time.Now().Add(c.config.TCPDeadline)
			err = c.ws.WriteControl(websocket.PongMessage, nil, deadline)
		case priorityNotification, priorityMessage:
			ev := e.(*event)
			c.responses.Store(ev.id, ev)

			err = c.ws.WriteMessage(websocket.BinaryMessage, ev.data)
			if err != nil {
				if ev.cb != nil {
					ev.cb(err)
				} else {
					ev.err <- err
				}
				continue
			}
		}

		if err != nil {
			c.close(err)
		}
	}
}

func (c *Websocket) reconnect(err error) {
	if !c.close(err) {
		return
	}

	switch e := err.(type) {
	case net.Error:
		if !e.Timeout() {
			log.Println("[websocket]", e)
		}
	case *websocket.CloseError:
		if e.Code != websocket.CloseAbnormalClosure {
			if e.Text != io.ErrUnexpectedEOF.Error() {
				return
			}
		}
	}

	for i := 0; i < 20; i++ {
		log.Println("[websocket] attempting reconnect")

		time.Sleep(c.config.TCPDeadline)

		err := c.Connect()
		if err == nil {
			return
		}

		log.Println("[websocket] failed to connect to messaging")
	}
}

func (c *Websocket) close(err error) bool {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return false
	}

	if c.config.OnDisconnect != nil {
		c.config.OnDisconnect(err)
	}

	c.queue.Push(priorityClose, sigclose(true))

	time.Sleep(time.Millisecond * 10)
	c.ws.Close()

	c.responses = sync.Map{}

	return true
}

func (c *Websocket) isClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func (c *Websocket) isShutdown() bool {
	return atomic.LoadInt32(&c.shutdown) == 1
}
