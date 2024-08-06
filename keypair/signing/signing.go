package signing

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair"
)

type PublicKey C.self_signing_public_key

// FromAddress converts a hex address to a public key
func FromAddress(hex string) *PublicKey {
	var public *C.self_signing_public_key
	publicPtr := &public

	hexBuf := (*C.uint8_t)(C.CBytes([]byte(hex)))
	hexLen := C.size_t(len(hex))

	defer func() {
		C.free(unsafe.Pointer(hexBuf))
	}()

	status := C.self_signing_public_key_decode(
		publicPtr,
		hexBuf,
		hexLen,
	)

	if status != 0 {
		return nil
	}

	runtime.SetFinalizer(publicPtr, func(public **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*public,
		)
	})

	return (*PublicKey)(*publicPtr)
}

// Type returns the type of key
func (p *PublicKey) Type() keypair.KeyType {
	return keypair.KeyTypeSigning
}

// String returns the hex encoded address of a public key
func (p *PublicKey) String() string {
	encoded := make([]byte, 66)

	status := C.self_signing_public_key_encode(
		(*C.self_signing_public_key)(p),
		(*C.uint8_t)(unsafe.Pointer(&encoded[0])),
		C.ulong(len(encoded)),
	)

	if status > 0 {
		panic("invalid key conversion!")
	}

	return *(*string)(unsafe.Pointer(&encoded))
}

type PublicKeyCollection C.self_collection_signing_public_key

func NewPublicKeyCollection() *PublicKeyCollection {
	collection := (*PublicKeyCollection)(C.self_collection_signing_public_key_init())

	runtime.SetFinalizer(collection, func(collection *PublicKeyCollection) {
		C.self_collection_signing_public_key_destroy(
			(*C.self_collection_signing_public_key)(collection),
		)
	})

	return collection
}

func (c *PublicKeyCollection) Length() int {
	return int(C.self_collection_signing_public_key_len(
		(*C.self_collection_signing_public_key)(c),
	))
}

func (c *PublicKeyCollection) Get(index int) *PublicKey {
	publicKey := (*PublicKey)(C.self_collection_signing_public_key_at(
		(*C.self_collection_signing_public_key)(c),
		C.ulong(index),
	))

	runtime.SetFinalizer(publicKey, func(publicKey *PublicKey) {
		C.self_signing_public_key_destroy(
			(*C.self_signing_public_key)(publicKey),
		)
	})

	return publicKey
}
