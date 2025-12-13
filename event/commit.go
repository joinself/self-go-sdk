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

type Commit struct {
	ptr *C.self_commit
}

func newCommit(ptr *C.self_commit) *Commit {
	e := &Commit{
		ptr: ptr,
	}

	runtime.AddCleanup(e, func(e *Commit) {
		C.self_commit_destroy(
			e.ptr,
		)
	}, e)

	return e
}

// ToAddress returns the address the event was addressed to
func (c *Commit) ToAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_commit_to_address(
		c.ptr,
	))
}

// FromAddress returns the address the event was sent by
func (c *Commit) FromAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_commit_from_address(
		c.ptr,
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Commit) Sequence() uint64 {
	return uint64(C.self_commit_sequence(
		c.ptr,
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Commit) Timestamp() time.Time {
	return time.Unix(int64(C.self_commit_timestamp(
		c.ptr,
	)), 0)
}
