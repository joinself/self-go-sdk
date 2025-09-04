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
	"unsafe"

	"github.com/joinself/self-go-sdk/status"
)

type Custom struct {
	ptr *C.self_message_content_custom
}

func newCustom(ptr *C.self_message_content_custom) *Custom {
	c := &Custom{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *Custom) {
		C.self_message_content_custom_destroy(
			c.ptr,
		)
	})

	return c
}

type CustomBuilder struct {
	ptr *C.self_message_content_custom_builder
}

func newCustomBuilder(ptr *C.self_message_content_custom_builder) *CustomBuilder {
	c := &CustomBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CustomBuilder) {
		C.self_message_content_custom_builder_destroy(
			c.ptr,
		)
	})

	return c
}

// DeocodeCustom decodes a custom message
func DecodeCustom(content *Content) (*Custom, error) {
	contentPtr := contentPtr(content)

	var customContent *C.self_message_content_custom

	result := C.self_message_content_as_custom(
		contentPtr,
		&customContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCustom(customContent), nil
}

// Payload returns the custom payload
func (c *Custom) Payload() []byte {
	payloadBuf := C.self_message_content_custom_payload(c.ptr)

	defer C.self_bytes_buffer_destroy(
		payloadBuf,
	)

	return C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(
			payloadBuf,
		)),
		C.int(C.self_bytes_buffer_len(
			payloadBuf,
		)),
	)
}

// NewCustom constructs a new custom message
func NewCustom() *CustomBuilder {
	return newCustomBuilder(C.self_message_content_custom_builder_init())
}

// Payload sets the messages payload
func (b *CustomBuilder) Payload(payload []byte) *CustomBuilder {
	payloadBuf := C.CBytes(payload)
	payloadLen := len(payload)

	C.self_message_content_custom_builder_payload(
		b.ptr,
		(*C.uint8_t)(payloadBuf),
		C.size_t(payloadLen),
	)

	C.free(unsafe.Pointer(payloadBuf))

	return b
}

// Finish finalizes the custom message and prepares it for sending
func (b *CustomBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_custom_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}
