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

	"github.com/joinself/self-go-sdk-next/keypair/signing"
)

type QREncoding int

const (
	QREncodingSVG     QREncoding = C.QR_SVG
	QREncodingUnicode QREncoding = C.QR_UNICODE
)

type AnonymousMessage C.self_anonymous_message
type Commit C.self_commit
type KeyPackage C.self_key_package
type Message C.self_message
type Proposal C.self_proposal
type Welcome C.self_welcome

func NewAnonymousMessage(content *Content) *AnonymousMessage {
	anonymousMessage := C.self_anonymous_message_init(
		(*C.self_message_content)(content),
	)

	runtime.SetFinalizer(&anonymousMessage, func(anonymousMessage **C.self_anonymous_message) {
		C.self_anonymous_message_destroy(
			*anonymousMessage,
		)
	})

	return (*AnonymousMessage)(anonymousMessage)
}

func DecodeAnonymousMessage(data []byte) (*AnonymousMessage, error) {
	var anonymousMessage *C.self_anonymous_message
	anonymousMessagePtr := &anonymousMessage

	dataBuf := C.CBytes(data)
	dataLen := len(data)
	defer C.free(dataBuf)

	status := C.self_anonymous_message_decode(
		anonymousMessagePtr,
		(*C.uint8_t)(dataBuf),
		C.size_t(dataLen),
	)

	if status > 0 {
		return nil, errors.New("failed to decode anonymous message")
	}

	runtime.SetFinalizer(anonymousMessagePtr, func(anonymousMessagePtr **C.self_anonymous_message) {
		C.self_anonymous_message_destroy(
			*anonymousMessagePtr,
		)
	})

	return (*AnonymousMessage)(*anonymousMessagePtr), nil
}

// ID returns the id of the messages content
func (a *AnonymousMessage) ID() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_anonymous_message_id((*C.self_anonymous_message)(a))),
		20,
	)
}

// Content returns the messages content
func (a *AnonymousMessage) Content() *Content {
	content := C.self_anonymous_message_message_content((*C.self_anonymous_message)(a))

	runtime.SetFinalizer(&content, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(content)
}

// Content returns the messages content
func (a *AnonymousMessage) EncodeToQR(encoding QREncoding) ([]byte, error) {
	var qrCode *C.self_encoded_buffer
	qrCodePtr := &qrCode

	status := C.self_anonymous_message_encode_as_qr(
		(*C.self_anonymous_message)(a),
		qrCodePtr,
		uint32(encoding),
	)

	if status > 0 {
		return nil, errors.New("failed to encode QR code")
	}

	encodedQR := C.GoBytes(
		unsafe.Pointer(C.self_encoded_buffer_buf(*qrCodePtr)),
		C.int(C.self_encoded_buffer_len(*qrCodePtr)),
	)

	C.self_encoded_buffer_destroy(
		*qrCodePtr,
	)

	return encodedQR, nil
}

// ToAddress returns the address the event was addressed to
func (c *Commit) ToAddress() *signing.PublicKey {
	address := C.self_commit_to_address(
		(*C.self_commit)(c),
	)

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// FromAddress returns the address the event was sent by
func (c *Commit) FromAddress() *signing.PublicKey {
	address := C.self_commit_from_address(
		(*C.self_commit)(c),
	)

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Commit) Sequence() uint64 {
	return uint64(C.self_commit_sequence(
		(*C.self_commit)(c),
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Commit) Timestamp() time.Time {
	return time.Unix(int64(C.self_commit_timestamp(
		(*C.self_commit)(c),
	)), 0)
}

// ToAddress returns the address the event was addressed to
func (c *KeyPackage) ToAddress() *signing.PublicKey {
	address := C.self_key_package_to_address(
		(*C.self_key_package)(c),
	)

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// FromAddress returns the address the event was sent by
func (c *KeyPackage) FromAddress() *signing.PublicKey {
	address := C.self_key_package_from_address(
		(*C.self_key_package)(c),
	)

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *KeyPackage) Sequence() uint64 {
	return uint64(C.self_key_package_sequence(
		(*C.self_key_package)(c),
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *KeyPackage) Timestamp() time.Time {
	return time.Unix(int64(C.self_key_package_timestamp(
		(*C.self_key_package)(c),
	)), 0)
}

// ID returns the id of the messages content
func (m *Message) ID() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_id((*C.self_message)(m))),
		20,
	)
}

// FromAddress returns the address the event was sent by
func (m *Message) FromAddress() *signing.PublicKey {
	address := C.self_message_from_address((*C.self_message)(m))

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// ToAddress returns the address the event was addressed to
func (m *Message) ToAddress() *signing.PublicKey {
	address := C.self_message_to_address((*C.self_message)(m))

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// Content returns the messages content
func (m *Message) Content() *Content {
	content := C.self_message_message_content((*C.self_message)(m))

	runtime.SetFinalizer(&content, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(content)
}

// ToAddress returns the address the event was addressed to
func (c *Proposal) ToAddress() *signing.PublicKey {
	address := C.self_proposal_to_address(
		(*C.self_proposal)(c),
	)

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// FromAddress returns the address the event was sent by
func (c *Proposal) FromAddress() *signing.PublicKey {
	address := C.self_proposal_from_address(
		(*C.self_proposal)(c),
	)

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Proposal) Sequence() uint64 {
	return uint64(C.self_proposal_sequence(
		(*C.self_proposal)(c),
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Proposal) Timestamp() time.Time {
	return time.Unix(int64(C.self_proposal_timestamp(
		(*C.self_proposal)(c),
	)), 0)
}

// ToAddress returns the address the event was addressed to
func (c *Welcome) ToAddress() *signing.PublicKey {
	address := C.self_welcome_to_address(
		(*C.self_welcome)(c),
	)

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// FromAddress returns the address the event was sent by
func (c *Welcome) FromAddress() *signing.PublicKey {
	address := C.self_welcome_from_address(
		(*C.self_welcome)(c),
	)

	runtime.SetFinalizer(&address, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(address)
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Welcome) Sequence() uint64 {
	return uint64(C.self_welcome_sequence(
		(*C.self_welcome)(c),
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Welcome) Timestamp() time.Time {
	return time.Unix(int64(C.self_welcome_timestamp(
		(*C.self_welcome)(c),
	)), 0)
}
