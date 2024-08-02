package credential

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"errors"
	"runtime"
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair/signing"
)

var (
	CredentialTypePassport             = NewCredentialTypeCollection().Append("VerifiableCredential").Append("PassportCredential")
	CredentialTypeLiveness             = NewCredentialTypeCollection().Append("VerifiableCredential").Append("LivenessCredential")
	CredentialTypeProfileImage         = NewCredentialTypeCollection().Append("VerifiableCredential").Append("ProfileImageCredential")
	CredentialTypeApplicationPublisher = NewCredentialTypeCollection().Append("VerifiableCredential").Append("ApplicationPublisherCredential")
)

type Credential C.self_credential
type CredentialBuilder C.self_credential_builder
type VerifiableCredential C.self_verifiable_credential
type VerifiableCredentialCollection C.self_collection_verifiable_credential
type CredentialPresentationDetail C.self_credential_presentation_detail
type CredentialPresentationDetailCollection C.self_collection_credential_presentation_detail
type CredentialVerificationEvidence C.self_credential_verification_evidence
type CredentialVerificationEvidenceCollection C.self_collection_credential_verification_evidence
type CredentialTypeCollection C.self_collection_credential_type

// NewCredential creates a new credential builder
func NewCredential() *CredentialBuilder {
	builder := (*CredentialBuilder)(C.self_credential_builder_init())

	runtime.SetFinalizer(builder, func(builder *CredentialBuilder) {
		C.self_credential_builder_destroy(
			(*C.self_credential_builder)(builder),
		)
	})

	return builder
}

// CredentialType sets the type of credential
func (b *CredentialBuilder) CredentialType(credentialType *CredentialTypeCollection) *CredentialBuilder {
	C.self_credential_builder_credential_type(
		(*C.self_credential_builder)(b),
		(*C.self_collection_credential_type)(credentialType),
	)
	return b
}

// CredentialSubject sets the address of the credential's subject
func (b *CredentialBuilder) CredentialSubject(subjectAddress *Address) *CredentialBuilder {
	C.self_credential_builder_credential_subject(
		(*C.self_credential_builder)(b),
		(*C.self_credential_address)(subjectAddress),
	)
	return b
}

// CredentialSubjectClaim adds a claim about the subject to the credential
func (b *CredentialBuilder) CredentialSubjectClaim(claimKey, claimValue string) *CredentialBuilder {
	key := C.CString(claimKey)
	value := C.CString(claimValue)

	defer func() {
		C.free(unsafe.Pointer(key))
		C.free(unsafe.Pointer(value))
	}()

	C.self_credential_builder_credential_subject_claim(
		(*C.self_credential_builder)(b),
		key,
		value,
	)

	return b
}

// CredentialSubjectClaims adds a collection of claims about the subject to te credential
func (b *CredentialBuilder) CredentialSubjectClaims(claims map[string]interface{}) *CredentialBuilder {
	claim, err := json.Marshal(claims)
	if err == nil {
		return b
	}

	claimBuffer := C.CBytes(claim)
	claimLength := len(claim)

	defer func() {
		C.free(claimBuffer)
	}()

	C.self_credential_builder_credential_subject_json(
		(*C.self_credential_builder)(b),
		(*C.uchar)(claimBuffer),
		(C.ulong)(claimLength),
	)

	return b
}

// Issuer sets the address of the credential's issuer
func (b *CredentialBuilder) Issuer(issuerAddress *Address) *CredentialBuilder {
	C.self_credential_builder_issuer(
		(*C.self_credential_builder)(b),
		(*C.self_credential_address)(issuerAddress),
	)
	return b
}

// ValidFrom sets the point of validity for the credential
func (b *CredentialBuilder) ValidFrom(timestamp time.Time) *CredentialBuilder {
	C.self_credential_builder_valid_from(
		(*C.self_credential_builder)(b),
		C.long(timestamp.Unix()),
	)
	return b
}

// SignWith signs the credential
func (b *CredentialBuilder) SignWith(signer *signing.PublicKey, issuedAt time.Time) *CredentialBuilder {
	C.self_credential_builder_sign_with(
		(*C.self_credential_builder)(b),
		(*C.self_signing_public_key)(signer),
		C.long(issuedAt.Unix()),
	)
	return b
}

// Finish generates and prepares the credential for being signed by an account
func (b *CredentialBuilder) Finish() (*Credential, error) {
	var credentialPtr *C.self_credential

	status := C.self_credential_builder_finish(
		(*C.self_credential_builder)(b),
		&credentialPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to create credential")
	}

	credential := (*Credential)(credentialPtr)

	runtime.SetFinalizer(credential, func(credential *Credential) {
		C.self_credential_destroy(
			(*C.self_credential)(credential),
		)
	})

	return credential, nil
}

// CredentialType returns the type of credential
func (c *VerifiableCredential) CredentialType() *CredentialTypeCollection {
	collection := (*CredentialTypeCollection)(C.self_verifiable_credential_credential_type(
		(*C.self_verifiable_credential)(c),
	))

	runtime.SetFinalizer(collection, func(collection *CredentialTypeCollection) {
		C.self_collection_credential_type_destroy(
			(*C.self_collection_credential_type)(collection),
		)
	})

	return collection
}

// CredentialSubject returns the subject of the credential's address
func (c *VerifiableCredential) CredentialSubject() *Address {
	subject := (*Address)(C.self_verifiable_credential_credential_subject(
		(*C.self_verifiable_credential)(c),
	))

	runtime.SetFinalizer(subject, func(address *Address) {
		C.self_credential_address_destroy(
			(*C.self_credential_address)(address),
		)
	})

	return subject
}

// CredentialSubject returns the subject of the credential's address
func (c *VerifiableCredential) CredentialSubjectClaim(claimKey string) (string, bool) {
	key := C.CString(claimKey)

	value := C.self_verifiable_credential_credential_subject_claim(
		(*C.self_verifiable_credential)(c),
		key,
	)

	C.free(unsafe.Pointer(key))

	if value == nil {
		return "", false
	}

	claimValue := C.GoBytes(
		unsafe.Pointer(C.self_credential_claim_value_buf(value)),
		C.int(C.self_credential_claim_value_len(value)),
	)

	C.self_credential_claim_value_destroy(value)

	return string(claimValue), true
}

// Issuer returns the address of the credential's issuer
func (c *VerifiableCredential) Issuer() *Address {
	issuer := (*Address)(C.self_verifiable_credential_issuer(
		(*C.self_verifiable_credential)(c),
	))

	runtime.SetFinalizer(issuer, func(address *Address) {
		C.self_credential_address_destroy(
			(*C.self_credential_address)(address),
		)
	})

	return issuer
}

// ValidFrom returns the time period that the credential is valid from
func (c *VerifiableCredential) ValidFrom() time.Time {
	validFrom := C.self_verifiable_credential_valid_from(
		(*C.self_verifiable_credential)(c),
	)

	return time.Unix(int64(validFrom), 0)
}

// Created returns the time that the credential was created
func (c *VerifiableCredential) Created() time.Time {
	created := C.self_verifiable_credential_created(
		(*C.self_verifiable_credential)(c),
	)

	return time.Unix(int64(created), 0)
}

// Validate validates the contents of the credential and it's signatures
func (c *VerifiableCredential) Validate() error {
	status := C.self_verifiable_credential_validate(
		(*C.self_verifiable_credential)(c),
	)

	if status > 0 {
		return errors.New("credential invalid")
	}

	return nil
}

func NewVerifiableCredentialCollection() *VerifiableCredentialCollection {
	collection := (*VerifiableCredentialCollection)(C.self_collection_verifiable_credential_init())

	runtime.SetFinalizer(collection, func(collection *VerifiableCredentialCollection) {
		C.self_collection_verifiable_credential_destroy(
			(*C.self_collection_verifiable_credential)(collection),
		)
	})

	return collection
}

func (c *VerifiableCredentialCollection) Length() int {
	return int(C.self_collection_verifiable_credential_len(
		(*C.self_collection_verifiable_credential)(c),
	))
}

func (c *VerifiableCredentialCollection) Get(index int) *VerifiableCredential {
	return (*VerifiableCredential)(C.self_collection_verifiable_credential_at(
		(*C.self_collection_verifiable_credential)(c),
		C.ulong(index),
	))
}

func NewCredentialVerificationEvidenceCollection() *CredentialVerificationEvidenceCollection {
	collection := (*CredentialVerificationEvidenceCollection)(C.self_collection_credential_verification_evidence_init())

	runtime.SetFinalizer(collection, func(collection *CredentialVerificationEvidenceCollection) {
		C.self_collection_credential_verification_evidence_destroy(
			(*C.self_collection_credential_verification_evidence)(collection),
		)
	})

	return collection
}

func (c *CredentialVerificationEvidenceCollection) Length() int {
	return int(C.self_collection_credential_verification_evidence_len(
		(*C.self_collection_credential_verification_evidence)(c),
	))
}

func (c *CredentialVerificationEvidenceCollection) Get(index int) *CredentialVerificationEvidence {
	return (*CredentialVerificationEvidence)(C.self_collection_credential_verification_evidence_at(
		(*C.self_collection_credential_verification_evidence)(c),
		C.ulong(index),
	))
}

func NewCredentialTypeCollection() *CredentialTypeCollection {
	collection := (*CredentialTypeCollection)(C.self_collection_credential_type_init())

	runtime.SetFinalizer(collection, func(collection *CredentialTypeCollection) {
		C.self_collection_credential_type_destroy(
			(*C.self_collection_credential_type)(collection),
		)
	})

	return collection
}

func (c *CredentialTypeCollection) Length() int {
	return int(C.self_collection_credential_type_len(
		(*C.self_collection_credential_type)(c),
	))
}

func (c *CredentialTypeCollection) Get(index int) string {
	return C.GoString(C.self_collection_credential_type_at(
		(*C.self_collection_credential_type)(c),
		C.ulong(index),
	))
}

func (c *CredentialTypeCollection) Append(element string) *CredentialTypeCollection {
	elementC := C.CString(element)

	C.self_collection_credential_type_append(
		(*C.self_collection_credential_type)(c),
		elementC,
	)

	C.free(unsafe.Pointer(elementC))

	return c
}
