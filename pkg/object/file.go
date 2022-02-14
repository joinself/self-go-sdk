// Copyright 2020 Self Group Ltd. All Rights Reserved.

package object

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
)

type restTransport interface {
	Get(path string) ([]byte, error)
	Post(path string, ctype string, data []byte) ([]byte, error)
	BuildURL(path string) string
}

type EncryptedObject struct {
	Link    string
	Key     string
	Nonce   string
	Content string
	Expires int64
}

type RemoteFileInteractor struct {
	api restTransport
}

func NewRemoteFileInteractor(api restTransport) *RemoteFileInteractor {
	return &RemoteFileInteractor{api: api}
}

// Encrypts and uploads an object.
func (r *RemoteFileInteractor) SetObject(data []byte) (*EncryptedObject, error) {
	// Encrypt input data
	encryptedData, key, err := r.encrypt(data)
	if err != nil {
		return nil, err
	}

	// Upload to server
	resp, err := r.sendHttpPostRequest(encryptedData)
	if err != nil {
		return nil, err
	}

	return &EncryptedObject{
		Link:    r.api.BuildURL("/v1/objects/" + resp.ID),
		Key:     string(key),
		Expires: resp.Expires,
	}, nil
}

// GetObject downloads a decrypts a remote object.
func (r *RemoteFileInteractor) GetObject(link, key string) ([]byte, error) {
	parts := strings.Split(link, "joinself.com")
	path := parts[len(parts)-1]
	content, err := r.api.Get(path)
	if err != nil {
		println("xoxo 1")
		return nil, err
	}

	keyNonce, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		println("xoxo 2")
		return nil, err
	}

	if len(keyNonce) != chacha20poly1305.KeySize+chacha20poly1305.NonceSizeX {
		println("xoxo 4")
		return nil, errors.New("object key and nonce are of an invalid size")
	}

	sk := keyNonce[:chacha20poly1305.KeySize]
	nonce := keyNonce[chacha20poly1305.KeySize:]

	aead, err := chacha20poly1305.NewX(sk)
	if err != nil {
		println("xoxo 4")
		return nil, err
	}

	return aead.Open(nil, nonce, content, nil)
}

type postResponse struct {
	ID      string `json:"id"`
	Expires int64  `json:"expires"`
}

func (r *RemoteFileInteractor) sendHttpPostRequest(data []byte) (payload *postResponse, err error) {
	content, err := r.api.Post("/v1/objects", "application/octet-stream", data)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &payload)

	return
}

func (r *RemoteFileInteractor) encrypt(data []byte) ([]byte, string, error) {
	key := make([]byte, chacha20poly1305.KeySize+chacha20poly1305.NonceSizeX)
	_, err := rand.Read(key)
	if err != nil {
		log.Println("error building shareable key " + err.Error())
		return []byte(""), "", err
	}

	sk := key[:chacha20poly1305.KeySize]
	nonce := key[chacha20poly1305.KeySize:]

	_, err = rand.Read(key)
	if err != nil {
		return []byte(""), "", err
	}

	aead, err := chacha20poly1305.NewX(sk)
	if err != nil {
		return []byte(""), "", err
	}

	encryptedData := aead.Seal(nil, nonce, data, nil)

	return encryptedData, base64.RawURLEncoding.EncodeToString(key), nil
}
