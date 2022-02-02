// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"time"

	"github.com/joinself/self-go-sdk/messaging"
	"golang.org/x/crypto/ed25519"
)

type messagingService interface {
	Subscribe(msgType string, h func(m *messaging.Message))
	PermitConnection(selfID string) error
}

type messagingClient interface {
	Send(recipients []string, data []byte) error
}

type restTransport interface {
	Get(path string) ([]byte, error)
}

type Service struct {
	messagingService messagingService
	messagingClient  messagingClient
	api              restTransport
	selfID           string
	deviceID         string
	keyID            string
	expiry           time.Duration
	sk               ed25519.PrivateKey
}

type Config struct {
	SelfID           string
	DeviceID         string
	KeyID            string
	MessagingClient  messagingClient
	MessagingService messagingService
	Rest             restTransport
	PrivateKey       ed25519.PrivateKey
}

func NewService(config Config) *Service {
	return &Service{
		selfID:           config.SelfID,
		deviceID:         config.DeviceID,
		keyID:            config.KeyID,
		messagingClient:  config.MessagingClient,
		messagingService: config.MessagingService,
		api:              config.Rest,
		expiry:           time.Minute,
		sk:               config.PrivateKey,
	}
}
