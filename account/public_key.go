package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

type PublicKey struct {
	public *C.self_signing_public_key
}

func (p *PublicKey) String() string {
	encoded := make([]byte, 66)

	result := C.self_signing_public_key_encode(
		p.public,
		(*C.uchar)(unsafe.Pointer(&encoded[0])),
		C.ulong(len(encoded)),
	)

	if result > 0 {
		panic("invalid key conversion!")
	}

	return *(*string)(unsafe.Pointer(&encoded))
}
