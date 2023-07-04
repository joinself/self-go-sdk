// Copyright 2020 Self Group Ltd. All Rights Reserved.

package transport

import "github.com/joinself/self-go-sdk/pkg/protos/msgprotov2"

// Encoder represents the interface an encoder needs to implement to work with a
// given protocol version
type Encoder interface {
	MarshalAuth(device, token string, offset int64) ([]byte, error)
	MarshalMessage(id, sender, recipient, mtype string, priority int, ciphertext []byte) ([]byte, error)
	UnmarshalHeader(data []byte) (Header, error)
	UnmarshalNotification(data []byte) (Notification, error)
	UnmarshalMessage(data []byte) (Message, int64, int64, error)
}

// Header represents an event header
type Header interface {
	Id() []byte
	Msgtype() msgprotov2.MsgType
}

// Notification represents a notification event
type Notification interface {
	Id() []byte
	Msgtype() msgprotov2.MsgType
	Error() []byte
}

// Message represents a message event
type Message interface {
	Id() []byte
	Sender() []byte
	Recipient() []byte
	CiphertextBytes() []byte
}

// Metadata represents a messages metadata
type Metadata interface {
	Offset() int64
	Timestamp() int64
}
