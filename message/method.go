package message

import (
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
)

var (
	methodAure = []byte("aure:")
	methodKey  = []byte("key:")
)

func aure(id []byte, key ...[]byte) string {
	var b strings.Builder

	if len(key) > 0 {
		b.Grow(116)
	} else {
		b.Grow(69)
	}

	b.Write(methodAure)
	b.WriteString(hex.EncodeToString(id))

	if len(key) > 0 {
		b.WriteRune('#')

		algorithmKey := make([]byte, len(key[0])+1)
		// TODO load these from the native SDK
		algorithmKey[0] = 0 //C.self_public_key_algorithm::ED25519
		copy(algorithmKey[1:], key[0])

		b.WriteString(base58.Encode(algorithmKey))
	}

	return b.String()
}
