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

	"github.com/joinself/self-go-sdk/keypair/signing"
)

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
func (c *Welcome) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_welcome_to_address(
		c.ptr,
	))
}

// FromAddress returns the address the event was sent by
func (c *Welcome) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_welcome_from_address(
		c.ptr,
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Welcome) Sequence() uint64 {
	return uint64(C.self_welcome_sequence(
		c.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Welcome) Timestamp() time.Time {
	return time.Unix(int64(C.self_welcome_timestamp(
		c.ptr,
	)), 0)
}
