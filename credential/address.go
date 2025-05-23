package credential

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

	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk/keypair/signing.newSigningPublicKey
func newSigningPublicKey(ptr *C.self_signing_public_key) *signing.PublicKey

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

type Method int

const (
	MethodUnknown        = 1<<63 - 1
	MethodAure    Method = C.METHOD_AURE
	MethodKey     Method = C.METHOD_KEY
)

type Address struct {
	ptr *C.self_credential_address
}

func newAddress(ptr *C.self_credential_address) *Address {
	a := &Address{
		ptr: ptr,
	}

	runtime.SetFinalizer(a, func(a *Address) {
		C.self_credential_address_destroy(
			a.ptr,
		)
	})

	return a
}

func credentialAddressPtr(a *Address) *C.self_credential_address {
	return a.ptr
}

// DecodeAddress decodes a did address
func DecodeAddress(did string) (*Address, error) {
	didPtr := C.CString(did)

	var address *C.self_credential_address

	result := C.self_credential_address_decode(
		didPtr,
		&address,
	)

	C.free(unsafe.Pointer(didPtr))

	if result > 0 {
		return nil, status.New(result)
	}

	return newAddress(address), nil
}

// AddressAure creates a new aure method address
func AddressAure(address *signing.PublicKey) *Address {
	return newAddress(C.self_credential_address_aure(
		signingPublicKeyPtr(address),
	))
}

// AddressAureWithKey creates a new aure method address with a signing key
func AddressAureWithKey(address, key *signing.PublicKey) *Address {
	return newAddress(C.self_credential_address_aure_with_key(
		signingPublicKeyPtr(address),
		signingPublicKeyPtr(key),
	))
}

// AddressKey creates a new key method address
func AddressKey(address *signing.PublicKey) *Address {
	return newAddress(C.self_credential_address_key(
		signingPublicKeyPtr(address),
	))
}

// Method returns the method of the address
func (a *Address) Method() Method {
	switch C.self_credential_address_method(a.ptr) {
	case C.METHOD_AURE:
		return MethodAure
	case C.METHOD_KEY:
		return MethodKey
	default:
		return MethodUnknown
	}
}

// Address returs the address key of the address
func (a *Address) Address() *signing.PublicKey {
	return newSigningPublicKey(C.self_credential_address_address(a.ptr))
}

// SigningKey returns the signing key of the address
func (a *Address) SigningKey() *signing.PublicKey {
	sk := C.self_credential_address_signing_key(a.ptr)
	if sk == nil {
		return nil
	}

	return newSigningPublicKey(sk)
}

// String encodes the address to a string
func (a *Address) String() string {
	encodedAddressBuffer := C.self_credential_address_encode(
		a.ptr,
	)

	encodedAddress := C.GoString(
		C.self_string_buffer_ptr(encodedAddressBuffer),
	)

	C.self_string_buffer_destroy(
		encodedAddressBuffer,
	)

	return encodedAddress
}
