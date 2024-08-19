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

	"github.com/joinself/self-go-sdk-next/keypair/signing"
	"github.com/joinself/self-go-sdk-next/object"
)

var (
	CredentialTypePassport             = newCredentialTypeCollection().Append("VerifiableCredential").Append("PassportCredential")
	CredentialTypeLiveness             = newCredentialTypeCollection().Append("VerifiableCredential").Append("LivenessCredential")
	CredentialTypeProfileImage         = newCredentialTypeCollection().Append("VerifiableCredential").Append("ProfileImageCredential")
	CredentialTypeApplicationPublisher = newCredentialTypeCollection().Append("VerifiableCredential").Append("ApplicationPublisherCredential")
)

type Credential C.self_credential
type CredentialBuilder C.self_credential_builder
type VerifiableCredential C.self_verifiable_credential
type VerifiableCredentialCollection C.self_collection_verifiable_credential
type CredentialPresentationDetail C.self_credential_presentation_detail
type CredentialPresentationDetailCollection C.self_collection_credential_presentation_detail
type CredentialVerificationEvidence C.self_credential_verification_evidence
type CredentialVerificationParameter C.self_credential_verification_parameter
type CredentialVerificationEvidenceCollection C.self_collection_credential_verification_evidence
type CredentialVerificationParameterCollection C.self_collection_credential_verification_parameter
type CredentialTypeCollection C.self_collection_credential_type

// NewCredential creates a new credential builder
func NewCredential() *CredentialBuilder {
	builder := C.self_credential_builder_init()

	runtime.SetFinalizer(&builder, func(builder **C.self_credential_builder) {
		C.self_credential_builder_destroy(
			*builder,
		)
	})

	return (*CredentialBuilder)(builder)
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
		(*C.uint8_t)(claimBuffer),
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
	var credential *C.self_credential
	credentialPtr := &credential

	status := C.self_credential_builder_finish(
		(*C.self_credential_builder)(b),
		credentialPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to create credential")
	}

	runtime.SetFinalizer(credentialPtr, func(credential **C.self_credential) {
		C.self_credential_destroy(
			*credential,
		)
	})

	return (*Credential)(*credentialPtr), nil
}

func DecodeVerifiableCredential(encodedCredential []byte) (*VerifiableCredential, error) {
	var verifiableCredential *C.self_verifiable_credential
	verifiableCredentialPtr := &verifiableCredential

	encodedBuf := C.CBytes(encodedCredential)
	encodedLen := len(encodedCredential)

	defer func() {
		C.free(encodedBuf)
	}()

	status := C.self_verifiable_credential_decode(
		verifiableCredentialPtr,
		(*C.uint8_t)(encodedBuf),
		(C.ulong)(encodedLen),
	)

	if status > 0 {
		return nil, errors.New("decode credential failed")
	}

	runtime.SetFinalizer(verifiableCredentialPtr, func(verifiableCredential **C.self_verifiable_credential) {
		C.self_verifiable_credential_destroy(
			*verifiableCredential,
		)
	})

	return (*VerifiableCredential)(*verifiableCredentialPtr), nil
}

// CredentialType returns the type of credential
func (c *VerifiableCredential) CredentialType() *CredentialTypeCollection {
	collection := C.self_verifiable_credential_credential_type(
		(*C.self_verifiable_credential)(c),
	)

	runtime.SetFinalizer(&collection, func(collection **C.self_collection_credential_type) {
		C.self_collection_credential_type_destroy(
			*collection,
		)
	})

	return (*CredentialTypeCollection)(collection)
}

// CredentialSubject returns the subject of the credential's address
func (c *VerifiableCredential) CredentialSubject() *Address {
	subject := C.self_verifiable_credential_credential_subject(
		(*C.self_verifiable_credential)(c),
	)

	runtime.SetFinalizer(&subject, func(address **C.self_credential_address) {
		C.self_credential_address_destroy(
			*address,
		)
	})

	return (*Address)(subject)
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
	issuer := C.self_verifiable_credential_issuer(
		(*C.self_verifiable_credential)(c),
	)

	runtime.SetFinalizer(&issuer, func(address **C.self_credential_address) {
		C.self_credential_address_destroy(
			*address,
		)
	})

	return (*Address)(issuer)
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

// Encode returns a json encoded verifiable credential
func (c *VerifiableCredential) Encode() ([]byte, error) {
	var encodedCredentialBuffer *C.self_encoded_buffer
	encodedCredentialBufferPtr := &encodedCredentialBuffer

	status := C.self_verifiable_credential_encode(
		(*C.self_verifiable_credential)(c),
		encodedCredentialBufferPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to encode credential")
	}

	encodedCredential := C.GoBytes(
		unsafe.Pointer(C.self_encoded_buffer_buf(*encodedCredentialBufferPtr)),
		C.int(C.self_encoded_buffer_len(*encodedCredentialBufferPtr)),
	)

	C.self_encoded_buffer_destroy(
		*encodedCredentialBufferPtr,
	)

	return encodedCredential, nil
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
	collection := C.self_collection_verifiable_credential_init()

	runtime.SetFinalizer(&collection, func(collection **C.self_collection_verifiable_credential) {
		C.self_collection_verifiable_credential_destroy(
			*collection,
		)
	})

	return (*VerifiableCredentialCollection)(collection)
}

func (c *VerifiableCredentialCollection) Length() int {
	return int(C.self_collection_verifiable_credential_len(
		(*C.self_collection_verifiable_credential)(c),
	))
}

func (c *VerifiableCredentialCollection) Get(index int) *VerifiableCredential {
	credential := C.self_collection_verifiable_credential_at(
		(*C.self_collection_verifiable_credential)(c),
		C.ulong(index),
	)

	runtime.SetFinalizer(&credential, func(credential **C.self_verifiable_credential) {
		C.self_verifiable_credential_destroy(
			*credential,
		)
	})

	return (*VerifiableCredential)(credential)
}

func NewCredentialVerificationEvidenceCollection() *CredentialVerificationEvidenceCollection {
	collection := C.self_collection_credential_verification_evidence_init()

	runtime.SetFinalizer(&collection, func(collection **C.self_collection_credential_verification_evidence) {
		C.self_collection_credential_verification_evidence_destroy(
			*collection,
		)
	})

	return (*CredentialVerificationEvidenceCollection)(collection)
}

func (c *CredentialVerificationEvidenceCollection) Length() int {
	return int(C.self_collection_credential_verification_evidence_len(
		(*C.self_collection_credential_verification_evidence)(c),
	))
}

func (c *CredentialVerificationEvidenceCollection) Get(index int) *CredentialVerificationEvidence {
	evidence := C.self_collection_credential_verification_evidence_at(
		(*C.self_collection_credential_verification_evidence)(c),
		C.ulong(index),
	)

	runtime.SetFinalizer(&evidence, func(evidence **C.self_credential_verification_evidence) {
		C.self_credential_verification_evidence_destroy(
			*evidence,
		)
	})

	return (*CredentialVerificationEvidence)(evidence)
}

func NewCredentialVerificationParameterCollection() *CredentialVerificationParameterCollection {
	collection := C.self_collection_credential_verification_parameter_init()

	runtime.SetFinalizer(&collection, func(collection **C.self_collection_credential_verification_parameter) {
		C.self_collection_credential_verification_parameter_destroy(
			*collection,
		)
	})

	return (*CredentialVerificationParameterCollection)(collection)
}

func (c *CredentialVerificationParameterCollection) Length() int {
	return int(C.self_collection_credential_verification_parameter_len(
		(*C.self_collection_credential_verification_parameter)(c),
	))
}

func (c *CredentialVerificationParameterCollection) Get(index int) *CredentialVerificationParameter {
	parameter := C.self_collection_credential_verification_parameter_at(
		(*C.self_collection_credential_verification_parameter)(c),
		C.ulong(index),
	)

	runtime.SetFinalizer(&parameter, func(parameter **C.self_credential_verification_parameter) {
		C.self_credential_verification_parameter_destroy(
			*parameter,
		)
	})

	return (*CredentialVerificationParameter)(parameter)
}

func NewCredentialTypeCollection() *CredentialTypeCollection {
	collection := C.self_collection_credential_type_init()

	runtime.SetFinalizer(&collection, func(collection **C.self_collection_credential_type) {
		C.self_collection_credential_type_destroy(
			*collection,
		)
	})

	return (*CredentialTypeCollection)(collection)
}

func newCredentialTypeCollection() *CredentialTypeCollection {
	return (*CredentialTypeCollection)(C.self_collection_credential_type_init())
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

// EvidenceType returns the evidence type
func (c *CredentialVerificationEvidence) EvidenceType() string {
	return C.GoString(
		C.self_credential_verification_evidence_evidence_type(
			(*C.self_credential_verification_evidence)(c),
		),
	)
}

// Object returns the object that makes up the content of the evidence
func (c *CredentialVerificationEvidence) Object() *object.Object {
	obj := C.self_credential_verification_evidence_object(
		(*C.self_credential_verification_evidence)(c),
	)

	runtime.SetFinalizer(&obj, func(obj **C.self_object) {
		C.self_object_destroy(
			*obj,
		)
	})

	return (*object.Object)(obj)
}

// ParameterType returns the parameter type
func (c *CredentialVerificationParameter) ParameterType() string {
	return C.GoString(
		C.self_credential_verification_parameter_parameter_type(
			(*C.self_credential_verification_parameter)(c),
		),
	)
}

// Value returns the value of the parameter
func (c *CredentialVerificationParameter) Value() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_credential_verification_parameter_value_buf(
			(*C.self_credential_verification_parameter)(c),
		)),
		(C.int)(C.self_credential_verification_parameter_value_len(
			(*C.self_credential_verification_parameter)(c),
		)),
	)
}
