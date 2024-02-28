// Copyright 2020 Self Group Ltd. All Rights Reserved.

package transport

import (
	"runtime"
	"testing"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/joinself/self-go-sdk/pkg/protos/msgprotov2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsocketConnect(t *testing.T) {
	s := newTestMessagingServer(t)
	defer s.s.Close()

	cfg := WebsocketConfig{
		SelfID:       "test",
		DeviceID:     "1",
		PrivateKey:   sk,
		MessagingURL: s.endpoint,
		TCPDeadline:  time.Millisecond * 100,
		InboxSize:    128,
	}

	c, err := NewWebsocket(cfg)
	require.Nil(t, err)
	err = c.Connect()
	require.Nil(t, err)
	require.Nil(t, c.Close())
}

func TestWebsocketReconnect(t *testing.T) {
	s := newTestMessagingServer(t)
	defer s.s.Close()

	connected := make(chan bool, 1)
	disconnected := make(chan bool, 1)

	cfg := WebsocketConfig{
		SelfID:       "test",
		DeviceID:     "1",
		PrivateKey:   sk,
		MessagingURL: s.endpoint,
		TCPDeadline:  time.Second,
		InboxSize:    128,
		OnConnect: func() {
			connected <- true
		},
		OnDisconnect: func(err error) {
			disconnected <- true
		},
	}

	c, err := NewWebsocket(cfg)
	require.Nil(t, err)
	err = c.Connect()
	require.Nil(t, err)
	defer c.Close()

	assert.True(t, <-connected)

	s.stop <- true

	assert.True(t, <-disconnected)

	s = newTestMessagingServer(t)
	defer s.s.Close()

	assert.True(t, <-connected)
}

func TestWebsocketSend(t *testing.T) {
	s := newTestMessagingServer(t)
	defer s.s.Close()

	cfg := WebsocketConfig{
		SelfID:       "test",
		DeviceID:     "1",
		PrivateKey:   sk,
		MessagingURL: s.endpoint,
		TCPDeadline:  time.Millisecond * 100,
		InboxSize:    128,
	}

	c, err := NewWebsocket(cfg)
	require.Nil(t, err)
	err = c.Connect()
	require.Nil(t, err)
	defer c.Close()

	err = c.Send([]string{"alice:1"}, "test.message", 1, []byte("test"))
	require.Nil(t, err)

	msg, err := wait(s.in)
	require.Nil(t, err)

	assert.NotEmpty(t, msg.Id)
	assert.Equal(t, []byte("test:1"), msg.Sender())
	assert.Equal(t, []byte("alice:1"), msg.Recipient())
	assert.Equal(t, []byte("test"), msg.CiphertextBytes())
}

func TestWebsocketReceive(t *testing.T) {
	s := newTestMessagingServer(t)
	defer s.s.Close()

	cfg := WebsocketConfig{
		SelfID:       "test",
		DeviceID:     "1",
		PrivateKey:   sk,
		MessagingURL: s.endpoint,
		TCPDeadline:  time.Millisecond * 100,
		InboxSize:    128,
	}

	c, err := NewWebsocket(cfg)
	require.Nil(t, err)
	err = c.Connect()
	require.Nil(t, err)
	defer c.Close()

	b := flatbuffers.NewBuilder(1024)

	mid := b.CreateString("test")
	msd := b.CreateString("alice:1")
	mrp := b.CreateString("test:1")
	mct := b.CreateByteVector([]byte("test"))

	msgprotov2.MessageStart(b)
	msgprotov2.MessageAddId(b, mid)
	msgprotov2.MessageAddMsgtype(b, msgprotov2.MsgTypeMSG)
	msgprotov2.MessageAddSender(b, msd)
	msgprotov2.MessageAddRecipient(b, mrp)
	msgprotov2.MessageAddMetadata(b, msgprotov2.CreateMetadata(
		b,
		0,
		0,
	))

	msgprotov2.MessageAddCiphertext(b, mct)
	msg := msgprotov2.MessageEnd(b)

	b.Finish(msg)

	s.out <- b.FinishedBytes()

	_, sender, _, m, err := c.Receive()
	require.Nil(t, err)

	assert.Equal(t, "alice:1", sender)
	assert.Equal(t, []byte("test"), m)
}

func TestWebsocketClose(t *testing.T) {
	s := newTestMessagingServer(t)
	defer s.s.Close()

	cfg := WebsocketConfig{
		SelfID:       "test",
		DeviceID:     "1",
		PrivateKey:   sk,
		MessagingURL: s.endpoint,
		TCPDeadline:  time.Millisecond * 100,
		InboxSize:    10240,
	}

	c, err := NewWebsocket(cfg)
	require.Nil(t, err)
	err = c.Connect()
	require.Nil(t, err)

	// handle received messages
	go func() {
		for i := 0; i < 10000; i++ {
			c.Receive()
		}
	}()

	go func() {
		// send some messages
		b := flatbuffers.NewBuilder(1024)

		for i := 0; i < 10000; i++ {
			b.Reset()

			mid := b.CreateString("test")
			msd := b.CreateString("alice:1")
			mrp := b.CreateString("test:1")
			mct := b.CreateByteVector([]byte("test"))

			msgprotov2.MessageStart(b)
			msgprotov2.MessageAddId(b, mid)
			msgprotov2.MessageAddMsgtype(b, msgprotov2.MsgTypeMSG)
			msgprotov2.MessageAddSender(b, msd)
			msgprotov2.MessageAddRecipient(b, mrp)
			msgprotov2.MessageAddMetadata(b, msgprotov2.CreateMetadata(
				b,
				0,
				0,
			))

			msgprotov2.MessageAddCiphertext(b, mct)
			msg := msgprotov2.MessageEnd(b)

			b.Finish(msg)

			fb := b.FinishedBytes()
			m := make([]byte, len(fb))
			copy(m, fb)

			s.out <- m
		}
	}()

	for len(c.inbox) < 1 {
		time.Sleep(time.Millisecond)
	}

	// close the connection
	c.Close()

	// check all messages have been processed
	assert.Len(t, c.inbox, 0)
}

func TestWebsocketCleanup(t *testing.T) {
	time.Sleep(time.Second)

	s := newTestMessagingServerUnresponsive(t)
	defer s.s.Close()

	cfg := WebsocketConfig{
		SelfID:       "test",
		DeviceID:     "1",
		PrivateKey:   sk,
		MessagingURL: s.endpoint,
		TCPDeadline:  time.Millisecond * 100,
		InboxSize:    10240,
	}

	goroutines := runtime.NumGoroutine()

	for i := 0; i < 5; i++ {
		cc, err := NewWebsocket(cfg)
		require.Nil(t, err)

		err = cc.Connect()
		require.NotNil(t, err)

		time.Sleep(time.Millisecond * 100)

		assert.Equal(t, goroutines, runtime.NumGoroutine())
	}

	rs := newTestMessagingServer(t)
	defer rs.s.Close()

	cfg.MessagingURL = rs.endpoint

	c, err := NewWebsocket(cfg)
	require.Nil(t, err)
	err = c.Connect()
	require.Nil(t, err)

	assert.NotEqual(t, goroutines, runtime.NumGoroutine())

	c.Close()
}
