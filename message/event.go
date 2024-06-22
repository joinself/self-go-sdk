package message

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"time"

	"github.com/joinself/self-go-sdk/keypair/signing"
)

type Commit C.self_commit
type KeyPackage C.self_key_package
type Proposal C.self_proposal
type Welcome C.self_welcome

// ToAddress returns the address the event was addressed to
func (c *Commit) ToAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_commit_to_address(
		(*C.self_commit)(c),
	))
}

// FromAddress returns the address the event was sent by
func (c *Commit) FromAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_commit_from_address(
		(*C.self_commit)(c),
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Commit) Sequence() uint64 {
	return uint64(C.self_commit_sequence(
		(*C.self_commit)(c),
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Commit) Timestamp() time.Time {
	return time.Unix(int64(C.self_commit_timestamp(
		(*C.self_commit)(c),
	)), 0)
}

// ToAddress returns the address the event was addressed to
func (c *KeyPackage) ToAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_key_package_to_address(
		(*C.self_key_package)(c),
	))
}

// FromAddress returns the address the event was sent by
func (c *KeyPackage) FromAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_key_package_from_address(
		(*C.self_key_package)(c),
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *KeyPackage) Sequence() uint64 {
	return uint64(C.self_key_package_sequence(
		(*C.self_key_package)(c),
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *KeyPackage) Timestamp() time.Time {
	return time.Unix(int64(C.self_key_package_timestamp(
		(*C.self_key_package)(c),
	)), 0)
}

// ToAddress returns the address the event was addressed to
func (c *Proposal) ToAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_proposal_to_address(
		(*C.self_proposal)(c),
	))
}

// FromAddress returns the address the event was sent by
func (c *Proposal) FromAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_proposal_from_address(
		(*C.self_proposal)(c),
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Proposal) Sequence() uint64 {
	return uint64(C.self_proposal_sequence(
		(*C.self_proposal)(c),
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Proposal) Timestamp() time.Time {
	return time.Unix(int64(C.self_proposal_timestamp(
		(*C.self_proposal)(c),
	)), 0)
}

// ToAddress returns the address the event was addressed to
func (c *Welcome) ToAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_welcome_to_address(
		(*C.self_welcome)(c),
	))
}

// FromAddress returns the address the event was sent by
func (c *Welcome) FromAddress() *signing.PublicKey {
	return (*signing.PublicKey)(C.self_welcome_from_address(
		(*C.self_welcome)(c),
	))
}

// Sequence returns the sequence of this event as determined by it's sender
func (c *Welcome) Sequence() uint64 {
	return uint64(C.self_welcome_sequence(
		(*C.self_welcome)(c),
	))
}

// Timestamp returns the timestamp the event was sent at
func (c *Welcome) Timestamp() time.Time {
	return time.Unix(int64(C.self_welcome_timestamp(
		(*C.self_welcome)(c),
	)), 0)
}
