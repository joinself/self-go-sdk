// Copyright 2020 Self Group Ltd. All Rights Reserved.

package messaging

import (
	"encoding/json"
	"time"

	"golang.org/x/crypto/ed25519"
)

// restTransport handles all interactions with the self api
type restTransport interface {
	Get(path string) ([]byte, error)
}

// messagingClient handles all interactions with self messaging and its users
type messagingClient interface {
	Start() bool
	Send(recipients []string, mtype string, data []byte) error
	SendAsync(recipients []string, mtype string, data []byte, callback func(error))
	Request(recipients []string, cid string, mtype string, data []byte, timeout time.Duration) (string, []byte, error)
	Subscribe(msgType string, sub func(sender string, payload []byte))
}

type pkiClient interface {
	GetHistory(selfID string) ([]json.RawMessage, error)
}

// Service handles all messaging operations
type Service struct {
	selfID    string
	deviceID  string
	keyID     string
	sk        ed25519.PrivateKey
	api       restTransport
	pki       pkiClient
	messaging messagingClient
}

// Config stores all configuration needed by the messaging service
type Config struct {
	SelfID     string
	DeviceID   string
	KeyID      string
	PrivateKey ed25519.PrivateKey
	PKI        pkiClient
	Messaging  messagingClient
	Rest       restTransport
}

type jwsPayload struct {
	ID           string    `json:"jti"`
	Conversation string    `json:"cid"`
	Issuer       string    `json:"iss"`
	Audience     string    `json:"aud"`
	Subject      string    `json:"sub"`
	IssuedAt     time.Time `json:"iat"`
	ExpiresAt    time.Time `json:"exp"`
}

// NewService creates a new client for interacting with messaging
func NewService(cfg Config) *Service {
	return &Service{
		selfID:    cfg.SelfID,
		deviceID:  cfg.DeviceID,
		keyID:     cfg.KeyID,
		sk:        cfg.PrivateKey,
		api:       cfg.Rest,
		pki:       cfg.PKI,
		messaging: cfg.Messaging,
	}
}
