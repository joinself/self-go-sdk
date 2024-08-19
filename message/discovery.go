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
)

type DiscoveryRequest C.self_message_content_discovery_request
type DiscoveryResponse C.self_message_content_discovery_response
type DiscoveryRequestBuilder C.self_message_content_discovery_request_builder
type DiscoveryResponseBuilder C.self_message_content_discovery_response_builder

// DecodeDiscoveryRequest decodes a message to a discovery request
func DecodeDiscoveryRequest(msg *Message) (*DiscoveryRequest, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var discoveryRequestContent *C.self_message_content_discovery_request
	discoveryRequestContentPtr := &discoveryRequestContent

	status := C.self_message_content_as_discovery_request(
		content,
		discoveryRequestContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to decode discovery request message")
	}

	runtime.SetFinalizer(discoveryRequestContentPtr, func(discoveryRequest **C.self_message_content_discovery_request) {
		C.self_message_content_discovery_request_destroy(
			*discoveryRequest,
		)
	})

	return (*DiscoveryRequest)(*discoveryRequestContentPtr), nil
}

// KeyPackage returns the embedded key package conntained in the discovery request
func (c *DiscoveryRequest) KeyPackage() *KeyPackage {
	keyPackage := C.self_message_content_discovery_request_key_package(
		(*C.self_message_content_discovery_request)(c),
	)

	runtime.SetFinalizer(keyPackage, func(keyPackage *C.self_key_package) {
		C.self_key_package_destroy(
			keyPackage,
		)
	})

	return (*KeyPackage)(keyPackage)
}

// Type returns the time the request expires at
func (c *DiscoveryRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_discovery_request_expires(
		(*C.self_message_content_discovery_request)(c),
	)), 0)
}

// NewDiscoveryRequest creates a new discovery request
func NewDiscoveryRequest() *DiscoveryRequestBuilder {
	builder := C.self_message_content_discovery_request_builder_init()

	runtime.SetFinalizer(builder, func(builder *C.self_message_content_discovery_request_builder) {
		C.self_message_content_discovery_request_builder_destroy(
			builder,
		)
	})

	return (*DiscoveryRequestBuilder)(builder)
}

// KeyPackage sets the key package that will be embedded in the request
func (b *DiscoveryRequestBuilder) KeyPackage(keyPackage *KeyPackage) *DiscoveryRequestBuilder {
	C.self_message_content_discovery_request_builder_key_package(
		(*C.self_message_content_discovery_request_builder)(b),
		(*C.self_key_package)(keyPackage),
	)
	return b
}

// Expires sets the time that the request expires at
func (b *DiscoveryRequestBuilder) Expires(expires time.Time) *DiscoveryRequestBuilder {
	C.self_message_content_discovery_request_builder_expires(
		(*C.self_message_content_discovery_request_builder)(b),
		C.long(expires.Unix()),
	)
	return b
}

// Finish finalises the request and builds the content
func (b *DiscoveryRequestBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content
	finishedContentPtr := &finishedContent

	status := C.self_message_content_discovery_request_builder_finish(
		(*C.self_message_content_discovery_request_builder)(b),
		finishedContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to build discovery request")
	}

	runtime.SetFinalizer(finishedContentPtr, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(*finishedContentPtr), nil
}

// DecodeDiscoveryResponse decodes a message to a discovery response
func DecodeDiscoveryResponse(msg *Message) (*DiscoveryResponse, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var discoveryResponseContent *C.self_message_content_discovery_response
	discoveryResponseContentPtr := &discoveryResponseContent

	status := C.self_message_content_as_discovery_response(
		content,
		discoveryResponseContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to decode discovery response message")
	}

	runtime.SetFinalizer(discoveryResponseContentPtr, func(discoveryResponse **C.self_message_content_discovery_response) {
		C.self_message_content_discovery_response_destroy(
			*discoveryResponse,
		)
	})

	return (*DiscoveryResponse)(*discoveryResponseContentPtr), nil
}

// ResponseTo returns the id of the request that is being responded to
func (c *DiscoveryResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_discovery_response_response_to(
			(*C.self_message_content_discovery_response)(c),
		)),
		20,
	)
}

// Status returns the status of the request
func (c *DiscoveryResponse) Status() ResponseStatus {
	return ResponseStatus(C.self_message_content_discovery_response_status(
		(*C.self_message_content_discovery_response)(c),
	))
}

// NewDiscoveryResponse creates a new discovery response
func NewDiscoveryResponse() *DiscoveryResponseBuilder {
	builder := C.self_message_content_discovery_response_builder_init()

	runtime.SetFinalizer(builder, func(builder *C.self_message_content_discovery_response_builder) {
		C.self_message_content_discovery_response_builder_destroy(
			builder,
		)
	})

	return (*DiscoveryResponseBuilder)(builder)
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
		(*C.self_message_content_discovery_response_builder)(b),
		(*C.uint8_t)(requestIDBuf),
	)

	C.free(requestIDBuf)

	return b
}

// ResponseTo sets the request id that is being responded to
func (b *DiscoveryResponseBuilder) Status(status ResponseStatus) *DiscoveryResponseBuilder {
	C.self_message_content_discovery_response_builder_status(
		(*C.self_message_content_discovery_response_builder)(b),
		uint32(status),
	)

	return b
}

// Finish finalises the response and builds the content
func (b *DiscoveryResponseBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content
	finishedContentPtr := &finishedContent

	status := C.self_message_content_discovery_response_builder_finish(
		(*C.self_message_content_discovery_response_builder)(b),
		finishedContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to build discovery response")
	}

	runtime.SetFinalizer(finishedContentPtr, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(*finishedContentPtr), nil
}
