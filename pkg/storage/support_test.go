package storage

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/joinself/self-go-sdk/pkg/siggraph"
	"github.com/square/go-jose"
	"github.com/stretchr/testify/require"
)

type testPKI struct {
	dkoff   map[string]int
	dkeys   map[string][]byte
	history map[string][]json.RawMessage
	mu      sync.Mutex
}

func newTestPKI(t testing.TB) *testPKI {
	return &testPKI{
		dkoff:   make(map[string]int),
		dkeys:   make(map[string][]byte),
		history: make(map[string][]json.RawMessage),
	}
}

func (p *testPKI) GetHistory(selfID string) ([]json.RawMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.history[selfID], nil
}

func (p *testPKI) GetDeviceKey(selfID, deviceID string) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var keys oneTimeKeys

	err := json.Unmarshal(p.dkeys[selfID+":"+deviceID], &keys)
	if err != nil {
		return nil, err
	}

	kid := p.dkoff[selfID+":"+deviceID]

	if kid > len(keys)-1 {
		return nil, errors.New("prekeys exhausted")
	}

	p.dkoff[selfID+":"+deviceID]++

	return json.Marshal(keys[kid])
}

func (p *testPKI) SetDeviceKeys(selfID, deviceID string, pkb []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.dkeys[selfID+":"+deviceID] = pkb
	return nil
}

func (p *testPKI) ListDeviceKeys(selfID, deviceID string) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.dkeys[selfID+":"+deviceID]) < 1 {
		return []byte("[]"), nil
	}
	return p.dkeys[selfID+":"+deviceID], nil
}

func (p *testPKI) purgeDeviceKeys(selfID, deviceID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.dkeys[selfID+":"+deviceID] = nil
}

func (p *testPKI) addpk(selfID string, sk ed25519.PrivateKey, pk ed25519.PublicKey) {
	now := time.Now().Unix()

	rpk, _, _ := ed25519.GenerateKey(rand.Reader)

	p.history[selfID] = []json.RawMessage{
		testop(sk, "1", &siggraph.Operation{
			Sequence:  0,
			Version:   "1.0.0",
			Previous:  "-",
			Timestamp: now,
			Actions: []siggraph.Action{
				{
					KID:           "1",
					DID:           "1",
					Type:          siggraph.TypeDeviceKey,
					Action:        siggraph.ActionKeyAdd,
					EffectiveFrom: now,
					Key:           base64.RawURLEncoding.EncodeToString(pk),
				},
				{
					KID:           "2",
					Type:          siggraph.TypeRecoveryKey,
					Action:        siggraph.ActionKeyAdd,
					EffectiveFrom: now,
					Key:           base64.RawURLEncoding.EncodeToString(rpk),
				},
			},
		}),
	}
}

func registerUser(t testing.TB, pki *testPKI, id string) ed25519.PrivateKey {
	// generate an identity keypair
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	identifier, _ := idsplit(id)

	pki.addpk(identifier, sk, pk)

	return sk
}

func testop(sk ed25519.PrivateKey, kid string, op *siggraph.Operation) json.RawMessage {
	data, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}

	opts := &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"kid": kid,
		},
	}

	s, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: sk}, opts)
	if err != nil {
		panic(err)
	}

	jws, err := s.Sign(data)
	if err != nil {
		panic(err)
	}

	return json.RawMessage(jws.FullSerialize())
}
