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
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/identity"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/status"
)

type SignedPayloadType int
type UnsignedPayloadType int

const (
	SignedPayloadUnknown                     SignedPayloadType   = 1<<63 - 1
	SignedPayloadIdentityDocumentOperation   SignedPayloadType   = C.SIGNED_PAYLOAD_IDENTITY_DOCUMENT_OPERATION
	UnsignedPayloadUnknown                   UnsignedPayloadType = 1<<63 - 1
	UnsignedPayloadIdentityDocumentOperation UnsignedPayloadType = C.UNSIGNED_PAYLOAD_IDENTITY_DOCUMENT_OPERATION
)

type UnsignedPayload struct {
	ptr *C.self_message_content_unsigned_payload
}

func newUnsignedPayload(ptr *C.self_message_content_unsigned_payload) *UnsignedPayload {
	c := &UnsignedPayload{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *UnsignedPayload) {
		C.self_message_content_unsigned_payload_destroy(
			c.ptr,
		)
	})

	return c
}

func unsignedPayloadPtr(p *UnsignedPayload) *C.self_message_content_unsigned_payload {
	return p.ptr
}

// PayloadType returns the type of unsigned payload
func (p *UnsignedPayload) PayloadType() UnsignedPayloadType {
	switch C.self_message_content_unsigned_payload_payload_type(unsignedPayloadPtr(p)) {
	case C.UNSIGNED_PAYLOAD_IDENTITY_DOCUMENT_OPERATION:
		return UnsignedPayloadIdentityDocumentOperation
	default:
		return UnsignedPayloadUnknown
	}
}

// AsIdentityDocumentOperation extracts the unsigned payload as an identity document operation
func (p *UnsignedPayload) AsIdentityDocumentOperation() (*UnsignedIdentityDocumentOperation, error) {
	var identityDocumentOperation *C.self_message_content_unsigned_payload_identity_document_operation

	result := C.self_message_content_unsigned_payload_as_identity_document_operation(
		unsignedPayloadPtr(p),
		&identityDocumentOperation,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newUnsignedIdentityDocumentOperation(identityDocumentOperation), nil
}

type SignedPayload struct {
	ptr *C.self_message_content_signed_payload
}

func newSignedPayload(ptr *C.self_message_content_signed_payload) *SignedPayload {
	c := &SignedPayload{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *SignedPayload) {
		C.self_message_content_signed_payload_destroy(
			c.ptr,
		)
	})

	return c
}

func signedPayloadPtr(p *SignedPayload) *C.self_message_content_signed_payload {
	return p.ptr
}

// PayloadType returns the type of signed payload
func (p *SignedPayload) PayloadType() SignedPayloadType {
	switch C.self_message_content_signed_payload_payload_type(signedPayloadPtr(p)) {
	case C.SIGNED_PAYLOAD_IDENTITY_DOCUMENT_OPERATION:
		return SignedPayloadIdentityDocumentOperation
	default:
		return SignedPayloadUnknown
	}
}

// AsIdentityDocumentOperation extracts the signed payload as an identity document operation
func (p *SignedPayload) AsIdentityDocumentOperation() (*SignedIdentityDocumentOperation, error) {
	var identityDocumentOperation *C.self_message_content_signed_payload_identity_document_operation

	result := C.self_message_content_signed_payload_as_identity_document_operation(
		signedPayloadPtr(p),
		&identityDocumentOperation,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSignedIdentityDocumentOperation(identityDocumentOperation), nil
}

// NewUnsignedIdentityDocumentOperation creates a new unsigned payload identity document operation
func NewUnsignedIdentityDocumentOperation(documentAddress *signing.PublicKey, operation *identity.Operation) *UnsignedPayload {
	return newUnsignedPayload(
		C.self_message_content_unsigned_payload_identity_document_operation_init(
			signingPublicKeyPtr(documentAddress),
			operationPtr(operation),
		),
	)
}

type UnsignedIdentityDocumentOperation struct {
	ptr *C.self_message_content_unsigned_payload_identity_document_operation
}

func newUnsignedIdentityDocumentOperation(ptr *C.self_message_content_unsigned_payload_identity_document_operation) *UnsignedIdentityDocumentOperation {
	c := &UnsignedIdentityDocumentOperation{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *UnsignedIdentityDocumentOperation) {
		C.self_message_content_unsigned_payload_identity_document_operation_destroy(
			c.ptr,
		)
	})

	return c
}

// DocumentAddress returns the address of the document the operation relates to
func (c *UnsignedIdentityDocumentOperation) DocumentAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_message_content_unsigned_payload_identity_document_operation_document_address(
		c.ptr,
	))
}

// Operation returns the operation that the signature is being requested for
func (c *UnsignedIdentityDocumentOperation) Operation() *identity.Operation {
	return newOperation(C.self_message_content_unsigned_payload_identity_document_operation_operation(
		c.ptr,
	))
}

func NewSignedIdentityDocumentOperation(documentAddress *signing.PublicKey, operation *identity.Operation) *SignedPayload {
	return newSignedPayload(
		C.self_message_content_signed_payload_identity_document_operation_init(
			signingPublicKeyPtr(documentAddress),
			operationPtr(operation),
		),
	)
}

type SignedIdentityDocumentOperation struct {
	ptr *C.self_message_content_signed_payload_identity_document_operation
}

func newSignedIdentityDocumentOperation(ptr *C.self_message_content_signed_payload_identity_document_operation) *SignedIdentityDocumentOperation {
	c := &SignedIdentityDocumentOperation{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *SignedIdentityDocumentOperation) {
		C.self_message_content_signed_payload_identity_document_operation_destroy(
			c.ptr,
		)
	})

	return c
}

// DocumentAddress returns the address of the document the operation relates to
func (c *SignedIdentityDocumentOperation) DocumentAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_message_content_signed_payload_identity_document_operation_document_address(
		c.ptr,
	))
}

// Operation returns the operation that the signature is being requested for
func (c *SignedIdentityDocumentOperation) Operation() *identity.Operation {
	return newOperation(C.self_message_content_signed_payload_identity_document_operation_operation(
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
func DecodeSigningRequest(msg *event.Message) (*SigningRequest, error) {
	content := contentPtr(msg.Content())

	var accountPairingRequestContent *C.self_message_content_signing_request

	result := C.self_message_content_as_signing_request(
		content,
		&accountPairingRequestContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningRequest(accountPairingRequestContent), nil
}

// UnsignedPayload returns the unsigned payload of the
func (c *SigningRequest) UnsignedPayload() *UnsignedPayload {
	return newUnsignedPayload(C.self_message_content_signing_request_unsigned_payload(
		c.ptr,
	))
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

// UnsignedPayload sets the unsigned payload of the request
func (b *SigningRequestBuilder) UnsignedPayload(payload *UnsignedPayload) *SigningRequestBuilder {
	C.self_message_content_signing_request_builder_unsigned_payload(
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
func (b *SigningRequestBuilder) Finish() (*event.Content, error) {
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
func DecodeSigningResponse(msg *event.Message) (*SigningResponse, error) {
	content := contentPtr(msg.Content())

	var accountPairingResponseContent *C.self_message_content_signing_response

	result := C.self_message_content_as_signing_response(
		content,
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

// SignedPayload returns the signed payload of the
func (c *SigningResponse) SignedPayload() *SignedPayload {
	return newSignedPayload(C.self_message_content_signing_response_signed_payload(
		c.ptr,
	))
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

// SignedPayload sets the signed payload of the response
func (b *SigningResponseBuilder) SignedPayload(signer *signing.PublicKey, payload *SignedPayload) *SigningResponseBuilder {
	C.self_message_content_signing_response_builder_signed_payload(
		b.ptr,
		signingPublicKeyPtr(signer),
		payload.ptr,
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
func (b *SigningResponseBuilder) Finish() (*event.Content, error) {
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
