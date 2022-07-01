// Copyright 2020 Self Group Ltd. All Rights Reserved.

package identity

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eteu-technologies/near-api-go/pkg/client"
	"github.com/eteu-technologies/near-api-go/pkg/client/block"
	"github.com/joinself/self-go-sdk/pkg/siggraph"
	"github.com/purehyperbole/merkletree"
)

var (
	// IdentityTypeIndividual an individual's identity
	IdentityTypeIndividual = "individual"
	// IdentityTypeApp an application's identity
	IdentityTypeApp = "app"
)

// User identity -> Individual?
// Would be nice to have a way to refer/represent different classes of app or user identity
// so Identity can be either an [Individual, App, Asset], etc.

// Identity represents all information about a self identity
type Identity struct {
	Name    string            `json:"name,omitempty"`
	SelfID  string            `json:"id"`
	Type    string            `json:"type"`
	History []json.RawMessage `json:"history"`
	Proofs  []PublicKeyProof  `json:"proofs"`
}

// PublicKeyProof a merkle proof for a public key operation
type PublicKeyProof struct {
	Operation  *Operation  `json:"operation,omitempty"`
	MerkleRoot *MerkleRoot `json:"merkle_root,omitempty"`
}

// Operation references a public key operation
type Operation struct {
	Sequence int    `json:"sequence"`
	Hash     string `json:"hash"`
	Alg      string `json:"alg"`
}

// MerkleRoot represents a merkle tree root hash that has been checkpointed to a public source
type MerkleRoot struct {
	Location  string `json:"location"`
	Sequence  int    `json:"sequence"`
	Timestamp int64  `json:"timestamp"`
	Hash      string `json:"hash"`
	Alg       string `json:"alg"`
}

// MerkleProof represents a merkle tree proof
type MerkleProof []MerkleProofEntry

// MerkleProofEntry represents a entry for a merkle tree proof
type MerkleProofEntry struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

// Device represents an identities device
type Device string

// GetIdentity gets an identity by its self ID
func (s Service) GetIdentity(selfID string) (*Identity, error) {
	var identity Identity
	var resp []byte
	var err error

	switch classifySelfID(selfID) {
	case IdentityTypeIndividual:
		resp, err = s.api.Get("/v1/identities/" + selfID)
	case IdentityTypeApp:
		resp, err = s.api.Get("/v1/apps/" + selfID)
	}

	if err != nil {
		return nil, err
	}

	fmt.Println(string(resp))

	err = json.Unmarshal(resp, &identity)
	if err != nil {
		return nil, err
	}

	// validate the signature graph
	_, err = siggraph.New(identity.History)
	if err != nil {
		return nil, err
	}

	// check that history entries have been checkpointed
	for i := 0; i < len(identity.History); i++ {
		if i >= len(identity.Proofs) {
			break
		}

		hash := sha256.Sum256(identity.History[i])
		target := hex.EncodeToString(hash[:])

		root, err := hex.DecodeString(identity.Proofs[i].MerkleRoot.Hash)
		if err != nil {
			return nil, err
		}

		if identity.Proofs[i].Operation.Hash != target {
			return nil, errors.New("operation hash does not match")
		}

		if identity.Proofs[i].Operation.Sequence != i {
			return nil, errors.New("operation sequence does not match")
		}

		fmt.Println("root:", identity.Proofs[i].MerkleRoot.Hash, "target:", target)

		endpoint := fmt.Sprintf("/v1/merkle/%s/proof/%s", identity.Proofs[i].MerkleRoot.Hash, target)

		resp, err = s.api.Get(endpoint)
		if err != nil {
			return nil, err
		}

		var proof MerkleProof

		fmt.Println(string(resp))

		err = json.Unmarshal(resp, &proof)
		if err != nil {
			return nil, err
		}

		var mp merkletree.Proof

		for _, p := range proof {
			if p.Left != "" {
				l, err := hex.DecodeString(p.Left)
				if err != nil {
					return nil, err
				}

				mp = append(mp, merkletree.Pair{
					Left: l,
				})
			} else {
				r, err := hex.DecodeString(p.Right)
				if err != nil {
					return nil, err
				}

				mp = append(mp, merkletree.Pair{
					Right: r,
				})
			}
		}

		fmt.Println("validating merkle proof")

		err = merkletree.Validate(sha256.New(), hash[:], root, mp)
		if err != nil {
			return nil, err
		}

		fmt.Println("checking merkle root exists on NEAR chain")

		rpc, err := client.NewClient("https://rpc.testnet.near.org")
		if err != nil {
			return nil, err
		}

		args, err := json.Marshal(map[string]string{
			"hash": identity.Proofs[i].MerkleRoot.Hash,
		})

		if err != nil {
			return nil, err
		}

		resp, err := rpc.ContractViewCallFunction(context.Background(), "identity.review.self.testnet", "get", base64.RawStdEncoding.EncodeToString(args), block.FinalityFinal())
		if err != nil {
			return nil, err
		}

		if bytes.Equal([]byte("null"), resp.Result) {
			return nil, errors.New("merkle root checkpoint does not exist!")
		}

		fmt.Println("checkpoint:", resp.Result)

		checkpoint, err := base64.RawURLEncoding.DecodeString(string(resp.Result[1 : len(resp.Result)-1]))
		if err != nil {
			return nil, err
		}

		// check checkpoint contains merkle root
		if !bytes.Equal(checkpoint[1:len(root)+1], root) {
			return nil, errors.New("checkpoint signed data does not contain merkle root hash")
		}

		// validate signature of checkpoint data
		pkd, err := base64.RawStdEncoding.DecodeString("dRglugY4mSrmXHPWC9b+NRCTBSZHrI/TV8tm6zQWHCo")
		if err != nil {
			return nil, err
		}

		pk := ed25519.PublicKey(pkd)

		signature := checkpoint[len(checkpoint)-ed25519.SignatureSize:]
		data := checkpoint[:len(checkpoint)-(ed25519.SignatureSize+1)]

		if !ed25519.Verify(pk, data, signature) {
			return nil, errors.New("checkpoint signature mismatch!")
		}

		fmt.Println("checkpoint validation passed!")
	}

	return &identity, nil
}

// GetDevices gets an identities devices
func (s Service) GetDevices(selfID string) ([]Device, error) {
	var devices []Device

	var resp []byte
	var err error

	switch classifySelfID(selfID) {
	case IdentityTypeIndividual:
		resp, err = s.api.Get("/v1/identities/" + selfID + "/devices/")
	case IdentityTypeApp:
		resp, err = s.api.Get("/v1/apps/" + selfID + "/devices/")
	}

	if err != nil {
		return nil, err
	}

	return devices, json.Unmarshal(resp, &devices)
}

// GetHistory gets the public key history of an identity
func (s Service) GetHistory(selfID string) ([]json.RawMessage, error) {
	var keys []json.RawMessage

	var resp []byte
	var err error

	switch classifySelfID(selfID) {
	case IdentityTypeIndividual:
		resp, err = s.api.Get("/v1/identities/" + selfID + "/history/")
	case IdentityTypeApp:
		resp, err = s.api.Get("/v1/apps/" + selfID + "/history/")
	}

	if err != nil {
		return nil, err
	}

	return keys, json.Unmarshal(resp, &keys)
}

func classifySelfID(selfID string) string {
	if len(selfID) > 11 {
		return IdentityTypeApp
	}

	return IdentityTypeIndividual
}
