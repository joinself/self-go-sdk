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

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk-next/keypair/signing.newSigningPublicKey
func newSigningPublicKey(*C.self_signing_public_key) *signing.PublicKey

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

type Commit struct {
	ptr *C.self_commit
}

func newCommit(ptr *C.self_commit) *Commit {
	e := &Commit{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *Commit) {
		C.self_commit_destroy(
			e.ptr,
		)
	})

	return e
}

type KeyPackage struct {
	ptr *C.self_key_package
}

func newKeyPackage(ptr *C.self_key_package) *KeyPackage {
	e := &KeyPackage{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *KeyPackage) {
		C.self_key_package_destroy(
			e.ptr,
		)
	})

	return e
}

func keyPackagePtr(k *KeyPackage) *C.self_key_package {
	return k.ptr
}

type Message struct {
	ptr *C.self_message
}

func newMessage(ptr *C.self_message) *Message {
	e := &Message{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *Message) {
		C.self_message_destroy(
			e.ptr,
		)
	})

	return e
}

type Proposal struct {
	ptr *C.self_proposal
}

func newProposal(ptr *C.self_proposal) *Proposal {
	e := &Proposal{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *Proposal) {
		C.self_proposal_destroy(
			e.ptr,
		)
	})

	return e
}

type Welcome struct {
	ptr *C.self_welcome
}

func newWelcome(ptr *C.self_welcome) *Welcome {
	e := &Welcome{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *Welcome) {
		C.self_welcome_destroy(
			e.ptr,
		)
	})

	return e
}

func welcomePtr(w *Welcome) *C.self_welcome {
	return w.ptr
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

	status := C.self_anonymous_message_decode(
		&anonymousMessage,
		(*C.uint8_t)(dataBuf),
		C.size_t(dataLen),
	)

	if status > 0 {
		return nil, errors.New("failed to decode anonymous message")
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
	qrCodePtr := &qrCode

	status := C.self_anonymous_message_encode_as_qr(
		a.ptr,
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
	return newSigningPublicKey(C.self_commit_to_address(
		c.ptr,
	))
}

// FromAddress returns the address the event was sent by
func (c *Commit) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_commit_from_address(
		c.ptr,
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Commit) Sequence() uint64 {
	return uint64(C.self_commit_sequence(
		c.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Commit) Timestamp() time.Time {
	return time.Unix(int64(C.self_commit_timestamp(
		c.ptr,
	)), 0)
}

// ToAddress returns the address the event was addressed to
func (c *KeyPackage) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_key_package_to_address(
		c.ptr,
	))
}

// FromAddress returns the address the event was sent by
func (c *KeyPackage) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_key_package_from_address(
		c.ptr,
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *KeyPackage) Sequence() uint64 {
	return uint64(C.self_key_package_sequence(
		c.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *KeyPackage) Timestamp() time.Time {
	return time.Unix(int64(C.self_key_package_timestamp(
		c.ptr,
	)), 0)
}

// ID returns the id of the messages content
func (m *Message) ID() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_id(m.ptr)),
		20,
	)
}

// FromAddress returns the address the event was sent by
func (m *Message) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(
		C.self_message_from_address(m.ptr),
	)
}

// ToAddress returns the address the event was addressed to
func (m *Message) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(
		C.self_message_to_address(m.ptr),
	)
}

// Content returns the messages content
func (m *Message) Content() *Content {
	return newContent(
		C.self_message_message_content(m.ptr),
	)
}

// ToAddress returns the address the event was addressed to
func (c *Proposal) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(
		C.self_proposal_to_address(c.ptr),
	)
}

// FromAddress returns the address the event was sent by
func (c *Proposal) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_proposal_from_address(
		c.ptr,
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Proposal) Sequence() uint64 {
	return uint64(C.self_proposal_sequence(
		c.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Proposal) Timestamp() time.Time {
	return time.Unix(int64(C.self_proposal_timestamp(
		c.ptr,
	)), 0)
}

// ToAddress returns the address the event was addressed to
func (c *Welcome) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_welcome_to_address(
		c.ptr,
	))
}

// FromAddress returns the address the event was sent by
func (c *Welcome) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_welcome_from_address(
		c.ptr,
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Welcome) Sequence() uint64 {
	return uint64(C.self_welcome_sequence(
		c.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Welcome) Timestamp() time.Time {
	return time.Unix(int64(C.self_welcome_timestamp(
		c.ptr,
	)), 0)
}
