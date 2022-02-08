package chat

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/joinself/self-go-sdk/chat/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDelivered(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", "c"}
	cids := []string{"cid1", "cid2"}
	gid := "gid"

	expectedDeviceURLs := []string{"/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/identities/c/devices"}
	expectedRecipients := recipients

	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"a:device1", "b:device1", "c:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "chat.message.delivered")
		assert.Equal(t, payload["aud"], expectedRecipients[0])
		expectedRecipients = expectedRecipients[1:]
		assert.Equal(t, payload["iss"], config.SelfID)

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Delivered(recipients, cids, gid)
}

func TestRead(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", "c"}
	cids := []string{"cid1", "cid2"}
	gid := "gid"

	expectedDeviceURLs := []string{"/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/identities/c/devices"}
	expectedRecipients := recipients

	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"a:device1", "b:device1", "c:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "chat.message.read")
		assert.Equal(t, payload["aud"], expectedRecipients[0])
		expectedRecipients = expectedRecipients[1:]
		assert.Equal(t, payload["iss"], config.SelfID)

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Read(recipients, cids, gid)
}

func TestEdit(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", "c"}
	cid := "cid"
	gid := "gid"
	newBody := "hello!"

	expectedDeviceURLs := []string{"/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/identities/c/devices"}
	expectedRecipients := recipients

	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"a:device1", "b:device1", "c:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "chat.message.edit")
		assert.Equal(t, payload["aud"], expectedRecipients[0])
		expectedRecipients = expectedRecipients[1:]
		assert.Equal(t, payload["iss"], config.SelfID)
		assert.Equal(t, payload["msg"], newBody)
		assert.Equal(t, payload["cid"], "cid")
		assert.Equal(t, payload["gid"], "gid")

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Edit(recipients, cid, newBody, gid)
}

func TestDelete(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", "c"}
	cids := []string{"cid1", "cid2"}
	gid := "gid"

	expectedDeviceURLs := []string{"/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/identities/c/devices"}
	expectedRecipients := recipients

	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"a:device1", "b:device1", "c:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "chat.message.delete")
		assert.Equal(t, payload["aud"], expectedRecipients[0])
		expectedRecipients = expectedRecipients[1:]
		assert.Equal(t, payload["iss"], config.SelfID)
		assert.Equal(t, 2, len(payload["cids"].([]interface{})))
		assert.Equal(t, payload["gid"], "gid")

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Delete(recipients, cids, gid)
}

func TestInviite(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", "c"}
	name := "name"
	gid := "gid"

	expectedDeviceURLs := []string{"/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/identities/c/devices"}
	expectedRecipients := recipients

	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"a:device1", "b:device1", "c:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "chat.invite")
		assert.Equal(t, payload["aud"], expectedRecipients[0])
		expectedRecipients = expectedRecipients[1:]
		assert.Equal(t, payload["iss"], config.SelfID)
		assert.Equal(t, 3, len(payload["members"].([]interface{})))
		assert.Equal(t, payload["gid"], "gid")

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Invite(gid, name, recipients, map[string]interface{}{})
}

func TestJoin(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", config.SelfID, "c"}
	gid := "gid"

	expectedDeviceURLs := []string{"/v1/identities/c/devices", "/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/apps/" + config.SelfID + "/devices", "/v1/identities/c/devices"}
	expectedRecipients := []string{"a", "b", "c"}
	expectedPermissions := expectedRecipients
	firstMessage := true

	// MessagingClient mock
	msgMock := mocks.MessagingClientMock{}

	msgMock.On("Send", mock.AnythingOfType("[]string"), mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		if firstMessage == true {
			// A sessions.create message will be sent to the users who had joined
			// the group after the current identity.
			assert.Equal(t, "sessions.create", payload["typ"])
			assert.Equal(t, "c", payload["aud"])
			firstMessage = false
		} else {
			assert.Equal(t, payload["typ"], "chat.join")
			assert.Equal(t, payload["aud"], expectedRecipients[0])
			expectedRecipients = expectedRecipients[1:]
			assert.Equal(t, payload["iss"], config.SelfID)
			assert.Equal(t, payload["gid"], "gid")
		}
		return true
	})).Return(nil)

	s.messagingClient = msgMock

	// MessagingService mock
	msgServiceMock := mocks.MessagingServiceMock{}
	msgServiceMock.On("PermitConnection", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedPermissions[0], input)
		expectedPermissions = expectedPermissions[1:]
		return true
	})).Return(nil)
	s.messagingService = msgServiceMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Join(gid, recipients)
}

func TestLeave(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", "c"}
	gid := "gid"

	expectedDeviceURLs := []string{"/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/identities/c/devices"}
	expectedRecipients := recipients

	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"a:device1", "b:device1", "c:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "chat.remove")
		assert.Equal(t, payload["aud"], expectedRecipients[0])
		expectedRecipients = expectedRecipients[1:]
		assert.Equal(t, payload["iss"], config.SelfID)
		assert.Equal(t, payload["gid"], "gid")

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Leave(gid, recipients)
}

func TestMessage(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", "c"}
	body := "hi"

	expectedDeviceURLs := []string{"/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/identities/c/devices"}
	expectedRecipients := recipients

	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"a:device1", "b:device1", "c:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "chat.message")
		assert.Equal(t, payload["aud"], expectedRecipients[0])
		expectedRecipients = expectedRecipients[1:]
		assert.Equal(t, payload["iss"], config.SelfID)
		assert.Equal(t, payload["msg"], body)

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Message(recipients, body, map[string]interface{}{})
}

func TestMessageWithOpts(t *testing.T) {
	s, config := setup(t)

	recipients := []string{"a", "b", "c"}
	body := "hi"

	expectedDeviceURLs := []string{"/v1/identities/a/devices", "/v1/identities/b/devices", "/v1/identities/c/devices"}
	expectedRecipients := recipients

	msgMock := mocks.MessagingClientMock{}
	msgMock.On("Send", []string{"a:device1", "b:device1", "c:device1"}, mock.MatchedBy(func(input []byte) bool {
		var envelope map[string]string
		assert.Nil(t, json.Unmarshal(input, &envelope))
		decoder := base64.RawStdEncoding
		payloadStr, err := decoder.DecodeString(envelope["payload"])
		assert.Nil(t, err)
		var payload map[string]interface{}
		assert.Nil(t, json.Unmarshal(payloadStr, &payload))

		assert.Equal(t, payload["typ"], "chat.message")
		assert.Equal(t, payload["aud"], expectedRecipients[0])
		expectedRecipients = expectedRecipients[1:]
		assert.Equal(t, payload["iss"], config.SelfID)
		assert.Equal(t, payload["msg"], body)
		assert.Equal(t, payload["gid"], "gid")
		assert.Equal(t, payload["rid"], "rid")

		return true
	})).Return(nil)
	s.messagingClient = msgMock

	// Mock API interactions
	restMock := mocks.RestTransportMock{}
	restMock.On("Get", mock.MatchedBy(func(input string) bool {
		assert.Equal(t, expectedDeviceURLs[0], input)
		expectedDeviceURLs = expectedDeviceURLs[1:]
		return true
	})).Return([]byte(`["device1"]`), nil)
	s.api = &restMock

	// Call public method
	s.Message(recipients, body, map[string]interface{}{
		"gid": "gid",
		"rid": "rid",
	})
}
