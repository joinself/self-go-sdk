package status

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"

type Error struct {
	status  uint32
	message string
}

func New(status uint32) *Error {

	return &Error{
		status: status,
		message: C.GoString(
			C.self_status_message(status),
		),
	}
}

func (e Error) Error() string {
	return e.message
}
