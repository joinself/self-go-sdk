package voice

import (
	"github.com/joinself/self-go-sdk/messaging"
)

// Sends a chat.voice.summary message Sending details about the call.
func (s *Service) Summary(recipient, cid, callID string) error {
	return s.send(recipient, map[string]interface{}{
		"typ":     "chat.voice.summary",
		"cid":     cid,
		"call_id": callID,
	})
}

// OnMessage subscribes to an incoming chat.voice.summary message
func (s *Service) OnSummary(callback func(iss, cid, callID string)) {
	s.messagingService.Subscribe("chat.voice.summary", func(m *messaging.Message) {
		payload, err := s.processMessage(m)
		if err == nil {
			callback(
				payload["iss"].(string),
				payload["cid"].(string),
				payload["call_id"].(string),
			)
		}
	})
}
