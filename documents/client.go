// Copyright 2020 Self Group Ltd. All Rights Reserved.

package documents

import (
	"crypto/ed25519"
	"encoding/json"
	"time"

	"github.com/joinself/self-go-sdk/pkg/object"
)

// messagingClient handles all interactions with self messaging and its users
type messagingClient interface {
	Send(recipients []string, data []byte) error
	Request(recipients []string, cid string, data []byte, timeout time.Duration) (string, []byte, error)
	Subscribe(msgType string, sub func(sender string, payload []byte))
	Command(command, selfID string, payload []byte) ([]byte, error)
	ListConnections() ([]string, error)
}
type pkiClient interface {
	GetHistory(selfID string) ([]json.RawMessage, error)
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

type requestHelper interface {
	FormatRecipients(recipients []string) ([]string, error)
}

// Service handles all messaging operations
type Service struct {
	selfID         string
	deviceID       string
	keyID          string
	api            restTransport
	sk             ed25519.PrivateKey
	URL            string
	messaging      messagingClient
	pki            pkiClient
	expiry         time.Duration
	fileInteractor remoteFile
	requestHelper  requestHelper
}

// Config stores all configuration needed by the messaging service
type Config struct {
	SelfID         string
	DeviceID       string
	PrivateKey     ed25519.PrivateKey
	Rest           restTransport
	KeyID          string
	Messaging      messagingClient
	PKI            pkiClient
	FileInteractor remoteFile
	RequestHelper  requestHelper
}

// NewService creates a new client for interacting with messaging
func NewService(cfg Config) *Service {
	return &Service{
		selfID:         cfg.SelfID,
		deviceID:       cfg.DeviceID,
		sk:             cfg.PrivateKey,
		keyID:          cfg.KeyID,
		messaging:      cfg.Messaging,
		api:            cfg.Rest,
		pki:            cfg.PKI,
		expiry:         time.Minute,
		fileInteractor: cfg.FileInteractor,
		requestHelper:  cfg.RequestHelper,
	}
}
