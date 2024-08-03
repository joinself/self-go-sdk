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
	RESPONSE_STATUS_ACCEPTED           ResponseStatus = C.RESPONSE_STATUS_ACCEPTED
	ResponseStatusCreated              ResponseStatus = C.RESPONSE_STATUS_CREATED
	ResponseStatusBadRequest           ResponseStatus = C.RESPONSE_STATUS_BAD_REQUEST
	ResponseStatusUnauthorized         ResponseStatus = C.RESPONSE_STATUS_UNAUTHORIZED
	ResponseStatusForbidden            ResponseStatus = C.RESPONSE_STATUS_FORBIDDEN
	ResponseStatusNotFound             ResponseStatus = C.RESPONSE_STATUS_NOT_FOUND
	ResponseStatusNotAcceptable        ResponseStatus = C.RESPONSE_STATUS_NOT_ACCEPTABLE
	ResponseStatusConflict             ResponseStatus = C.RESPONSE_STATUS_CONFLICT
)
