package voice

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/messaging"
	"github.com/joinself/self-go-sdk/pkg/helpers"
	"github.com/joinself/self-go-sdk/pkg/ntp"
)

// Message represents a chat message.
type Message struct {
	Recipients []string
	JTI        string
	ISS        string
	Payload    map[string]interface{}
}

func (s *Service) processMessage(m *messaging.Message) (map[string]interface{}, error) {
	println("message received from " + m.Sender)
	var payload map[string]interface{}
	err := json.Unmarshal(m.Payload, &payload)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return payload, nil
}

func (s *Service) send(recipient string, payload map[string]interface{}) error {
	recs, err := helpers.PrepareRecipients([]string{recipient}, []string{s.selfID + ":" + s.deviceID}, s.api)
	if err != nil {
		return err
	}

	typ := "chat.voice.start"
	body, err := s.buildPayload(recipient, payload)

	return s.messagingClient.Send(recs, typ, body)
}

func (s *Service) buildPayload(recipient string, input map[string]interface{}) ([]byte, error) {
	req := map[string]interface{}{
		"jti": uuid.New().String(),
		"cid": uuid.New().String(),
		"aud": recipient,
		"sub": recipient,
		"iss": s.selfID,
		"iat": ntp.TimeFunc().Format(time.RFC3339),
		"exp": ntp.TimeFunc().Add(s.expiry).Format(time.RFC3339),
	}
	for k, v := range input {
		req[k] = v
	}

	return helpers.PrepareJWS(req, s.keyID, s.sk)
}
