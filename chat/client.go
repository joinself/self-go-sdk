// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"time"

	"github.com/joinself/self-go-sdk/messaging"
	"github.com/joinself/self-go-sdk/pkg/object"
	"golang.org/x/crypto/ed25519"
)

// messagingService handles all interactions with the messaging service
type messagingService interface {
	Subscribe(msgType string, h func(m *messaging.Message))
	PermitConnection(selfID string) error
}

// messagingClient handles all interactions with self messaging and its users
type messagingClient interface {
	Send(recipients []string, mtype string, data []byte) error
}

// restTransport handles all interactions with the self api
type restTransport interface {
	Get(path string) ([]byte, error)
	Post(path string, ctype string, data []byte) ([]byte, error)
	BuildURL(path string) string
}

// remoteFile manages interactions with the remote filles
type remoteFile interface {
	SetObject(data []byte) (*object.EncryptedObject, error)
	GetObject(link, key string) ([]byte, error)
}

// Service handles all chat interactions.
type Service struct {
	messagingService messagingService
	messagingClient  messagingClient
	api              restTransport
	selfID           string
	deviceID         string
	keyID            string
	expiry           time.Duration
	sk               ed25519.PrivateKey
	fileInteractor   remoteFile
	environment      string
}

// Config stores all configuration needed by the chat service.
type Config struct {
	SelfID           string
	DeviceID         string
	KeyID            string
	MessagingClient  messagingClient
	MessagingService messagingService
	Rest             restTransport
	PrivateKey       ed25519.PrivateKey
	FileInteractor   remoteFile
	Environment      string
}

// NewService creates a new client for interacting with facts.
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
		fileInteractor:   config.FileInteractor,
		environment:      config.Environment,
	}
}
