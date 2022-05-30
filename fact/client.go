// Copyright 2020 Self Group Ltd. All Rights Reserved.

package fact

import (
	"github.com/joinself/self-go-sdk/request"
)

// Service handles all fact operations
type Service struct {
	requester *request.Service
}

// Config stores all configuration needed by the fact service
type Config struct {
	Requester *request.Service
}

// NewService creates a new client for interacting with facts
func NewService(cfg Config) *Service {
	return &Service{
		requester: cfg.Requester,
	}
}
