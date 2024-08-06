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

type Chat C.self_message_content_chat
type ChatBuilder C.self_message_content_chat_builder

// DeocodeChat decodes a chat message
func DecodeChat(msg *Message) (*Chat, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var chatContent *C.self_message_content_chat
	chatContentPtr := &chatContent

	status := C.self_message_content_as_chat(
		content,
		chatContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to decode chat message")
	}

	runtime.SetFinalizer(chatContentPtr, func(chat **C.self_message_content_chat) {
		C.self_message_content_chat_destroy(
			*chat,
		)
	})

	return (*Chat)(*chatContentPtr), nil
}

// Message returns the chat message
func (c *Chat) Message() string {
	return C.GoString(C.self_message_content_chat_message((*C.self_message_content_chat)(c)))
}

// NewChat constructs a new chat message
func NewChat() *ChatBuilder {
	builder := (*ChatBuilder)(C.self_message_content_chat_builder_init())

	runtime.SetFinalizer(builder, func(builder *ChatBuilder) {
		C.self_message_content_chat_builder_destroy(
			(*C.self_message_content_chat_builder)(builder),
		)
	})

	return builder
}

// Message sets the message on the chat message
func (b *ChatBuilder) Message(msg string) *ChatBuilder {
	cMsg := C.CString(msg)

	C.self_message_content_chat_builder_message(
		(*C.self_message_content_chat_builder)(b),
		cMsg,
	)

	C.free(unsafe.Pointer(cMsg))

	return b
}

// Finish finalizes the chat message and prepares it for sending
func (b *ChatBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content
	finishedContentPtr := &finishedContent

	status := C.self_message_content_chat_builder_finish(
		(*C.self_message_content_chat_builder)(b),
		finishedContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to build chat request")
	}

	runtime.SetFinalizer(finishedContentPtr, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(*finishedContentPtr), nil
}
