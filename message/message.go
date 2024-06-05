package message

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "github.com/joinself/self-go-sdk/keypair/signing"

type Type int

const (
	TypeChat Type = iota
	TypeChatReceiptDelivered
	TypeChatReceiptRead
	TypeConnectionRequest
	TypeConnectionResponse
	TypeCredentialVerificationRequest
	TypeCredentialVerificationResponse
	TypeCredentialPresentationRequest
	TypeCredentialPresentationResponse
)

func ContentType(message *Message) Type {
	content := C.self_message_message_content((*C.self_message)(message))

	switch C.self_message_content_type_of(content) {
	case C.CONTENT_CHAT:
		return TypeChat
	default:
		return C.CONTENT_CUSTOM
	}
}

type Message C.self_message

func (m *Message) FromAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_message_from_address((*C.self_message)(m)))
}

func (m *Message) ToAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_message_to_address((*C.self_message)(m)))
}

func (m *Message) Content() *Content {
	return (*Content)(C.self_message_message_content((*C.self_message)(m)))
}
