package voice

import (
	"crypto/ed25519"
	"time"

	"github.com/joinself/self-go-sdk/messaging"
)

// messagingClient handles all interactions with self messaging and its users
type messagingClient interface {
	Send(recipients []string, mtype string, data []byte) error
}

// messagingService handles all interactions with the messaging service
type messagingService interface {
	Subscribe(msgType string, h func(m *messaging.Message))
}

// restTransport handles all interactions with the self api
type restTransport interface {
	Get(path string) ([]byte, error)
	Post(path string, ctype string, data []byte) ([]byte, error)
	BuildURL(path string) string
}

type Service struct {
	messagingClient  messagingClient
	messagingService messagingService
	selfID           string
	deviceID         string
	keyID            string
	api              restTransport
	sk               ed25519.PrivateKey
	expiry           time.Duration
}

type Config struct {
	SelfID           string
	DeviceID         string
	PrivateKey       ed25519.PrivateKey
	Rest             restTransport
	KeyID            string
	MessagingClient  messagingClient
	MessagingService messagingService
}

func NewService(cfg Config) *Service {
	return &Service{
		messagingClient:  cfg.MessagingClient,
		messagingService: cfg.MessagingService,
		selfID:           cfg.SelfID,
		deviceID:         cfg.DeviceID,
		sk:               cfg.PrivateKey,
		keyID:            cfg.KeyID,
		api:              cfg.Rest,
		expiry:           time.Minute,
	}
}
