package message

import "github.com/joinself/self-go-sdk/account"

type MessageType int

const (
	MessageTypeChat MessageType = iota
	MessageTypeChatReceiptDelivered
	MessageTypeChatReceiptRead
	MessageTypeConnectionRequest
	MessageTypeConnectionResponse
	MessageTypeCredentialVerificationRequest
	MessageTypeCredentialVerificationResponse
	MessageTypeCredentialPresentationRequest
	MessageTypeCredentialPresentationResponse
)

func Classify(message *account.Message) MessageType {
	return MessageTypeChat
}
