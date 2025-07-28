package identity

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration
#cgo linux LDFLAGS: -lself_sdk -Wl,--allow-multiple-definition
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair"
	"github.com/joinself/self-go-sdk/keypair/exchange"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/status"
)

type Role uint64
type Method uint16
type Action int
type Description int

const (
	RoleAssertion        Role        = C.KEY_ROLE_ASSERTION
	RoleAuthentication   Role        = C.KEY_ROLE_AUTHENTICATION
	RoleVerification     Role        = C.KEY_ROLE_VERIFICATION
	RoleInvocation       Role        = C.KEY_ROLE_INVOCATION
	RoleDelegation       Role        = C.KEY_ROLE_DELEGATION
	RoleMessaging        Role        = C.KEY_ROLE_MESSAGING
	MethodAure           Method      = C.METHOD_AURE
	MethodKey            Method      = C.METHOD_KEY
	ActionGrant          Action      = C.OPERATION_ACTION_GRANT
	ActionRevoke         Action      = C.OPERATION_ACTION_REVOKE
	ActionModify         Action      = C.OPERATION_ACTION_MODIFY
	ActionRecover        Action      = C.OPERATION_ACTION_RECOVER
	ActionDeactivate     Action      = C.OPERATION_ACTION_DEACTIVATE
	DescriptionNone      Description = C.OPERATION_DESCRIPTION_NONE
	DescriptionEmbedded  Description = C.OPERATION_DESCRIPTION_EMBEDDED
	DescriptionReference Description = C.OPERATION_DESCRIPTION_REFERENCE
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

func DecodeOperation(documentAddress *signing.PublicKey, encodedOperation []byte) (*Operation, error) {
	var operation *C.self_identity_operation

	operationPtr := (*C.uint8_t)(C.CBytes(encodedOperation))
	operationLen := C.size_t(len(encodedOperation))

	result := C.self_identity_operation_decode(
		signingPublicKeyPtr(documentAddress),
		operationPtr,
		operationLen,
		&operation,
	)

	C.free(unsafe.Pointer(operationPtr))

	if result > 0 {
		return nil, status.New(result)
	}

	return newOperation(operation), nil
}

// Sequence the sequence number of the operation
func (o *Operation) Sequence(address *signing.PublicKey) uint32 {
	return uint32(C.self_identity_operation_sequence(
		operationPtr(o),
	))
}

// Hash returns the 32 byte hash of the operations content
func (o *Operation) Hash() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_identity_operation_hash(
			operationPtr(o),
		)),
		32,
	)
}

// SignedBy checks if the operation has been signed with a given key
func (o *Operation) SignedBy(address *signing.PublicKey) bool {
	return bool(C.self_identity_operation_signed_by(
		operationPtr(o),
		signingPublicKeyPtr(address),
	))
}

// Actions returns a summary of the operations actions
func (o *Operation) Actions() []*ActionSummary {
	collection := C.self_identity_operation_actions(
		operationPtr(o),
	)

	collectionLen := int(C.self_collection_identity_operation_action_len(
		collection,
	))

	actions := make([]*ActionSummary, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_identity_operation_action_at(
			collection,
			C.size_t(i),
		)

		actions[i] = newOperationAction(ptr)
	}

	C.self_collection_identity_operation_action_destroy(
		collection,
	)

	return actions
}

func (o *Operation) Encode() ([]byte, error) {
	var buf *C.self_bytes_buffer

	result := C.self_identity_operation_encode(o.ptr, &buf)
	if result > 0 {
		return nil, status.New(result)
	}

	operation := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(buf)),
		C.int(C.self_bytes_buffer_len(buf)),
	)

	C.self_bytes_buffer_destroy(buf)

	return operation, nil
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
func (b *OperationBuilder) GrantReferenced(method Method, controller *signing.PublicKey, key *signing.PublicKey, roles Role) *OperationBuilder {
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

// Threshold sets a threshold that must be met by keys that are performing an operation related to a specific role
func (b *OperationBuilder) Threshold(role Role, threshold uint8) *OperationBuilder {
	C.self_identity_operation_builder_threshold(
		b.ptr,
		C.self_identity_key_role(role),
		C.uint8_t(threshold),
	)

	return b
}

// Weight sets a weight for a key and a given role
func (b *OperationBuilder) Weight(key *signing.PublicKey, role Role, weight uint8) *OperationBuilder {
	C.self_identity_operation_builder_weight(
		b.ptr,
		signingPublicKeyPtr(key),
		C.self_identity_key_role(role),
		C.uint8_t(weight),
	)

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
