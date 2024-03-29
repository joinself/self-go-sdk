// Copyright 2020 Self Group Ltd. All Rights Reserved.

package transport

import (
	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/pkg/protos/msgprotov1"
	"github.com/joinself/self-go-sdk/pkg/protos/msgprotov2"
	"google.golang.org/protobuf/proto"
)

type encoderV1 struct {
}

func newEncoderV1() Encoder {
	return &encoderV1{}
}

// MarshalAuth creates a protocol v1 auth message
func (e *encoderV1) MarshalAuth(device, token string, offset int64) ([]byte, error) {
	return proto.Marshal(&msgprotov1.Auth{
		Id:     uuid.New().String(),
		Type:   msgprotov1.MsgType_AUTH,
		Device: device,
		Token:  token,
		Offset: offset,
	})
}

// MarshalMessage creates a protocol v1 message
func (e *encoderV1) MarshalMessage(id, sender, recipient, mtype string, priority int, ciphertext []byte) ([]byte, error) {
	return proto.Marshal(&msgprotov1.Message{
		Id:         id,
		Type:       msgprotov1.MsgType_MSG,
		Sender:     sender,
		Recipient:  recipient,
		Ciphertext: ciphertext,
	})
}

type header struct {
	id []byte
	mt msgprotov2.MsgType
}

// Id implements the id method for a header
func (h *header) Id() []byte {
	return h.id
}

// Msgtype implements the msg type method for a header
func (h *header) Msgtype() msgprotov2.MsgType {
	return h.mt
}

// UnmarshalHeader reads a protocol v1 header event
func (e *encoderV1) UnmarshalHeader(data []byte) (Header, error) {
	var h msgprotov1.Header

	err := proto.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}

	return &header{
		id: []byte(h.Id),
		mt: msgprotov2.MsgType(h.Type),
	}, nil
}

type notification struct {
	id []byte
	mt msgprotov2.MsgType
	er []byte
}

// Id implements the id method for a notification
func (n *notification) Id() []byte {
	return n.id
}

// Msgtype implements the msg type method for a notification
func (n *notification) Msgtype() msgprotov2.MsgType {
	return n.mt
}

// Msgtype implements the msg type method for a notification
func (n *notification) Error() []byte {
	return n.er
}

// UnmarshalNotification reads a protocol v1 notification event
func (e *encoderV1) UnmarshalNotification(data []byte) (Notification, error) {
	var n msgprotov1.Notification

	err := proto.Unmarshal(data, &n)
	if err != nil {
		return nil, err
	}

	return &notification{
		id: []byte(n.Id),
		mt: msgprotov2.MsgType(n.Type),
		er: []byte(n.Error),
	}, nil
}

type message struct {
	id []byte
	sd []byte
	rt []byte
	ct []byte
}

// Id implements the id method for a message
func (m *message) Id() []byte {
	return m.id
}

// Sender implements the sender method for a message
func (m *message) Sender() []byte {
	return m.sd
}

// Recipient implements the recipient method for a message
func (m *message) Recipient() []byte {
	return m.rt
}

// CiphertextBytes implements the ciphertext method for a message
func (m *message) CiphertextBytes() []byte {
	return m.ct
}

// UnmarshalMessage reads a protocol v1 message event
func (e *encoderV1) UnmarshalMessage(data []byte) (Message, int64, int64, error) {
	var m msgprotov1.Message

	err := proto.Unmarshal(data, &m)
	if err != nil {
		return nil, 0, 0, err
	}

	return &message{
		id: []byte(m.Id),
		sd: []byte(m.Sender),
		rt: []byte(m.Recipient),
		ct: m.Ciphertext,
	}, m.Timestamp.AsTime().Unix(), m.Offset, nil
}
