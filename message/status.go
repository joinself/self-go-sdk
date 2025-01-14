package message

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"

const (
	ResponseStatusUnknown       ResponseStatus = C.RESPONSE_STATUS_UNKNOWN
	ResponseStatusOk            ResponseStatus = C.RESPONSE_STATUS_OK
	ResponseStatusAccepted      ResponseStatus = C.RESPONSE_STATUS_ACCEPTED
	ResponseStatusCreated       ResponseStatus = C.RESPONSE_STATUS_CREATED
	ResponseStatusBadRequest    ResponseStatus = C.RESPONSE_STATUS_BAD_REQUEST
	ResponseStatusUnauthorized  ResponseStatus = C.RESPONSE_STATUS_UNAUTHORIZED
	ResponseStatusForbidden     ResponseStatus = C.RESPONSE_STATUS_FORBIDDEN
	ResponseStatusNotFound      ResponseStatus = C.RESPONSE_STATUS_NOT_FOUND
	ResponseStatusNotAcceptable ResponseStatus = C.RESPONSE_STATUS_NOT_ACCEPTABLE
	ResponseStatusConflict      ResponseStatus = C.RESPONSE_STATUS_CONFLICT
)

type ResponseStatus int

func (s ResponseStatus) String() string {
	switch s {
	case ResponseStatusOk:
		return "Ok"
	case ResponseStatusAccepted:
		return "Accepted"
	case ResponseStatusCreated:
		return "Created"
	case ResponseStatusBadRequest:
		return "Bad Request"
	case ResponseStatusUnauthorized:
		return "Unauthorized"
	case ResponseStatusForbidden:
		return "Forbidden"
	case ResponseStatusNotFound:
		return "Not Found"
	case ResponseStatusNotAcceptable:
		return "Not Acceptable"
	case ResponseStatusConflict:
		return "Conflict"
	default:
		return "Unknown"
	}
}
