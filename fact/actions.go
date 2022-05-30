// Copyright 2020 Self Group Ltd. All Rights Reserved.

package fact

import (
	"encoding/json"
	"time"

	"github.com/joinself/self-go-sdk/request"
)

// Attest creates an attested fact about a self identity
func (s Service) Attest(selfID string, facts []request.Fact) ([]json.RawMessage, error) {
	return s.requester.Attest(selfID, facts)
}

// Request requests a fact from a given identity
func (s Service) Request(req *request.FactRequest) (*request.FactResponse, error) {
	return s.requester.Request(req)
}

// RequestAsync requests a fact from a given identity and does not
// wait for the response
func (s Service) RequestAsync(req *request.FactRequestAsync) error {
	return s.requester.RequestAsync(req)
}

// RequestViaIntermediary requests a fact from a given identity via a trusted
// intermediary. The intermediary verifies that the identity has a given fact
// and that it meets the requested requirements.
func (s Service) RequestViaIntermediary(req *request.IntermediaryFactRequest) (*request.IntermediaryFactResponse, error) {
	return s.requester.RequestViaIntermediary(req)
}

// GenerateQRCode generates a qr code containing an fact request
func (s Service) GenerateQRCode(req *request.QRFactRequest) ([]byte, error) {
	return s.requester.GenerateQRCode(req)
}

// GenerateDeepLink generates a qr code containing an fact request
func (s Service) GenerateDeepLink(req *request.DeepLinkFactRequest) (string, error) {
	return s.requester.GenerateDeepLink(req)
}

// WaitForResponse waits for completion of a fact request that was initiated by qr code
func (s Service) WaitForResponse(cid string, exp time.Duration) (*request.QRFactResponse, error) {
	return s.requester.WaitForResponse(cid, exp)
}

// Subscribe subscribes to fact request responses
func (s Service) Subscribe(sub func(sender string, res *request.QRFactResponse)) {
	s.requester.Subscribe(false, func(sender string, res *request.StandardResponse) {
		sub(sender, &request.QRFactResponse{
			Responder: sender,
			Facts:     res.Facts,
		})
	})
}
