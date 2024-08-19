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

type Address = C.self_credential_address

func AddressAure(address *signing.PublicKey) *Address {
	a := C.self_credential_address_aure(
		(*C.self_signing_public_key)(address),
	)

	runtime.SetFinalizer(&a, func(a **C.self_credential_address) {
		C.self_credential_address_destroy(
			*a,
		)
	})

	return (*Address)(a)
}

func AddressAureWithKey(address, key *signing.PublicKey) *Address {
	a := C.self_credential_address_aure_with_key(
		(*C.self_signing_public_key)(address),
		(*C.self_signing_public_key)(key),
	)

	runtime.SetFinalizer(&a, func(a **C.self_credential_address) {
		C.self_credential_address_destroy(
			*a,
		)
	})

	return (*Address)(a)
}

func AddressKey(address *signing.PublicKey) *Address {
	a := C.self_credential_address_key(
		(*C.self_signing_public_key)(address),
	)

	runtime.SetFinalizer(&a, func(a **C.self_credential_address) {
		C.self_credential_address_destroy(
			*a,
		)
	})

	return (*Address)(a)
}

func (a *Address) String() string {
	encodedAddressBuffer := C.self_credential_address_encode(
		(*C.self_credential_address)(a),
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
