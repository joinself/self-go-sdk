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

	"github.com/joinself/self-go-sdk/keypair"
	"github.com/joinself/self-go-sdk/keypair/exchange"
	"github.com/joinself/self-go-sdk/keypair/signing"
)

// Document a collection of public keys tied to an identity
type Document C.self_identity_document

// NewDocument creates a new identity document
func NewDocument() *Document {
	document := (*Document)(C.self_identity_document_init())

	runtime.SetFinalizer(document, func(document *Document) {
		C.self_identity_document_destroy(
			(*C.self_identity_document)(document),
		)
	})

	return document
}

// HasRolesAt returns true if a key had a given set of roles at a time
func (d *Document) HasRolesAt(key keypair.PublicKey, roles Role, at time.Time) bool {
	switch pk := key.(type) {
	case *signing.PublicKey:
		return bool(C.self_identity_document_signing_key_roles_at(
			(*C.self_identity_document)(d),
			(*C.self_signing_public_key)(pk),
			C.ulong(roles),
			C.long(at.Unix()),
		))
	case *exchange.PublicKey:
		return bool(C.self_identity_document_exchange_key_roles_at(
			(*C.self_identity_document)(d),
			(*C.self_exchange_public_key)(pk),
			C.ulong(roles),
			C.long(at.Unix()),
		))
	default:
		return false
	}
}

// ValidAt returns true if a key was valid at a given time
func (d *Document) ValidAt(key keypair.PublicKey, at time.Time) bool {
	switch pk := key.(type) {
	case *signing.PublicKey:
		return bool(C.self_identity_document_signing_key_valid_at(
			(*C.self_identity_document)(d),
			(*C.self_signing_public_key)(pk),
			C.long(at.Unix()),
		))
	case *exchange.PublicKey:
		return bool(C.self_identity_document_exchange_key_valid_at(
			(*C.self_identity_document)(d),
			(*C.self_exchange_public_key)(pk),
			C.long(at.Unix()),
		))
	default:
		return false
	}
}

// Create creates a new operation to update the document
func (d *Document) Create() *OperationBuilder {
	builder := (*OperationBuilder)(C.self_identity_document_create(
		(*C.self_identity_document)(d),
	))

	runtime.SetFinalizer(builder, func(builder *OperationBuilder) {
		C.self_identity_operation_builder_destroy(
			(*C.self_identity_operation_builder)(builder),
		)
	})

	return builder
}
