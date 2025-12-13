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
)

type Reference struct {
	ptr *C.self_reference
}

func newReference(ptr *C.self_reference) *Reference {
	e := &Reference{
		ptr: ptr,
	}

	runtime.AddCleanup(e, func(e *Reference) {
		C.self_reference_destroy(
			e.ptr,
		)
	}, e)

	return e
}

// ID returns the id of the messages content
func (r *Reference) ID() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_reference_id(r.ptr)),
		20,
	)
}
