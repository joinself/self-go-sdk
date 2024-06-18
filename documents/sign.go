// Copyright 2020 Self Group Ltd. All Rights Reserved.

package documents

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/pkg/helpers"
	"github.com/joinself/self-go-sdk/pkg/ntp"
	"github.com/joinself/self-go-sdk/pkg/object"
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

// RequestSignature sends a signature request to the specified user.
func (s *Service) RequestSignature(recipient string, body string, objects []InputObject) (Response, error) {
	var resp Response
	jti := uuid.New().String()
	oo := make([]map[string]interface{}, 0)
	for _, o := range objects {
		fo := object.New(s.fileInteractor)
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

	payload, err := helpers.PrepareJWS(req, s.keyID, s.sk)
	if err != nil {
		return resp, err
	}

	recs, err := helpers.PrepareRecipients([]string{recipient}, []string{s.selfID + ":" + s.deviceID}, s.api)
	if err != nil {
		return resp, err
	}

	issuer, response, err := s.messaging.Request(recs, jti, req["typ"].(string), payload, 0)
	if err != nil {
		return resp, err
	}

	resp, err = s.response(issuer, response)
	if err != nil {
		log.Println(err.Error())
		return resp, err
	}

	return resp, nil
}

func (s *Service) response(issuer string, response []byte) (resp Response, err error) {
	selfID := strings.Split(issuer, ":")[0]

	history, err := s.pki.GetHistory(selfID)
	if err != nil {
		return resp, err
	}

	msg, err := helpers.ParseJWS(response, history)
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
