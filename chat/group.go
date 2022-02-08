// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

type Group struct {
	service *Service
	Payload map[string]interface{}
	GID     string
	Members []string
	Name    string
	Link    string
	Key     string
	Mime    string
}

func NewGroup(service *Service, payload map[string]interface{}) *Group {
	members := []string{}
	for _, v := range payload["members"].([]interface{}) {
		members = append(members, v.(string))
	}

	g := Group{
		service: service,
		Payload: payload,
		GID:     payload["gid"].(string),
		Members: members,
		Name:    payload["name"].(string),
	}
	if payload["link"] != nil {
		g.Link = payload["link"].(string)
	}
	if payload["key"] != nil {
		g.Key = payload["key"].(string)
	}
	if payload["mime"] != nil {
		g.Mime = payload["mime"].(string)
	}

	return &g
}

func (m *Group) Invite(user string) {
	if len(user) == 0 {
		return
	}

	m.Members = append(m.Members, user)
	m.service.Invite(m.GID, m.Name, m.Members)
}

func (m *Group) Leave() {
	m.service.Leave(m.GID, m.Members)
}

func (m *Group) Join() {
	m.service.Join(m.GID, m.Members)
}

func (m *Group) Message(body string, opts ...MessageOptions) (*Message, error) {
	if len(opts) == 0 {
		opts = append(opts, MessageOptions{})
	}
	opts[0].GID = m.GID
	return m.service.Message(m.Members, body, opts[0])
}
