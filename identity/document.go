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

	"github.com/joinself/self-go-sdk/keypair"
	"github.com/joinself/self-go-sdk/keypair/exchange"
	"github.com/joinself/self-go-sdk/keypair/signing"
)

type AddressDescription interface {
	Address() keypair.PublicKey
	Controller() *signing.PublicKey
}

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

//go:linkname exchangePublicKeyPtr github.com/joinself/self-go-sdk/keypair/exchange.exchangePublicKeyPtr
func exchangePublicKeyPtr(p *exchange.PublicKey) *C.self_exchange_public_key

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk/keypair/signing.newSigningPublicKey
func newSigningPublicKey(*C.self_signing_public_key) *signing.PublicKey

//go:linkname newExchangePublicKey github.com/joinself/self-go-sdk/keypair/exchange.newExchangePublicKey
func newExchangePublicKey(*C.self_exchange_public_key) *exchange.PublicKey

//go:linkname toSigningPublicKeyCollection github.com/joinself/self-go-sdk/keypair/signing.toSigningPublicKeyCollection
func toSigningPublicKeyCollection(c []*signing.PublicKey) *C.self_collection_signing_public_key

//go:linkname fromSigningPublicKeyCollection github.com/joinself/self-go-sdk/keypair/signing.fromSigningPublicKeyCollection
func fromSigningPublicKeyCollection(ptr *C.self_collection_signing_public_key) []*signing.PublicKey

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

// SigningKeys returns all signing keys that have been added to the document
func (d *Document) SigningKeys() []*signing.PublicKey {
	collection := C.self_identity_document_signing_keys(
		d.ptr,
	)

	keys := fromSigningPublicKeyCollection(
		collection,
	)

	C.self_collection_signing_public_key_destroy(
		collection,
	)

	return keys
}

// SigningKeysWithRoles returns all signing keys that have been added to the document with a given set of roles
func (d *Document) SigningKeysWithRoles(roles Role) []*signing.PublicKey {
	collection := C.self_identity_document_signing_keys_with_roles(
		d.ptr,
		C.uint64_t(roles),
	)

	keys := fromSigningPublicKeyCollection(
		collection,
	)

	C.self_collection_signing_public_key_destroy(
		collection,
	)

	return keys
}

// SigningKeysWithRolesAt returns all signing keys that have been added to the document with a given set of roles at a given timeframe
func (d *Document) SigningKeysWithRolesAt(roles Role, at time.Time) []*signing.PublicKey {
	collection := C.self_identity_document_signing_keys_with_roles_at(
		d.ptr,
		C.uint64_t(roles),
		C.int64_t(at.Unix()),
	)

	keys := fromSigningPublicKeyCollection(
		collection,
	)

	C.self_collection_signing_public_key_destroy(
		collection,
	)

	return keys
}

// ThresholdMet checks if the provided signers have the required weight for a role's threshold
func (d *Document) ThresholdMet(role Role, signers []*signing.PublicKey) bool {
	collection := toSigningPublicKeyCollection(signers)

	defer C.self_collection_signing_public_key_destroy(collection)

	return bool(C.self_identity_document_threshold_met(
		d.ptr,
		C.self_identity_key_role(role),
		collection,
	))
}

// ThresholdMet checks if the provided signers have the required weight for a role's threshold at a given time
func (d *Document) ThresholdMetAt(role Role, at time.Time, signers []*signing.PublicKey) bool {
	collection := toSigningPublicKeyCollection(signers)

	defer C.self_collection_signing_public_key_destroy(collection)

	return bool(C.self_identity_document_threshold_met_at(
		d.ptr,
		C.self_identity_key_role(role),
		C.int64_t(at.Unix()),
		collection,
	))
}

// DescriptionsAt returns all key descriptions at a given time
func (d *Document) DescriptionsAt(at time.Time) []AddressDescription {
	collection := C.self_identity_document_descriptions_at(
		d.ptr,
		C.int64_t(at.Unix()),
	)

	collectionLen := int(C.self_collection_identity_operation_description_len(
		collection,
	))

	descriptions := make([]AddressDescription, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_identity_operation_description_at(
			collection,
			C.size_t(i),
		)

		switch C.self_identity_operation_description_type_of(ptr) {
		case C.OPERATION_DESCRIPTION_EMBEDDED:
			descriptions[i] = newEmbeddedDescription(
				C.self_identity_operation_description_as_embedded(ptr),
			)
		case C.OPERATION_DESCRIPTION_REFERENCE:
			descriptions[i] = newReferenceDescription(
				C.self_identity_operation_description_as_reference(ptr),
			)
		}
	}

	C.self_collection_identity_operation_description_destroy(
		collection,
	)

	return descriptions
}

// Create creates a new operation to update the document
func (d *Document) Create() *OperationBuilder {
	return newOperationBuilder(C.self_identity_document_create(
		d.ptr,
	))
}
