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

//go:linkname operationPtr github.com/joinself/self-go-sdk/identity.operationPtr
func operationPtr(o *identity.Operation) *C.self_identity_operation

//go:linkname newOperation github.com/joinself/self-go-sdk/identity.newOperation
func newOperation(ptr *C.self_identity_operation) *identity.Operation

//go:linkname fromObjectCollection github.com/joinself/self-go-sdk/object.fromObjectCollection
func fromObjectCollection(ptr *C.self_collection_object) []*object.Object

type AccountPairingRequest struct {
	ptr *C.self_message_content_account_pairing_request
}

func newAccountPairingRequest(ptr *C.self_message_content_account_pairing_request) *AccountPairingRequest {
	c := &AccountPairingRequest{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(ptr *C.self_message_content_account_pairing_request) {
		C.self_message_content_account_pairing_request_destroy(
			ptr,
		)
	}, c.ptr)

	return c
}

type AccountPairingResponse struct {
	ptr *C.self_message_content_account_pairing_response
}

func newAccountPairingResponse(ptr *C.self_message_content_account_pairing_response) *AccountPairingResponse {
	c := &AccountPairingResponse{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(ptr *C.self_message_content_account_pairing_response) {
		C.self_message_content_account_pairing_response_destroy(
			ptr,
		)
	}, c.ptr)

	return c
}

type AccountPairingRequestBuilder struct {
	ptr *C.self_message_content_account_pairing_request_builder
}

func newAccountPairingRequestBuilder(ptr *C.self_message_content_account_pairing_request_builder) *AccountPairingRequestBuilder {
	c := &AccountPairingRequestBuilder{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(ptr *C.self_message_content_account_pairing_request_builder) {
		C.self_message_content_account_pairing_request_builder_destroy(
			ptr,
		)
	}, c.ptr)

	return c
}

type AccountPairingResponseBuilder struct {
	ptr *C.self_message_content_account_pairing_response_builder
}

func newAccountPairingResponseBuilder(ptr *C.self_message_content_account_pairing_response_builder) *AccountPairingResponseBuilder {
	c := &AccountPairingResponseBuilder{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(ptr *C.self_message_content_account_pairing_response_builder) {
		C.self_message_content_account_pairing_response_builder_destroy(
			ptr,
		)
	}, c.ptr)

	return c
}

// DecodeAccountPairingRequest decodes a message to a account pairing request
func DecodeAccountPairingRequest(content *Content) (*AccountPairingRequest, error) {
	contentPtr := contentPtr(content)

	var accountPairingRequestContent *C.self_message_content_account_pairing_request

	result := C.self_message_content_as_account_pairing_request(
		contentPtr,
		&accountPairingRequestContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newAccountPairingRequest(accountPairingRequestContent), nil
}

// Address returns the key that the sender wishes to be linked
func (c *AccountPairingRequest) Address() *signing.PublicKey {
	return newSigningPublicKey(C.self_message_content_account_pairing_request_address(
		c.ptr,
	))
}

// Roles returns the roles that the sender wishes the specified key to have
func (c *AccountPairingRequest) Roles() identity.Role {
	return identity.Role(C.self_message_content_account_pairing_request_roles(
		c.ptr,
	))
}

// Expires returns the time the request expires at
func (c *AccountPairingRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_account_pairing_request_expires(
		c.ptr,
	)), 0)
}

// NewAccountPairingRequest creates a new account pairing request
func NewAccountPairingRequest() *AccountPairingRequestBuilder {
	return newAccountPairingRequestBuilder(
		C.self_message_content_account_pairing_request_builder_init(),
	)
}

// Address sets the key that is to be linked to an identity document
func (b *AccountPairingRequestBuilder) Address(address *signing.PublicKey) *AccountPairingRequestBuilder {
	C.self_message_content_account_pairing_request_builder_address(
		b.ptr,
		signingPublicKeyPtr(address),
	)
	return b
}

// Roles sets the requested roles that the linked key should have
func (b *AccountPairingRequestBuilder) Roles(roles identity.Role) *AccountPairingRequestBuilder {
	C.self_message_content_account_pairing_request_builder_roles(
		b.ptr,
		C.uint64_t(roles),
	)
	return b
}

// Expires sets the time that the request expires at
func (b *AccountPairingRequestBuilder) Expires(expires time.Time) *AccountPairingRequestBuilder {
	C.self_message_content_account_pairing_request_builder_expires(
		b.ptr,
		C.int64_t(expires.Unix()),
	)
	return b
}

// Finish finalises the request and builds the content
func (b *AccountPairingRequestBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_account_pairing_request_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}

// DecodeAccountPairingResponse decodes a message to a account pairing response
func DecodeAccountPairingResponse(content *Content) (*AccountPairingResponse, error) {
	contentPtr := contentPtr(content)

	var accountPairingResponseContent *C.self_message_content_account_pairing_response

	result := C.self_message_content_as_account_pairing_response(
		contentPtr,
		&accountPairingResponseContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newAccountPairingResponse(accountPairingResponseContent), nil
}

// ResponseTo returns the id of the request that is being responded to
func (c *AccountPairingResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_account_pairing_response_response_to(
			c.ptr,
		)),
		20,
	)
}

// Status returns the status of the request
func (c *AccountPairingResponse) Status() ResponseStatus {
	return ResponseStatus(C.self_message_content_account_pairing_response_status(
		c.ptr,
	))
}

// DocumentAddress returns the address of the document the key will be linked to
func (c *AccountPairingResponse) DocumentAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_message_content_account_pairing_response_document_address(
		c.ptr,
	))
}

// Operation returns the signed operation that can be executed to add the key to the document
func (c *AccountPairingResponse) Operation() *identity.Operation {
	return newOperation(C.self_message_content_account_pairing_response_operation(
		c.ptr,
	))
}

// Presentations returns any presentations that can be used by the linked by
func (c *AccountPairingResponse) Presentations() []*credential.VerifiablePresentation {
	collection := C.self_message_content_account_pairing_response_presentations(
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
func (c *AccountPairingResponse) Assets() []*object.Object {
	collection := C.self_message_content_account_pairing_response_assets(
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

// NewAccountPairingResponse creates a new account pairing response
func NewAccountPairingResponse() *AccountPairingResponseBuilder {
	return newAccountPairingResponseBuilder(
		C.self_message_content_account_pairing_response_builder_init(),
	)
}

// ResponseTo sets the request id that is being responded to
func (b *AccountPairingResponseBuilder) ResponseTo(requestID []byte) *AccountPairingResponseBuilder {
	if len(requestID) != 20 {
		return b
	}

	requestIDBuf := C.CBytes(
		requestID,
	)

	C.self_message_content_account_pairing_response_builder_response_to(
		b.ptr,
		(*C.uint8_t)(requestIDBuf),
	)

	C.free(requestIDBuf)

	return b
}

// ResponseTo sets the request id that is being responded to
func (b *AccountPairingResponseBuilder) Status(status ResponseStatus) *AccountPairingResponseBuilder {
	C.self_message_content_account_pairing_response_builder_status(
		b.ptr,
		uint32(status),
	)

	return b
}

// DocumentAddress sets the address of the document that the key will be linked to
func (b *AccountPairingResponseBuilder) DocumentAddress(address *signing.PublicKey) *AccountPairingResponseBuilder {
	C.self_message_content_account_pairing_response_builder_document_address(
		b.ptr,
		signingPublicKeyPtr(address),
	)

	return b
}

// Operation sets the operation that will add the key to the identity document
func (b *AccountPairingResponseBuilder) Operation(operation *identity.Operation) *AccountPairingResponseBuilder {
	C.self_message_content_account_pairing_response_builder_operation(
		b.ptr,
		operationPtr(operation),
	)

	return b
}

// Presentation adds a presentation that can be used by the key that has been added to the document
func (b *AccountPairingResponseBuilder) Presentation(presentation *credential.VerifiablePresentation) *AccountPairingResponseBuilder {
	C.self_message_content_account_pairing_response_builder_presentation(
		b.ptr,
		verifiablePresentationPtr(presentation),
	)

	return b
}

// Asset adds an asset that can be used in support of an attached presentation
func (b *AccountPairingResponseBuilder) Asset(asset *object.Object) *AccountPairingResponseBuilder {
	C.self_message_content_account_pairing_response_builder_asset(
		b.ptr,
		objectPtr(asset),
	)

	return b
}

// Finish finalises the response and builds the content
func (b *AccountPairingResponseBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_account_pairing_response_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}
