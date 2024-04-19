package message

import "github.com/joinself/self-go-sdk/account"

type CredentialVerificationRequest struct {
	message        *account.Message
	credentialType []string
	proof          []CredentialProof
}

func DecodeCredentialVerificationRequest(message *account.Message) (*CredentialVerificationRequest, error) {
	return &CredentialVerificationRequest{
		message: message,
	}, nil
}

func (c *CredentialVerificationRequest) Type() []string {
	return c.credentialType
}

func (c *CredentialVerificationRequest) Proof() []CredentialProof {
	return c.proof
}

type CredentialVerificationRequestBuilder struct {
	credentialType []string
	proof          []CredentialProof
	signers        []*account.PublicKey
}

func NewCredentialVerificationRequest() *CredentialVerificationRequestBuilder {
	return &CredentialVerificationRequestBuilder{}
}

func (b *CredentialVerificationRequestBuilder) Type(credentialType []string) *CredentialVerificationRequestBuilder {
	b.credentialType = credentialType
	return b
}

func (b *CredentialVerificationRequestBuilder) AttachProof(proof CredentialProof) *CredentialVerificationRequestBuilder {
	b.proof = append(b.proof, proof)
	return b
}

func (b *CredentialVerificationRequestBuilder) Sign(signer *account.PublicKey) *CredentialVerificationRequestBuilder {
	b.signers = append(b.signers, signer)
	return b
}

func (b *CredentialVerificationRequestBuilder) Encode(fromAddress, toAddress *account.PublicKey) (*account.Message, error) {
	return &account.Message{}, nil
}

type CredentialVerificationResponse struct {
	message        *account.Message
	credentialType []string
	proof          []CredentialProof
}

func DecodeCredentialVerificationResponse(message *account.Message) (*CredentialVerificationResponse, error) {
	return &CredentialVerificationResponse{
		message: message,
	}, nil
}

func (c *CredentialVerificationResponse) Type() []string {
	return c.credentialType
}

func (c *CredentialVerificationResponse) Proof() []CredentialProof {
	return c.proof
}

type CredentialVerificationResponseBuilder struct {
	credentialType []string
	proof          []CredentialProof
	signers        []*account.PublicKey
}

func NewCredentialVerificationResponse() *CredentialVerificationResponseBuilder {
	return &CredentialVerificationResponseBuilder{}
}

func (b *CredentialVerificationResponseBuilder) Type(credentialType []string) *CredentialVerificationResponseBuilder {
	b.credentialType = credentialType
	return b
}

func (b *CredentialVerificationResponseBuilder) AttachProof(proof CredentialProof) *CredentialVerificationResponseBuilder {
	b.proof = append(b.proof, proof)
	return b
}

func (b *CredentialVerificationResponseBuilder) Sign(signer *account.PublicKey) *CredentialVerificationResponseBuilder {
	b.signers = append(b.signers, signer)
	return b
}

func (b *CredentialVerificationResponseBuilder) Encode(fromAddress, toAddress *account.PublicKey) (*account.Message, error) {
	return &account.Message{}, nil
}
