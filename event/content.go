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

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk/keypair/signing.newSigningPublicKey
func newSigningPublicKey(*C.self_signing_public_key) *signing.PublicKey

//go:linkname newPlatformAttestation github.com/joinself/self-go-sdk/platform.newPlatformAttestation
func newPlatformAttestation(*C.self_platform_attestation) *platform.Attestation

type Type int

type Content struct {
	ptr *C.self_message_content
}

const (
	TypeUnknown                        Type = 1<<63 - 1
	TypeCustom                         Type = C.CONTENT_CUSTOM
	TypeChat                           Type = C.CONTENT_CHAT
	TypeReceipt                        Type = C.CONTENT_RECEIPT
	TypeIntroduction                   Type = C.CONTENT_INTRODUCTION
	TypeDiscoveryRequest               Type = C.CONTENT_DISCOVERY_REQUEST
	TypeDiscoveryResponse              Type = C.CONTENT_DISCOVERY_RESPONSE
	TypeSigningRequest                 Type = C.CONTENT_SIGNING_REQUEST
	TypeSigningResponse                Type = C.CONTENT_SIGNING_RESPONSE
	TypeAccountPairingRequest          Type = C.CONTENT_ACCOUNT_PAIRING_REQUEST
	TypeAccountPairingResponse         Type = C.CONTENT_ACCOUNT_PAIRING_RESPONSE
	TypeCredentialVerificationRequest  Type = C.CONTENT_CREDENTIAL_VERIFICATION_REQUEST
	TypeCredentialVerificationResponse Type = C.CONTENT_CREDENTIAL_VERIFICATION_RESPONSE
	TypeCredentialPresentationRequest  Type = C.CONTENT_CREDENTIAL_PRESENTATION_REQUEST
	TypeCredentialPresentationResponse Type = C.CONTENT_CREDENTIAL_PRESENTATION_RESPONSE
)

func newContent(ptr *C.self_message_content) *Content {
	c := &Content{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *Content) {
		C.self_message_content_destroy(
			c.ptr,
		)
	})

	return c
}

func contentPtr(c *Content) *C.self_message_content {
	return c.ptr
}

// ContentType get the content type of the message
func ContentType(message *Message) Type {
	content := C.self_message_message_content(message.ptr)

	switch C.self_message_content_type_of(content) {
	case C.CONTENT_CUSTOM:
		return TypeCustom
	case C.CONTENT_CHAT:
		return TypeChat
	case C.CONTENT_RECEIPT:
		return TypeReceipt
	case C.CONTENT_INTRODUCTION:
		return TypeIntroduction
	case C.CONTENT_DISCOVERY_REQUEST:
		return TypeDiscoveryRequest
	case C.CONTENT_DISCOVERY_RESPONSE:
		return TypeDiscoveryResponse
	case C.CONTENT_SIGNING_REQUEST:
		return TypeSigningRequest
	case C.CONTENT_SIGNING_RESPONSE:
		return TypeSigningResponse
	case C.CONTENT_ACCOUNT_PAIRING_REQUEST:
		return TypeAccountPairingRequest
	case C.CONTENT_ACCOUNT_PAIRING_RESPONSE:
		return TypeAccountPairingResponse
	case C.CONTENT_CREDENTIAL_VERIFICATION_REQUEST:
		return TypeCredentialVerificationRequest
	case C.CONTENT_CREDENTIAL_VERIFICATION_RESPONSE:
		return TypeCredentialVerificationResponse
	case C.CONTENT_CREDENTIAL_PRESENTATION_REQUEST:
		return TypeCredentialPresentationRequest
	case C.CONTENT_CREDENTIAL_PRESENTATION_RESPONSE:
		return TypeCredentialPresentationResponse
	default:
		return TypeUnknown
	}
}

// ID returns the unique id of the message
func (c *Content) ID() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_id(c.ptr)),
		20,
	)
}
