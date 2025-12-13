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

	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/platform"
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk/keypair/signing.newSigningPublicKey
func newSigningPublicKey(*C.self_signing_public_key) *signing.PublicKey

//go:linkname newPlatformAttestation github.com/joinself/self-go-sdk/platform.newPlatformAttestation
func newPlatformAttestation(*C.self_platform_attestation) *platform.Attestation

type Chat struct {
	ptr *C.self_message_content_chat
}

func newChat(ptr *C.self_message_content_chat) *Chat {
	c := &Chat{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(c *Chat) {
		C.self_message_content_chat_destroy(
			c.ptr,
		)
	}, c)

	return c
}

type ChatBuilder struct {
	ptr *C.self_message_content_chat_builder
}

func newChatBuilder(ptr *C.self_message_content_chat_builder) *ChatBuilder {
	c := &ChatBuilder{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(c *ChatBuilder) {
		C.self_message_content_chat_builder_destroy(
			c.ptr,
		)
	}, c)

	return c
}

// DeocodeChat decodes a chat message
func DecodeChat(content *Content) (*Chat, error) {
	contentPtr := contentPtr(content)

	var chatContent *C.self_message_content_chat

	result := C.self_message_content_as_chat(
		contentPtr,
		&chatContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newChat(chatContent), nil
}

// Message returns the chat message
func (c *Chat) Message() string {
	return C.GoString(C.self_message_content_chat_message(c.ptr))
}

// Referencing returns the id the message is replying to
func (c *Chat) Referencing() []byte {
	referencing := C.self_message_content_chat_referencing(c.ptr)
	if referencing == nil {
		return nil
	}

	return C.GoBytes(
		unsafe.Pointer(referencing),
		20,
	)
}

// Attachments returns the attachments
func (c *Chat) Attachments() []*object.Object {
	collection := C.self_message_content_chat_attachments(c.ptr)

	var attachments []*object.Object

	for i := 0; i < int(C.self_collection_object_len(collection)); i++ {
		attachments = append(attachments, newObject(
			C.self_collection_object_at(collection, C.size_t(i)),
		))
	}

	C.self_collection_object_destroy(collection)

	return attachments
}

// NewChat constructs a new chat message
func NewChat() *ChatBuilder {
	return newChatBuilder(C.self_message_content_chat_builder_init())
}

// Message sets the message on the chat message
func (b *ChatBuilder) Message(msg string) *ChatBuilder {
	cMsg := C.CString(msg)

	C.self_message_content_chat_builder_message(
		b.ptr,
		cMsg,
	)

	C.free(unsafe.Pointer(cMsg))

	return b
}

// Reference references a message in the chat
func (b *ChatBuilder) Reference(messageID []byte) *ChatBuilder {
	cID := C.CBytes(messageID)

	C.self_message_content_chat_builder_reference(
		b.ptr,
		(*C.uint8_t)(cID),
	)

	C.free(unsafe.Pointer(cID))

	return b
}

// Attach attaches an object to the message
func (b *ChatBuilder) Attach(attachment *object.Object) *ChatBuilder {
	C.self_message_content_chat_builder_attach(
		b.ptr,
		objectPtr(attachment),
	)

	return b
}

// Finish finalizes the chat message and prepares it for sending
func (b *ChatBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_chat_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}
