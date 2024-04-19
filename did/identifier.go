package did

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
)

type Method int

const (
	MethodAure Method = iota
	MethodKey
)

const (
	AlgorithmEd25519   = byte(0)
	AlgorithmX25519    = byte(1)
	AlgorithmSecp256k1 = byte(2)
)

var (
	methodAure = []byte("aure:")
	methodKey  = []byte("key:")
)

type Identifier struct {
	method Method
	id     []byte
	pk     []byte
}

func NewIdentifier(method Method, id []byte, key ...[]byte) *Identifier {
	var pk []byte

	if len(key) > 0 {
		pk = key[0]
	}

	return &Identifier{
		method: method,
		id:     id,
		pk:     pk,
	}
}

func (i *Identifier) Build() string {
	switch i.method {
	case MethodAure:
		if i.pk != nil {
			return Aure(i.id, i.pk)
		} else {
			return Aure(i.id)
		}
	default:
		panic("")
	}
}

func (i *Identifier) Subject() string {
	switch i.method {
	case MethodAure:
		return Aure(i.id)
	default:
		panic("")
	}
}

func (i *Identifier) PublicKey() ed25519.PublicKey {
	return i.pk
}

func Aure(id []byte, key ...[]byte) string {
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
		algorithmKey[0] = AlgorithmEd25519
		copy(algorithmKey[1:], key[0])

		b.WriteString(base58.Encode(algorithmKey))
	}

	return b.String()
}

func Key(key []byte) string {
	var b strings.Builder

	b.Grow(50)

	b.Write(methodKey)

	algorithmKey := make([]byte, len(key)+1)
	algorithmKey[0] = AlgorithmEd25519
	copy(algorithmKey[1:], key)

	b.WriteString(base58.Encode(algorithmKey))

	return b.String()
}

func Extract(id string) ed25519.PublicKey {
	if len(id) != 116 {
		return nil
	}

	key := id[len(id)-46:]
	return ed25519.PublicKey(base58.Decode(key)[1:])
}

func Decompose(id string) (string, string, string, error) {
	parts := strings.Split(id, ":")

	switch len(parts) {
	case 2:
		return parts[0], parts[1], "", nil
	case 3:
		return parts[0], parts[1], parts[3], nil
	default:
		return "", "", "", errors.New("invalid did")
	}
}
