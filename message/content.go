package message

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"

type Content C.self_message_content
type Type int
type ResponseStatus int

const (
	TypeUnknown                        Type           = 1<<63 - 1
	TypeCustom                         Type           = C.CONTENT_CUSTOM
	TypeChat                           Type           = C.CONTENT_CHAT
	TypeReceipt                        Type           = C.CONTENT_RECEIPT
	TypeDiscoveryRequest               Type           = C.CONTENT_DISCOVERY_REQUEST
	TypeDiscoveryResponse              Type           = C.CONTENT_DISCOVERY_RESPONSE
	TypeCredentialVerificationRequest  Type           = C.CONTENT_CREDENTIAL_VERIFICATION_REQUEST
	TypeCredentialVerificationResponse Type           = C.CONTENT_CREDENTIAL_VERIFICATION_RESPONSE
	TypeCredentialPresentationRequest  Type           = C.CONTENT_CREDENTIAL_PRESENTATION_REQUEST
	TypeCredentialPresentationResponse Type           = C.CONTENT_CREDENTIAL_PRESENTATION_RESPONSE
	ResponseStatusUnknown              ResponseStatus = C.RESPONSE_STATUS_UNKNOWN
	ResponseStatusOk                   ResponseStatus = C.RESPONSE_STATUS_OK
	ResponseStatusAccepted             ResponseStatus = C.RESPONSE_STATUS_ACCEPTED
	ResponseStatusCreated              ResponseStatus = C.RESPONSE_STATUS_CREATED
	ResponseStatusBadRequest           ResponseStatus = C.RESPONSE_STATUS_BAD_REQUEST
	ResponseStatusUnauthorized         ResponseStatus = C.RESPONSE_STATUS_UNAUTHORIZED
	ResponseStatusForbidden            ResponseStatus = C.RESPONSE_STATUS_FORBIDDEN
	ResponseStatusNotFound             ResponseStatus = C.RESPONSE_STATUS_NOT_FOUND
	ResponseStatusNotAcceptable        ResponseStatus = C.RESPONSE_STATUS_NOT_ACCEPTABLE
	ResponseStatusConflict             ResponseStatus = C.RESPONSE_STATUS_CONFLICT
)

func ContentType(message *Message) Type {
	content := C.self_message_message_content((*C.self_message)(message))

	switch C.self_message_content_type_of(content) {
	case C.CONTENT_CUSTOM:
		return TypeCustom
	case C.CONTENT_CHAT:
		return TypeChat
	case C.CONTENT_RECEIPT:
		return TypeReceipt
	case C.CONTENT_DISCOVERY_REQUEST:
		return TypeDiscoveryRequest
	case C.CONTENT_DISCOVERY_RESPONSE:
		return TypeDiscoveryResponse
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
