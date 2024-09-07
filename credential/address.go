package credential

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/joinself/self-go-sdk-next/keypair/signing"
)

//go:linkname signingPublicKeyPtr signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

type Address struct {
	ptr *C.self_credential_address
}

func newAddress(ptr *C.self_credential_address) *Address {
	a := &Address{
		ptr: ptr,
	}

	runtime.SetFinalizer(&a, func(a *Address) {
		C.self_credential_address_destroy(
			a.ptr,
		)
	})

	return a
}

func AddressAure(address *signing.PublicKey) *Address {
	return newAddress(C.self_credential_address_aure(
		signingPublicKeyPtr(address),
	))
}

func AddressAureWithKey(address, key *signing.PublicKey) *Address {
	return newAddress(C.self_credential_address_aure_with_key(
		signingPublicKeyPtr(address),
		signingPublicKeyPtr(key),
	))
}

func AddressKey(address *signing.PublicKey) *Address {
	return newAddress(C.self_credential_address_key(
		signingPublicKeyPtr(address),
	))
}

func (a *Address) String() string {
	encodedAddressBuffer := C.self_credential_address_encode(
		a.ptr,
	)

	encodedAddress := C.GoBytes(
		unsafe.Pointer(C.self_encoded_buffer_buf(encodedAddressBuffer)),
		C.int(C.self_encoded_buffer_len(encodedAddressBuffer)),
	)

	C.self_encoded_buffer_destroy(
		encodedAddressBuffer,
	)

	return string(encodedAddress)
}
