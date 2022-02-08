// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"log"
)

type Message struct {
	service    *Service
	Body       string
	Recipients []string
	JTI        string
	GID        string
	ISS        string
	Payload    map[string]interface{}
	Objects    []*Object
}

// TODO: objects

func NewMessage(chat *Service, recipients []string, payload map[string]interface{}) *Message {
	// TODO: process incomming objects
	m := Message{
		service:    chat,
		ISS:        payload["iss"].(string),
		Body:       payload["msg"].(string),
		Recipients: recipients,
		JTI:        payload["jti"].(string),
		Payload:    payload,
	}

	if payload["gid"] != nil {
		m.GID = payload["gid"].(string)
	}

	if payload["objects"] != nil {
		for _, oo := range payload["objects"].([]interface{}) {
			o := oo.(map[string]interface{})
			if _, ok := o["key"]; ok {
				obj := NewObject(chat.FileInteractor)
				err := obj.BuildFromObject(o)
				if err != nil {
					log.Println(err)
					continue
				}

				m.Objects = append(m.Objects, obj)
			} else {
				// TODO: implement public object here
				println("TODO received a public object " + o["link"].(string))
			}
		}
	}

	return &m
}

func (m *Message) Delete() {
	m.service.Delete(m.Recipients, []string{m.JTI}, m.GID)
}

func (m *Message) Edit(Body string) {
	if m.amITheRecipient() {
		return
	}

	m.Body = Body
	m.service.Edit(m.Recipients, m.JTI, m.Body, m.GID)
}

func (m *Message) MarkAsDelivered() {
	if !m.amITheRecipient() {
		return
	}

	m.service.Delivered([]string{m.ISS}, []string{m.JTI}, m.GID)
}

func (m *Message) MarkAsRead() {
	if !m.amITheRecipient() {
		return
	}

	m.service.Read([]string{m.ISS}, []string{m.JTI}, m.GID)
}

func (m *Message) Respond(body string) *Message {
	return m.Message(body, MessageOptions{RID: m.JTI})
}

func (m *Message) Message(body string, opts ...MessageOptions) *Message {
	if len(opts) == 0 {
		opts = append(opts, MessageOptions{})
	}

	if len(m.GID) > 0 {
		opts[0].AUD = m.GID
		opts[0].GID = m.GID
	}

	to := m.Recipients
	if m.amITheRecipient() {
		to = []string{m.ISS}
	}

	return m.service.Message(to, body, opts[0])
}

func (m *Message) amITheRecipient() bool {
	return (len(m.Recipients) == 1 && m.Recipients[0] == m.service.SelfID())
}
