package documents

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/joinself/self-go-sdk/documents/mocks"
	"github.com/joinself/self-go-sdk/pkg/ntp"
	"github.com/joinself/self-go-sdk/pkg/siggraph"
	"github.com/square/go-jose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSign(t *testing.T) {
	s, config := setup(t)

	recipient := "1112223334"
	body := "testing"
	objects := []InputObject{}

	// Mock API interactions
	expectedDeviceURLs := []string{"/v1/identities/" + recipient + "/devices"}
	restMock := mocks.RestTransport{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Mock messaging interactions
	var duration time.Duration = 0
	msgMock := mocks.MessagingClient{}
	msgMock.On("Request", []string{recipient + ":device1"}, mock.AnythingOfType("string"), mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "document.sign.req")
		assert.Equal(t, payload["aud"], recipient)
		assert.Equal(t, payload["iss"], config.SelfID)
		assert.Equal(t, payload["msg"], body)

		return true
	}), duration).Return(recipient, []byte("hello"), nil)
	s.messaging = msgMock

	// Mock pki interactions
	pki := mocks.PkiClient{}
	pki.On("GetHistory", recipient).Return(getHistory(config), nil)
	s.pki = &pki

	s.RequestSignature(recipient, body, objects)
}

func getHistory(config *Config) []json.RawMessage {
	now := ntp.TimeFunc().Add(-(time.Hour * 356 * 24)).Unix()

	return []json.RawMessage{
		testop(config.PrivateKey, "1", &siggraph.Operation{
			Sequence:  0,
			Version:   "1.0.0",
			Previous:  "-",
			Timestamp: now,
			Actions: []siggraph.Action{
				{
					KID:           "1",
					DID:           "1",
					Type:          siggraph.TypeDeviceKey,
					Action:        siggraph.ActionKeyAdd,
					EffectiveFrom: now,
					Key:           base64.RawURLEncoding.EncodeToString(config.PrivateKey),
				},
			},
		}),
	}
}
func testop(sk ed25519.PrivateKey, kid string, op *siggraph.Operation) json.RawMessage {
	data, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}

	opts := &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"kid": kid,
		},
	}

	s, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: sk}, opts)
	if err != nil {
		panic(err)
	}

	jws, err := s.Sign(data)
	if err != nil {
		panic(err)
	}

	return json.RawMessage(jws.FullSerialize())
}
