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

	"github.com/joinself/self-go-sdk/crypto"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname cryptoKeyPackagePtr github.com/joinself/self-go-sdk/crypto.cryptoKeyPackagePtr
func cryptoKeyPackagePtr(k *crypto.KeyPackage) *C.self_crypto_key_package

//go:linkname newCryptoKeyPackage github.com/joinself/self-go-sdk/crypto.newCryptoKeyPackage
func newCryptoKeyPackage(ptr *C.self_crypto_key_package) *crypto.KeyPackage

type DiscoveryRequest struct {
	ptr *C.self_message_content_discovery_request
}

func newDiscoveryRequest(ptr *C.self_message_content_discovery_request) *DiscoveryRequest {
	c := &DiscoveryRequest{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(ptr *C.self_message_content_discovery_request) {
		C.self_message_content_discovery_request_destroy(
			ptr,
		)
	}, c.ptr)

	return c
}

type DiscoveryResponse struct {
	ptr *C.self_message_content_discovery_response
}

func newDiscoveryResponse(ptr *C.self_message_content_discovery_response) *DiscoveryResponse {
	c := &DiscoveryResponse{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(ptr *C.self_message_content_discovery_response) {
		C.self_message_content_discovery_response_destroy(
			ptr,
		)
	}, c.ptr)

	return c
}

type DiscoveryRequestBuilder struct {
	ptr *C.self_message_content_discovery_request_builder
}

func newDiscoveryRequestBuilder(ptr *C.self_message_content_discovery_request_builder) *DiscoveryRequestBuilder {
	c := &DiscoveryRequestBuilder{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(ptr *C.self_message_content_discovery_request_builder) {
		C.self_message_content_discovery_request_builder_destroy(
			ptr,
		)
	}, c.ptr)

	return c
}

type DiscoveryResponseBuilder struct {
	ptr *C.self_message_content_discovery_response_builder
}

func newDiscoveryResponseBuilder(ptr *C.self_message_content_discovery_response_builder) *DiscoveryResponseBuilder {
	c := &DiscoveryResponseBuilder{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(ptr *C.self_message_content_discovery_response_builder) {
		C.self_message_content_discovery_response_builder_destroy(
			ptr,
		)
	}, c.ptr)

	return c
}

// DecodeDiscoveryRequest decodes a message to a discovery request
func DecodeDiscoveryRequest(content *Content) (*DiscoveryRequest, error) {
	contentPtr := contentPtr(content)

	var discoveryRequestContent *C.self_message_content_discovery_request

	result := C.self_message_content_as_discovery_request(
		contentPtr,
		&discoveryRequestContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newDiscoveryRequest(discoveryRequestContent), nil
}

// DocumentAddress returns the document address of the requester, or nil if not provided
func (c *DiscoveryRequest) DocumentAddress() *signing.PublicKey {
	documentAddress := C.self_message_content_discovery_request_document_address(
		c.ptr,
	)

	if documentAddress == nil {
		return nil
	}

	return newSigningPublicKey(documentAddress)
}

// InboxAddress returns the inbox address of the requester
func (c *DiscoveryRequest) InboxAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_message_content_discovery_request_from_address(
		c.ptr,
	))
}

// KeyPackage returns the embedded key package contained in the discovery request, or nil if not provided
func (c *DiscoveryRequest) KeyPackage() *crypto.KeyPackage {
	keyPackage := C.self_message_content_discovery_request_key_package(
		c.ptr,
	)

	if keyPackage == nil {
		return nil
	}

	return newCryptoKeyPackage(keyPackage)
}

// Type returns the time the request expires at
func (c *DiscoveryRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_discovery_request_expires(
		c.ptr,
	)), 0)
}

// NewDiscoveryRequest creates a new discovery request
func NewDiscoveryRequest() *DiscoveryRequestBuilder {
	return newDiscoveryRequestBuilder(
		C.self_message_content_discovery_request_builder_init(),
	)
}

// DocumentAddress sets an optional document address
func (b *DiscoveryRequestBuilder) DocumentAddress(documentAddress *signing.PublicKey) *DiscoveryRequestBuilder {
	C.self_message_content_discovery_request_builder_document_address(
		b.ptr,
		signingPublicKeyPtr(documentAddress),
	)
	return b
}

// InboxAddress sets the inbox address of the requester. This can be omitted if a key package is provided
func (b *DiscoveryRequestBuilder) InboxAddress(inboxAddress *signing.PublicKey) *DiscoveryRequestBuilder {
	C.self_message_content_discovery_request_builder_from_address(
		b.ptr,
		signingPublicKeyPtr(inboxAddress),
	)
	return b
}

// KeyPackage sets an optional key package that will be embedded in the request
func (b *DiscoveryRequestBuilder) KeyPackage(keyPackage *crypto.KeyPackage) *DiscoveryRequestBuilder {
	C.self_message_content_discovery_request_builder_key_package(
		b.ptr,
		cryptoKeyPackagePtr(keyPackage),
	)
	return b
}

// Expires sets the time that the request expires at
func (b *DiscoveryRequestBuilder) Expires(expires time.Time) *DiscoveryRequestBuilder {
	C.self_message_content_discovery_request_builder_expires(
		b.ptr,
		C.int64_t(expires.Unix()),
	)
	return b
}

// Finish finalises the request and builds the content
func (b *DiscoveryRequestBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_discovery_request_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}

// DecodeDiscoveryResponse decodes a message to a discovery response
func DecodeDiscoveryResponse(content *Content) (*DiscoveryResponse, error) {
	contentPtr := contentPtr(content)

	var discoveryResponseContent *C.self_message_content_discovery_response

	result := C.self_message_content_as_discovery_response(
		contentPtr,
		&discoveryResponseContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newDiscoveryResponse(discoveryResponseContent), nil
}

// ResponseTo returns the id of the request that is being responded to
func (c *DiscoveryResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_discovery_response_response_to(
			c.ptr,
		)),
		20,
	)
}

// Status returns the status of the request
func (c *DiscoveryResponse) Status() ResponseStatus {
	return ResponseStatus(C.self_message_content_discovery_response_status(
		c.ptr,
	))
}

// NewDiscoveryResponse creates a new discovery response
func NewDiscoveryResponse() *DiscoveryResponseBuilder {
	return newDiscoveryResponseBuilder(
		C.self_message_content_discovery_response_builder_init(),
	)
}

// ResponseTo sets the request id that is being responded to
func (b *DiscoveryResponseBuilder) ResponseTo(requestID []byte) *DiscoveryResponseBuilder {
	if len(requestID) != 20 {
		return b
	}

	requestIDBuf := C.CBytes(
		requestID,
	)

	C.self_message_content_discovery_response_builder_response_to(
		b.ptr,
		(*C.uint8_t)(requestIDBuf),
	)

	C.free(requestIDBuf)

	return b
}

// ResponseTo sets the request id that is being responded to
func (b *DiscoveryResponseBuilder) Status(status ResponseStatus) *DiscoveryResponseBuilder {
	C.self_message_content_discovery_response_builder_status(
		b.ptr,
		uint32(status),
	)

	return b
}

// Finish finalises the response and builds the content
func (b *DiscoveryResponseBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_discovery_response_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}
