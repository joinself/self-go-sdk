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
	RoleAssertion      Role = C.ASSERTION
	RoleAuthentication Role = C.AUTHENTICATION
	RoleVerification   Role = C.VERIFICATION
	RoleInvocation     Role = C.INVOCATION
	RoleDelegation     Role = C.DELEGATION
	RoleMessaging      Role = C.MESSAGING
)

type Operation C.self_identity_operation
type OperationBuilder C.self_identity_operation_builder

// NewOperation creates a new operation
func NewOperation() *OperationBuilder {
	builder := C.self_identity_operation_builder_init()

	runtime.SetFinalizer(&builder, func(builder **C.self_identity_operation_builder) {
		C.self_identity_operation_builder_destroy(
			*builder,
		)
	})

	return (*OperationBuilder)(builder)
}

// Identifier sets the identifier of the document to target
func (b *OperationBuilder) Identifier(address *signing.PublicKey) *OperationBuilder {
	C.self_identity_operation_builder_id(
		(*C.self_identity_operation_builder)(b),
		(*C.self_signing_public_key)(address),
	)

	return b
}

// Identifier sets the identifier of the document to target
func (b *OperationBuilder) Sequence(sequence uint32) *OperationBuilder {
	C.self_identity_operation_builder_sequence(
		(*C.self_identity_operation_builder)(b),
		C.uint(sequence),
	)

	return b
}

// Timestamp sets the timestamp of the operation
func (b *OperationBuilder) Timestamp(timestamp time.Time) *OperationBuilder {
	C.self_identity_operation_builder_timestamp(
		(*C.self_identity_operation_builder)(b),
		C.long(timestamp.Unix()),
	)

	return b
}

// Previous sets the hash of the previous operation
func (b *OperationBuilder) Previous(previousHash []byte) *OperationBuilder {
	previousBuf := C.CBytes(previousHash)
	previousLen := len(previousHash)

	C.self_identity_operation_builder_previous(
		(*C.self_identity_operation_builder)(b),
		(*C.uint8_t)(previousBuf),
		C.ulong(previousLen),
	)

	C.free(previousBuf)

	return b
}

// GrantEmbedded grants an embedded key with a given set of roles
func (b *OperationBuilder) GrantEmbedded(key keypair.PublicKey, roles Role) *OperationBuilder {
	switch pk := key.(type) {
	case *signing.PublicKey:
		C.self_identity_operation_builder_signing_grant_embedded(
			(*C.self_identity_operation_builder)(b),
			(*C.self_signing_public_key)(pk),
			C.ulong(roles),
		)
	case *exchange.PublicKey:
		C.self_identity_operation_builder_exchange_grant_embedded(
			(*C.self_identity_operation_builder)(b),
			(*C.self_exchange_public_key)(pk),
			C.ulong(roles),
		)
	}

	return b
}

// GrantReferenced grants roles to a key controlled by another identity
func (b *OperationBuilder) GrantReferenced(method uint16, controller *signing.PublicKey, key *signing.PublicKey, roles Role) *OperationBuilder {
	C.self_identity_operation_builder_signing_grant_referenced(
		(*C.self_identity_operation_builder)(b),
		C.ushort(method),
		(*C.self_signing_public_key)(controller),
		(*C.self_signing_public_key)(key),
		C.ulong(roles),
	)

	return b
}

// Modify modifies the roles of an existing key
func (b *OperationBuilder) Modify(key keypair.PublicKey, roles Role) *OperationBuilder {
	switch pk := key.(type) {
	case *signing.PublicKey:
		C.self_identity_operation_builder_signing_modify(
			(*C.self_identity_operation_builder)(b),
			(*C.self_signing_public_key)(pk),
			C.ulong(roles),
		)
	case *exchange.PublicKey:
		C.self_identity_operation_builder_exchange_modify(
			(*C.self_identity_operation_builder)(b),
			(*C.self_exchange_public_key)(pk),
			C.ulong(roles),
		)
	}

	return b
}

// Revoke revokes a key from a given point in time
func (b *OperationBuilder) Revoke(key keypair.PublicKey, effectiveFrom time.Time) *OperationBuilder {
	switch pk := key.(type) {
	case *signing.PublicKey:
		C.self_identity_operation_builder_signing_revoke(
			(*C.self_identity_operation_builder)(b),
			(*C.self_signing_public_key)(pk),
			C.long(effectiveFrom.Unix()),
		)
	case *exchange.PublicKey:
		C.self_identity_operation_builder_exchange_revoke(
			(*C.self_identity_operation_builder)(b),
			(*C.self_exchange_public_key)(pk),
			C.long(effectiveFrom.Unix()),
		)
	}

	return b
}

// Recover recovers an identity, revoking all existing keys
func (b *OperationBuilder) Recover(effectiveFrom time.Time) *OperationBuilder {
	C.self_identity_operation_builder_recover(
		(*C.self_identity_operation_builder)(b),
		C.long(effectiveFrom.Unix()),
	)

	return b
}

// Deactivate permanently deactivates the identity
func (b *OperationBuilder) Deactivate(effectiveFrom time.Time) *OperationBuilder {
	C.self_identity_operation_builder_deactivate(
		(*C.self_identity_operation_builder)(b),
		C.long(effectiveFrom.Unix()),
	)

	return b
}

// SignWith specifies which key to sign the operation with
func (b *OperationBuilder) SignWith(signer *signing.PublicKey) *OperationBuilder {
	C.self_identity_operation_builder_sign_with(
		(*C.self_identity_operation_builder)(b),
		(*C.self_signing_public_key)(signer),
	)

	return b
}

// Finish finalizes the operation and prepares it for execution
func (b *OperationBuilder) Finish() *Operation {
	var operation *C.self_identity_operation
	operationPtr := &operation

	C.self_identity_operation_builder_finish(
		(*C.self_identity_operation_builder)(b),
		operationPtr,
	)

	runtime.SetFinalizer(operationPtr, func(operation **C.self_identity_operation) {
		C.self_identity_operation_destroy(
			*operation,
		)
	})

	return (*Operation)(*operationPtr)
}
