// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"encoding/base64"
	"errors"
	"net/http"

	"golang.org/x/crypto/chacha20poly1305"
)

type RemoteFileInteractor struct {
	token string
	url   string
}

func NewRemoteFileInteractor(token, url string) *RemoteFileInteractor {
	return &RemoteFileInteractor{
		token: token,
		url:   url,
	}
}

func (r *RemoteFileInteractor) Set(link, key string) ([]byte, error) {
	return []byte{}, errors.New("todo")

}

func (r *RemoteFileInteractor) Get(objectID, key string) ([]byte, error) {
	resp, err := http.Get(r.url + "/v1/objects/" + objectID)
	if err != nil {
		return nil, err
	}

	content := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(content)
	if err != nil {
		return nil, err
	}

	keyNonce, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	if len(keyNonce) != chacha20poly1305.KeySize+chacha20poly1305.NonceSizeX {
		return nil, errors.New("object key and nonce are of an invalid size")
	}

	sk := keyNonce[:chacha20poly1305.KeySize]
	nonce := keyNonce[chacha20poly1305.KeySize:]

	aead, err := chacha20poly1305.NewX(sk)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, nonce, content, nil)
}
