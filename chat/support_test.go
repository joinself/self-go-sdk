package chat

import (
	"crypto/ed25519"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/joinself/self-go-sdk/pkg/request"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*Service, *Config) {
	pk := "1:56qJGhYCJmTHsYChCp3sPSjmiGlN2yG0KakYDquMAD0"
	kp := strings.Split(pk, ":")

	decoder := base64.RawStdEncoding

	skData, err := decoder.DecodeString(kp[1])
	assert.Nil(t, err)

	rh := request.New(request.Config{
		SelfID:   "c4f81d86-9dac-40fd-9830-13c66a0b2345",
		DeviceID: "1",
	})
	config := Config{
		SelfID:        "c4f81d86-9dac-40fd-9830-13c66a0b2345",
		KeyID:         kp[0],
		PrivateKey:    ed25519.NewKeyFromSeed(skData),
		RequestHelper: rh,
	}
	s := NewService(config)

	return s, &config
}
