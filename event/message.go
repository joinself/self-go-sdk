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
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/platform"
	"github.com/joinself/self-go-sdk/token"
)

//go:linkname newPlatformAttestation github.com/joinself/self-go-sdk/platform.newPlatformAttestation
func newPlatformAttestation(ptr *C.self_platform_attestation) *platform.Attestation

//go:linkname newToken github.com/joinself/self-go-sdk/token.newToken
func newToken(ptr *C.self_token) *token.Token

//go:linkname contentType github.com/joinself/self-go-sdk/message.contentType
func contentType(ptr *C.self_message_content) message.ContentType

type Message struct {
	ptr *C.self_message
}

func newMessage(ptr *C.self_message) *Message {
	e := &Message{
		ptr: ptr,
	}

	runtime.AddCleanup(e, func(ptr *C.self_message) {
		C.self_message_destroy(
			ptr,
		)
	}, e.ptr)

	return e
}

// ContentTypeOf get the content type of the message
func ContentTypeOf(message *Message) message.ContentType {
	return contentType(C.self_message_message_content(message.ptr))
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
func (m *Message) Content() *message.Content {
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

// Integrity returns an integrity check performed over the contents of the message
func (m *Message) Integrity() (*platform.Attestation, bool) {
	integrity := C.self_message_message_integrity(m.ptr)
	if integrity == nil {
		return nil, false
	}

	return newPlatformAttestation(
		integrity,
	), true
}

// Tokens returns tokens attached to the message
func (m *Message) Tokens() []*token.Token {
	collection := C.self_message_tokens(
		m.ptr,
	)

	var tokens []*token.Token

	for i := 0; i < int(C.self_collection_token_len(collection)); i++ {
		tokens = append(tokens, newToken(
			C.self_collection_token_at(collection, C.size_t(i)),
		))
	}

	C.self_collection_token_destroy(collection)

	return tokens
}

// MerkleRoot returns a merkle root from an attached merkle proof, if provided
func (m *Message) MerkleRoot() []byte {
	buf := C.self_message_merkle_root(m.ptr)
	if buf == nil {
		return nil
	}

	merkleRoot := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(buf)),
		C.int(C.self_bytes_buffer_len(buf)),
	)

	C.self_bytes_buffer_destroy(buf)

	return merkleRoot
}
