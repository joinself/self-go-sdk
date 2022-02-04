// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"encoding/json"
	"log"

	"github.com/joinself/self-go-sdk/messaging"
)

// OnMessage subscribes to an incoming chat message.
func (s *Service) OnMessage(callback func(cm *Message)) {
	s.messagingService.Subscribe("chat.message", func(m *messaging.Message) {
		println("message received from " + m.Sender)
		var payload map[string]interface{}
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			log.Println(err)
			return
		}

		nm := NewMessage(s, []string{payload["aud"].(string)}, payload)
		nm.MarkAsDelivered()
		nm.MarkAsRead()
		callback(nm)
	})
}

// OnInvite subscribes to group invitations.
func (s *Service) OnInvite(callback func(m *Group)) {
	s.messagingService.Subscribe("chat.invite", func(m *messaging.Message) {
		println("invited to a group by " + m.Sender)

		var payload map[string]interface{}
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			log.Println(err)
			return
		}

		callback(NewGroup(s, payload))
	})
}

// OnJoin subscribes to people joining a group
func (s *Service) OnJoin(callback func(iss, gid string)) {
	s.messagingService.Subscribe("chat.join", func(m *messaging.Message) {
		println(m.Sender + " joined a group you're in")

		var payload map[string]interface{}
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			log.Println(err)
			return
		}

		callback(m.Sender, payload["gid"].(string))
	})
}

// OnLeave subscribes to people leaving the specified group.
func (s *Service) OnLeave(callback func(iss, gid string)) {
	s.messagingService.Subscribe("chat.remove", func(m *messaging.Message) {
		println(m.Sender + " left a group you're in")
		var payload map[string]interface{}
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			log.Println(err)
			return
		}

		callback(m.Sender, payload["gid"].(string))
	})
}
