// Copyright 2020 Self Group Ltd. All Rights Reserved.

package request

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"

	"github.com/joinself/self-go-sdk/pkg/kidhelper"
	"github.com/joinself/self-go-sdk/pkg/siggraph"
	"github.com/square/go-jose"
)

var (
	ErrBadJSONPayload       = errors.New("bad json payload")
	ErrResponseBadSignature = errors.New("bad response signature")
)

type RestTransport interface {
	Get(path string) ([]byte, error)
}

func FormatRecipients(recipients []string, selfID, deviceID string, api RestTransport) ([]string, error) {
	devices := make([]string, 0)
	for _, sID := range recipients {
		dds, err := getDevices(api, sID)
		if err != nil {
			return nil, err
		}

		for i := range dds {
			// if is not the current device
			if sID != selfID && dds[i] != deviceID {
				devices = append(devices, sID+":"+dds[i])
			}
		}
	}

	return devices, nil
}

func getDevices(api RestTransport, selfID string) ([]string, error) {
	var resp []byte
	var err error

	if len(selfID) > 11 {
		resp, err = api.Get("/v1/apps/" + selfID + "/devices")
	} else {
		resp, err = api.Get("/v1/identities/" + selfID + "/devices")
	}
	if err != nil {
		return nil, err
	}

	var devices []string
	err = json.Unmarshal(resp, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func Serialize(req map[string]interface{}, kid string, sk ed25519.PrivateKey) ([]byte, error) {
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

func ParseResponse(response []byte, history []json.RawMessage) (msg []byte, err error) {
	sg, err := siggraph.New(history)
	if err != nil {
		return msg, err
	}

	kid, err := kidhelper.GetJWSKID(response)
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
