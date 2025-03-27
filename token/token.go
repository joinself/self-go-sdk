package token

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

type Token struct {
	ptr *C.self_token
}

func newToken(ptr *C.self_token) *Token {
	t := &Token{
		ptr: ptr,
	}

	runtime.SetFinalizer(t, func(t *Token) {
		C.self_token_destroy(
			t.ptr,
		)
	})

	return t
}

func tokenPtr(t *Token) *C.self_token {
	return t.ptr
}
