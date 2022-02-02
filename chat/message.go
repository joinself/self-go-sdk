// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

type Message struct {
	service    *Service
	Body       string
	Recipients []string
	JTI        string
	GID        string
	ISS        string
	Payload    map[string]interface{}
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
	return m.Message(body, map[string]interface{}{"rid": m.JTI})
}

func (m *Message) Message(body string, opts map[string]interface{}) *Message {
	if len(m.GID) > 0 {
		opts["aud"] = m.GID
		opts["gid"] = m.GID
	}

	to := m.Recipients
	if m.amITheRecipient() {
		to = []string{m.ISS}
	}

	return m.service.Message(to, body, opts)
}

func (m *Message) amITheRecipient() bool {
	return (len(m.Recipients) == 1 && m.Recipients[0] == m.service.SelfID())
}
