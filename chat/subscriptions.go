// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"encoding/json"
	"log"

	"github.com/joinself/self-go-sdk/messaging"
)

type OnMessageOptions struct {
	MarkAsDelivered bool
	MarkAsRead      bool
}

// OnMessage subscribes to an incoming chat message.
func (s *Service) OnMessage(callback func(cm *Message), opts ...OnMessageOptions) {
	s.messagingService.Subscribe("chat.message", func(m *messaging.Message) {
		nm, err := s.processChatMessage(m, opts...)
		if err == nil {
			callback(nm)
		}
	})
}

func (s *Service) processChatMessage(m *messaging.Message, opts ...OnMessageOptions) (*Message, error) {
	println("message received from " + m.Sender)
	var payload map[string]interface{}
	err := json.Unmarshal(m.Payload, &payload)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	nm := NewMessage(s, []string{payload["aud"].(string)}, payload)

	if len(opts) > 0 {
		if opts[0].MarkAsRead == true {
			nm.MarkAsRead()
		}
		if opts[0].MarkAsDelivered != false {
			nm.MarkAsDelivered()
		}
	} else {
		nm.MarkAsDelivered()
	}

	return nm, err
}

// OnInvite subscribes to group invitations.
func (s *Service) OnInvite(callback func(m *Group)) {
	s.messagingService.Subscribe("chat.invite", func(m *messaging.Message) {
		println("invited to a group by " + m.Sender)

		g, err := s.processChatInvite(m)
		if err == nil {
			callback(g)
		}
	})
}

func (s *Service) processChatInvite(m *messaging.Message) (*Group, error) {
	var payload map[string]interface{}
	err := json.Unmarshal(m.Payload, &payload)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return NewGroup(s, payload), nil
}

// OnJoin subscribes to people joining a group
func (s *Service) OnJoin(callback func(iss, gid string)) {
	s.messagingService.Subscribe("chat.join", func(m *messaging.Message) {
		println(m.Sender + " joined a group you're in")

		gid, err := s.getMessageGID(m)
		if err == nil {
			callback(m.Sender, gid)
		}
	})
}

// OnLeave subscribes to people leaving the specified group.
func (s *Service) OnLeave(callback func(iss, gid string)) {
	s.messagingService.Subscribe("chat.remove", func(m *messaging.Message) {
		println(m.Sender + " left a group you're in")
		gid, err := s.getMessageGID(m)
		if err == nil {
			callback(m.Sender, gid)
		}
	})
}

// OnConnection subscribes to people connecting to your app.
func (s *Service) OnConnection(callback func(iss, status string)) {
	s.messagingService.Subscribe("identities.connections.resp", func(m *messaging.Message) {
		println(m.Sender + " connected to your app")

		var payload map[string]interface{}
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			log.Println(err)
			callback(m.Sender, "errored")
		}

		if err == nil {
			callback(m.Sender, payload["status"].(string))
		}
	})
}

func (s *Service) getMessageGID(m *messaging.Message) (string, error) {
	var payload map[string]interface{}
	err := json.Unmarshal(m.Payload, &payload)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return payload["gid"].(string), nil

}
