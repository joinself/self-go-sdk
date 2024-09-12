package identity

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"time"

	"github.com/joinself/self-go-sdk-next/keypair"
	"github.com/joinself/self-go-sdk-next/keypair/exchange"
	"github.com/joinself/self-go-sdk-next/keypair/signing"
)

// TODO define directly from C types...
type Role uint64

const (
	RoleAssertion      Role = C.KEY_ROLE_ASSERTION
	RoleAuthentication Role = C.KEY_ROLE_AUTHENTICATION
	RoleVerification   Role = C.KEY_ROLE_VERIFICATION
	RoleInvocation     Role = C.KEY_ROLE_INVOCATION
	RoleDelegation     Role = C.KEY_ROLE_DELEGATION
	RoleMessaging      Role = C.KEY_ROLE_MESSAGING
)

type Operation struct {
	ptr *C.self_identity_operation
}

func newOperation(ptr *C.self_identity_operation) *Operation {
	o := &Operation{
		ptr: ptr,
	}

	runtime.SetFinalizer(o, func(o *Operation) {
		C.self_identity_operation_destroy(
			o.ptr,
		)
	})

	return o
}

func operationPtr(o *Operation) *C.self_identity_operation {
	return o.ptr
}

type OperationBuilder struct {
	ptr *C.self_identity_operation_builder
}

func newOperationBuilder(ptr *C.self_identity_operation_builder) *OperationBuilder {
	b := &OperationBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(b, func(b *OperationBuilder) {
		C.self_identity_operation_builder_destroy(
			b.ptr,
		)
	})

	return b
}

// NewOperation creates a new operation
func NewOperation() *OperationBuilder {
	return newOperationBuilder(C.self_identity_operation_builder_init())
}

// Identifier sets the identifier of the document to target
func (b *OperationBuilder) Identifier(address *signing.PublicKey) *OperationBuilder {
	C.self_identity_operation_builder_id(
		b.ptr,
		signingPublicKeyPtr(address),
	)

	return b
}

// Identifier sets the identifier of the document to target
func (b *OperationBuilder) Sequence(sequence uint32) *OperationBuilder {
	C.self_identity_operation_builder_sequence(
		b.ptr,
		C.uint32_t(sequence),
	)

	return b
}

// Timestamp sets the timestamp of the operation
func (b *OperationBuilder) Timestamp(timestamp time.Time) *OperationBuilder {
	C.self_identity_operation_builder_timestamp(
		b.ptr,
		C.int64_t(timestamp.Unix()),
	)

	return b
}

// Previous sets the hash of the previous operation
func (b *OperationBuilder) Previous(previousHash []byte) *OperationBuilder {
	previousBuf := C.CBytes(previousHash)
	previousLen := len(previousHash)

	C.self_identity_operation_builder_previous(
		b.ptr,
		(*C.uint8_t)(previousBuf),
		C.size_t(previousLen),
	)

	C.free(previousBuf)

	return b
}

// GrantEmbedded grants an embedded key with a given set of roles
func (b *OperationBuilder) GrantEmbedded(key keypair.PublicKey, roles Role) *OperationBuilder {
	switch pk := key.(type) {
	case *signing.PublicKey:
		C.self_identity_operation_builder_signing_grant_embedded(
			b.ptr,
			signingPublicKeyPtr(pk),
			C.uint64_t(roles),
		)
	case *exchange.PublicKey:
		C.self_identity_operation_builder_exchange_grant_embedded(
			b.ptr,
			exchangePublicKeyPtr(pk),
			C.uint64_t(roles),
		)
	}

	return b
}

// GrantReferenced grants roles to a key controlled by another identity
func (b *OperationBuilder) GrantReferenced(method uint16, controller *signing.PublicKey, key *signing.PublicKey, roles Role) *OperationBuilder {
	C.self_identity_operation_builder_signing_grant_referenced(
		b.ptr,
		C.uint16_t(method),
		signingPublicKeyPtr(controller),
		signingPublicKeyPtr(key),
		C.uint64_t(roles),
	)

	return b
}

// Modify modifies the roles of an existing key
func (b *OperationBuilder) Modify(key keypair.PublicKey, roles Role) *OperationBuilder {
	switch pk := key.(type) {
	case *signing.PublicKey:
		C.self_identity_operation_builder_signing_modify(
			b.ptr,
			signingPublicKeyPtr(pk),
			C.uint64_t(roles),
		)
	case *exchange.PublicKey:
		C.self_identity_operation_builder_exchange_modify(
			b.ptr,
			exchangePublicKeyPtr(pk),
			C.uint64_t(roles),
		)
	}

	return b
}

// Revoke revokes a key from a given point in time
func (b *OperationBuilder) Revoke(key keypair.PublicKey, effectiveFrom time.Time) *OperationBuilder {
	switch pk := key.(type) {
	case *signing.PublicKey:
		C.self_identity_operation_builder_signing_revoke(
			b.ptr,
			signingPublicKeyPtr(pk),
			C.int64_t(effectiveFrom.Unix()),
		)
	case *exchange.PublicKey:
		C.self_identity_operation_builder_exchange_revoke(
			b.ptr,
			exchangePublicKeyPtr(pk),
			C.int64_t(effectiveFrom.Unix()),
		)
	}

	return b
}

// Recover recovers an identity, revoking all existing keys
func (b *OperationBuilder) Recover(effectiveFrom time.Time) *OperationBuilder {
	C.self_identity_operation_builder_recover(
		b.ptr,
		C.int64_t(effectiveFrom.Unix()),
	)

	return b
}

// Deactivate permanently deactivates the identity
func (b *OperationBuilder) Deactivate(effectiveFrom time.Time) *OperationBuilder {
	C.self_identity_operation_builder_deactivate(
		b.ptr,
		C.int64_t(effectiveFrom.Unix()),
	)

	return b
}

// SignWith specifies which key to sign the operation with
func (b *OperationBuilder) SignWith(signer *signing.PublicKey) *OperationBuilder {
	C.self_identity_operation_builder_sign_with(
		b.ptr,
		signingPublicKeyPtr(signer),
	)

	return b
}

// Finish finalizes the operation and prepares it for execution
func (b *OperationBuilder) Finish() *Operation {
	var ptr *C.self_identity_operation

	C.self_identity_operation_builder_finish(
		b.ptr,
		&ptr,
	)

	return newOperation(ptr)
}
