package message

import "github.com/joinself/self-go-sdk/account"

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

func Classify(message *account.Message) Type {
	return TypeChat
}
