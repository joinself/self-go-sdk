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

	"github.com/joinself/self-go-sdk-next/status"
)

type QREncoding int

const (
	QREncodingSVG     QREncoding = C.QR_SVG
	QREncodingUnicode QREncoding = C.QR_UNICODE
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
func NewAnonymousMessage(content *Content) *AnonymousMessage {
	return newAnonymousMessage(C.self_anonymous_message_init(
		content.ptr,
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
func (a *AnonymousMessage) Content() *Content {
	return newContent(
		C.self_anonymous_message_message_content(a.ptr),
	)
}

// Content returns the messages content
func (a *AnonymousMessage) EncodeToQR(encoding QREncoding) ([]byte, error) {
	var qrCode *C.self_encoded_buffer

	result := C.self_anonymous_message_encode_as_qr(
		a.ptr,
		&qrCode,
		uint32(encoding),
	)

	if result > 0 {
		return nil, status.New(result)
	}

	encodedQR := C.GoBytes(
		unsafe.Pointer(C.self_encoded_buffer_buf(qrCode)),
		C.int(C.self_encoded_buffer_len(qrCode)),
	)

	C.self_encoded_buffer_destroy(
		qrCode,
	)

	return encodedQR, nil
}

// EncodeToString encodes a message to an encoded string
func (a *AnonymousMessage) EncodeToString() (string, error) {
	var encodeBuffer *C.self_encoded_buffer

	result := C.self_anonymous_message_encode(
		a.ptr,
		&encodeBuffer,
	)

	if result > 0 {
		return "", status.New(result)
	}

	encodedMessage := C.GoBytes(
		unsafe.Pointer(C.self_encoded_buffer_buf(encodeBuffer)),
		C.int(C.self_encoded_buffer_len(encodeBuffer)),
	)

	C.self_encoded_buffer_destroy(
		encodeBuffer,
	)

	return base64.RawURLEncoding.EncodeToString(encodedMessage), nil
}

// AnonymousMessageDecode decodes a message from an encoded string
func AnonymousMessageDecode(encoded string) (*AnonymousMessage, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	decodedBuf := C.CBytes(decoded)

	var anonymousMessage *C.self_anonymous_message

	result := C.self_anonymous_message_decode(
		&anonymousMessage,
		(*C.uint8_t)(decodedBuf),
		C.size_t(len(decoded)),
	)

	C.free(decodedBuf)

	if result > 0 {
		return nil, status.New(result)
	}

	return newAnonymousMessage(anonymousMessage), nil
}
