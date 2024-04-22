package message

import (
	"crypto/ed25519"
	"encoding/json"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/gowebpki/jcs"
)

var (
	DefaultContext                 = []string{"VerifiableCredential"}
	DefaultCredential              = []string{"VerifiableCredential"}
	PassportCredential             = []string{"VerifiableCredential", "PassportCredential"}
	LivenessCredential             = []string{"VerifiableCredential", "LivenessCredential"}
	ProfileImageCredential         = []string{"VerifiableCredential", "ProfileImageCredential"}
	ApplicationPublisherCredential = []string{"VerifiableCredential", "ApplicationPublisherCredential"}
	DefaultCryptoSuite             = "jcs-eddsa-2022"
	TypeDataIntegrityProof         = "DataIntegrityProof"
	PurposeAssertion               = "assertionMethod"
)

type VerifiableCredential struct {
	Context   []string               `json:"@context"`
	ID        string                 `json:"id,omitempty"`
	Type      []string               `json:"type"`
	Issuer    string                 `json:"issuer"`
	ValidFrom string                 `json:"validFrom"`
	Subject   map[string]interface{} `json:"credentialSubject"`
	Proof     *Proof                 `json:"proof,omitempty"`
}

func NewVerifiableCredential(credentialType []string, issuer []byte, validFrom time.Time, subject map[string]interface{}) *VerifiableCredential {
	return &VerifiableCredential{
		Context:   DefaultContext,
		Type:      credentialType,
		Issuer:    aure(issuer),
		ValidFrom: validFrom.UTC().Format("2006-01-02T15:04:05Z07:00"),
		Subject:   subject,
	}
}

type Proof struct {
	Type               string `json:"type"`
	Cryptosuite        string `json:"cryptosuite"`
	Created            string `json:"created"`
	VerificationMethod string `json:"verificationMethod"`
	ProofPurpose       string `json:"proofPurpose"`
	ProofValue         string `json:"proofValue"`
}

func (vc *VerifiableCredential) SignDataIntegrity(by []byte, at time.Time, sk ed25519.PrivateKey) ([]byte, error) {
	return vc.SignDataIntegrityFunc(
		by,
		sk.Public().(ed25519.PublicKey),
		at,
		func(message []byte) []byte {
			return ed25519.Sign(sk, message)
		},
	)
}

func (vc *VerifiableCredential) SignDataIntegrityFunc(by []byte, pk ed25519.PublicKey, at time.Time, signFunc func([]byte) []byte) ([]byte, error) {
	data, err := json.Marshal(vc)
	if err != nil {
		return nil, err
	}

	// transform the data into JCS format for signing
	jcsData, err := jcs.Transform(data)
	if err != nil {
		return nil, err
	}

	vc.Proof = &Proof{
		Type:        TypeDataIntegrityProof,
		Cryptosuite: DefaultCryptoSuite,
		Created:     at.UTC().Format("2006-01-02T15:04:05Z07:00"),
		VerificationMethod: aure(
			by,
			pk,
		),
		ProofPurpose: PurposeAssertion,
		ProofValue:   base58.Encode(signFunc(jcsData)),
	}

	return json.Marshal(vc)
}

type CredentialProof struct {
	Type  []string
	Proof []byte
}
