package exchange

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair"
)

type PublicKey C.self_exchange_public_key

// FromAddress converts a hex address to a public key
func FromAddress(hex string) *PublicKey {
	var public *C.self_exchange_public_key

	hexBuf := (*C.uint8_t)(C.CBytes([]byte(hex)))
	hexLen := C.size_t(len(hex))

	defer func() {
		C.free(unsafe.Pointer(hexBuf))
	}()

	status := C.self_exchange_public_key_decode(
		&public,
		hexBuf,
		hexLen,
	)

	if status != 0 {
		return nil
	}

	return (*PublicKey)(public)
}

// Type returns the type of key
func (p *PublicKey) Type() keypair.KeyType {
	return keypair.KeyTypeExchange
}

// String returns the hex encoded address of a public key
func (p *PublicKey) String() string {
	encoded := make([]byte, 66)

	status := C.self_exchange_public_key_encode(
		(*C.self_exchange_public_key)(p),
		(*C.uint8_t)(unsafe.Pointer(&encoded[0])),
		C.ulong(len(encoded)),
	)

	if status > 0 {
		panic("invalid key conversion!")
	}

	return *(*string)(unsafe.Pointer(&encoded))
}
