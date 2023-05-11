// Copyright 2020 Self Group Ltd. All Rights Reserved.

package authentication

import (
	"time"

	"github.com/joinself/self-go-sdk/fact"
)

type Requester interface {
	Request(req *fact.FactRequest) (*fact.FactResponse, error)
	RequestAsync(req *fact.FactRequestAsync) error
	GenerateQRCode(req *fact.QRFactRequest) ([]byte, error)
	GenerateDeepLink(req *fact.DeepLinkFactRequest) (string, error)
	WaitForResponse(cid string, exp time.Duration) (*fact.QRFactResponse, error)
	Subscribe(auth bool, sub func(sender string, res *fact.StandardResponse))
}

// Service handles all fact operations
type Service struct {
	requester Requester
	expiry    time.Duration
}

// Config stores all configuration needed by the authentication service
type Config struct {
	Requester *fact.Service
}

type QRConfig = fact.QRConfig

// NewService creates a new client for interacting with facts
func NewService(cfg Config) *Service {
	return &Service{
		requester: *cfg.Requester,
	}
}
