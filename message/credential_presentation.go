package message

import (
	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/keypair/signing"
)

type CredentialPresentationRequest struct {
	message     *account.Message
	credentials [][]string
	challenge   []byte
}

func DecodeCredentialPresentationRequest(message *account.Message) (*CredentialPresentationRequest, error) {
	return &CredentialPresentationRequest{
		message: message,
	}, nil
}

func (c *CredentialPresentationRequest) Credentials() [][]string {
	return c.credentials
}

func (c *CredentialPresentationRequest) Challenge() []byte {
	return c.challenge
}

type CredentialPresentationRequestBuilder struct {
	credentialType []string
	proof          []CredentialProof
	signers        []*signing.PublicKey
}

func NewCredentialPresentationRequest() *CredentialPresentationRequestBuilder {
	return &CredentialPresentationRequestBuilder{}
}

func (b *CredentialPresentationRequestBuilder) Type(credentialType []string) *CredentialPresentationRequestBuilder {
	b.credentialType = credentialType
	return b
}

func (b *CredentialPresentationRequestBuilder) AttachProof(proof CredentialProof) *CredentialPresentationRequestBuilder {
	b.proof = append(b.proof, proof)
	return b
}

func (b *CredentialPresentationRequestBuilder) Sign(signer *signing.PublicKey) *CredentialPresentationRequestBuilder {
	b.signers = append(b.signers, signer)
	return b
}

func (b *CredentialPresentationRequestBuilder) Build(fromAddress, toAddress *signing.PublicKey) (*account.Message, error) {
	return &account.Message{}, nil
}

type CredentialPresentationResponse struct {
	message        *account.Message
	credentialType []string
	proof          []CredentialProof
}

func DecodeCredentialPresentationResponse(message *account.Message) (*CredentialPresentationResponse, error) {
	return &CredentialPresentationResponse{
		message: message,
	}, nil
}

func (c *CredentialPresentationResponse) Type() []string {
	return c.credentialType
}

func (c *CredentialPresentationResponse) Proof() []CredentialProof {
	return c.proof
}

type CredentialPresentationResponseBuilder struct {
	credentialType []string
	proof          []CredentialProof
	signers        []*signing.PublicKey
}

func NewCredentialPresentationResponse() *CredentialPresentationResponseBuilder {
	return &CredentialPresentationResponseBuilder{}
}

func (b *CredentialPresentationResponseBuilder) Type(credentialType []string) *CredentialPresentationResponseBuilder {
	b.credentialType = credentialType
	return b
}

func (b *CredentialPresentationResponseBuilder) AttachProof(proof CredentialProof) *CredentialPresentationResponseBuilder {
	b.proof = append(b.proof, proof)
	return b
}

func (b *CredentialPresentationResponseBuilder) Sign(signer *signing.PublicKey) *CredentialPresentationResponseBuilder {
	b.signers = append(b.signers, signer)
	return b
}

func (b *CredentialPresentationResponseBuilder) Encode(fromAddress, toAddress *signing.PublicKey) (*account.Message, error) {
	return &account.Message{}, nil
}
