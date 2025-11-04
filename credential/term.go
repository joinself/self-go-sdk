package credential

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

type Term struct {
	ptr *C.self_credential_term
}

func newTerm(ptr *C.self_credential_term) *Term {
	t := &Term{
		ptr: ptr,
	}

	runtime.SetFinalizer(t, func(t *Term) {
		C.self_credential_term_destroy(
			t.ptr,
		)
	})

	return t
}

func termPtr(t *Term) *C.self_credential_term {
	return t.ptr
}
