// Copyright 2020 Self Group Ltd. All Rights Reserved.

package messaging

import (
	"errors"

	"github.com/joinself/self-go-sdk/pkg/transport"
)

type testEvent struct {
	sender    string
	recipient string
	data      []byte
}

type testWebsocket struct {
	in     chan *testEvent
	out    chan *testEvent
	offset int64
}

func newTestWebsocket() *testWebsocket {
	return &testWebsocket{
		in:  make(chan *testEvent, 10),
		out: make(chan *testEvent, 10),
	}
}

func (c *testWebsocket) Send(recipients []string, mtype string, priority int, data []byte) error {
	for _, r := range recipients {
		if r == "non-existent" {
			return errors.New("recipient does not exist")
		}
		c.out <- &testEvent{recipient: r, data: data}
	}
	return nil
}

func (c *testWebsocket) SendAsync(recipients []string, mtype string, priority int, data []byte, callback func(error)) {
	for _, r := range recipients {
		if r == "non-existent" {
			callback(errors.New("recipient does not exist"))
			return
		}
		c.out <- &testEvent{recipient: r, data: data}
	}
	callback(nil)
}

func (c *testWebsocket) Receive() (string, int64, []byte, error) {
	c.offset++

	e, ok := <-c.in
	if !ok {
		return "", c.offset, nil, transport.ErrChannelClosed
	}

	if e.recipient == "failure" {
		return "", c.offset, nil, errors.New("transport failure")
	}
	return e.sender, c.offset, e.data, nil
}

func (c *testWebsocket) Command(command string, payload []byte) ([]byte, error) {
	return []byte(`["*"]`), nil
}

func (c *testWebsocket) Connect() error {
	return nil
}

func (c *testWebsocket) Close() error {
	return nil
}

type testStorage struct {
	offset int64
}

func newTestStorage() *testStorage {
	return &testStorage{}
}

func (c *testStorage) AccountOffset(inboxID string) (int64, error) {
	return c.offset, nil
}

func (c *testStorage) Encrypt(from string, to []string, data []byte) ([]byte, error) {
	// fake encrypt the payload
	return []byte(decoder.EncodeToString(data)), nil
}

func (c *testStorage) Decrypt(from string, to string, offset int64, data []byte) ([]byte, error) {
	// fake decrypt the payload
	c.offset = offset
	return decoder.DecodeString(string(data))
}
