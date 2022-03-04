package kidhelper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

func GetJWSKID(payload []byte) (string, error) {
	var jws struct {
		Protected string `json:"protected"`
	}

	err := json.Unmarshal(payload, &jws)
	if err != nil {
		return "", err
	}

	return getKID(jws.Protected)
}

func getKID(token string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(strings.Split(token, ".")[0])
	if err != nil {
		return "", err
	}

	hdr := make(map[string]string)

	err = json.Unmarshal(data, &hdr)
	if err != nil {
		return "", err
	}

	kid := hdr["kid"]
	if kid == "" {
		return "", errors.New("token must specify an identifier for the signing key")
	}

	return kid, nil
}
