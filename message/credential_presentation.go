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

	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/credential/predicate"
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname fromCredentialTypeCollection github.com/joinself/self-go-sdk/credential.fromCredentialTypeCollection
func fromCredentialTypeCollection(c *C.self_collection_credential_type) []string

//go:linkname fromPresentationTypeCollection github.com/joinself/self-go-sdk/credential.fromPresentationTypeCollection
func fromPresentationTypeCollection(c *C.self_collection_presentation_type) []string

//go:linkname toCredentialTypeCollection github.com/joinself/self-go-sdk/credential.toCredentialTypeCollection
func toCredentialTypeCollection(credentialType []string) *C.self_collection_credential_type

//go:linkname toPresentationTypeCollection github.com/joinself/self-go-sdk/credential.toPresentationTypeCollection
func toPresentationTypeCollection(presentationType []string) *C.self_collection_presentation_type

//go:linkname fromVerifiableCredentialCollection github.com/joinself/self-go-sdk/credential.fromVerifiableCredentialCollection
func fromVerifiableCredentialCollection(c *C.self_collection_verifiable_credential) []*credential.VerifiableCredential

//go:linkname fromVerifiablePresentationCollection github.com/joinself/self-go-sdk/credential.fromVerifiablePresentationCollection
func fromVerifiablePresentationCollection(c *C.self_collection_verifiable_presentation) []*credential.VerifiablePresentation

//go:linkname verifiableCredentialPtr github.com/joinself/self-go-sdk/credential.verifiableCredentialPtr
func verifiableCredentialPtr(ptr *credential.VerifiableCredential) *C.self_verifiable_credential

//go:linkname verifiablePresentationPtr github.com/joinself/self-go-sdk/credential.verifiablePresentationPtr
func verifiablePresentationPtr(ptr *credential.VerifiablePresentation) *C.self_verifiable_presentation

//go:linkname newCredentialPredicateTree github.com/joinself/self-go-sdk/credential/predicate.newCredentialPredicateTree
func newCredentialPredicateTree(f *C.self_credential_predicate_tree) *predicate.Tree

//go:linkname credentialPredicateTreePtr github.com/joinself/self-go-sdk/credential/predicate.credentialPredicateTreePtr
func credentialPredicateTreePtr(ptr *predicate.Tree) *C.self_credential_predicate_tree

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
func DecodeCredentialPresentationRequest(content *Content) (*CredentialPresentationRequest, error) {
	contentPtr := contentPtr(content)

	var credentialPresentationRequestContent *C.self_message_content_credential_presentation_request

	result := C.self_message_content_as_credential_presentation_request(
		contentPtr,
		&credentialPresentationRequestContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCredentialPresentationRequest(credentialPresentationRequestContent), nil
}

// PresentationType returns the type of credential that presentation is being requested for
func (c *CredentialPresentationRequest) PresentationType() []string {
	collection := C.self_message_content_credential_presentation_request_presentation_type(
		c.ptr,
	)

	presentationType := fromPresentationTypeCollection(collection)

	C.self_collection_presentation_type_destroy(
		collection,
	)

	return presentationType
}

// Predicates returns a predicate tree that defines the requirements that returned credentials must match
func (c *CredentialPresentationRequest) Predicates() *predicate.Tree {
	return newCredentialPredicateTree(
		C.self_message_content_credential_presentation_request_predicates(
			c.ptr,
		),
	)
}

// Proof returns associated verifiable credential proof to support the presentation request
func (c *CredentialPresentationRequest) Proof() []*credential.VerifiablePresentation {
	collection := C.self_message_content_credential_presentation_request_proof(
		c.ptr,
	)

	presentations := fromVerifiablePresentationCollection(
		collection,
	)

	C.self_collection_verifiable_presentation_destroy(
		collection,
	)

	return presentations
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

// PresentationType sets the type of presentation being requested
func (b *CredentialPresentationRequestBuilder) PresentationType(presentationType ...string) *CredentialPresentationRequestBuilder {
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

// Predicates specifies the predicates that returned credentials must match
func (b *CredentialPresentationRequestBuilder) Predicates(tree *predicate.Tree) *CredentialPresentationRequestBuilder {
	C.self_message_content_credential_presentation_request_builder_predicates(
		b.ptr,
		credentialPredicateTreePtr(tree),
	)

	return b
}

// Proof attaches proof to the credential presentation request
func (b *CredentialPresentationRequestBuilder) Proof(proof ...*credential.VerifiablePresentation) *CredentialPresentationRequestBuilder {
	for i := range proof {
		C.self_message_content_credential_presentation_request_builder_proof(
			b.ptr,
			verifiablePresentationPtr(proof[i]),
		)
	}
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
func (b *CredentialPresentationRequestBuilder) Finish() (*Content, error) {
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
func DecodeCredentialPresentationResponse(content *Content) (*CredentialPresentationResponse, error) {
	contentPtr := contentPtr(content)

	var credentialPresentationResponseContent *C.self_message_content_credential_presentation_response

	result := C.self_message_content_as_credential_presentation_response(
		contentPtr,
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
func (b *CredentialPresentationResponseBuilder) VerifiablePresentation(presentations ...*credential.VerifiablePresentation) *CredentialPresentationResponseBuilder {
	for i := range presentations {
		C.self_message_content_credential_presentation_response_builder_verifiable_presentation(
			b.ptr,
			verifiablePresentationPtr(presentations[i]),
		)
	}
	return b
}

// Finish finalises the response and builds the content
func (b *CredentialPresentationResponseBuilder) Finish() (*Content, error) {
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
