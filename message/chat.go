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

	chat := (*Chat)(chatContent)

	runtime.SetFinalizer(chat, func(chat *Chat) {
		C.self_message_content_chat_destroy(
			(*C.self_message_content_chat)(chat),
		)
	})

	return chat, nil
}

// Message returns the chat message
func (c *Chat) Message() string {
	return C.GoString(C.self_message_content_chat_message((*C.self_message_content_chat)(c)))
}

type ChatBuilder C.self_message_content_chat_builder

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
func (b *ChatBuilder) Finish() *Content {
	var finishedContent *C.self_message_content

	C.self_message_content_chat_builder_finish(
		(*C.self_message_content_chat_builder)(b),
		&finishedContent,
	)

	content := (*Content)(finishedContent)

	runtime.SetFinalizer(content, func(content *Content) {
		C.self_message_content_destroy(
			(*C.self_message_content)(content),
		)
	})

	return content
}
