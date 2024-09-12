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

	"github.com/joinself/self-go-sdk-next/keypair"
	"github.com/joinself/self-go-sdk-next/keypair/exchange"
	"github.com/joinself/self-go-sdk-next/keypair/signing"
)

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk-next/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

//go:linkname exchangePublicKeyPtr github.com/joinself/self-go-sdk-next/keypair/exchange.exchangePublicKeyPtr
func exchangePublicKeyPtr(p *exchange.PublicKey) *C.self_exchange_public_key

// Document a collection of public keys tied to an identity
type Document struct {
	ptr *C.self_identity_document
}

func newIdentityDocument(ptr *C.self_identity_document) *Document {
	d := &Document{
		ptr: ptr,
	}

	runtime.SetFinalizer(d, func(d *Document) {
		C.self_identity_document_destroy(
			d.ptr,
		)
	})

	return d
}

func identityDocumentPtr(d *Document) *C.self_identity_document {
	return d.ptr
}

// NewDocument creates a new identity document
func NewDocument() *Document {
	return newIdentityDocument(C.self_identity_document_init())
}

// HasRolesAt returns true if a key had a given set of roles at a time
func (d *Document) HasRolesAt(key keypair.PublicKey, roles Role, at time.Time) bool {
	switch pk := key.(type) {
	case *signing.PublicKey:
		return bool(C.self_identity_document_signing_key_roles_at(
			d.ptr,
			signingPublicKeyPtr(pk),
			C.uint64_t(roles),
			C.int64_t(at.Unix()),
		))
	case *exchange.PublicKey:
		return bool(C.self_identity_document_exchange_key_roles_at(
			d.ptr,
			exchangePublicKeyPtr(pk),
			C.uint64_t(roles),
			C.int64_t(at.Unix()),
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
			d.ptr,
			signingPublicKeyPtr(pk),
			C.int64_t(at.Unix()),
		))
	case *exchange.PublicKey:
		return bool(C.self_identity_document_exchange_key_valid_at(
			d.ptr,
			exchangePublicKeyPtr(pk),
			C.int64_t(at.Unix()),
		))
	default:
		return false
	}
}

// Create creates a new operation to update the document
func (d *Document) Create() *OperationBuilder {
	return newOperationBuilder(C.self_identity_document_create(
		d.ptr,
	))
}
