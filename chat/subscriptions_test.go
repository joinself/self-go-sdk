package chat

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/joinself/self-go-sdk/chat/mocks"
	"github.com/joinself/self-go-sdk/messaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/ed25519"
)

func setup(t *testing.T) (*Service, *Config) {
	pk := "1:56qJGhYCJmTHsYChCp3sPSjmiGlN2yG0KakYDquMAD0"
	kp := strings.Split(pk, ":")

	decoder := base64.RawStdEncoding

	skData, err := decoder.DecodeString(kp[1])
	assert.Nil(t, err)

	config := Config{
		SelfID:     "c4f81d86-9dac-40fd-9830-13c66a0b2345",
		KeyID:      kp[0],
		PrivateKey: ed25519.NewKeyFromSeed(skData),
	}
	s := NewService(config)

	return s, &config
}

func TestProcessChatMessage(t *testing.T) {
	s, config := setup(t)

	expectedOutgoingMessages := []string{"chat.message.delivered", "chat.message.read"}

	// Mock messaging interactions
	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"iss:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], expectedOutgoingMessages[0])
		expectedOutgoingMessages = expectedOutgoingMessages[1:]
		assert.Equal(t, payload["aud"], "iss")
		assert.Equal(t, payload["iss"], "c4f81d86-9dac-40fd-9830-13c66a0b2345")

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", "/v1/identities/iss/devices").Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	payload, err := json.Marshal(map[string]interface{}{
		"iss": "iss",
		"msg": "hi",
		"aud": config.SelfID,
		"jti": "jti",
	})
	assert.Nil(t, err)

	nm, err := s.processChatMessage(&messaging.Message{
		Sender:  "bob",
		Payload: payload,
	})

	assert.Equal(t, "iss", nm.ISS)
	assert.Equal(t, "hi", nm.Body)
	assert.Equal(t, []string{config.SelfID}, nm.Recipients)
	assert.Equal(t, "jti", nm.JTI)
	assert.Equal(t, "", nm.GID)
	assert.Nil(t, err)
}

func TestProcessChatInvite(t *testing.T) {
	members := []string{"a", "b", "c"}
	s, config := setup(t)

	payload, err := json.Marshal(map[string]interface{}{
		"iss":     "iss",
		"aud":     config.SelfID,
		"jti":     "jti",
		"gid":     "gid",
		"name":    "group name",
		"members": members,
	})
	assert.Nil(t, err)

	g, err := s.processChatInvite(&messaging.Message{
		Sender:  "bob",
		Payload: payload,
	})

	assert.Nil(t, err)
	assert.Equal(t, g.GID, "gid")
	assert.Equal(t, g.Members, members)
	assert.Equal(t, g.Name, "group name")
}

func TestGetMessageGID(t *testing.T) {
	members := []string{"a", "b", "c"}
	s, config := setup(t)

	payload, err := json.Marshal(map[string]interface{}{
		"iss":     "iss",
		"aud":     config.SelfID,
		"jti":     "jti",
		"gid":     "gid",
		"name":    "group name",
		"members": members,
	})
	assert.Nil(t, err)

	gid, err := s.getMessageGID(&messaging.Message{
		Sender:  "bob",
		Payload: payload,
	})

	assert.Nil(t, err)
	assert.Equal(t, gid, "gid")
}
