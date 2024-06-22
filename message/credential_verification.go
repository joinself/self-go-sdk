package message

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/object"
)

type CredentialVerificationRequest C.self_message_content_credential_verification_request
type CredentialVerificationResponse C.self_message_content_credential_verification_response
type CredentialVerificationRequestBuilder C.self_message_content_credential_verification_request_builder
type CredentialVerificationResponseBuilder C.self_message_content_credential_verification_response_builder

// DecodeCredentialVerificationRequest decodes a message to a credential verification request
func DecodeCredentialVerificationRequest(msg *Message) (*CredentialVerificationRequest, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var credentialVerificationRequestContent *C.self_message_content_credential_verification_request

	status := C.self_message_content_as_credential_verification_request(
		content,
		&credentialVerificationRequestContent,
	)

	if status > 0 {
		return nil, errors.New("failed to decode credential verification request message")
	}

	credentialVerificationRequest := (*CredentialVerificationRequest)(credentialVerificationRequestContent)

	runtime.SetFinalizer(credentialVerificationRequest, func(credentialVerificationRequest *CredentialVerificationRequest) {
		C.self_message_content_credential_verification_request_destroy(
			(*C.self_message_content_credential_verification_request)(credentialVerificationRequest),
		)
	})

	return credentialVerificationRequest, nil
}

// Type returns the type of credential that verification is being requested for
func (c *CredentialVerificationRequest) Type() credential.CredentialType {
	return credential.CredentialType(C.self_message_content_credential_verification_request_credential_type(
		(*C.self_message_content_credential_verification_request)(c),
	))
}

// Proof returns associated verifiable credential proof to support the verification request
func (c *CredentialVerificationRequest) Proof() *credential.VerifiableCredentialCollection {
	return (*credential.VerifiableCredentialCollection)(C.self_message_content_credential_verification_request_proof(
		(*C.self_message_content_credential_verification_request)(c),
	))
}

// Evidence returns associated data to be used as evidence to support the verification request
func (c *CredentialVerificationRequest) Evidence() *credential.CredentialVerificationEvidenceCollection {
	return (*credential.CredentialVerificationEvidenceCollection)(C.self_message_content_credential_verification_request_evidence(
		(*C.self_message_content_credential_verification_request)(c),
	))
}

// Type returns the time the request expires at
func (c *CredentialVerificationRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_credential_verification_request_expires(
		(*C.self_message_content_credential_verification_request)(c),
	)), 0)
}

// NewCredentialVerificationRequest creates a new credential verification request
func NewCredentialVerificationRequest() *CredentialVerificationRequestBuilder {
	builder := (*CredentialVerificationRequestBuilder)(C.self_message_content_credential_verification_request_builder_init())

	runtime.SetFinalizer(builder, func(builder *CredentialVerificationRequestBuilder) {
		C.self_message_content_credential_verification_request_builder_destroy(
			(*C.self_message_content_credential_verification_request_builder)(builder),
		)
	})

	return builder
}

// Type sets the type of credential being requested
func (b *CredentialVerificationRequestBuilder) Type(credentialType credential.CredentialType) *CredentialVerificationRequestBuilder {
	C.self_message_content_credential_verification_request_builder_credential_type(
		(*C.self_message_content_credential_verification_request_builder)(b),
		uint32(credentialType),
	)
	return b
}

// Proof attaches proof to the credential verification request
func (b *CredentialVerificationRequestBuilder) Proof(proof *credential.VerifiableCredential) *CredentialVerificationRequestBuilder {
	C.self_message_content_credential_verification_request_builder_proof(
		(*C.self_message_content_credential_verification_request_builder)(b),
		(*C.self_verifiable_credential)(proof),
	)
	return b
}

// Evidence attaches evidence to the credential verification request
func (b *CredentialVerificationRequestBuilder) Evidence(evidenceType string, evidence *object.Object) *CredentialVerificationRequestBuilder {
	evidenceTypeC := C.CString(evidenceType)

	C.self_message_content_credential_verification_request_builder_evidence(
		(*C.self_message_content_credential_verification_request_builder)(b),
		evidenceTypeC,
		(*C.self_object)(evidence),
	)

	C.free(unsafe.Pointer(evidenceTypeC))

	return b
}

// Expires sets the time that the request expires at
func (b *CredentialVerificationRequestBuilder) Expires(expires time.Time) *CredentialVerificationRequestBuilder {
	C.self_message_content_credential_verification_request_builder_expires(
		(*C.self_message_content_credential_verification_request_builder)(b),
		C.long(expires.Unix()),
	)
	return b
}

// Finish finalises the request and builds the content
func (b *CredentialVerificationRequestBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	status := C.self_message_content_credential_verification_request_builder_finish(
		(*C.self_message_content_credential_verification_request_builder)(b),
		&finishedContent,
	)

	if status > 0 {
		return nil, errors.New("failed to build credential verificaiton request")
	}

	content := (*Content)(finishedContent)

	runtime.SetFinalizer(content, func(content *Content) {
		C.self_message_content_destroy(
			(*C.self_message_content)(content),
		)
	})

	return content, nil
}

// DecodeCredentialVerificationResponse decodes a message to a credential verification response
func DecodeCredentialVerificationResponse(msg *Message) (*CredentialVerificationResponse, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var credentialVerificationResponseContent *C.self_message_content_credential_verification_response

	status := C.self_message_content_as_credential_verification_response(
		content,
		&credentialVerificationResponseContent,
	)

	if status > 0 {
		return nil, errors.New("failed to decode credential verification response message")
	}

	credentialVerificationResponse := (*CredentialVerificationResponse)(credentialVerificationResponseContent)

	runtime.SetFinalizer(credentialVerificationResponse, func(credentialVerificationResponse *CredentialVerificationResponse) {
		C.self_message_content_credential_verification_response_destroy(
			(*C.self_message_content_credential_verification_response)(credentialVerificationResponse),
		)
	})

	return credentialVerificationResponse, nil
}

// Status returns the status of the request
func (c *CredentialVerificationResponse) Status() int {
	return 0
}

// Credentials returns verified credentials that have been asserted by the responder
func (c *CredentialVerificationResponse) Credentials() *credential.VerifiableCredentialCollection {
	return (*credential.VerifiableCredentialCollection)(C.self_message_content_credential_verification_response_verifiable_credentials(
		(*C.self_message_content_credential_verification_response)(c),
	))
}

// NewCredentialVerificationResponse creates a new credential verification response
func NewCredentialVerificationResponse() *CredentialVerificationResponseBuilder {
	builder := (*CredentialVerificationResponseBuilder)(C.self_message_content_credential_verification_response_builder_init())

	runtime.SetFinalizer(builder, func(builder *CredentialVerificationResponseBuilder) {
		C.self_message_content_credential_verification_response_builder_destroy(
			(*C.self_message_content_credential_verification_response_builder)(builder),
		)
	})

	return builder
}

// VerifiableCredential attaches a verified credential to the response
func (b *CredentialVerificationResponseBuilder) VerifiableCredential(proof *credential.VerifiableCredential) *CredentialVerificationResponseBuilder {
	C.self_message_content_credential_verification_response_builder_verifiable_credential(
		(*C.self_message_content_credential_verification_response_builder)(b),
		(*C.self_verifiable_credential)(proof),
	)
	return b
}

// Finish finalises the response and builds the content
func (b *CredentialVerificationResponseBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	status := C.self_message_content_credential_verification_response_builder_finish(
		(*C.self_message_content_credential_verification_response_builder)(b),
		&finishedContent,
	)

	if status > 0 {
		return nil, errors.New("failed to build credential verificaiton response")
	}

	content := (*Content)(finishedContent)

	runtime.SetFinalizer(content, func(content *Content) {
		C.self_message_content_destroy(
			(*C.self_message_content)(content),
		)
	})

	return content, nil
}
