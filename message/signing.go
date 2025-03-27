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
	"github.com/joinself/self-go-sdk/identity"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname toSigningPublicKeyCollection github.com/joinself/self-go-sdk/keypair/signing.toSigningPublicKeyCollection
func toSigningPublicKeyCollection(p []*signing.PublicKey) *C.self_collection_signing_public_key

type SigningPayloadType int

const (
	SigningPayloadUnknown                   SigningPayloadType = 1<<63 - 1
	SigningPayloadIdentityDocumentOperation SigningPayloadType = C.SIGNING_PAYLOAD_IDENTITY_DOCUMENT_OPERATION
)

type SigningPayload struct {
	ptr *C.self_message_content_signing_payload
}

func newSigningPayload(ptr *C.self_message_content_signing_payload) *SigningPayload {
	c := &SigningPayload{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *SigningPayload) {
		C.self_message_content_signing_payload_destroy(
			c.ptr,
		)
	})

	return c
}

func signingPayloadPtr(p *SigningPayload) *C.self_message_content_signing_payload {
	return p.ptr
}

func fromSigningPayloadCollection(collection *C.self_collection_message_content_signing_payload) []*SigningPayload {
	collectionLen := int(C.self_collection_message_content_signing_payload_len(
		collection,
	))

	payloads := make([]*SigningPayload, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_message_content_signing_payload_at(
			collection,
			C.size_t(i),
		)

		payloads[i] = newSigningPayload(ptr)
	}

	return payloads
}

// PayloadType returns the type of signing payload
func (p *SigningPayload) PayloadType() SigningPayloadType {
	switch C.self_message_content_signing_payload_type_of(signingPayloadPtr(p)) {
	case C.SIGNING_PAYLOAD_IDENTITY_DOCUMENT_OPERATION:
		return SigningPayloadIdentityDocumentOperation
	default:
		return SigningPayloadUnknown
	}
}

// AsIdentityDocumentOperation extracts the signing payload as an identity document operation
func (p *SigningPayload) AsIdentityDocumentOperation() (*SigningIdentityDocumentOperation, error) {
	var identityDocumentOperation *C.self_message_content_signing_payload_identity_document_operation

	result := C.self_message_content_signing_payload_as_identity_document_operation(
		signingPayloadPtr(p),
		&identityDocumentOperation,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningIdentityDocumentOperation(identityDocumentOperation), nil
}

// NewSigningIdentityDocumentOperation creates a new Signing payload identity document operation
func NewSigningIdentityDocumentOperation(documentAddress *signing.PublicKey, operation *identity.Operation) *SigningPayload {
	return newSigningPayload(
		C.self_message_content_signing_payload_identity_document_operation_init(
			signingPublicKeyPtr(documentAddress),
			operationPtr(operation),
		),
	)
}

type SigningIdentityDocumentOperation struct {
	ptr *C.self_message_content_signing_payload_identity_document_operation
}

func newSigningIdentityDocumentOperation(ptr *C.self_message_content_signing_payload_identity_document_operation) *SigningIdentityDocumentOperation {
	c := &SigningIdentityDocumentOperation{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *SigningIdentityDocumentOperation) {
		C.self_message_content_signing_payload_identity_document_operation_destroy(
			c.ptr,
		)
	})

	return c
}

// DocumentAddress returns the address of the document the operation relates to
func (c *SigningIdentityDocumentOperation) DocumentAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_message_content_signing_payload_identity_document_operation_document_address(
		c.ptr,
	))
}

// Operation returns the operation that the signature is being requested for
func (c *SigningIdentityDocumentOperation) Operation() *identity.Operation {
	return newOperation(C.self_message_content_signing_payload_identity_document_operation_operation(
		c.ptr,
	))
}

type SigningRequest struct {
	ptr *C.self_message_content_signing_request
}

func newSigningRequest(ptr *C.self_message_content_signing_request) *SigningRequest {
	c := &SigningRequest{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *SigningRequest) {
		C.self_message_content_signing_request_destroy(
			c.ptr,
		)
	})

	return c
}

type SigningResponse struct {
	ptr *C.self_message_content_signing_response
}

func newSigningResponse(ptr *C.self_message_content_signing_response) *SigningResponse {
	c := &SigningResponse{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *SigningResponse) {
		C.self_message_content_signing_response_destroy(
			c.ptr,
		)
	})

	return c
}

type SigningRequestBuilder struct {
	ptr *C.self_message_content_signing_request_builder
}

func newSigningRequestBuilder(ptr *C.self_message_content_signing_request_builder) *SigningRequestBuilder {
	c := &SigningRequestBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *SigningRequestBuilder) {
		C.self_message_content_signing_request_builder_destroy(
			c.ptr,
		)
	})

	return c
}

type SigningResponseBuilder struct {
	ptr *C.self_message_content_signing_response_builder
}

func newSigningResponseBuilder(ptr *C.self_message_content_signing_response_builder) *SigningResponseBuilder {
	c := &SigningResponseBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *SigningResponseBuilder) {
		C.self_message_content_signing_response_builder_destroy(
			c.ptr,
		)
	})

	return c
}

// DecodeSigningRequest decodes a message to a account pairing request
func DecodeSigningRequest(content *Content) (*SigningRequest, error) {
	contentPtr := contentPtr(content)

	var accountPairingRequestContent *C.self_message_content_signing_request

	result := C.self_message_content_as_signing_request(
		contentPtr,
		&accountPairingRequestContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningRequest(accountPairingRequestContent), nil
}

// Payload returns the Signing payload of the
func (c *SigningRequest) Payloads() ([]*SigningPayload, error) {
	var payloads *C.self_collection_message_content_signing_payload

	result := C.self_message_content_signing_request_payloads(
		c.ptr,
		&payloads,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return fromSigningPayloadCollection(payloads), nil
}

// RequiresLiveness returns true if the request requires an accompanying liveness check
func (c *SigningRequest) RequiresLiveness() bool {
	return bool(C.self_message_content_signing_request_requires_liveness(
		c.ptr,
	))
}

// Expires returns the time the request expires at
func (c *SigningRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_signing_request_expires(
		c.ptr,
	)), 0)
}

// NewSigningRequest creates a new account pairing request
func NewSigningRequest() *SigningRequestBuilder {
	return newSigningRequestBuilder(
		C.self_message_content_signing_request_builder_init(),
	)
}

// Expires sets the time that the request expires at
func (b *SigningRequestBuilder) Expires(expires time.Time) *SigningRequestBuilder {
	C.self_message_content_signing_request_builder_expires(
		b.ptr,
		C.int64_t(expires.Unix()),
	)
	return b
}

// Payload sets the Signing payload of the request
func (b *SigningRequestBuilder) Payload(payload *SigningPayload) *SigningRequestBuilder {
	C.self_message_content_signing_request_builder_payload(
		b.ptr,
		payload.ptr,
	)
	return b
}

// RequireLiveness specifies the signature is required to be accompanied by a linked liveness check
func (b *SigningRequestBuilder) RequireLiveness() *SigningRequestBuilder {
	C.self_message_content_signing_request_builder_requires_liveness(
		b.ptr,
	)
	return b
}

// Finish finalises the request and builds the content
func (b *SigningRequestBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_signing_request_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}

// DecodeSigningResponse decodes a message to a account pairing response
func DecodeSigningResponse(content *Content) (*SigningResponse, error) {
	contentPtr := contentPtr(content)

	var accountPairingResponseContent *C.self_message_content_signing_response

	result := C.self_message_content_as_signing_response(
		contentPtr,
		&accountPairingResponseContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningResponse(accountPairingResponseContent), nil
}

// ResponseTo returns the id of the request that is being responded to
func (c *SigningResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_signing_response_response_to(
			c.ptr,
		)),
		20,
	)
}

// Status returns the status of the request
func (c *SigningResponse) Status() ResponseStatus {
	return ResponseStatus(C.self_message_content_signing_response_status(
		c.ptr,
	))
}

// Payload returns the signing payload of the
func (c *SigningResponse) Payloads() ([]*SigningPayload, error) {
	var payloads *C.self_collection_message_content_signing_payload

	result := C.self_message_content_signing_response_payloads(
		c.ptr,
		&payloads,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return fromSigningPayloadCollection(payloads), nil
}

// Presentations returns any presentations that can be used by the linked by
func (c *SigningResponse) Presentations() []*credential.VerifiablePresentation {
	collection := C.self_message_content_signing_response_presentations(
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

// Assets returns any supporting objects needed to support claims in the provided presentations
func (c *SigningResponse) Assets() []*object.Object {
	collection := C.self_message_content_signing_response_assets(
		c.ptr,
	)

	objects := fromObjectCollection(
		collection,
	)

	C.self_collection_object_destroy(
		collection,
	)

	return objects
}

// NewSigningResponse creates a new account pairing response
func NewSigningResponse() *SigningResponseBuilder {
	return newSigningResponseBuilder(
		C.self_message_content_signing_response_builder_init(),
	)
}

// ResponseTo sets the request id that is being responded to
func (b *SigningResponseBuilder) ResponseTo(requestID []byte) *SigningResponseBuilder {
	if len(requestID) != 20 {
		return b
	}

	requestIDBuf := C.CBytes(
		requestID,
	)

	C.self_message_content_signing_response_builder_response_to(
		b.ptr,
		(*C.uint8_t)(requestIDBuf),
	)

	C.free(requestIDBuf)

	return b
}

// ResponseTo sets the request id that is being responded to
func (b *SigningResponseBuilder) Status(status ResponseStatus) *SigningResponseBuilder {
	C.self_message_content_signing_response_builder_status(
		b.ptr,
		uint32(status),
	)

	return b
}

// Payload sets the signing payload of the response and which keys to sign it with
func (b *SigningResponseBuilder) Payload(payload *SigningPayload, signers []*signing.PublicKey) *SigningResponseBuilder {
	C.self_message_content_signing_response_builder_payload(
		b.ptr,
		payload.ptr,
		toSigningPublicKeyCollection(signers),
	)
	return b
}

// Presentation adds a presentation that can be used by the key that has been added to the document
func (b *SigningResponseBuilder) Presentation(presentation *credential.VerifiablePresentation) *SigningResponseBuilder {
	C.self_message_content_signing_response_builder_presentation(
		b.ptr,
		verifiablePresentationPtr(presentation),
	)

	return b
}

// Asset adds an asset that can be used in support of an attached presentation
func (b *SigningResponseBuilder) Asset(asset *object.Object) *SigningResponseBuilder {
	C.self_message_content_signing_response_builder_asset(
		b.ptr,
		objectPtr(asset),
	)

	return b
}

// Finish finalises the response and builds the content
func (b *SigningResponseBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_signing_response_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}
