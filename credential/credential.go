package credential

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"runtime"
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/status"
)

var (
	CredentialTypeEmail              = []string{"VerifiableCredential", "EmailCredential"}
	CredentialTypePhone              = []string{"VerifiableCredential", "PhoneCredential"}
	CredentialTypePassport           = []string{"VerifiableCredential", "PassportCredential"}
	CredentialTypeLiveness           = []string{"VerifiableCredential", "LivenessCredential"}
	CredentialTypeProfileName        = []string{"VerifiableCredential", "ProfileNameCredential"}
	CredentialTypeProfileImage       = []string{"VerifiableCredential", "ProfileImageCredential"}
	CredentialTypeOrganisation       = []string{"VerifiableCredential", "OrganisationCredential"}
	CredentialTypeApplication        = []string{"VerifiableCredential", "ApplicationCredential"}
	CredentialTypeApplicationInstall = []string{"VerifiableCredential", "ApplicationInstallCredential"}
)

//go:linkname newObject github.com/joinself/self-go-sdk/object.newObject
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

func credentialPtr(c *Credential) *C.self_credential {
	return c.ptr
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

func verifiableCredentialPtr(v *VerifiableCredential) *C.self_verifiable_credential {
	return v.ptr
}

type CredentialPresentationDetail struct {
	ptr *C.self_credential_presentation_detail
}

func newCredentialPresentationDetail(ptr *C.self_credential_presentation_detail) *CredentialPresentationDetail {
	c := &CredentialPresentationDetail{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialPresentationDetail) {
		C.self_credential_presentation_detail_destroy(
			c.ptr,
		)
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
	if err != nil {
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

	result := C.self_credential_builder_finish(
		b.ptr,
		&credential,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCredential(credential), nil
}

// DecodeVerifiableCredential decodes a verifiable credential from it's json form
func DecodeVerifiableCredential(encodedCredential []byte) (*VerifiableCredential, error) {
	var verifiableCredential *C.self_verifiable_credential

	encodedBuf := C.CBytes(encodedCredential)
	encodedLen := len(encodedCredential)

	defer func() {
		C.free(encodedBuf)
	}()

	result := C.self_verifiable_credential_decode(
		&verifiableCredential,
		(*C.uint8_t)(encodedBuf),
		(C.size_t)(encodedLen),
	)

	if result > 0 {
		return nil, status.New(result)
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

// CredentialSubjectClaim returns the one of the claims about the subject
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

// CredentialSubjectClaims returns all of the claims about the subject
func (c *VerifiableCredential) CredentialSubjectClaims() (map[string]interface{}, error) {
	var claims map[string]interface{}

	value := C.self_verifiable_credential_credential_subject_json(
		c.ptr,
	)

	claimValue := C.GoBytes(
		unsafe.Pointer(C.self_credential_claim_value_buf(value)),
		C.int(C.self_credential_claim_value_len(value)),
	)

	C.self_credential_claim_value_destroy(value)

	return claims, json.Unmarshal(claimValue, &claims)
}

// Issuer returns the address of the credential's issuer
func (c *VerifiableCredential) Issuer() *Address {
	return newAddress(C.self_verifiable_credential_issuer(
		c.ptr,
	))
}

// Signer returns the address of the credential's signer
func (c *VerifiableCredential) Signer() (*Address, error) {
	var signerAddress *C.self_credential_address

	result := C.self_verifiable_credential_signer(
		c.ptr,
		&signerAddress,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newAddress(signerAddress), nil
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
	var encodedCredentialBuffer *C.self_bytes_buffer
	encodedCredentialBufferPtr := &encodedCredentialBuffer

	result := C.self_verifiable_credential_encode(
		c.ptr,
		encodedCredentialBufferPtr,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	encodedCredential := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(*encodedCredentialBufferPtr)),
		C.int(C.self_bytes_buffer_len(*encodedCredentialBufferPtr)),
	)

	C.self_bytes_buffer_destroy(
		*encodedCredentialBufferPtr,
	)

	return encodedCredential, nil
}

// Validate validates the contents of the credential and it's signatures
func (c *VerifiableCredential) Validate() error {
	result := C.self_verifiable_credential_validate(
		c.ptr,
	)

	if result > 0 {
		return status.New(result)
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

// Key returns the parameters key
func (c *CredentialVerificationParameter) Key() string {
	return C.GoString(
		C.self_credential_verification_parameter_parameter_key(
			c.ptr,
		),
	)
}

// Value returns the value of the parameter
func (c *CredentialVerificationParameter) Value() any {
	ptr := C.self_credential_verification_parameter_value(
		c.ptr,
	)

	defer func() {
		C.self_message_content_parameter_value_destroy(ptr)
	}()

	switch C.self_message_content_parameter_value_value_type(
		ptr,
	) {
	case C.PARAMETER_VALUE_BYTES:
		buf := C.self_message_content_parameter_value_as_bytes(
			ptr,
		)

		bytes := C.GoBytes(
			unsafe.Pointer(C.self_bytes_buffer_buf(
				buf,
			)),
			C.int(C.self_bytes_buffer_len(
				buf,
			)),
		)

		C.self_bytes_buffer_destroy(
			buf,
		)

		return bytes
	case C.PARAMETER_VALUE_STRING:
		buf := C.self_message_content_parameter_value_as_string(
			ptr,
		)

		str := C.GoString(
			C.self_string_buffer_ptr(
				buf,
			),
		)

		C.self_string_buffer_destroy(
			buf,
		)

		return str
	case C.PARAMETER_VALUE_INT8:
		return C.self_message_content_parameter_value_as_int8(
			ptr,
		)
	case C.PARAMETER_VALUE_INT16:
		return C.self_message_content_parameter_value_as_int16(
			ptr,
		)
	case C.PARAMETER_VALUE_INT32:
		return C.self_message_content_parameter_value_as_int32(
			ptr,
		)
	case C.PARAMETER_VALUE_INT64:
		return C.self_message_content_parameter_value_as_int64(
			ptr,
		)
	case C.PARAMETER_VALUE_UINT8:
		return C.self_message_content_parameter_value_as_uint8(
			ptr,
		)
	case C.PARAMETER_VALUE_UINT16:
		return C.self_message_content_parameter_value_as_uint16(
			ptr,
		)
	case C.PARAMETER_VALUE_UINT32:
		return C.self_message_content_parameter_value_as_uint32(
			ptr,
		)
	case C.PARAMETER_VALUE_UINT64:
		return C.self_message_content_parameter_value_as_uint64(
			ptr,
		)
	case C.PARAMETER_VALUE_FLOAT32:
		return C.self_message_content_parameter_value_as_float32(
			ptr,
		)
	case C.PARAMETER_VALUE_FLOAT64:
		return C.self_message_content_parameter_value_as_float64(
			ptr,
		)
	case C.PARAMETER_VALUE_ARRAY_BYTES:
		collection := C.self_message_content_parameter_value_as_array_bytes(
			ptr,
		)

		collectionLen := int(C.self_collection_bytes_buffer_len(collection))

		values := make([][]byte, collectionLen)

		for i := 0; i < collectionLen; i++ {
			buf := C.self_collection_bytes_buffer_at(collection, C.ulong(i))

			values[i] = C.GoBytes(
				unsafe.Pointer(C.self_bytes_buffer_buf(
					buf,
				)),
				C.int(C.self_bytes_buffer_len(
					buf,
				)),
			)

			C.self_bytes_buffer_destroy(
				buf,
			)
		}

		C.self_collection_bytes_buffer_destroy(
			collection,
		)

		return values
	case C.PARAMETER_VALUE_ARRAY_STRING:
		collection := C.self_message_content_parameter_value_as_array_string(
			ptr,
		)

		collectionLen := int(C.self_collection_string_buffer_len(collection))

		values := make([]string, collectionLen)

		for i := 0; i < collectionLen; i++ {
			buf := C.self_collection_string_buffer_at(collection, C.ulong(i))

			values[i] = C.GoString(C.self_string_buffer_ptr(
				buf,
			))

			C.self_string_buffer_destroy(
				buf,
			)
		}

		C.self_collection_string_buffer_destroy(
			collection,
		)

		return values
	default:
		return nil
	}
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

func fromCredentialVerificationEvidenceCollection(collection *C.self_collection_credential_verification_evidence) []*CredentialVerificationEvidence {
	collectionLen := int(C.self_collection_credential_verification_evidence_len(
		collection,
	))

	details := make([]*CredentialVerificationEvidence, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_credential_verification_evidence_at(
			collection,
			C.size_t(i),
		)

		details[i] = newCredentialVerificationEvidence(ptr)
	}

	return details
}

func fromCredentialVerificationParameterCollection(collection *C.self_collection_credential_verification_parameter) []*CredentialVerificationParameter {
	collectionLen := int(C.self_collection_credential_verification_parameter_len(
		collection,
	))

	details := make([]*CredentialVerificationParameter, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_credential_verification_parameter_at(
			collection,
			C.size_t(i),
		)

		details[i] = newCredentialVerificationParameter(ptr)
	}

	return details
}
