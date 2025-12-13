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

	"github.com/joinself/self-go-sdk/status"
)

type ContentType int

type Content struct {
	ptr *C.self_message_content
}

const (
	ContentTypeUnknown                        ContentType = 1<<63 - 1
	ContentTypeCustom                         ContentType = C.CONTENT_CUSTOM
	ContentTypeChat                           ContentType = C.CONTENT_CHAT
	ContentTypeReceipt                        ContentType = C.CONTENT_RECEIPT
	ContentTypeCredential                     ContentType = C.CONTENT_CREDENTIAL
	ContentTypeIntroduction                   ContentType = C.CONTENT_INTRODUCTION
	ContentTypeDiscoveryRequest               ContentType = C.CONTENT_DISCOVERY_REQUEST
	ContentTypeDiscoveryResponse              ContentType = C.CONTENT_DISCOVERY_RESPONSE
	ContentTypeSigningRequest                 ContentType = C.CONTENT_SIGNING_REQUEST
	ContentTypeSigningResponse                ContentType = C.CONTENT_SIGNING_RESPONSE
	ContentTypeAccountPairingRequest          ContentType = C.CONTENT_ACCOUNT_PAIRING_REQUEST
	ContentTypeAccountPairingResponse         ContentType = C.CONTENT_ACCOUNT_PAIRING_RESPONSE
	ContentTypeCredentialVerificationRequest  ContentType = C.CONTENT_CREDENTIAL_VERIFICATION_REQUEST
	ContentTypeCredentialVerificationResponse ContentType = C.CONTENT_CREDENTIAL_VERIFICATION_RESPONSE
	ContentTypeCredentialPresentationRequest  ContentType = C.CONTENT_CREDENTIAL_PRESENTATION_REQUEST
	ContentTypeCredentialPresentationResponse ContentType = C.CONTENT_CREDENTIAL_PRESENTATION_RESPONSE
)

func (t ContentType) String() string {
	switch t {
	case ContentTypeCustom:
		return "Custom"
	case ContentTypeChat:
		return "Chat"
	case ContentTypeReceipt:
		return "Receipt"
	case ContentTypeCredential:
		return "Credential"
	case ContentTypeIntroduction:
		return "Introduction"
	case ContentTypeDiscoveryRequest:
		return "DiscoveryRequest"
	case ContentTypeDiscoveryResponse:
		return "DiscoveryResponse"
	case ContentTypeSigningRequest:
		return "SigningRequest"
	case ContentTypeSigningResponse:
		return "SigningResponse"
	case ContentTypeAccountPairingRequest:
		return "AccountPairingRequest"
	case ContentTypeAccountPairingResponse:
		return "AccountPairingResponse"
	case ContentTypeCredentialVerificationRequest:
		return "CredentialVerificationRequest"
	case ContentTypeCredentialVerificationResponse:
		return "CredentialVerificationResponse"
	case ContentTypeCredentialPresentationRequest:
		return "CredentialPresentationRequest"
	case ContentTypeCredentialPresentationResponse:
		return "CredentialPresentationResponse"
	default:
		return "Unknown"
	}
}

func newContent(ptr *C.self_message_content) *Content {
	c := &Content{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(c *Content) {
		C.self_message_content_destroy(
			c.ptr,
		)
	}, c)

	return c
}

func contentPtr(c *Content) *C.self_message_content {
	return c.ptr
}

func contentType(content *C.self_message_content) ContentType {
	switch C.self_message_content_type_of(content) {
	case C.CONTENT_CUSTOM:
		return ContentTypeCustom
	case C.CONTENT_CHAT:
		return ContentTypeChat
	case C.CONTENT_RECEIPT:
		return ContentTypeReceipt
	case C.CONTENT_CREDENTIAL:
		return ContentTypeCredential
	case C.CONTENT_INTRODUCTION:
		return ContentTypeIntroduction
	case C.CONTENT_DISCOVERY_REQUEST:
		return ContentTypeDiscoveryRequest
	case C.CONTENT_DISCOVERY_RESPONSE:
		return ContentTypeDiscoveryResponse
	case C.CONTENT_SIGNING_REQUEST:
		return ContentTypeSigningRequest
	case C.CONTENT_SIGNING_RESPONSE:
		return ContentTypeSigningResponse
	case C.CONTENT_ACCOUNT_PAIRING_REQUEST:
		return ContentTypeAccountPairingRequest
	case C.CONTENT_ACCOUNT_PAIRING_RESPONSE:
		return ContentTypeAccountPairingResponse
	case C.CONTENT_CREDENTIAL_VERIFICATION_REQUEST:
		return ContentTypeCredentialVerificationRequest
	case C.CONTENT_CREDENTIAL_VERIFICATION_RESPONSE:
		return ContentTypeCredentialVerificationResponse
	case C.CONTENT_CREDENTIAL_PRESENTATION_REQUEST:
		return ContentTypeCredentialPresentationRequest
	case C.CONTENT_CREDENTIAL_PRESENTATION_RESPONSE:
		return ContentTypeCredentialPresentationResponse
	default:
		return ContentTypeUnknown
	}
}

// ID returns the unique id of the message
func (c *Content) ID() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_id(c.ptr)),
		20,
	)
}

// ContentType get the content type
func (c *Content) ContentType() ContentType {
	return contentType(c.ptr)
}

// Summary creates a summary of the content used for sending a push notification
func (c *Content) Summary() (*ContentSummary, error) {
	var summary *C.self_message_content_summary

	result := C.self_message_content_summary_of(
		c.ptr,
		&summary,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContentSummary(summary), nil
}
