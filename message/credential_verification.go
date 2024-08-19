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

	"github.com/joinself/self-go-sdk-next/credential"
	"github.com/joinself/self-go-sdk-next/object"
)

type CredentialVerificationRequest C.self_message_content_credential_verification_request
type CredentialVerificationResponse C.self_message_content_credential_verification_response
type CredentialVerificationRequestBuilder C.self_message_content_credential_verification_request_builder
type CredentialVerificationResponseBuilder C.self_message_content_credential_verification_response_builder

// DecodeCredentialVerificationRequest decodes a message to a credential verification request
func DecodeCredentialVerificationRequest(msg *Message) (*CredentialVerificationRequest, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var credentialVerificationRequestContent *C.self_message_content_credential_verification_request
	credentialVerificationRequestContentPtr := &credentialVerificationRequestContent

	status := C.self_message_content_as_credential_verification_request(
		content,
		credentialVerificationRequestContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to decode credential verification request message")
	}

	runtime.SetFinalizer(credentialVerificationRequestContentPtr, func(credentialVerificationRequest **C.self_message_content_credential_verification_request) {
		C.self_message_content_credential_verification_request_destroy(
			*credentialVerificationRequest,
		)
	})

	return (*CredentialVerificationRequest)(*credentialVerificationRequestContentPtr), nil
}

// Type returns the type of credential that verification is being requested for
func (c *CredentialVerificationRequest) Type() *credential.CredentialTypeCollection {
	collection := C.self_message_content_credential_verification_request_credential_type(
		(*C.self_message_content_credential_verification_request)(c),
	)

	runtime.SetFinalizer(collection, func(collection *C.self_collection_credential_type) {
		C.self_collection_credential_type_destroy(
			collection,
		)
	})

	return (*credential.CredentialTypeCollection)(collection)
}

// Proof returns associated verifiable credential proof to support the verification request
func (c *CredentialVerificationRequest) Proof() *credential.VerifiableCredentialCollection {
	collection := C.self_message_content_credential_verification_request_proof(
		(*C.self_message_content_credential_verification_request)(c),
	)

	return (*credential.VerifiableCredentialCollection)(collection)
}

// Evidence returns associated data to be used as evidence to support the verification request
func (c *CredentialVerificationRequest) Evidence() *credential.CredentialVerificationEvidenceCollection {
	collection := C.self_message_content_credential_verification_request_evidence(
		(*C.self_message_content_credential_verification_request)(c),
	)

	runtime.SetFinalizer(collection, func(collection *C.self_collection_credential_verification_evidence) {
		C.self_collection_credential_verification_evidence_destroy(
			collection,
		)
	})

	return (*credential.CredentialVerificationEvidenceCollection)(collection)
}

// Parameters returns associated data to be used as parameters to support the verification request
func (c *CredentialVerificationRequest) Parameters() *credential.CredentialVerificationParameterCollection {
	collection := C.self_message_content_credential_verification_request_parameters(
		(*C.self_message_content_credential_verification_request)(c),
	)

	runtime.SetFinalizer(collection, func(collection *C.self_collection_credential_verification_parameter) {
		C.self_collection_credential_verification_parameter_destroy(
			collection,
		)
	})

	return (*credential.CredentialVerificationParameterCollection)(collection)
}

// Type returns the time the request expires at
func (c *CredentialVerificationRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_credential_verification_request_expires(
		(*C.self_message_content_credential_verification_request)(c),
	)), 0)
}

// NewCredentialVerificationRequest creates a new credential verification request
func NewCredentialVerificationRequest() *CredentialVerificationRequestBuilder {
	builder := C.self_message_content_credential_verification_request_builder_init()

	runtime.SetFinalizer(builder, func(builder *C.self_message_content_credential_verification_request_builder) {
		C.self_message_content_credential_verification_request_builder_destroy(
			builder,
		)
	})

	return (*CredentialVerificationRequestBuilder)(builder)
}

// Type sets the type of credential being requested
func (b *CredentialVerificationRequestBuilder) Type(credentialType *credential.CredentialTypeCollection) *CredentialVerificationRequestBuilder {
	C.self_message_content_credential_verification_request_builder_credential_type(
		(*C.self_message_content_credential_verification_request_builder)(b),
		(*C.self_collection_credential_type)(credentialType),
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

// Evidence attaches evidence to the credential verification request
func (b *CredentialVerificationRequestBuilder) Parameter(parameterType string, value []byte) *CredentialVerificationRequestBuilder {
	parameterTypeC := C.CString(parameterType)
	valueBuf := C.CBytes(value)
	valueLen := len(value)

	C.self_message_content_credential_verification_request_builder_parameter(
		(*C.self_message_content_credential_verification_request_builder)(b),
		parameterTypeC,
		(*C.uint8_t)(valueBuf),
		(C.ulong)(valueLen),
	)

	C.free(unsafe.Pointer(parameterTypeC))
	C.free(valueBuf)

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
	finishedContentPtr := &finishedContent

	status := C.self_message_content_credential_verification_request_builder_finish(
		(*C.self_message_content_credential_verification_request_builder)(b),
		finishedContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to build credential verificaiton request")
	}

	runtime.SetFinalizer(finishedContentPtr, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(*finishedContentPtr), nil
}

// DecodeCredentialVerificationResponse decodes a message to a credential verification response
func DecodeCredentialVerificationResponse(msg *Message) (*CredentialVerificationResponse, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var credentialVerificationResponseContent *C.self_message_content_credential_verification_response
	credentialVerificationResponseContentPtr := &credentialVerificationResponseContent

	status := C.self_message_content_as_credential_verification_response(
		content,
		credentialVerificationResponseContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to decode credential verification response message")
	}

	runtime.SetFinalizer(credentialVerificationResponseContentPtr, func(credentialVerificationResponse **C.self_message_content_credential_verification_response) {
		C.self_message_content_credential_verification_response_destroy(
			*credentialVerificationResponse,
		)
	})

	return (*CredentialVerificationResponse)(*credentialVerificationResponseContentPtr), nil
}

// ResponseTo returns the id of the request that is being responded to
func (c *CredentialVerificationResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_credential_verification_response_response_to(
			(*C.self_message_content_credential_verification_response)(c),
		)),
		20,
	)
}

// Status returns the status of the request
func (c *CredentialVerificationResponse) Status() ResponseStatus {
	return ResponseStatus(C.self_message_content_credential_verification_response_status(
		(*C.self_message_content_credential_verification_response)(c),
	))
}

// Credentials returns verified credentials that have been asserted by the responder
func (c *CredentialVerificationResponse) Credentials() *credential.VerifiableCredentialCollection {
	collection := C.self_message_content_credential_verification_response_verifiable_credentials(
		(*C.self_message_content_credential_verification_response)(c),
	)

	runtime.SetFinalizer(collection, func(collection *C.self_collection_verifiable_credential) {
		C.self_collection_verifiable_credential_destroy(
			collection,
		)
	})

	return (*credential.VerifiableCredentialCollection)(collection)
}

// NewCredentialVerificationResponse creates a new credential verification response
func NewCredentialVerificationResponse() *CredentialVerificationResponseBuilder {
	builder := C.self_message_content_credential_verification_response_builder_init()

	runtime.SetFinalizer(builder, func(builder *C.self_message_content_credential_verification_response_builder) {
		C.self_message_content_credential_verification_response_builder_destroy(
			builder,
		)
	})

	return (*CredentialVerificationResponseBuilder)(builder)
}

// ResponseTo sets the request id that is being responded to
func (b *CredentialVerificationResponseBuilder) ResponseTo(requestID []byte) *CredentialVerificationResponseBuilder {
	if len(requestID) != 20 {
		return b
	}

	requestIDBuf := C.CBytes(
		requestID,
	)

	C.self_message_content_credential_verification_response_builder_response_to(
		(*C.self_message_content_credential_verification_response_builder)(b),
		(*C.uint8_t)(requestIDBuf),
	)

	C.free(requestIDBuf)

	return b
}

// ResponseTo sets the request id that is being responded to
func (b *CredentialVerificationResponseBuilder) Status(status ResponseStatus) *CredentialVerificationResponseBuilder {
	C.self_message_content_credential_verification_response_builder_status(
		(*C.self_message_content_credential_verification_response_builder)(b),
		uint32(status),
	)

	return b
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
	finishedContentPtr := &finishedContent

	status := C.self_message_content_credential_verification_response_builder_finish(
		(*C.self_message_content_credential_verification_response_builder)(b),
		finishedContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to build credential verificaiton response")
	}

	runtime.SetFinalizer(finishedContentPtr, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(*finishedContentPtr), nil
}
