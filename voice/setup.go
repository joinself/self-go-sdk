package voice

import (
	"github.com/joinself/self-go-sdk/messaging"
)

func (s *Service) Setup(recipient, name, cid string) error {
	return s.send(recipient, map[string]interface{}{
		"typ":  "chat.voice.setup",
		"cid":  cid,
		"data": map[string]string{"name": name},
	})
}

// OnMessage subscribes to an incoming chat.voice.setup message
func (s *Service) OnSetup(callback func(iss, cid string, data interface{})) {
	s.messagingService.Subscribe("chat.voice.setup", func(m *messaging.Message) {
		payload, err := s.processMessage(m)
		if err == nil {
			callback(
				payload["iss"].(string),
				payload["cid"].(string),
				payload["data"],
			)
		}
	})
}
