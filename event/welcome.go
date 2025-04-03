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

//go:linkname newCryptoWelcome github.com/joinself/self-go-sdk/crypto.newWelcome
func newCryptoWelcome(ptr *C.self_welcome, owned bool) *crypto.Welcome

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

// Sequence returns the sequence of this event as determined by it's sender
func (w *Welcome) Sequence() uint64 {
	return uint64(C.self_welcome_sequence(
		w.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (w *Welcome) Timestamp() time.Time {
	return time.Unix(int64(C.self_welcome_timestamp(
		w.ptr,
	)), 0)
}

// Welcome returns the event's welcome
func (w *Welcome) Welcome() *crypto.Welcome {
	return newCryptoWelcome(w.ptr, false)
}
