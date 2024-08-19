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

	"github.com/joinself/self-go-sdk-next/keypair"
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

// FromBytes constructs a public key from bytes
func FromBytes(data []byte) *PublicKey {
	var public *C.self_signing_public_key
	publicPtr := &public

	dataBuf := (*C.uint8_t)(C.CBytes(data))
	dataLen := C.size_t(len(data))

	defer func() {
		C.free(unsafe.Pointer(dataBuf))
	}()

	status := C.self_signing_public_key_from_bytes(
		publicPtr,
		dataBuf,
		dataLen,
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
	encodedPtr := C.CBytes(make([]byte, 66))
	defer C.free(encodedPtr)

	status := C.self_signing_public_key_encode(
		(*C.self_signing_public_key)(p),
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

	status := C.self_signing_public_key_as_bytes(
		(*C.self_signing_public_key)(p),
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

type PublicKeyCollection C.self_collection_signing_public_key

func NewPublicKeyCollection() *PublicKeyCollection {
	collection := C.self_collection_signing_public_key_init()

	runtime.SetFinalizer(collection, func(collection *C.self_collection_signing_public_key) {
		C.self_collection_signing_public_key_destroy(
			collection,
		)
	})

	return (*PublicKeyCollection)(collection)
}

func (c *PublicKeyCollection) Length() int {
	return int(C.self_collection_signing_public_key_len(
		(*C.self_collection_signing_public_key)(c),
	))
}

func (c *PublicKeyCollection) Get(index int) *PublicKey {
	publicKey := C.self_collection_signing_public_key_at(
		(*C.self_collection_signing_public_key)(c),
		C.ulong(index),
	)

	runtime.SetFinalizer(publicKey, func(publicKey *C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			publicKey,
		)
	})

	return (*PublicKey)(publicKey)
}
