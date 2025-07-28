package event

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
	"time"

	"github.com/joinself/self-go-sdk/crypto"
	"github.com/joinself/self-go-sdk/keypair/signing"
)

//go:linkname newCryptoKeyPackage github.com/joinself/self-go-sdk/crypto.newCryptoKeyPackage
func newCryptoKeyPackage(ptr *C.self_crypto_key_package, owned bool) *crypto.KeyPackage

type KeyPackage struct {
	ptr *C.self_key_package
}

func newKeyPackage(ptr *C.self_key_package) *KeyPackage {
	e := &KeyPackage{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *KeyPackage) {
		C.self_key_package_destroy(
			e.ptr,
		)
	})

	return e
}

func keyPackagePtr(k *KeyPackage) *C.self_key_package {
	return k.ptr
}

// ToAddress returns the address the event was addressed to
func (c *KeyPackage) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_key_package_to_address(
		c.ptr,
	))
}

// FromAddress returns the address the event was sent by
func (c *KeyPackage) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_key_package_from_address(
		c.ptr,
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *KeyPackage) Sequence() uint64 {
	return uint64(C.self_key_package_sequence(
		c.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *KeyPackage) Timestamp() time.Time {
	return time.Unix(int64(C.self_key_package_timestamp(
		c.ptr,
	)), 0)
}

// KeyPackage returns the events key package
func (c *KeyPackage) KeyPackage() *crypto.KeyPackage {
	return newCryptoKeyPackage(C.self_key_package_crypto_key_package(c.ptr), false)
}
