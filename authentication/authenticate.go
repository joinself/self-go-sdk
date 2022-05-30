// Copyright 2020 Self Group Ltd. All Rights Reserved.

package authentication

import (
	"errors"
	"time"

	"github.com/joinself/self-go-sdk/request"
)

var (
	RequestAuthentication  = "identities.authenticate.req"
	ResponseAuthentication = "identities.authenticate.resp"

	ErrMissingConversationID      = errors.New("qr request must specify a unique conversation id")
	ErrRequestTimeout             = errors.New("request timeout")
	ErrResponseBadType            = errors.New("received response is not an authentication response")
	ErrResponseBadIssuer          = errors.New("bad response issuer")
	ErrResponseBadAudience        = errors.New("bad response audience")
	ErrResponseBadSubject         = errors.New("bad response subject")
	ErrResponseBadSignature       = errors.New("bad response signature")
	ErrResponseBadStatus          = errors.New("bad response status")
	ErrInvalidExpiry              = errors.New("invalid expiry format")
	ErrInvalidIssuedAt            = errors.New("invalid issued at format")
	ErrResponseExpired            = errors.New("response has expired")
	ErrResponseIssuedTooSoon      = errors.New("response was issued in the future")
	ErrResponseStatusRejected     = errors.New("authentication was rejected")
	ErrMissingConversationIDForDL = errors.New("deep link request must specify a unique conversation id")
	ErrMissingCallback            = errors.New("deep link request must specify a callback url")
	ErrNotConnected               = errors.New("you're not permitting connections from the specifed recipient")
	ErrNotEnoughCredits           = errors.New("your credits have expired, please log in to the developer portal and top up your account")
)

// QRAuthenticationRequest specifies options in a qr code authentication request
type QRAuthenticationRequest struct {
	ConversationID string
	Expiry         time.Duration
	QRConfig       request.QRConfig
}

// DeepLinkAuthenticationRequest specifies options in a deep link authentication request
type DeepLinkAuthenticationRequest struct {
	Callback       string
	ConversationID string
	Expiry         time.Duration
}

// Response is returned on an asynchronous authentication
// from either a qr code or deep link authentication
type Response struct {
	CID      string
	SelfID   string
	DeviceID string
	Accepted bool
}

type AuthRequest struct {
	SelfID string
	Facts  []request.Fact
}

type AuthRequestAsync struct {
	SelfID string
	Facts  []request.Fact
	CID    string
}

// Request prompts a user to authenticate via biometrics
func (s Service) Request(req AuthRequest) error {
	resp, err := s.requester.Request(&request.FactRequest{
		SelfID: req.SelfID,
		Facts:  req.Facts,
		Auth:   true,
		Expiry: s.expiry,
	})

	if err != nil {
		return err
	}

	if resp.Status == "rejected" {
		return ErrResponseStatusRejected
	} else if resp.Status != "accepted" {
		return ErrResponseBadStatus
	}

	return nil
}

// RequestAsync prompts a user to authenticate via biometrics but
// does not wait for the response.
func (s Service) RequestAsync(req AuthRequestAsync) error {
	return s.requester.RequestAsync(&request.FactRequestAsync{
		SelfID: req.SelfID,
		Facts:  req.Facts,
		Auth:   true,
		CID:    req.CID,
	})
}

// GenerateQRCode generates an authentication request as a qr code
func (s *Service) GenerateQRCode(req *QRAuthenticationRequest) ([]byte, error) {
	return s.requester.GenerateQRCode(&request.QRFactRequest{
		ConversationID: req.ConversationID,
		Expiry:         req.Expiry,
		QRConfig:       req.QRConfig,
		Auth:           true,
	})
}

// GenerateDeepLink generates an authentication request as a deep link
func (s *Service) GenerateDeepLink(req *DeepLinkAuthenticationRequest) (string, error) {
	return s.requester.GenerateDeepLink(&request.DeepLinkFactRequest{
		ConversationID: req.ConversationID,
		Expiry:         req.Expiry,
		Auth:           true,
		Callback:       req.Callback,
	})
}

// WaitForResponse waits for a response from a qr code authentication request
func (s *Service) WaitForResponse(cid string, exp time.Duration) (*request.QRFactResponse, error) {
	return s.requester.WaitForResponse(cid, exp)
}

// Subscribe subscribes to fact request responses
func (s *Service) Subscribe(sub func(sender, cid string, authenticated bool)) {
	s.requester.Subscribe(true, func(sender string, res *request.StandardResponse) {
		sub(sender, res.Conversation, res.Auth)
	})
}
