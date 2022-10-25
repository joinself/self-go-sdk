package voice

import (
	"github.com/joinself/self-go-sdk/messaging"
)

// Sends a chat.voice.busy message finishing the call.
func (s *Service) Busy(recipient, cid, callID string) error {
	return s.send(recipient, map[string]interface{}{
		"typ":     "chat.voice.busy",
		"cid":     cid,
		"call_id": callID,
	})
}

// OnMessage subscribes to an incoming chat.voice.busy message
func (s *Service) OnBusy(callback func(iss, cid, callID string)) {
	s.messagingService.Subscribe("chat.voice.busy", func(m *messaging.Message) {
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
