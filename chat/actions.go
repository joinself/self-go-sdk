// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/pkg/ntp"
	"github.com/joinself/self-go-sdk/pkg/request"
)

func (s *Service) SelfID() string {
	return s.selfID
}

type MessageObject struct {
	Link string
	Name string
	Data []byte
	Mime string
}
type MessageOptions struct {
	GID     string
	AUD     string
	JTI     string
	RID     string
	Objects []MessageObject
}

// Message sends a message to a list of recipients.
func (s *Service) Message(recipients []string, body string, opts ...MessageOptions) (*Message, error) {
	payload := map[string]interface{}{
		"typ": "chat.message",
		"msg": body,
	}

	payload["aud"] = recipients[0]
	if len(opts) > 0 {
		if len(opts[0].GID) > 0 {
			payload["gid"] = opts[0].GID
			payload["aud"] = opts[0].GID
		}
		if len(opts[0].JTI) > 0 {
			payload["jti"] = opts[0].JTI
		}
		if len(opts[0].RID) > 0 {
			payload["rid"] = opts[0].RID
		}

		objects := make([]interface{}, 0)
		for _, o := range opts[0].Objects {
			if len(o.Link) > 0 {
				// Is a public image, just append it
				objects = append(objects, map[string]interface{}{
					"link": o.Link,
					"name": o.Name,
					"mime": o.Mime,
				})
			} else {
				fo := NewObject(s.fileInteractor)
				err := fo.BuildFromData(o.Data, o.Name, o.Mime)
				if err != nil {
					return nil, err
				}
				objects = append(objects, fo.ToPayload())
			}
		}
		payload["objects"] = objects
	}

	err := s.send(recipients, payload)
	if err != nil {
		return nil, err
	}

	return NewMessage(s, recipients, payload), nil
}

// Delivered sends a message to confirm a list of messages (identified by it's cids)
// have been delivered.
func (s *Service) Delivered(recipients []string, cids []string, gid string) {
	s.confirm("delivered", recipients, cids, gid)
}

// Read sends a message to confirm a list of messages (identified by it's cids)
// have been read.
func (s *Service) Read(recipients []string, cids []string, gid string) {
	s.confirm("read", recipients, cids, gid)
}

// Edit changes the body of a previous message.
func (s *Service) Edit(recipients []string, cid string, body string, gid string) {
	p := map[string]interface{}{
		"typ": "chat.message.edit",
		"cid": cid,
		"msg": body,
	}

	if len(gid) > 0 {
		p["gid"] = gid
	}

	s.send(recipients, p)
}

// Delete deletes previous messages.
func (s *Service) Delete(recipients []string, cids []string, gid string) {
	p := map[string]interface{}{
		"typ":  "chat.message.delete",
		"cids": cids,
	}

	if gid != "" {
		p["gid"] = gid
	}

	s.send(recipients, p)
}

type InviteOptions struct {
	Data []byte
	Mime string
}

// Invite sends a group invitation to a list of members.
func (s *Service) Invite(gid string, name string, members []string, opts ...InviteOptions) error {
	p := map[string]interface{}{
		"typ":     "chat.invite",
		"gid":     gid,
		"name":    name,
		"members": members,
		"aud":     gid,
	}

	if len(opts) > 0 {
		fo := NewObject(s.fileInteractor)
		err := fo.BuildFromData(opts[0].Data, "", opts[0].Mime)
		if err != nil {
			return err
		}
		objPayload := fo.ToPayload()
		p["link"] = objPayload["link"]
		p["key"] = objPayload["key"]
		p["expires"] = objPayload["expires"]
	}

	return s.send(members, p)
}

// Join joins a group.
func (s *Service) Join(gid string, members []string) {
	// Allow incoming connections from the given members.
	for _, m := range members {
		if m != s.selfID {
			s.messagingService.PermitConnection(m)
		}
	}

	// Create missing sessions with group members.
	s.createMissingSessions(members)

	// Send joining confirmation.
	s.send(members, map[string]interface{}{
		"typ": "chat.join", "gid": gid, "aud": gid,
	})
}

// Leave leaves a group.
func (s *Service) Leave(gid string, members []string) {
	s.send(members, map[string]interface{}{
		"typ": "chat.remove",
		"gid": gid,
	})
}

func (s *Service) confirm(action string, recipients []string, cids []string, gid string) error {
	req := map[string]interface{}{
		"cids": cids,
		"typ":  "chat.message." + action,
	}
	if gid != "" {
		req["gid"] = gid
	}

	return s.send(recipients, req)
}

func (s *Service) send(recipients []string, req map[string]interface{}) error {
	recs, err := s.requestHelper.FormatRecipients(recipients)
	if err != nil {
		return err
	}

	req["jti"] = uuid.New().String()
	req["iss"] = s.selfID
	req["iat"] = ntp.TimeFunc().Format(time.RFC3339)
	req["exp"] = ntp.TimeFunc().Add(s.expiry).Format(time.RFC3339)
	req["device_id"] = s.deviceID

	if gid, ok := req["gid"]; ok && len(gid.(string)) > 0 {
		req["aud"] = gid
		req["sub"] = gid

		body, err := request.Serialize(req, s.keyID, s.sk)
		if err != nil {
			return err
		}

		err = s.messagingClient.Send(recs, body)
		if err != nil {
			return err
		}
	} else {
		for _, recipient := range recs {
			r := strings.Split(recipient, ":")[0]
			req["aud"] = r
			req["sub"] = r

			body, err := request.Serialize(req, s.keyID, s.sk)
			if err != nil {
				return err
			}

			err = s.messagingClient.Send([]string{recipient}, body)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) createMissingSessions(members []string) error {
	println("creating missing sessions 1")
	sw := false
	unconnectedMembers := make([]string, 0)

	for _, m := range members {
		if sw {
			unconnectedMembers = append(unconnectedMembers, m)
		}
		if m == s.selfID {
			sw = true
		}
	}

	println("creating missing sessions")
	return s.send(unconnectedMembers, map[string]interface{}{"typ": "sessions.create"})
}
