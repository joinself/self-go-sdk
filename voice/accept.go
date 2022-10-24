package voice

import (
	"github.com/joinself/self-go-sdk/messaging"
)

// Sends a chat.voice.accept message accepting a specific call.
func (s *Service) Accept(recipient, cid, callID, peerInfo string, data map[string]interface{}) error {
	return s.send(recipient, map[string]interface{}{
		"typ":       "chat.voice.accept",
		"cid":       cid,
		"call_id":   callID,
		"peer_info": peerInfo,
		"data":      data,
	})
}

// OnMessage subscribes to an incoming chat.voice.accept message
func (s *Service) OnAccept(callback func(iss, cid, callID, peerInfo string, data interface{})) {
	s.messagingService.Subscribe("chat.voice.accept", func(m *messaging.Message) {
		payload, err := s.processMessage(m)
		if err == nil {
			callback(
				payload["iss"].(string),
				payload["cid"].(string),
				payload["callID"].(string),
				payload["peerInfo"].(string),
				payload["data"],
			)
		}
	})
}
