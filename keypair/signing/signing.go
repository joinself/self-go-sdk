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

type PublicKey struct {
	ptr *C.self_signing_public_key
}

func newSigningPublicKey(ptr *C.self_signing_public_key) *PublicKey {
	p := &PublicKey{
		ptr: ptr,
	}

	runtime.SetFinalizer(p, func(p *PublicKey) {
		C.self_signing_public_key_destroy(
			p.ptr,
		)
	})

	return p
}

func signingPublicKeyPtr(p *PublicKey) *C.self_signing_public_key {
	return p.ptr
}

// FromAddress converts a hex address to a public key
func FromAddress(hex string) *PublicKey {
	var ptr *C.self_signing_public_key

	hexBuf := (*C.uint8_t)(C.CBytes([]byte(hex)))
	hexLen := C.size_t(len(hex))

	defer func() {
		C.free(unsafe.Pointer(hexBuf))
	}()

	status := C.self_signing_public_key_decode(
		&ptr,
		hexBuf,
		hexLen,
	)

	if status != 0 {
		return nil
	}

	return newSigningPublicKey(ptr)
}

// FromBytes constructs a public key from bytes
func FromBytes(data []byte) *PublicKey {
	var ptr *C.self_signing_public_key

	dataBuf := (*C.uint8_t)(C.CBytes(data))
	dataLen := C.size_t(len(data))

	defer func() {
		C.free(unsafe.Pointer(dataBuf))
	}()

	status := C.self_signing_public_key_from_bytes(
		&ptr,
		dataBuf,
		dataLen,
	)

	if status != 0 {
		return nil
	}

	return newSigningPublicKey(ptr)
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

	status := C.self_signing_public_key_as_bytes(
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

// Matches compares two public keys
func (p *PublicKey) Matches(target *PublicKey) bool {
	return bool(C.self_signing_public_key_matches(
		p.ptr,
		target.ptr,
	))
}

/*
func toPublicKeyCollection(keys []*PublicKey) *C.self_collection_signing_public_key {
	collection := C.self_collection_signing_public_key_init()

	for i := 0; i < len(keys); i++ {
		C.self_collection_signing_public_key_append(
			collection,
			keys[i].ptr,
		)
	}

	return collection
}
*/

func fromSigningPublicKeyCollection(collection *C.self_collection_signing_public_key) []*PublicKey {
	collectionLen := int(C.self_collection_signing_public_key_len(
		collection,
	))

	keys := make([]*PublicKey, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_signing_public_key_at(
			collection,
			C.size_t(i),
		)

		keys[i] = newSigningPublicKey(ptr)
	}

	return keys
}
