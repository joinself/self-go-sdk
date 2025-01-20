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
	"runtime"
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/platform"
)

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

// Content returns the sha3 hash of the encoded content
func (m *Message) ContentHash() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_message_content_hash(m.ptr)),
		32,
	)
}

// Content returns the sha3 hash of the encoded content
func (m *Message) Integrity() (*platform.Attestation, bool) {
	integrity := C.self_message_message_integrity(m.ptr)
	if integrity == nil {
		return nil, false
	}

	return newPlatformAttestation(
		integrity,
	), true
}
