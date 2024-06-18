// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"log"

	"github.com/joinself/self-go-sdk/pkg/object"
)

// Message represents a chat message.
type Message struct {
	service    *Service
	Body       string
	Recipients []string
	JTI        string
	GID        string
	ISS        string
	Payload    map[string]interface{}
	Objects    []*object.Object
}

// NewMessage creates a chat message object.
func NewMessage(chat *Service, recipients []string, payload map[string]interface{}) *Message {
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
			obj := object.New(chat.fileInteractor)
			if _, ok := o["key"]; ok {
				err := obj.BuildFromObject(o)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				// TODO: implement public object here
				obj.Link = o["link"].(string)
				obj.Name = o["name"].(string)
				obj.Mime = o["mime"].(string)
			}
			m.Objects = append(m.Objects, obj)
		}
	}

	return &m
}

// Sends deletes the current message.
func (m *Message) Delete() {
	m.service.Delete(m.Recipients, []string{m.JTI}, m.GID)
}

// Edit edit the current message.
func (m *Message) Edit(Body string) {
	if m.amITheRecipient() {
		return
	}

	m.Body = Body
	m.service.Edit(m.Recipients, m.JTI, m.Body, m.GID)
}

// MarkAsDelivered marks the current message as delivered.
func (m *Message) MarkAsDelivered() {
	if !m.amITheRecipient() {
		return
	}

	m.service.Delivered([]string{m.ISS}, []string{m.JTI}, m.GID)
}

// MarkAsRead marks the current message as read.
func (m *Message) MarkAsRead() {
	if !m.amITheRecipient() {
		return
	}

	m.service.Read([]string{m.ISS}, []string{m.JTI}, m.GID)
}

// Respond sends a response to the current message.
func (m *Message) Respond(body string) (*Message, error) {
	return m.Message(body, MessageOptions{RID: m.JTI})
}

// Message sends a message to the current conversation.
func (m *Message) Message(body string, opts ...MessageOptions) (*Message, error) {
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
