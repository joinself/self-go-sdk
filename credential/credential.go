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

const (
	CredentialTypeEmail                   = "EmailCredential"
	CredentialTypePhone                   = "PhoneCredential"
	CredentialTypePassport                = "PassportCredential"
	CredentialTypeLiveness                = "LivenessCredential"
	CredentialTypeProfileName             = "ProfileNameCredential"
	CredentialTypeProfileImage            = "ProfileImageCredential"
	CredentialTypeOrganisation            = "OrganisationCredential"
	CredentialTypeApplication             = "ApplicationCredential"
	CredentialTypeApplicationInstall      = "ApplicationInstallCredential"
	FieldType                             = "/type"
	FieldIssuer                           = "/issuer"
	FieldValidFrom                        = "/validFrom"
	FieldValidUntil                       = "/validUntil"
	FieldSubject                          = "/credentialSubject/id"
	FieldSubjectClaims                    = "/credentialSubject"
	FieldSubjectEmailAddress              = "/credentialSubject/emailAddress"
	FieldSubjectPhoneNumber               = "/credentialSubject/phoneNumber"
	FieldSubjectLivenessSourceImageHash   = "/credentialSubject/sourceImageHash"
	FieldSubjectLivenessTargetImageHash   = "/credentialSubject/targetImageHash"
	FieldSubjectPassportDocumentNumber    = "/credentialSubject/documentNumber"
	FieldSubjectPassportGivenNames        = "/credentialSubject/givenNames"
	FieldSubjectPassportSurname           = "/credentialSubject/surname"
	FieldSubjectPassportSex               = "/credentialSubject/sex"
	FieldSubjectPassportNationality       = "/credentialSubject/nationality"
	FieldSubjectPassportDateOfBirth       = "/credentialSubject/dateOfBirth"
	FieldSubjectPassportDateOfExpiration  = "/credentialSubject/dateOfExpiration"
	FieldSubjectPassportCountryOfIssuance = "/credentialSubject/countryOfIssuance"
	FieldSubjectPassportDocumentMrz       = "/credentialSubject/mrz"
	FieldSubjectPassportImageType         = "/credentialSubject/imageType"
	FieldSubjectPassportimageHash         = "/credentialSubject/imageHash"
	FieldSubjectOrganisationName          = "/credentialSubject/organisationName"
	FieldSubjectApplicationName           = "/credentialSubject/applicationName"
	FieldSubjectApplicationSubsidiaryOf   = "/credentialSubject/subsidiaryOf"
)

// DateTime returns a verifiable credential date time
func DateTime(instant time.Time) string {
	return instant.UTC().Format("2006-01-02T15:04:05Z07:00")
}

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

// NewCredential creates a new credential builder
func NewCredential() *CredentialBuilder {
	return newCredentialBuilder(C.self_credential_builder_init())
}

// CredentialType sets the type of credential
func (b *CredentialBuilder) CredentialType(credentialType ...string) *CredentialBuilder {
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

// ValidUntil sets the end of the validity period for the credential
func (b *CredentialBuilder) ValidUntil(timestamp time.Time) *CredentialBuilder {
	C.self_credential_builder_valid_until(
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
	collection := C.self_verifiable_credential_type_of(
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

	claimValue := C.GoString(C.self_string_buffer_ptr(value))

	C.self_string_buffer_destroy(value)

	return claimValue, true
}

// CredentialSubjectClaims returns all of the claims about the subject
func (c *VerifiableCredential) CredentialSubjectClaims() (map[string]interface{}, error) {
	var claims map[string]interface{}

	value := C.self_verifiable_credential_credential_subject_json(
		c.ptr,
	)

	claimValue := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(value)),
		C.int(C.self_bytes_buffer_len(value)),
	)

	C.self_bytes_buffer_destroy(value)

	return claims, json.Unmarshal(claimValue, &claims)
}

// CredentialHash returns a hash of the complete verifiable credential
func (c *VerifiableCredential) CredentialHash() ([]byte, error) {
	hash := C.CBytes(make([]byte, 32))
	defer C.free(hash)

	result := C.self_verifiable_credential_credential_hash(
		c.ptr,
		(*C.uint8_t)(hash),
		32,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return C.GoBytes(
		unsafe.Pointer(hash),
		32,
	), nil
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

// PrimaryType returns the primary type of a credential or presentation
func PrimaryType(typ []string) string {
	return typ[1]
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

func toVerifiableCredentialCollection(credentials []*VerifiableCredential) *C.self_collection_verifiable_credential {
	collection := C.self_collection_verifiable_credential_init()

	for i := 0; i < len(credentials); i++ {
		C.self_collection_verifiable_credential_append(
			collection,
			verifiableCredentialPtr(credentials[i]),
		)
	}

	return collection
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
