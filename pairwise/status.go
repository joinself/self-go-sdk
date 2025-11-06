package pairwise

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration
#cgo linux LDFLAGS: -lself_sdk -Wl,--allow-multiple-definition
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"

type Status int

const (
	StatusPending     Status = C.CONNECTION_STATUS_PENDING
	StatusNegotiating Status = C.CONNECTION_STATUS_NEGOTIATING
	StatusEstablished Status = C.CONNECTION_STATUS_ESTABLISHED
)
