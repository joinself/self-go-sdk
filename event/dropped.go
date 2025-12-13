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

	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/status"
)

type Dropped struct {
	ptr *C.self_dropped_event
}

func newDropped(ptr *C.self_dropped_event) *Dropped {
	e := &Dropped{
		ptr: ptr,
	}

	runtime.AddCleanup(e, func(e *Dropped) {
		C.self_dropped_event_destroy(
			e.ptr,
		)
	}, e)

	return e
}

// ToAddress returns the address the event was addressed to
func (c *Dropped) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_dropped_event_to_address(
		c.ptr,
	))
}

// FromAddress returns the address the event was sent by
func (c *Dropped) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_dropped_event_from_address(
		c.ptr,
	))
}

// FromSequence returns the first sequence that was dropped
func (c *Dropped) FromSequence() uint64 {
	return uint64(C.self_dropped_event_from_sequence(
		c.ptr,
	))
}

// ToSequence returns the last sequence that was dropped
func (c *Dropped) ToSequence() uint64 {
	return uint64(C.self_dropped_event_to_sequence(
		c.ptr,
	))
}

// Reason returns the reason the events were dropped
func (c *Dropped) Reason() error {
	return status.New(C.self_dropped_event_reason(
		c.ptr,
	))
}
