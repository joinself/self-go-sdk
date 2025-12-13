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
)

// NOTE this is here specifically to standardize the api surface for account
// as a result of moving KeyPackage to this package
type Welcome struct {
	ptr *C.self_crypto_welcome
}

func newCryptoWelcome(ptr *C.self_crypto_welcome, owned bool) *Welcome {
	e := &Welcome{
		ptr: ptr,
	}

	if owned {
		runtime.AddCleanup(e, func(e *Welcome) {
			C.self_crypto_welcome_destroy(
				e.ptr,
			)
		}, e)
	}

	return e
}

func cryptoWelcomePtr(w *Welcome) *C.self_crypto_welcome {
	return w.ptr
}
