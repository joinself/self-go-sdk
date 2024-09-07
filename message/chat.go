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
	"unsafe"
)

type Chat struct {
	ptr *C.self_message_content_chat
}

func newChat(ptr *C.self_message_content_chat) *Chat {
	c := &Chat{
		ptr: ptr,
	}

	runtime.SetFinalizer(&c, func(c *Chat) {
		C.self_message_content_chat_destroy(
			c.ptr,
		)
	})

	return c
}

type ChatBuilder struct {
	ptr *C.self_message_content_chat_builder
}

func newChatBuilder(ptr *C.self_message_content_chat_builder) *ChatBuilder {
	c := &ChatBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(&c, func(c *ChatBuilder) {
		C.self_message_content_chat_builder_destroy(
			c.ptr,
		)
	})

	return c
}

// DeocodeChat decodes a chat message
func DecodeChat(msg *Message) (*Chat, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var chatContent *C.self_message_content_chat

	status := C.self_message_content_as_chat(
		content,
		&chatContent,
	)

	if status > 0 {
		return nil, errors.New("failed to decode chat message")
	}

	return newChat(chatContent), nil
}

// Message returns the chat message
func (c *Chat) Message() string {
	return C.GoString(C.self_message_content_chat_message(c.ptr))
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

// Finish finalizes the chat message and prepares it for sending
func (b *ChatBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	status := C.self_message_content_chat_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if status > 0 {
		return nil, errors.New("failed to build chat request")
	}

	return newContent(finishedContent), nil
}
