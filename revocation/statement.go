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
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk/keypair/signing.newSigningPublicKey
func newSigningPublicKey(ptr *C.self_signing_public_key) *signing.PublicKey

// Statement a revocation statement
type Statement struct {
	ptr *C.self_revocation_statement
}

func newStatement(ptr *C.self_revocation_statement) *Statement {
	s := &Statement{
		ptr: ptr,
	}

	runtime.AddCleanup(s, func(ptr *C.self_revocation_statement) {
		C.self_revocation_statement_destroy(
			ptr,
		)
	}, s.ptr)

	return s
}

func statementPtr(s *Statement) *C.self_revocation_statement {
	return s.ptr
}

// StatementBuilder a builder for revocation statements
type StatementBuilder struct {
	ptr *C.self_revocation_statement_builder
}

func newStatementBuilder(ptr *C.self_revocation_statement_builder) *StatementBuilder {
	b := &StatementBuilder{
		ptr: ptr,
	}

	runtime.AddCleanup(b, func(ptr *C.self_revocation_statement_builder) {
		C.self_revocation_statement_builder_destroy(
			ptr,
		)
	}, b.ptr)

	return b
}

// NewStatement creates a new revocation statement builder
func NewStatement() *StatementBuilder {
	return newStatementBuilder(C.self_revocation_statement_builder_init())
}

// Issuer sets the issuer of the revocation statement
func (b *StatementBuilder) Issuer(issuer *signing.PublicKey) *StatementBuilder {
	C.self_revocation_statement_builder_issuer(
		b.ptr,
		signingPublicKeyPtr(issuer),
	)

	return b
}

// Sequence sets the sequence of the revocation statement
func (b *StatementBuilder) Sequence(sequence uint64) *StatementBuilder {
	C.self_revocation_statement_builder_sequence(
		b.ptr,
		C.uint64_t(sequence),
	)

	return b
}

// Timestamp sets the timestamp of the revocation statement
func (b *StatementBuilder) Timestamp(timestamp time.Time) *StatementBuilder {
	C.self_revocation_statement_builder_timestamp(
		b.ptr,
		C.int64_t(timestamp.Unix()),
	)

	return b
}

// RevokeBy marks a credential hash as revoked at a given time
func (b *StatementBuilder) RevokeBy(hash []byte, revokedAt time.Time) *StatementBuilder {
	hashBuf := C.CBytes(hash)

	C.self_revocation_statement_builder_revoke_by(
		b.ptr,
		(*C.uint8_t)(hashBuf),
		C.size_t(len(hash)),
		C.int64_t(revokedAt.Unix()),
	)

	C.free(hashBuf)

	return b
}

// SignWith specifies the key to sign the revocation statement with
func (b *StatementBuilder) SignWith(signer *signing.PublicKey, issuedAt time.Time) *StatementBuilder {
	C.self_revocation_statement_builder_sign_with(
		b.ptr,
		signingPublicKeyPtr(signer),
		C.int64_t(issuedAt.Unix()),
	)

	return b
}

// Finish finalizes the revocation statement and prepares it for issuing
func (b *StatementBuilder) Finish() (*Statement, error) {
	var statement *C.self_revocation_statement

	result := C.self_revocation_statement_builder_finish(
		b.ptr,
		&statement,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newStatement(statement), nil
}

// Issuer returns the issuer of the revocation statement
func (s *Statement) Issuer() *signing.PublicKey {
	return newSigningPublicKey(C.self_revocation_statement_issuer(
		s.ptr,
	))
}

// Sequence returns the sequence of the revocation statement
func (s *Statement) Sequence() uint64 {
	return uint64(C.self_revocation_statement_sequence(
		s.ptr,
	))
}

// Timestamp returns the timestamp of the revocation statement
func (s *Statement) Timestamp() time.Time {
	return time.Unix(int64(C.self_revocation_statement_timestamp(
		s.ptr,
	)), 0)
}

// MerkleRoot returns the merkle root of the revocation statement
func (s *Statement) MerkleRoot() []byte {
	buf := C.self_revocation_statement_merkle_root(
		s.ptr,
	)

	merkleRoot := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(buf)),
		C.int(C.self_bytes_buffer_len(buf)),
	)

	C.self_bytes_buffer_destroy(buf)

	return merkleRoot
}

// SignedBy returns true if the statement has been signed by the given signer
func (s *Statement) SignedBy(signer *signing.PublicKey) bool {
	return bool(C.self_revocation_statement_signed_by(
		s.ptr,
		signingPublicKeyPtr(signer),
	))
}

// Signers returns all signers of the revocation statement
func (s *Statement) Signers() []*Signer {
	collection := C.self_revocation_statement_signers(
		s.ptr,
	)

	signers := fromSignerCollection(collection)

	C.self_collection_revocation_statement_signer_destroy(
		collection,
	)

	return signers
}

// RevokedAtBy returns the revocation timestamp for a given hash.
// Returns false if the hash is not in the statement.
func (s *Statement) RevokedAtBy(hash []byte) (time.Time, bool) {
	var timestamp C.int64_t

	hashBuf := C.CBytes(hash)
	defer C.free(hashBuf)

	found := C.self_revocation_statement_revoked_at_by(
		s.ptr,
		(*C.uint8_t)(hashBuf),
		C.size_t(len(hash)),
		&timestamp,
	)

	if !bool(found) {
		return time.Time{}, false
	}

	return time.Unix(int64(timestamp), 0), true
}

// Revocations returns all revocations in the statement
func (s *Statement) Revocations() []*Revocation {
	collection := C.self_revocation_statement_revocations(
		s.ptr,
	)

	revocations := fromRevocationCollection(collection)

	C.self_collection_revocation_statement_revocation_destroy(
		collection,
	)

	return revocations
}

// Merge merges the signatures from another statement into this one
func (s *Statement) Merge(other *Statement) error {
	result := C.self_revocation_statement_merge(
		s.ptr,
		other.ptr,
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// Encode encodes the revocation statement as bytes
func (s *Statement) Encode() ([]byte, error) {
	var buf *C.self_bytes_buffer

	result := C.self_revocation_statement_encode(
		s.ptr,
		&buf,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	encoded := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(buf)),
		C.int(C.self_bytes_buffer_len(buf)),
	)

	C.self_bytes_buffer_destroy(buf)

	return encoded, nil
}

// DecodeStatement decodes a revocation statement from bytes
func DecodeStatement(encodedStatement []byte) (*Statement, error) {
	var statement *C.self_revocation_statement

	encodedBuf := C.CBytes(encodedStatement)
	defer C.free(encodedBuf)

	result := C.self_revocation_statement_decode(
		(*C.uint8_t)(encodedBuf),
		C.size_t(len(encodedStatement)),
		&statement,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newStatement(statement), nil
}

// Signer a signer of a revocation statement
type Signer struct {
	ptr *C.self_revocation_statement_signer
}

func newSigner(ptr *C.self_revocation_statement_signer) *Signer {
	s := &Signer{
		ptr: ptr,
	}

	runtime.AddCleanup(s, func(ptr *C.self_revocation_statement_signer) {
		C.self_revocation_statement_signer_destroy(
			ptr,
		)
	}, s.ptr)

	return s
}

// Address returns the signing key address of the signer
func (s *Signer) Address() *signing.PublicKey {
	return newSigningPublicKey(C.self_revocation_statement_signer_address(
		s.ptr,
	))
}

// Issued returns the timestamp that the signing key signed the statement
func (s *Signer) Issued() time.Time {
	return time.Unix(int64(C.self_revocation_statement_signer_issued(
		s.ptr,
	)), 0)
}

// Revocation a single revocation entry in a statement
type Revocation struct {
	ptr *C.self_revocation_statement_revocation
}

func newRevocation(ptr *C.self_revocation_statement_revocation) *Revocation {
	r := &Revocation{
		ptr: ptr,
	}

	runtime.AddCleanup(r, func(ptr *C.self_revocation_statement_revocation) {
		C.self_revocation_statement_revocation_destroy(
			ptr,
		)
	}, r.ptr)

	return r
}

// Hash returns the hash of the revoked item
func (r *Revocation) Hash() []byte {
	ptr := C.self_revocation_statement_revocation_hash(
		r.ptr,
	)

	return C.GoBytes(unsafe.Pointer(ptr), 32)
}

// Timestamp returns the timestamp the revocation occurred
func (r *Revocation) Timestamp() time.Time {
	return time.Unix(int64(C.self_revocation_statement_revocation_timestamp(
		r.ptr,
	)), 0)
}

func fromSignerCollection(collection *C.self_collection_revocation_statement_signer) []*Signer {
	length := int(C.self_collection_revocation_statement_signer_len(
		collection,
	))

	signers := make([]*Signer, length)

	for i := 0; i < length; i++ {
		ptr := C.self_collection_revocation_statement_signer_at(
			collection,
			C.size_t(i),
		)

		signers[i] = newSigner(ptr)
	}

	return signers
}

func fromRevocationCollection(collection *C.self_collection_revocation_statement_revocation) []*Revocation {
	length := int(C.self_collection_revocation_statement_revocation_len(
		collection,
	))

	revocations := make([]*Revocation, length)

	for i := 0; i < length; i++ {
		ptr := C.self_collection_revocation_statement_revocation_at(
			collection,
			C.size_t(i),
		)

		revocations[i] = newRevocation(ptr)
	}

	return revocations
}
