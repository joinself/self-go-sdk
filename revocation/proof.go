package revocation

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
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair/signing"
)

// Proof a proof that a credential was revoked
type Proof struct {
	ptr *C.self_revocation_proof
}

func newProof(ptr *C.self_revocation_proof) *Proof {
	p := &Proof{
		ptr: ptr,
	}

	runtime.AddCleanup(p, func(ptr *C.self_revocation_proof) {
		C.self_revocation_proof_destroy(
			ptr,
		)
	}, p.ptr)

	return p
}

// Issuer returns the issuer of the revocation proof
func (p *Proof) Issuer() *signing.PublicKey {
	return newSigningPublicKey(C.self_revocation_proof_issuer(
		p.ptr,
	))
}

// Sequence returns the sequence of the statement the proof was created from
func (p *Proof) Sequence() uint64 {
	return uint64(C.self_revocation_proof_sequence(
		p.ptr,
	))
}

// Timestamp returns the timestamp of the statement the proof was created from
func (p *Proof) Timestamp() time.Time {
	return time.Unix(int64(C.self_revocation_proof_timestamp(
		p.ptr,
	)), 0)
}

// Revoked returns the timestamp for when the revocation took effect
func (p *Proof) Revoked() time.Time {
	return time.Unix(int64(C.self_revocation_proof_revoked(
		p.ptr,
	)), 0)
}

// RevocationHash returns the hash of the entity that was revoked
func (p *Proof) RevocationHash() []byte {
	ptr := C.self_revocation_proof_revocation_hash(
		p.ptr,
	)

	return C.GoBytes(unsafe.Pointer(ptr), 32)
}

// Signers returns the signers of the statement the proof was generated from
func (p *Proof) Signers() []*Signer {
	collection := C.self_revocation_proof_signers(
		p.ptr,
	)

	signers := fromSignerCollection(collection)

	C.self_collection_revocation_statement_signer_destroy(
		collection,
	)

	return signers
}
