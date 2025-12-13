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

type Proposal struct {
	ptr *C.self_proposal
}

func newProposal(ptr *C.self_proposal) *Proposal {
	e := &Proposal{
		ptr: ptr,
	}

	runtime.AddCleanup(e, func(ptr *C.self_proposal) {
		C.self_proposal_destroy(
			ptr,
		)
	}, e.ptr)

	return e
}

// ToAddress returns the address the event was addressed to
func (c *Proposal) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(
		C.self_proposal_to_address(c.ptr),
	)
}

// FromAddress returns the address the event was sent by
func (c *Proposal) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_proposal_from_address(
		c.ptr,
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Proposal) Sequence() uint64 {
	return uint64(C.self_proposal_sequence(
		c.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Proposal) Timestamp() time.Time {
	return time.Unix(int64(C.self_proposal_timestamp(
		c.ptr,
	)), 0)
}
