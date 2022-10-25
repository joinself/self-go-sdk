package voice

import (
	"github.com/joinself/self-go-sdk/messaging"
)

// Sends a chat.voice.accept message finishing the call.
func (s *Service) Stop(recipient, cid, callID string) error {
	return s.send(recipient, map[string]interface{}{
		"typ":     "chat.voice.stop",
		"cid":     cid,
		"call_id": callID,
	})
}

// OnMessage subscribes to an incoming chat.voice.stop message
func (s *Service) OnStop(callback func(iss, cid, callID string)) {
	s.messagingService.Subscribe("chat.voice.stop", func(m *messaging.Message) {
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
