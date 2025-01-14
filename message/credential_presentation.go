package message

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk-next/credential"
	"github.com/joinself/self-go-sdk-next/event"
	"github.com/joinself/self-go-sdk-next/status"
)

//go:linkname fromCredentialTypeCollection github.com/joinself/self-go-sdk-next/credential.fromCredentialTypeCollection
func fromCredentialTypeCollection(c *C.self_collection_credential_type) []string

//go:linkname fromPresentationTypeCollection github.com/joinself/self-go-sdk-next/credential.fromPresentationTypeCollection
func fromPresentationTypeCollection(c *C.self_collection_presentation_type) []string

//go:linkname fromPresentationDetailCollection github.com/joinself/self-go-sdk-next/credential.fromPresentationDetailCollection
func fromPresentationDetailCollection(c *C.self_collection_credential_presentation_detail) []*credential.CredentialPresentationDetail

//go:linkname toCredentialTypeCollection github.com/joinself/self-go-sdk-next/credential.toCredentialTypeCollection
func toCredentialTypeCollection(credentialType []string) *C.self_collection_credential_type

//go:linkname toPresentationTypeCollection github.com/joinself/self-go-sdk-next/credential.toPresentationTypeCollection
func toPresentationTypeCollection(presentationType []string) *C.self_collection_presentation_type

//go:linkname fromVerifiableCredentialCollection github.com/joinself/self-go-sdk-next/credential.fromVerifiableCredentialCollection
func fromVerifiableCredentialCollection(c *C.self_collection_verifiable_credential) []*credential.VerifiableCredential

//go:linkname fromVerifiablePresentationCollection github.com/joinself/self-go-sdk-next/credential.fromVerifiablePresentationCollection
func fromVerifiablePresentationCollection(c *C.self_collection_verifiable_presentation) []*credential.VerifiablePresentation

//go:linkname fromCredentialVerificationEvidenceCollection github.com/joinself/self-go-sdk-next/credential.fromCredentialVerificationEvidenceCollection
func fromCredentialVerificationEvidenceCollection(c *C.self_collection_credential_verification_evidence) []*credential.CredentialVerificationEvidence

//go:linkname fromCredentialVerificationParameterCollection github.com/joinself/self-go-sdk-next/credential.fromCredentialVerificationParameterCollection
func fromCredentialVerificationParameterCollection(c *C.self_collection_credential_verification_parameter) []*credential.CredentialVerificationParameter

//go:linkname verifiableCredentialPtr github.com/joinself/self-go-sdk-next/credential.verifiableCredentialPtr
func verifiableCredentialPtr(ptr *credential.VerifiableCredential) *C.self_verifiable_credential

//go:linkname verifiablePresentationPtr github.com/joinself/self-go-sdk-next/credential.verifiablePresentationPtr
func verifiablePresentationPtr(ptr *credential.VerifiablePresentation) *C.self_verifiable_presentation

type CredentialPresentationRequest struct {
	ptr *C.self_message_content_credential_presentation_request
}

func newCredentialPresentationRequest(ptr *C.self_message_content_credential_presentation_request) *CredentialPresentationRequest {
	c := &CredentialPresentationRequest{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialPresentationRequest) {
		C.self_message_content_credential_presentation_request_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialPresentationResponse struct {
	ptr *C.self_message_content_credential_presentation_response
}

func newCredentialPresentationResponse(ptr *C.self_message_content_credential_presentation_response) *CredentialPresentationResponse {
	c := &CredentialPresentationResponse{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialPresentationResponse) {
		C.self_message_content_credential_presentation_response_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialPresentationRequestBuilder struct {
	ptr *C.self_message_content_credential_presentation_request_builder
}

func newCredentialPresentationRequestBuilder(ptr *C.self_message_content_credential_presentation_request_builder) *CredentialPresentationRequestBuilder {
	c := &CredentialPresentationRequestBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialPresentationRequestBuilder) {
		C.self_message_content_credential_presentation_request_builder_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialPresentationResponseBuilder struct {
	ptr *C.self_message_content_credential_presentation_response_builder
}

func newCredentialPresentationResponseBuilder(ptr *C.self_message_content_credential_presentation_response_builder) *CredentialPresentationResponseBuilder {
	c := &CredentialPresentationResponseBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialPresentationResponseBuilder) {
		C.self_message_content_credential_presentation_response_builder_destroy(
			c.ptr,
		)
	})

	return c
}

// DecodeCredentialPresentationRequest decodes a message to a credential presentation request
func DecodeCredentialPresentationRequest(msg *event.Message) (*CredentialPresentationRequest, error) {
	content := contentPtr(msg.Content())

	var credentialPresentationRequestContent *C.self_message_content_credential_presentation_request

	result := C.self_message_content_as_credential_presentation_request(
		content,
		&credentialPresentationRequestContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCredentialPresentationRequest(credentialPresentationRequestContent), nil
}

// Type returns the type of credential that presentation is being requested for
func (c *CredentialPresentationRequest) Type() []string {
	collection := C.self_message_content_credential_presentation_request_presentation_type(
		c.ptr,
	)

	presentationType := fromPresentationTypeCollection(collection)

	C.self_collection_presentation_type_destroy(
		collection,
	)

	return presentationType
}

// Details returns details of the requested credential presentations
func (c *CredentialPresentationRequest) Details() []*credential.CredentialPresentationDetail {
	collection := C.self_message_content_credential_presentation_request_details(
		c.ptr,
	)

	details := fromPresentationDetailCollection(collection)

	C.self_collection_credential_presentation_detail_destroy(
		collection,
	)

	return details
}

// Type returns the time the request expires at
func (c *CredentialPresentationRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_credential_presentation_request_expires(
		c.ptr,
	)), 0)
}

// NewCredentialPresentationRequest creates a new credential presentation request
func NewCredentialPresentationRequest() *CredentialPresentationRequestBuilder {
	return newCredentialPresentationRequestBuilder(
		C.self_message_content_credential_presentation_request_builder_init(),
	)
}

// Type sets the type of presentation being requested
func (b *CredentialPresentationRequestBuilder) Type(presentationType []string) *CredentialPresentationRequestBuilder {
	collection := toPresentationTypeCollection(presentationType)

	C.self_message_content_credential_presentation_request_builder_presentation_type(
		b.ptr,
		collection,
	)

	C.self_collection_presentation_type_destroy(
		collection,
	)

	return b
}

// Details specifies the details of the credentials being requested for presentation
func (b *CredentialPresentationRequestBuilder) Details(credentialType []string, subject string) *CredentialPresentationRequestBuilder {
	subjectC := C.CString(subject)

	collection := toCredentialTypeCollection(credentialType)

	C.self_message_content_credential_presentation_request_builder_details(
		b.ptr,
		collection,
		subjectC,
	)

	C.free(unsafe.Pointer(subjectC))
	C.self_collection_credential_type_destroy(
		collection,
	)

	return b
}

// Expires sets the time that the request expires at
func (b *CredentialPresentationRequestBuilder) Expires(expires time.Time) *CredentialPresentationRequestBuilder {
	C.self_message_content_credential_presentation_request_builder_expires(
		b.ptr,
		C.int64_t(expires.Unix()),
	)
	return b
}

// Finish finalises the request and builds the content
func (b *CredentialPresentationRequestBuilder) Finish() (*event.Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_credential_presentation_request_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}

// DecodeCredentialPresentationResponse decodes a message to a credential presentation response
func DecodeCredentialPresentationResponse(msg *event.Message) (*CredentialPresentationResponse, error) {
	content := contentPtr(msg.Content())

	var credentialPresentationResponseContent *C.self_message_content_credential_presentation_response

	result := C.self_message_content_as_credential_presentation_response(
		content,
		&credentialPresentationResponseContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCredentialPresentationResponse(credentialPresentationResponseContent), nil
}

// ResponseTo returns the id of the request that is being responded to
func (c *CredentialPresentationResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_credential_presentation_response_response_to(
			c.ptr,
		)),
		20,
	)
}

// Status returns the status of the request
func (c *CredentialPresentationResponse) Status() ResponseStatus {
	return ResponseStatus(C.self_message_content_credential_presentation_response_status(
		c.ptr,
	))
}

// Presentations returns veriable presentations that have been asserted by the responder
func (c *CredentialPresentationResponse) Presentations() []*credential.VerifiablePresentation {
	collection := C.self_message_content_credential_presentation_response_verifiable_presentations(
		c.ptr,
	)

	credentials := fromVerifiablePresentationCollection(
		collection,
	)

	C.self_collection_verifiable_presentation_destroy(
		collection,
	)

	return credentials
}

// NewCredentialPresentationResponse creates a new credential presentation response
func NewCredentialPresentationResponse() *CredentialPresentationResponseBuilder {
	return newCredentialPresentationResponseBuilder(
		C.self_message_content_credential_presentation_response_builder_init(),
	)
}

// ResponseTo sets the request id that is being responded to
func (b *CredentialPresentationResponseBuilder) ResponseTo(requestID []byte) *CredentialPresentationResponseBuilder {
	if len(requestID) != 20 {
		return b
	}

	requestIDBuf := C.CBytes(
		requestID,
	)

	C.self_message_content_credential_presentation_response_builder_response_to(
		b.ptr,
		(*C.uint8_t)(requestIDBuf),
	)

	C.free(requestIDBuf)

	return b
}

// ResponseTo sets the request id that is being responded to
func (b *CredentialPresentationResponseBuilder) Status(status ResponseStatus) *CredentialPresentationResponseBuilder {
	C.self_message_content_credential_presentation_response_builder_status(
		b.ptr,
		uint32(status),
	)

	return b
}

// VerifiablePresentation attaches a verified presentation of credentails to the response
func (b *CredentialPresentationResponseBuilder) VerifiablePresentation(presentation *credential.VerifiablePresentation) *CredentialPresentationResponseBuilder {
	C.self_message_content_credential_presentation_response_builder_verifiable_presentation(
		b.ptr,
		verifiablePresentationPtr(presentation),
	)
	return b
}

// Finish finalises the response and builds the content
func (b *CredentialPresentationResponseBuilder) Finish() (*event.Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_credential_presentation_response_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}
