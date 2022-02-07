// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joinself/self-go-sdk/pkg/ntp"
	"github.com/square/go-jose"
)

func (s *Service) SelfID() string {
	return s.selfID
}

// Message sends a message to a list of recipients.
func (s *Service) Message(recipients []string, body string, opts map[string]interface{}) *Message {
	payload := map[string]interface{}{
		"typ": "chat.message",
		"msg": body,
	}

	payload["aud"] = recipients[0]
	if _, ok := opts["gid"]; ok {
		payload["gid"] = opts["gid"]
		payload["aud"] = opts["gid"]
	}
	if jti, ok := opts["jti"]; ok {
		payload["jti"] = jti
	}
	if rid, ok := opts["rid"]; ok {
		payload["rid"] = rid
	}

	if _, ok := opts["objects"]; ok {
		// fi := NewRemoteFileInteractor(s.api)
		objects := make([]interface{}, 0)
		for _, o := range opts["objects"].([]map[string]interface{}) {
			if _, ok := o["link"]; ok {
				// Is a public image, just append it
				objects = append(objects, o)
			} else {
				fo := NewObject(s.FileInteractor)
				err := fo.BuildFromData(o["data"].([]byte), o["name"].(string), o["mime"].(string))
				if err == nil {
					objects = append(objects, fo.ToPayload())
				} else {
					log.Println(err.Error())
				}
			}
		}
		payload["objects"] = objects
	}

	s.send(recipients, payload)

	return NewMessage(s, recipients, payload)
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

// Invite sends a group invitation to a list of members.
func (s *Service) Invite(gid string, name string, members []string, opts map[string]interface{}) {
	p := map[string]interface{}{
		"typ":     "chat.invite",
		"gid":     gid,
		"name":    name,
		"members": members,
		"aud":     gid,
	}

	if _, ok := opts["data"]; ok {
		fo := NewObject(s.FileInteractor)
		err := fo.BuildFromData(opts["data"].([]byte), "", opts["mime"].(string))
		if err == nil {
			objPayload := fo.ToPayload()
			opts["link"] = objPayload["link"]
			opts["key"] = objPayload["key"]
			opts["expires"] = objPayload["expires"]
		}
	}

	s.send(members, p)
}

// Join joins a group.
func (s *Service) Join(gid string, members []string) {
	// Allow incoming connections from the given members.
	for _, m := range members {
		s.messagingService.PermitConnection(m)
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
	recs, err := s.recipients(recipients)
	if err != nil {
		return err
	}

	req["jti"] = uuid.New().String()
	req["iss"] = s.selfID
	req["iat"] = ntp.TimeFunc().Format(time.RFC3339)
	req["exp"] = ntp.TimeFunc().Add(s.expiry).Format(time.RFC3339)
	req["device_id"] = s.deviceID

	for _, recipient := range recipients {
		req["aud"] = recipient
		req["sub"] = recipient

		payload, err := json.Marshal(req)
		if err != nil {
			return err
		}
		println(string(payload))

		opts := &jose.SignerOptions{
			ExtraHeaders: map[jose.HeaderKey]interface{}{
				"kid": s.keyID,
			},
		}

		signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: s.sk}, opts)
		if err != nil {
			return err
		}

		signature, err := signer.Sign(payload)
		if err != nil {
			return err
		}

		body := []byte(signature.FullSerialize())

		err = s.messagingClient.Send(recs, body)
		if err != nil {
			return err
		}
	}
	return nil
}

// builds a list of all devices associated with an identity
func (s Service) recipients(recipients []string) ([]string, error) {
	devices := make([]string, 0)
	for _, selfID := range recipients {
		dds, err := s.getDevices(selfID)
		if err != nil {
			return nil, err
		}

		for i := range dds {
			if selfID != s.selfID && dds[i] != s.deviceID {
				devices = append(devices, selfID+":"+dds[i])
			}
		}
	}

	return devices, nil
}

func (s *Service) createMissingSessions(members []string) error {
	println("creating missing sessions 1")
	sw := false
	posteriorMembers := make([]string, 0)

	for _, m := range members {
		if sw {
			posteriorMembers = append(posteriorMembers, m)
		}
		if m == s.selfID {
			sw = true
		}
	}

	println("creating missing sessions")
	return s.send(posteriorMembers, map[string]interface{}{"typ": "sessions.create"})
}

func (s Service) getDevices(selfID string) ([]string, error) {
	var resp []byte
	var err error

	if len(selfID) > 11 {
		resp, err = s.api.Get("/v1/apps/" + selfID + "/devices")
	} else {
		resp, err = s.api.Get("/v1/identities/" + selfID + "/devices")
	}
	if err != nil {
		return nil, err
	}

	var devices []string
	err = json.Unmarshal(resp, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
