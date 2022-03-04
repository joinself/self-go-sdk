// Copyright 2020 Self Group Ltd. All Rights Reserved.

package documents

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/chat"
	"github.com/joinself/self-go-sdk/pkg/ntp"
	"github.com/joinself/self-go-sdk/pkg/siggraph"
	"github.com/square/go-jose"
)

var (
	ErrBadJSONPayload       = errors.New("bad json payload")
	ErrResponseBadSignature = errors.New("bad response signature")
)

type InputObject struct {
	Name string
	Data []byte
	Mime string
}

type SignedObject struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Hash string `json:"hash"`
}

type Response struct {
	Signature     string
	SignedObjects []SignedObject `json:"signed_objects"`
	Status        string         `json:"status"`
}

func (s *Service) RequestSignature(recipient string, body string, objects []InputObject) (Response, error) {
	var resp Response
	jti := uuid.New().String()
	oo := make([]map[string]interface{}, 0)
	for _, o := range objects {
		fo := chat.NewObject(s.fileInteractor)
		err := fo.BuildFromData(o.Data, o.Name, o.Mime)
		if err != nil {
			return resp, err
		}
		oo = append(oo, fo.ToPayload())
	}

	req := map[string]interface{}{
		"jti":     jti,
		"cid":     uuid.New().String(),
		"typ":     "document.sign.req",
		"aud":     recipient,
		"sub":     recipient,
		"msg":     body,
		"objects": oo,
		"iss":     s.selfID,
		"iat":     ntp.TimeFunc().Format(time.RFC3339),
		"exp":     ntp.TimeFunc().Add(s.expiry).Format(time.RFC3339),
	}

	payload, err := s.serialize(req)
	if err != nil {
		return resp, err
	}

	recs, err := s.recipients([]string{recipient})
	if err != nil {
		return resp, err
	}

	log.Println("sending request")
	issuer, response, err := s.messaging.Request(recs, jti, payload, 0)
	if err != nil {
		return resp, err
	}
	log.Println("response received")

	resp, err = s.response(issuer, response)
	if err != nil {
		log.Println(err.Error())
		return resp, err
	}

	return resp, nil
}

func (s *Service) serialize(req map[string]interface{}) ([]byte, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return []byte(""), err
	}

	opts := &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"kid": s.keyID,
		},
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: s.sk}, opts)
	if err != nil {
		return []byte(""), err
	}

	signature, err := signer.Sign(payload)
	if err != nil {
		return []byte(""), err
	}

	return []byte(signature.FullSerialize()), nil
}

// builds a list of all devices associated with an identity
func (s Service) recipients(recipients []string) ([]string, error) {
	devices := make([]string, 0)
	for _, selfID := range recipients {
		dds, err := s.getDevices(selfID)
		if err != nil {
			return nil, err
		}

		for i := range dds {
			if selfID != s.selfID && dds[i] != s.deviceID {
				devices = append(devices, selfID+":"+dds[i])
			}
		}
	}

	return devices, nil
}

func (s Service) getDevices(selfID string) ([]string, error) {
	var resp []byte
	var err error

	if len(selfID) > 11 {
		resp, err = s.api.Get("/v1/apps/" + selfID + "/devices")
	} else {
		resp, err = s.api.Get("/v1/identities/" + selfID + "/devices")
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

func (s *Service) response(issuer string, response []byte) (resp Response, err error) {
	selfID := strings.Split(issuer, ":")[0]

	history, err := s.pki.GetHistory(selfID)
	if err != nil {
		return resp, err
	}

	sg, err := siggraph.New(history)
	if err != nil {
		return resp, err
	}

	kid, err := getJWSKID(response)
	if err != nil {
		return resp, err
	}

	pk, err := sg.ActiveKey(kid)
	if err != nil {
		return resp, err
	}

	jws, err := jose.ParseSigned(string(response))
	if err != nil {
		return resp, err
	}

	msg, err := jws.Verify(pk)
	if err != nil {
		return resp, ErrResponseBadSignature
	}

	err = json.Unmarshal(msg, &resp)
	if err != nil {
		return resp, ErrBadJSONPayload
	}
	resp.Signature = string(response)

	return resp, nil
}

func getJWSKID(payload []byte) (string, error) {
	var jws struct {
		Protected string `json:"protected"`
	}

	err := json.Unmarshal(payload, &jws)
	if err != nil {
		return "", err
	}

	return getKID(jws.Protected)
}

func getKID(token string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(strings.Split(token, ".")[0])
	if err != nil {
		return "", err
	}

	hdr := make(map[string]string)

	err = json.Unmarshal(data, &hdr)
	if err != nil {
		return "", err
	}

	kid := hdr["kid"]
	if kid == "" {
		return "", errors.New("token must specify an identifier for the signing key")
	}

	return kid, nil
}
