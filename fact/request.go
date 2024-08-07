// Copyright 2020 Self Group Ltd. All Rights Reserved.

package fact

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/pkg/helpers"
	"github.com/joinself/self-go-sdk/pkg/ntp"
	"github.com/joinself/self-go-sdk/pkg/object"
	"github.com/joinself/self-go-sdk/pkg/siggraph"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/skip2/go-qrcode"
	"github.com/square/go-jose"
	"github.com/tidwall/gjson"
)

var (
	ErrBadJSONPayload               = errors.New("bad json payload")
	ErrResponseBadSignature         = errors.New("bad response signature")
	ErrRequestTimeout               = errors.New("request timeout")
	ErrMessageBadIssuer             = errors.New("bad response issuer")
	ErrMessageBadSubject            = errors.New("bad response subject")
	ErrMessageBadAudience           = errors.New("bad response audience")
	ErrMessageBadStatus             = errors.New("bad response status")
	ErrMessageExpired               = errors.New("response has expired")
	ErrMessageIssuedTooSoon         = errors.New("response was issued in the future")
	ErrStatusRejected               = errors.New("fact request was rejected")
	ErrStatusUnauthorized           = errors.New("you are not authorized to interact with this user")
	ErrFactRequestBadIdentity       = errors.New("fact request must specify a valid self id")
	ErrFactRequestBadFacts          = errors.New("fact request must specify one or more facts")
	ErrFactQRRequestBadConversation = errors.New("fact qr request must specify a valid conversation id")
	ErrFactQRRequestBadFacts        = errors.New("fact qr request must specify one or more facts")
	ErrFactResultMismatch           = errors.New("fact has differing attested values")
	ErrFactNotAttested              = errors.New("fact has attestations with empty or invalid values")
	ErrBadAttestationSubject        = errors.New("attestation is not related to the responder")
	ErrMissingConversationID        = errors.New("deep link request must specify a unique conversation id")
	ErrMissingCallback              = errors.New("deep link request must specify a callback url")
	ErrFactRequestCID               = errors.New("cid not provided")
	ErrSigningKeyInvalid            = errors.New("signing key was invalid at the time the attestation was issued")
	ErrNotConnected                 = errors.New("you're not permitting connections from the specifed recipient")
	ErrNotEnoughCredits             = errors.New("your credits have expired, please log in to the developer portal and top up your account")
	ErrEmptyFacts                   = errors.New("facts not provided")
	ErrEmptySource                  = errors.New("empty source provided")

	ServiceSelfIntermediary = "self_intermediary"
)

// FactRequest specifies the parameters of an information request
type FactRequest struct {
	SelfID      string
	Description string
	Facts       []Fact
	Expiry      time.Duration
	AllowedFor  time.Duration
	Callback    json.RawMessage
	Auth        bool
}

// FactRequestAsync specifies the parameters of an information requestAsync
type FactRequestAsync struct {
	SelfID      string
	Description string
	Facts       []Fact
	Expiry      time.Duration
	AllowedFor  time.Duration
	CID         string
	Callback    json.RawMessage
	Auth        bool
}

// FactResponse contains the details of the requested facts
type FactResponse struct {
	Status  string
	Facts   []Fact
	Objects map[string]*object.Object
}

// QRFactRequest contains the details of the requested facts
type QRFactRequest struct {
	ConversationID string
	Description    string
	Facts          []Fact
	Options        map[string]string
	Expiry         time.Duration
	QRConfig       QRConfig
	Auth           bool
}

// QRFactResponse contains the details of the requested facts
type QRFactResponse struct {
	Responder      string
	Facts          []Fact
	Options        map[string]string
	Accepted       bool
	ConversationID string
	DeviceID       string
}

// DeepLinkFactRequest contains the details of the requested facts
type DeepLinkFactRequest struct {
	ConversationID string
	Description    string
	Callback       string
	Facts          []Fact
	Expiry         time.Duration
	Auth           bool
}

type QRConfig struct {
	Size            int
	ForegroundColor string
	BackgroundColor string
}

// IntermediaryFactRequest specifies the paramters on an information request via an intermediary
type IntermediaryFactRequest struct {
	SelfID       string
	Description  string
	Intermediary string
	Facts        []Fact
	Expiry       time.Duration
}

// IntermediaryFactResponse contains the details of the requested facts
type IntermediaryFactResponse struct {
	Facts []Fact
}

type remoteFile interface {
	SetObject(data []byte) (*object.EncryptedObject, error)
	GetObject(link, key string) ([]byte, error)
}

type StandardResponse struct {
	ID             string                   `json:"jti"`
	Type           string                   `json:"typ"`
	Conversation   string                   `json:"cid"`
	Issuer         string                   `json:"iss"`
	Audience       string                   `json:"aud"`
	Subject        string                   `json:"sub"`
	IssuedAt       time.Time                `json:"iat"`
	ExpiresAt      time.Time                `json:"exp"`
	DeviceID       string                   `json:"device_id"`
	Status         string                   `json:"status"`
	Description    string                   `json:"description"`
	Facts          []Fact                   `json:"facts"`
	Objects        []map[string]interface{} `json:"objects,omitempty"`
	Auth           bool                     `json:"auth"`
	FileInteractor *object.RemoteFileInteractor
}

// Request requests a fact from a given identity
func (s Service) Request(req *FactRequest) (*FactResponse, error) {
	if req.SelfID == "" {
		return nil, ErrFactRequestBadIdentity
	}

	for _, fact := range req.Facts {
		err := fact.validate()
		if err != nil {
			return nil, err
		}
	}

	if req.Expiry == 0 {
		req.Expiry = defaultRequestTimeout
	}

	if !s.paidActions() {
		return nil, ErrNotEnoughCredits
	}

	cid := uuid.New().String()

	payload, err := s.factPayload(cid, req.SelfID, req.SelfID, req.Description, req.Facts, nil, req.Expiry, &req.AllowedFor, req.Auth, req.Callback)
	if err != nil {
		return nil, err
	}

	recipients, err := helpers.PrepareRecipients([]string{req.SelfID}, []string{s.selfID + ":" + s.deviceID}, s.api)
	if err != nil {
		return nil, err
	}

	responder, response, err := s.messaging.Request(recipients, cid, RequestInformation, payload, req.Expiry)
	if err != nil {
		return nil, err
	}

	selfID := strings.Split(responder, ":")[0]

	if selfID != req.SelfID {
		return nil, ErrMessageBadIssuer
	}

	resp, err := s.factResponse(selfID, selfID, response)
	if err != nil {
		return nil, err
	}

	objects := map[string]*object.Object{}
	for _, o := range resp.Objects {
		fo := object.New(s.fileInteractor)
		o["name"] = o["id"]
		if err := fo.BuildFromObject(o); err == nil {
			if _, ok := o["image_hash"]; ok {
				objects[o["image_hash"].(string)] = fo
			} else if _, ok := o["object_hash"]; ok {
				objects[o["object_hash"].(string)] = fo
			}
		}
	}

	return &FactResponse{Facts: resp.Facts, Objects: objects, Status: resp.Status}, nil
}

// RequestAsync requests a fact from a given identity and does not
// wait for the response
func (s Service) RequestAsync(req *FactRequestAsync) error {
	if req.SelfID == "" {
		return ErrFactRequestBadIdentity
	}

	if req.Expiry == 0 {
		req.Expiry = defaultRequestTimeout
	}

	if req.CID == "" {
		return ErrFactRequestCID
	}

	if !s.paidActions() {
		return ErrNotEnoughCredits
	}

	payload, err := s.factPayload(req.CID, req.SelfID, req.SelfID, req.Description, req.Facts, nil, req.Expiry, &req.AllowedFor, req.Auth, req.Callback)
	if err != nil {
		return err
	}

	recipients, err := helpers.PrepareRecipients([]string{req.SelfID}, []string{s.selfID + ":" + s.deviceID}, s.api)
	if err != nil {
		return err
	}

	return s.messaging.Send(recipients, RequestInformation, payload)
}

// RequestViaIntermediary requests a fact from a given identity via a trusted
// intermediary. The intermediary verifies that the identity has a given fact
// and that it meets the requested requirements.
func (s Service) RequestViaIntermediary(req *IntermediaryFactRequest) (*IntermediaryFactResponse, error) {
	if req.Expiry == 0 {
		req.Expiry = defaultRequestTimeout
	}

	if req.Intermediary == "" {
		req.Intermediary = ServiceSelfIntermediary
	}

	cid := uuid.New().String()

	payload, err := s.factPayload(cid, req.SelfID, req.Intermediary, req.Description, req.Facts, nil, req.Expiry, nil, false, nil)
	if err != nil {
		return nil, err
	}

	recipients, err := helpers.PrepareRecipients([]string{req.Intermediary}, []string{s.selfID + ":" + s.deviceID}, s.api)
	if err != nil {
		return nil, err
	}

	responder, response, err := s.messaging.Request(recipients, cid, RequestInformation, payload, req.Expiry)
	if err != nil {
		return nil, err
	}

	selfID := strings.Split(responder, ":")[0]

	if selfID != req.Intermediary {
		return nil, ErrMessageBadIssuer
	}

	resp, err := jose.ParseSigned(string(response))
	if err != nil {
		return nil, err
	}

	sub := gjson.GetBytes(resp.UnsafePayloadWithoutVerification(), "sub").String()

	if sub != req.SelfID {
		return nil, ErrMessageBadSubject
	}

	res, err := s.factResponse(req.Intermediary, req.SelfID, response)
	if err != nil {
		return nil, err
	}

	return &IntermediaryFactResponse{Facts: res.Facts}, nil
}

// GenerateQRCode generates a qr code containing an fact request
func (s Service) GenerateQRCode(req *QRFactRequest) ([]byte, error) {
	if req.ConversationID == "" {
		return nil, ErrFactQRRequestBadConversation
	}

	if req.Expiry == 0 {
		req.Expiry = defaultRequestTimeout
	}

	if req.QRConfig.ForegroundColor == "" {
		req.QRConfig.ForegroundColor = "#0E1C42"
	}

	if req.QRConfig.BackgroundColor == "" {
		req.QRConfig.BackgroundColor = "#FFFFFF"
	}

	if req.QRConfig.Size == 0 {
		req.QRConfig.Size = 400
	}

	payload, err := s.factPayload(req.ConversationID, "-", "-", req.Description, req.Facts, req.Options, req.Expiry, nil, req.Auth, nil)
	if err != nil {
		return nil, err
	}

	q, err := qrcode.New(string(payload), qrcode.Low)
	if err != nil {
		return nil, err
	}

	q.BackgroundColor, _ = colorful.Hex(req.QRConfig.BackgroundColor)
	q.ForegroundColor, _ = colorful.Hex(req.QRConfig.ForegroundColor)

	return q.PNG(req.QRConfig.Size)
}

// GenerateDeepLink generates a qr code containing an fact request
func (s Service) GenerateDeepLink(req *DeepLinkFactRequest) (string, error) {
	if req.ConversationID == "" {
		return "", ErrMissingConversationID
	}

	if req.Callback == "" {
		return "", ErrMissingCallback
	}

	payload, err := s.factPayload(req.ConversationID, "-", "-", req.Description, req.Facts, nil, req.Expiry, nil, req.Auth, nil)
	if err != nil {
		return "", err
	}

	body := base64.RawStdEncoding.EncodeToString(payload)
	baseURL := fmt.Sprintf("https://%s.links.joinself.com", s.environment)
	portalURL := fmt.Sprintf("https://developer.%s.joinself.com", s.environment)
	apn := fmt.Sprintf("com.joinself.app.%s", s.environment)

	if s.environment == "" || s.environment == "development" {
		baseURL = "https://links.joinself.com"
		portalURL = "https://developer.joinself.com"
		apn = "com.joinself.app"
		if s.environment == "development" {
			apn = "com.joinself.com.dev"
		}
	}
	return fmt.Sprintf("%s?link=%s/callback/%s%%3Fqr=%s&apn=%s", baseURL, portalURL, req.Callback, body, apn), nil
}

// WaitForResponse waits for completion of a fact request that was initiated by qr code
func (s Service) WaitForResponse(cid string, exp time.Duration) (*QRFactResponse, error) {
	responder, response, err := s.messaging.Wait(cid, exp)
	if err != nil {
		return nil, err
	}

	selfID := strings.Split(responder, ":")[0]

	resp, err := s.factResponse(selfID, selfID, response)
	if err != nil {
		return nil, err
	}

	return &QRFactResponse{
		Responder: responder,
		Facts:     resp.Facts,
		Accepted:  (resp.Status == "accepted"),
		DeviceID:  resp.DeviceID,
	}, nil
}

// Subscribe subscribes to fact request responses
func (s Service) Subscribe(auth bool, sub func(sender string, res *StandardResponse)) {
	s.messaging.Subscribe(ResponseInformation, func(sender string, payload []byte) {
		selfID := strings.Split(sender, ":")[0]

		resp, err := s.factResponse(selfID, selfID, payload)
		if err != nil {
			if !errors.Is(err, ErrStatusRejected) {
				log.Println("fact response error:", err.Error())
				return
			}
		}

		sub(selfID, resp)
	})
}

func (s *Service) factResponse(issuer, subject string, response []byte) (*StandardResponse, error) {
	history, err := s.pki.GetHistory(issuer)
	if err != nil {
		return nil, err
	}

	msg, err := helpers.ParseJWS(response, history)
	if err != nil {
		return nil, ErrResponseBadSignature
	}

	return s.parseFactResponse(issuer, subject, msg)
}

func (s *Service) parseFactResponse(issuer, subject string, response []byte) (*StandardResponse, error) {
	var resp StandardResponse

	err := json.Unmarshal(response, &resp)
	if err != nil {
		return nil, ErrBadJSONPayload
	}
	resp.FileInteractor = s.fileInteractor

	if resp.Audience != s.selfID {
		return nil, ErrMessageBadAudience
	}

	if resp.Issuer != issuer {
		return nil, ErrMessageBadIssuer
	}

	if ntp.After(resp.ExpiresAt) {
		return nil, ErrMessageExpired
	}

	if ntp.Before(resp.IssuedAt) {
		return nil, ErrMessageIssuedTooSoon
	}

	for i, f := range resp.Facts {
		resp.Facts[i].payloads = make([][]byte, len(f.Attestations))

		for x, adata := range f.Attestations {
			jws, err := jose.ParseSigned(string(adata))
			if err != nil {
				return nil, err
			}

			apayload := jws.UnsafePayloadWithoutVerification()

			iss := gjson.GetBytes(apayload, "iss").String()
			iatRFC3999 := gjson.GetBytes(apayload, "iat").String()

			history, err := s.pki.GetHistory(iss)
			if err != nil {
				return nil, err
			}

			sg, err := siggraph.New(history)
			if err != nil {
				return nil, err
			}

			kid, err := helpers.GetJWSKID(adata)
			if err != nil {
				return nil, err
			}

			iat, err := time.Parse(time.RFC3339, iatRFC3999)
			if err != nil {
				return nil, err
			}

			if !sg.IsKeyValid(kid, iat) {
				return nil, ErrSigningKeyInvalid
			}

			pk, err := sg.Key(kid)
			if err != nil {
				return nil, err
			}

			msg, err := jws.Verify(pk)
			if err != nil {
				return nil, ErrResponseBadSignature
			}

			sub := gjson.GetBytes(msg, "sub").String()

			if strings.Split(sub, ":")[0] != subject {
				return nil, ErrBadAttestationSubject
			}

			resp.Facts[i].payloads[x] = msg
		}
	}

	switch resp.Status {
	case StatusAccepted:
		return &resp, nil
	case StatusRejected:
		return &resp, ErrStatusRejected
	case StatusUnauthorized:
		return &resp, ErrStatusUnauthorized
	default:
		return &resp, ErrMessageBadStatus
	}
}

// FactResponse validate and process a fact response
func (s *Service) FactResponse(issuer, subject string, response []byte) ([]Fact, error) {
	resp, err := s.parseFactResponse(issuer, subject, response)
	if err != nil {
		return nil, err
	}

	return resp.Facts, nil
}

func (s *Service) factPayload(cid, selfID, intermediary, description string, facts []Fact, options map[string]string, exp time.Duration, au *time.Duration, auth bool, callback json.RawMessage) ([]byte, error) {
	if facts == nil {
		facts = make([]Fact, 0)
	}
	req := map[string]interface{}{
		"typ":         RequestInformation,
		"cid":         cid,
		"jti":         uuid.New().String(),
		"iss":         s.selfID,
		"sub":         selfID,
		"aud":         intermediary,
		"iat":         ntp.TimeFunc().Format(time.RFC3339),
		"exp":         ntp.TimeFunc().Add(exp).Format(time.RFC3339),
		"device_id":   s.deviceID,
		"description": description,
		"facts":       facts,
		"auth":        auth,
	}

	if au != nil {
		req["allowed_until"] = ntp.TimeFunc().Add(*au).Format(time.RFC3339)
	}

	if options != nil {
		req["options"] = options
	}
	if callback != nil {
		req["callback"] = callback
	}

	return helpers.PrepareJWS(req, s.keyID, s.sk)
}

// builds a list of all devices associated with an identity
func (s Service) paidActions() bool {
	var resp []byte
	var err error

	resp, err = s.api.Get("/v1/apps/" + s.selfID)
	if err != nil {
		return false
	}

	var app struct {
		PaidActions bool `json:"paid_actions"`
	}

	err = json.Unmarshal(resp, &app)
	if err != nil {
		return false
	}

	return app.PaidActions
}
