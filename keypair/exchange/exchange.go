package exchange

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/joinself/self-go-sdk-next/keypair"
)

type PublicKey struct {
	ptr *C.self_exchange_public_key
}

func newExchangePublicKey(ptr *C.self_exchange_public_key) *PublicKey {
	p := &PublicKey{
		ptr: ptr,
	}

	runtime.SetFinalizer(p, func(p *PublicKey) {
		C.self_exchange_public_key_destroy(
			p.ptr,
		)
	})

	return p
}

func exchangePublicKeyPtr(p *PublicKey) *C.self_exchange_public_key {
	return p.ptr
}

// FromAddress converts a hex address to a public key
func FromAddress(hex string) *PublicKey {
	var publicKey *C.self_exchange_public_key

	hexBuf := (*C.uint8_t)(C.CBytes([]byte(hex)))
	hexLen := C.size_t(len(hex))

	defer func() {
		C.free(unsafe.Pointer(hexBuf))
	}()

	status := C.self_exchange_public_key_decode(
		&publicKey,
		hexBuf,
		hexLen,
	)

	if status != 0 {
		return nil
	}

	return newExchangePublicKey(publicKey)
}

// Type returns the type of key
func (p *PublicKey) Type() keypair.KeyType {
	return keypair.KeyTypeExchange
}

// String returns the hex encoded address of a public key
func (p *PublicKey) String() string {
	encodedPtr := C.CBytes(make([]byte, 66))
	defer C.free(encodedPtr)

	status := C.self_exchange_public_key_encode(
		p.ptr,
		(*C.uint8_t)(encodedPtr),
		66,
	)

	if status > 0 {
		panic("invalid key conversion!")
	}

	encoded := C.GoBytes(
		encodedPtr,
		66,
	)

	return string(encoded)
}

// Bytes returns the raw bytes of the address
func (p *PublicKey) Bytes() []byte {
	bytesPtr := C.CBytes(make([]byte, 33))
	defer C.free(bytesPtr)

	status := C.self_exchange_public_key_as_bytes(
		p.ptr,
		(*C.uint8_t)(bytesPtr),
		33,
	)

	if status > 0 {
		panic("invalid key conversion!")
	}

	bytes := C.GoBytes(
		bytesPtr,
		33,
	)

	return bytes
}
