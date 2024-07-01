// Copyright 2020 Self Group Ltd. All Rights Reserved.

package object

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"strconv"
)

// RemoteFile manages interactions with the remote filles
type RemoteFile interface {
	SetObject(data []byte) (*EncryptedObject, error)
	GetObject(link, key string) ([]byte, error)
}

type Object struct {
	fi         RemoteFile
	Link       string
	Name       string
	Mime       string
	Expires    int64
	Key        string
	Nonce      string
	Ciphertext string
	Content    []byte
	Digest     string
}

// New creates an object.
func New(fi RemoteFile) *Object {
	return &Object{
		fi: fi,
	}
}

// BuildFromData builds an object from data.
func (o *Object) BuildFromData(data []byte, name, mime string) error {
	println("building an object")
	// Encrypt message
	obj, err := o.fi.SetObject(data)
	if err != nil {
		log.Println("error uploading the object")
		return err
	}

	o.Name = name
	o.Link = obj.Link
	o.Key = obj.Key
	o.Mime = mime
	o.Expires = obj.Expires
	o.Content = data
	o.Digest = o.calculateDigest(data)

	return nil
}

// BuildFromObject builds an object from a message object payload.
func (o *Object) BuildFromObject(obj map[string]interface{}) error {
	o.Name = obj["name"].(string)
	o.Link = obj["link"].(string)
	if _, ok := obj["key"]; ok {
		o.Key = obj["key"].(string)
		o.Mime = obj["mime"].(string)
		if val, ok := obj["object_hash"]; ok {
			o.Digest = val.(string)
		} else if val, ok := obj["image_hash"]; ok {
			o.Digest = val.(string)
		}
		if exp, exists := obj["expires"]; exists {
			switch v := exp.(type) {
			case string:
				if expires, err := strconv.ParseInt(obj["expires"].(string), 10, 64); err == nil {
					o.Expires = expires
				}
			case int64:
				o.Expires = v
			}
		}
	}

	return nil
}

// ToPayload translates the current object to payload.
func (o *Object) ToPayload() map[string]interface{} {
	return map[string]interface{}{
		"name":        o.Name,
		"link":        o.Link,
		"key":         o.Key,
		"mime":        o.Mime,
		"expires":     strconv.FormatInt(o.Expires, 10),
		"public":      false,
		"object_hash": o.Digest,
	}
}

// GetContent gets the current object content.
func (o *Object) GetContent() ([]byte, error) {
	content, err := o.fi.GetObject(o.Link, o.Key)
	if err != nil {
		println("error getting the object")
		println(err.Error())
		return []byte(""), err
	}
	if o.calculateDigest(content) != o.Digest {
		return content, errors.New("object digest and content digest do not match")
	}

	return content, nil
}

func (o *Object) calculateDigest(ct []byte) string {
	h := sha256.Sum256(ct)
	return base64.RawURLEncoding.EncodeToString(h[:])
}
