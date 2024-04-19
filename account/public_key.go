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

// PublicKey a signing public key
type PublicKey struct {
	public *C.self_signing_public_key
}

// Address converts a hex address to a public key
func Address(hex string) *PublicKey {
	var public *C.self_signing_public_key

	hexBuf := (*C.uint8_t)(C.CBytes([]byte(hex)))
	hexLen := C.size_t(len(hex))

	defer func() {
		C.free(unsafe.Pointer(hexBuf))
	}()

	status := C.self_signing_public_key_decode(
		&public,
		hexBuf,
		hexLen,
	)

	if status != 0 {
		return nil
	}

	return &PublicKey{
		public: public,
	}
}

// String returns the hex encoded address of a public key
func (p *PublicKey) String() string {
	encoded := make([]byte, 66)

	status := C.self_signing_public_key_encode(
		p.public,
		(*C.uchar)(unsafe.Pointer(&encoded[0])),
		C.ulong(len(encoded)),
	)

	if status > 0 {
		panic("invalid key conversion!")
	}

	return *(*string)(unsafe.Pointer(&encoded))
}
