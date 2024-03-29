// Copyright 2020 Self Group Ltd. All Rights Reserved.

package helpers

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"

	"github.com/joinself/self-go-sdk/pkg/siggraph"
	"github.com/square/go-jose"
)

var (
	ErrBadJSONPayload       = errors.New("bad json payload")
	ErrResponseBadSignature = errors.New("bad response signature")
)

func PrepareJWS(req map[string]interface{}, kid string, sk ed25519.PrivateKey) ([]byte, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return []byte(""), err
	}

	opts := &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"kid": kid,
		},
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: sk}, opts)
	if err != nil {
		return []byte(""), err
	}

	signature, err := signer.Sign(payload)
	if err != nil {
		return []byte(""), err
	}

	return []byte(signature.FullSerialize()), nil
}

func ParseJWS(response []byte, history []json.RawMessage) (msg []byte, err error) {
	sg, err := siggraph.New(history)
	if err != nil {
		return msg, err
	}

	kid, err := GetJWSKID(response)
	if err != nil {
		return msg, err
	}

	pk, err := sg.ActiveKey(kid)
	if err != nil {
		return msg, err
	}

	jws, err := jose.ParseSigned(string(response))
	if err != nil {
		return msg, err
	}

	msg, err = jws.Verify(pk)
	if err != nil {
		return msg, ErrResponseBadSignature
	}

	return msg, nil
}
