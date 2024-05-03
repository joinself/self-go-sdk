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
	"time"

	"github.com/joinself/self-go-sdk/keypair"
	"github.com/joinself/self-go-sdk/keypair/signing"
)

type Role int

const (
	RoleAssertion      Role = 1
	RoleAuthentication Role = 2
	RoleVerification   Role = 3
	RoleInvocation     Role = 4
	RoleDelegation     Role = 5
	RoleIdentifier     Role = 6
	RoleMessaging      Role = 7
)

type Operation C.self_identity_operation_builder
type OperationBuilder C.self_identity_operation_builder

func NewOperation() *OperationBuilder {
	return &OperationBuilder{}
}

func (b *OperationBuilder) Identifier(address *signing.PublicKey) *OperationBuilder {
	return b
}

func (b *OperationBuilder) Timestamp(timestamp time.Time) *OperationBuilder {
	return b
}

func (b *OperationBuilder) GrantEmbedded(key keypair.PublicKey, roles Role) *OperationBuilder {
	return b
}

func (b *OperationBuilder) GrantReferenced(controller keypair.PublicKey, key keypair.PublicKey, roles Role) *OperationBuilder {
	return b
}

func (b *OperationBuilder) Modify(key keypair.PublicKey, roles Role) *OperationBuilder {
	return b
}

func (b *OperationBuilder) Revoke(key keypair.PublicKey, at time.Time) *OperationBuilder {
	return b
}

func (b *OperationBuilder) Sign(key keypair.PublicKey) *OperationBuilder {
	return b
}

func (b *OperationBuilder) Build() *Operation {
	return (*Operation)(b)
}
