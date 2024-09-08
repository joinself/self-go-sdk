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
	CredentialTypeEmail                = []string{"VerifiableCredential", "EmailCredential"}
	CredentialTypePassport             = []string{"VerifiableCredential", "PassportCredential"}
	CredentialTypeLiveness             = []string{"VerifiableCredential", "LivenessCredential"}
	CredentialTypeProfileName          = []string{"VerifiableCredential", "ProfileNameCredential"}
	CredentialTypeProfileImage         = []string{"VerifiableCredential", "ProfileImageCredential"}
	CredentialTypeApplicationPublisher = []string{"VerifiableCredential", "ApplicationPublisherCredential"}
)

//go:linkname newObject object.newObject
func newObject(ptr *C.self_object) *object.Object

// Credential an unsigned credential
type Credential struct {
	ptr *C.self_credential
}

func newCredential(ptr *C.self_credential) *Credential {
	c := &Credential{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *Credential) {
		C.self_credential_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialBuilder struct {
	ptr *C.self_credential_builder
}

func newCredentialBuilder(ptr *C.self_credential_builder) *CredentialBuilder {
	b := &CredentialBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(b, func(b *CredentialBuilder) {
		C.self_credential_builder_destroy(
			b.ptr,
		)
	})

	return b
}

type VerifiableCredential struct {
	ptr *C.self_verifiable_credential
}

func newVerifiableCredential(ptr *C.self_verifiable_credential) *VerifiableCredential {
	c := &VerifiableCredential{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *VerifiableCredential) {
		C.self_verifiable_credential_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialPresentationDetail struct {
	ptr *C.self_credential_presentation_detail
}

func newCredentialPresentationDetail(ptr *C.self_credential_presentation_detail) *CredentialPresentationDetail {
	c := &CredentialPresentationDetail{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialPresentationDetail) {
		// TODO add this
		/*
			C.self_credential_presentation_detail_destroy(
				c.ptr,
			)
		*/
	})

	return c
}

type CredentialVerificationEvidence struct {
	ptr *C.self_credential_verification_evidence
}

func newCredentialVerificationEvidence(ptr *C.self_credential_verification_evidence) *CredentialVerificationEvidence {
	c := &CredentialVerificationEvidence{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialVerificationEvidence) {
		C.self_credential_verification_evidence_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialVerificationParameter struct {
	ptr *C.self_credential_verification_parameter
}

func newCredentialVerificationParameter(ptr *C.self_credential_verification_parameter) *CredentialVerificationParameter {
	c := &CredentialVerificationParameter{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialVerificationParameter) {
		C.self_credential_verification_parameter_destroy(
			c.ptr,
		)
	})

	return c
}

// NewCredential creates a new credential builder
func NewCredential() *CredentialBuilder {
	return newCredentialBuilder(C.self_credential_builder_init())
}

// CredentialType sets the type of credential
func (b *CredentialBuilder) CredentialType(credentialType []string) *CredentialBuilder {
	collection := toCredentialTypeCollection(credentialType)

	C.self_credential_builder_credential_type(
		b.ptr,
		collection,
	)

	C.self_collection_credential_type_destroy(
		collection,
	)

	return b
}

// CredentialSubject sets the address of the credential's subject
func (b *CredentialBuilder) CredentialSubject(subjectAddress *Address) *CredentialBuilder {
	C.self_credential_builder_credential_subject(
		b.ptr,
		subjectAddress.ptr,
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
		b.ptr,
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
		b.ptr,
		(*C.uint8_t)(claimBuffer),
		(C.size_t)(claimLength),
	)

	return b
}

// Issuer sets the address of the credential's issuer
func (b *CredentialBuilder) Issuer(issuerAddress *Address) *CredentialBuilder {
	C.self_credential_builder_issuer(
		b.ptr,
		issuerAddress.ptr,
	)
	return b
}

// ValidFrom sets the point of validity for the credential
func (b *CredentialBuilder) ValidFrom(timestamp time.Time) *CredentialBuilder {
	C.self_credential_builder_valid_from(
		b.ptr,
		C.int64_t(timestamp.Unix()),
	)
	return b
}

// SignWith signs the credential
func (b *CredentialBuilder) SignWith(signer *signing.PublicKey, issuedAt time.Time) *CredentialBuilder {
	C.self_credential_builder_sign_with(
		b.ptr,
		signingPublicKeyPtr(signer),
		C.int64_t(issuedAt.Unix()),
	)
	return b
}

// Finish generates and prepares the credential for being signed by an account
func (b *CredentialBuilder) Finish() (*Credential, error) {
	var credential *C.self_credential

	status := C.self_credential_builder_finish(
		b.ptr,
		&credential,
	)

	if status > 0 {
		return nil, errors.New("failed to create credential")
	}

	return newCredential(credential), nil
}

func DecodeVerifiableCredential(encodedCredential []byte) (*VerifiableCredential, error) {
	var verifiableCredential *C.self_verifiable_credential

	encodedBuf := C.CBytes(encodedCredential)
	encodedLen := len(encodedCredential)

	defer func() {
		C.free(encodedBuf)
	}()

	status := C.self_verifiable_credential_decode(
		&verifiableCredential,
		(*C.uint8_t)(encodedBuf),
		(C.size_t)(encodedLen),
	)

	if status > 0 {
		return nil, errors.New("decode credential failed")
	}

	return newVerifiableCredential(verifiableCredential), nil
}

// CredentialType returns the type of credential
func (c *VerifiableCredential) CredentialType() []string {
	collection := C.self_verifiable_credential_credential_type(
		c.ptr,
	)

	credentials := fromCredentialTypeCollection(collection)

	C.self_collection_credential_type_destroy(
		collection,
	)

	return credentials
}

// CredentialSubject returns the subject of the credential's address
func (c *VerifiableCredential) CredentialSubject() *Address {
	return newAddress(C.self_verifiable_credential_credential_subject(
		c.ptr,
	))
}

// CredentialSubject returns the subject of the credential's address
func (c *VerifiableCredential) CredentialSubjectClaim(claimKey string) (string, bool) {
	key := C.CString(claimKey)

	value := C.self_verifiable_credential_credential_subject_claim(
		c.ptr,
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
	return newAddress(C.self_verifiable_credential_issuer(
		c.ptr,
	))
}

// ValidFrom returns the time period that the credential is valid from
func (c *VerifiableCredential) ValidFrom() time.Time {
	validFrom := C.self_verifiable_credential_valid_from(
		c.ptr,
	)

	return time.Unix(int64(validFrom), 0)
}

// Created returns the time that the credential was created
func (c *VerifiableCredential) Created() time.Time {
	created := C.self_verifiable_credential_created(
		c.ptr,
	)

	return time.Unix(int64(created), 0)
}

// Encode returns a json encoded verifiable credential
func (c *VerifiableCredential) Encode() ([]byte, error) {
	var encodedCredentialBuffer *C.self_encoded_buffer
	encodedCredentialBufferPtr := &encodedCredentialBuffer

	status := C.self_verifiable_credential_encode(
		c.ptr,
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
		c.ptr,
	)

	if status > 0 {
		return errors.New("credential invalid")
	}

	return nil
}

// EvidenceType returns the evidence type
func (c *CredentialVerificationEvidence) EvidenceType() string {
	return C.GoString(
		C.self_credential_verification_evidence_evidence_type(
			c.ptr,
		),
	)
}

// Object returns the object that makes up the content of the evidence
func (c *CredentialVerificationEvidence) Object() *object.Object {
	return newObject(C.self_credential_verification_evidence_object(
		c.ptr,
	))
}

// ParameterType returns the parameter type
func (c *CredentialVerificationParameter) ParameterType() string {
	return C.GoString(
		C.self_credential_verification_parameter_parameter_type(
			c.ptr,
		),
	)
}

// Value returns the value of the parameter
func (c *CredentialVerificationParameter) Value() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_credential_verification_parameter_value_buf(
			c.ptr,
		)),
		(C.int)(C.self_credential_verification_parameter_value_len(
			c.ptr,
		)),
	)
}

func toCredentialTypeCollection(credentialType []string) *C.self_collection_credential_type {
	collection := C.self_collection_credential_type_init()

	for i := 0; i < len(credentialType); i++ {
		typ := C.CString(credentialType[i])

		C.self_collection_credential_type_append(
			collection,
			typ,
		)

		C.free(unsafe.Pointer(typ))
	}

	return collection
}

func fromCredentialTypeCollection(collection *C.self_collection_credential_type) []string {
	collectionLen := int(C.self_collection_credential_type_len(
		collection,
	))

	credentialType := make([]string, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_credential_type_at(
			collection,
			C.size_t(i),
		)

		credentialType[i] = C.GoString(ptr)
	}

	return credentialType
}

func fromVerifiableCredentialCollection(collection *C.self_collection_verifiable_credential) []*VerifiableCredential {
	collectionLen := int(C.self_collection_verifiable_credential_len(
		collection,
	))

	credentials := make([]*VerifiableCredential, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_verifiable_credential_at(
			collection,
			C.size_t(i),
		)

		credentials[i] = newVerifiableCredential(ptr)
	}

	return credentials
}

func fromPresentationDetailCollection(collection *C.self_collection_credential_presentation_detail) []*CredentialPresentationDetail {
	collectionLen := int(C.self_collection_credential_presentation_detail_len(
		collection,
	))

	details := make([]*CredentialPresentationDetail, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_credential_presentation_detail_at(
			collection,
			C.size_t(i),
		)

		details[i] = newCredentialPresentationDetail(ptr)
	}

	return details
}
