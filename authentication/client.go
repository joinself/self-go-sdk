// Copyright 2020 Self Group Ltd. All Rights Reserved.

package authentication

import (
	"time"

	"github.com/joinself/self-go-sdk/fact"
)

// Service handles all fact operations
type Service struct {
	requester fact.Service
	expiry    time.Duration
}

// Config stores all configuration needed by the authentication service
type Config struct {
	Requester *fact.Service
}

// NewService creates a new client for interacting with facts
func NewService(cfg Config) *Service {
	return &Service{
		requester: *cfg.Requester,
	}
}
