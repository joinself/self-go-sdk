package event

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"encoding/base64"
	"runtime"
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname contentPtr github.com/joinself/self-go-sdk/message.contentPtr
func contentPtr(c *message.Content) *C.self_message_content

//go:linkname newContent github.com/joinself/self-go-sdk/message.newContent
func newContent(ptr *C.self_message_content) *message.Content

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk/keypair/signing.newSigningPublicKey
func newSigningPublicKey(ptr *C.self_signing_public_key) *signing.PublicKey

type QREncoding int
type MessageFlag uint64

const (
	QREncodingSVG            QREncoding  = C.QR_SVG
	QREncodingUnicode        QREncoding  = C.QR_UNICODE
	MessageFlagTargetSandbox MessageFlag = C.MESSAGE_FLAG_TARGET_SANDBOX
)

type AnonymousMessage struct {
	ptr *C.self_anonymous_message
}

func newAnonymousMessage(ptr *C.self_anonymous_message) *AnonymousMessage {
	e := &AnonymousMessage{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *AnonymousMessage) {
		C.self_anonymous_message_destroy(
			e.ptr,
		)
	})

	return e
}

// NewAnonymousMessage creates a new anonymous message from content
func NewAnonymousMessage(content *message.Content) *AnonymousMessage {
	return newAnonymousMessage(C.self_anonymous_message_init(
		contentPtr(content),
	))
}

// DecodeAnonymousMessage decodes an anonymous message
func DecodeAnonymousMessage(data []byte) (*AnonymousMessage, error) {
	var anonymousMessage *C.self_anonymous_message

	dataBuf := C.CBytes(data)
	dataLen := len(data)
	defer C.free(dataBuf)

	result := C.self_anonymous_message_decode(
		&anonymousMessage,
		(*C.uint8_t)(dataBuf),
		C.size_t(dataLen),
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newAnonymousMessage(anonymousMessage), nil
}

// ID returns the id of the messages content
func (a *AnonymousMessage) ID() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_anonymous_message_id(a.ptr)),
		20,
	)
}

// Content returns the messages content
func (a *AnonymousMessage) Content() *message.Content {
	return newContent(
		C.self_anonymous_message_message_content(a.ptr),
	)
}

// Flags returns the messages content
func (a *AnonymousMessage) Flags() MessageFlag {
	return MessageFlag(C.self_anonymous_message_flags(a.ptr))
}

// SetFlags sets the messages flags
func (a *AnonymousMessage) SetFlags(flags MessageFlag) *AnonymousMessage {
	C.self_anonymous_message_set_flags(a.ptr, C.uint64_t(flags))
	return a
}

// HasFlags returns true if the message has flags
func (a *AnonymousMessage) HasFlags(flags MessageFlag) bool {
	return bool(C.self_anonymous_message_has_flags(a.ptr, C.uint64_t(flags)))
}

// Content returns the messages content
func (a *AnonymousMessage) EncodeToQR(encoding QREncoding) ([]byte, error) {
	var qrCode *C.self_bytes_buffer

	result := C.self_anonymous_message_encode_as_qr(
		a.ptr,
		&qrCode,
		uint32(encoding),
	)

	if result > 0 {
		return nil, status.New(result)
	}

	encodedQR := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(qrCode)),
		C.int(C.self_bytes_buffer_len(qrCode)),
	)

	C.self_bytes_buffer_destroy(
		qrCode,
	)

	return encodedQR, nil
}

// EncodeToString encodes a message to an encoded string
func (a *AnonymousMessage) EncodeToString() (string, error) {
	var encodeBuffer *C.self_bytes_buffer

	result := C.self_anonymous_message_encode(
		a.ptr,
		&encodeBuffer,
	)

	if result > 0 {
		return "", status.New(result)
	}

	encodedMessage := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(encodeBuffer)),
		C.int(C.self_bytes_buffer_len(encodeBuffer)),
	)

	C.self_bytes_buffer_destroy(
		encodeBuffer,
	)

	return base64.RawURLEncoding.EncodeToString(encodedMessage), nil
}

// AnonymousMessageDecodeFromString decodes a message from an encoded string
func AnonymousMessageDecodeFromString(encoded string) (*AnonymousMessage, error) {
	var anonymousMessage *C.self_anonymous_message

	encodedPtr := C.CString(encoded)

	result := C.self_anonymous_message_decode_from_string(
		&anonymousMessage,
		encodedPtr,
	)

	C.free(unsafe.Pointer(encodedPtr))

	if result > 0 {
		return nil, status.New(result)
	}

	return newAnonymousMessage(anonymousMessage), nil
}
