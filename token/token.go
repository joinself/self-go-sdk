package token

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

func encodeToken(t *Token) ([]byte, error) {
	var buf *C.self_bytes_buffer

	result := C.self_token_encode(
		t.ptr,
		&buf,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	encodedToken := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(buf)),
		C.int(C.self_bytes_buffer_len(buf)),
	)

	C.self_bytes_buffer_destroy(buf)

	return encodedToken, nil
}
