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

// NOTE this is here specifically to standardize the api surface for account
// as a result of moving KeyPackage to this package
type Welcome struct {
	ptr *C.self_welcome
}

func newWelcome(ptr *C.self_welcome) *Welcome {
	e := &Welcome{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *Welcome) {
		C.self_welcome_destroy(
			e.ptr,
		)
	})

	return e
}

func welcomePtr(w *Welcome) *C.self_welcome {
	return w.ptr
}

// ToAddress returns the address the event was addressed to
func (w *Welcome) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_welcome_to_address(
		w.ptr,
	))
}

// FromAddress returns the address the event was sent by
func (w *Welcome) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_welcome_from_address(
		w.ptr,
	))
}
