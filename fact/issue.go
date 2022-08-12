// Copyright 2020 Self Group Ltd. All Rights Reserved.

package fact

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/pkg/helpers"
	"github.com/joinself/self-go-sdk/pkg/ntp"
	"gopkg.in/square/go-jose.v2"
)

type FactGroup struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type FactToIssue struct {
	Key    string     `json:"key"`
	Value  string     `json:"value"`
	Source string     `json:"-"`
	Group  *FactGroup `json:"group,omitempty"`
}

func (f *FactToIssue) validate() error {
	if f.Key == "" {
		return errors.New("fact key not provided")
	}

	if f.Value == "" {
		return errors.New("fact value not provided")
	}

	if f.Source == "" {
		return errors.New("fact source not provided")
	}

	return nil
}

// Issues a fact to a specific user.
func (s *Service) Issue(selfID string, facts []FactToIssue, viewers []string) error {
	if selfID == "" {
		return ErrFactRequestBadIdentity
	}

	if len(facts) == 0 {
		return ErrEmptyFacts
	}

	for _, fact := range facts {
		err := fact.validate()
		if err != nil {
			return err
		}
	}

	return s.sendIssuedFacts(selfID, facts, viewers)
}

func (s *Service) sendIssuedFacts(selfID string, facts []FactToIssue, viewers []string) error {
	opts := &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"kid": s.keyID,
		},
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: s.sk}, opts)
	if err != nil {
		return err
	}

	attestations := make([]json.RawMessage, len(facts))
	for i, f := range facts {
		payload, err := json.Marshal(map[string]interface{}{
			"sub":      selfID,
			"iss":      s.selfID,
			"iat":      ntp.TimeFunc().Format(time.RFC3339),
			"exp":      ntp.TimeFunc().Add(defaultRequestTimeout).Format(time.RFC3339),
			"source":   f.Source,
			"verified": true,
			"facts":    []FactToIssue{f},
		})
		if err != nil {
			return err
		}

		attestation, err := signer.Sign(payload)
		if err != nil {
			return err
		}

		attestations[i] = json.RawMessage(attestation.FullSerialize())
	}

	req := map[string]interface{}{
		"typ":          "identities.facts.issue",
		"iss":          s.selfID,
		"sub":          selfID,
		"aud":          selfID,
		"iat":          ntp.TimeFunc().Format(time.RFC3339),
		"exp":          ntp.TimeFunc().Add(defaultRequestTimeout).Format(time.RFC3339),
		"cid":          uuid.New().String(),
		"jti":          uuid.New().String(),
		"status":       "verified",
		"attestations": attestations,
		"viewers":      viewers,
	}

	jws, err := helpers.PrepareJWS(req, s.keyID, s.sk)
	if err != nil {
		return err
	}

	recipients, err := helpers.PrepareRecipients([]string{selfID}, []string{s.selfID + ":" + s.deviceID}, s.api)
	if err != nil {
		return err
	}

	return s.messaging.Send(recipients, req["type"].(string), jws)
}
