package crypto

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

	"github.com/joinself/self-go-sdk/keypair/signing"
)

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk/keypair/signing.newSigningPublicKey
func newSigningPublicKey(ptr *C.self_signing_public_key) *signing.PublicKey

// NOTE this serves to provide a means of breaking an import cycle between message and event packages
// the underlying C type is the same as the `event.KeyPackage` type
type KeyPackage struct {
	ptr *C.self_crypto_key_package
}

func newCryptoKeyPackage(ptr *C.self_crypto_key_package, owned bool) *KeyPackage {
	e := &KeyPackage{
		ptr: ptr,
	}

	if owned {
		runtime.AddCleanup(e, func(ptr *C.self_crypto_key_package) {
			C.self_crypto_key_package_destroy(
				ptr,
			)
		}, e.ptr)
	}

	return e
}

func cryptoKeyPackagePtr(k *KeyPackage) *C.self_crypto_key_package {
	return k.ptr
}

// FromAddress returns the address the event was sent by
func (c *KeyPackage) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_crypto_key_package_from_address(
		c.ptr,
	))
}

func toCryptoKeyPackageCollection(packages []*KeyPackage) *C.self_collection_crypto_key_package {
	collection := C.self_collection_crypto_key_package_init()

	for i := 0; i < len(packages); i++ {
		C.self_collection_crypto_key_package_append(
			collection,
			packages[i].ptr,
		)
	}

	return collection
}
