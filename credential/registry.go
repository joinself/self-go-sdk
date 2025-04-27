package credential

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

	"github.com/joinself/self-go-sdk/status"
)

// TrustedIssuerRegistry a registry of trusted credential issuers
type TrustedIssuerRegistry struct {
	ptr *C.self_trusted_issuer_registry
}

func newTrustedIssuerRegistry(ptr *C.self_trusted_issuer_registry) *TrustedIssuerRegistry {
	r := &TrustedIssuerRegistry{
		ptr: ptr,
	}

	runtime.SetFinalizer(r, func(r *TrustedIssuerRegistry) {
		C.self_trusted_issuer_registry_destroy(
			r.ptr,
		)
	})

	return r
}

func trustedIssuerRegistryPtr(c *TrustedIssuerRegistry) *C.self_trusted_issuer_registry {
	return c.ptr
}

// NewTrustedIssuerRegistry creates a new empty issuer registry
func NewTrustedIssuerRegistry() *TrustedIssuerRegistry {
	return newTrustedIssuerRegistry(
		C.self_trusted_issuer_registry_init(),
	)
}

// ProductionTrustedIssuerRegistry creates a default issuer registry
// for use with Self's production environment
func ProductionTrustedIssuerRegistry() *TrustedIssuerRegistry {
	return newTrustedIssuerRegistry(
		C.self_trusted_issuer_registry_default_production(),
	)
}

// SandboxTrustedIssuerRegistry creates a default issuer registry
// for use with Self's sandbox environment
func SandboxTrustedIssuerRegistry() *TrustedIssuerRegistry {
	return newTrustedIssuerRegistry(
		C.self_trusted_issuer_registry_default_sandbox(),
	)
}

// DefaultCredentialTypes returns the default credential types issued
// by self's registry.
func DefaultCredentialTypes() []string {
	collection := C.self_trusted_issuer_registry_default_credential_types()

	collectionLen := int(C.self_collection_string_buffer_len(
		collection,
	))

	credentialTypes := make([]string, collectionLen)

	for i := 0; i < collectionLen; i++ {
		buf := C.self_collection_string_buffer_at(
			collection,
			C.size_t(i),
		)

		credentialTypes[i] = C.GoString(
			C.self_string_buffer_ptr(buf),
		)

		C.self_string_buffer_destroy(buf)
	}

	C.self_collection_string_buffer_destroy(
		collection,
	)

	return credentialTypes
}

// DefaultIssuerEpoch returns the default epoch from when Self
// was granted permission to issue credentials
func DefaultIssuerEpoch() time.Time {
	return time.Unix(
		int64(C.self_trusted_issuer_registry_default_issuer_epoch()),
		0,
	)
}

// AddIssuer adds an issuer to the trusted issuer registry
func (r *TrustedIssuerRegistry) AddIssuer(issuer *Address) bool {
	return bool(C.self_trusted_issuer_registry_issuer_add(
		r.ptr,
		issuer.ptr,
	))
}

// RemoveIssuer removes an issuer from the trusted issuer registry
func (r *TrustedIssuerRegistry) RemoveIssuer(issuer *Address) bool {
	return bool(C.self_trusted_issuer_registry_issuer_remove(
		r.ptr,
		issuer.ptr,
	))
}

// GrantAuthority grants authority to an issuer for a given credential type and time period
func (r *TrustedIssuerRegistry) GrantAuthority(issuer *Address, credentialType string, granted time.Time, revoked *time.Time) error {
	var revokedAt *C.int64_t

	if revoked != nil {
		r := C.int64_t(revoked.Unix())
		revokedAt = &r
	}

	credentialTypePtr := C.CString(credentialType)

	result := C.self_trusted_issuer_registry_authority_grant(
		r.ptr,
		issuer.ptr,
		credentialTypePtr,
		C.int64_t(granted.Unix()),
		revokedAt,
	)

	C.free(unsafe.Pointer(credentialTypePtr))

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// RevokeAuthority revokes authority given to an issuer for a given credential type and time period
func (r *TrustedIssuerRegistry) RevokeAuthority(issuer *Address, credentialType string, revoked time.Time) error {
	credentialTypePtr := C.CString(credentialType)

	result := C.self_trusted_issuer_registry_authority_revoke(
		r.ptr,
		issuer.ptr,
		credentialTypePtr,
		C.int64_t(revoked.Unix()),
	)

	C.free(unsafe.Pointer(credentialTypePtr))

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// AuthorityAt returns true if a issuer had authority to issue a given credential type at the specified time
func (r *TrustedIssuerRegistry) AuthorityAt(issuer *Address, credentialType string, at time.Time) bool {
	credentialTypePtr := C.CString(credentialType)
	defer C.free(unsafe.Pointer(credentialTypePtr))

	return bool(C.self_trusted_issuer_registry_authority_at(
		r.ptr,
		issuer.ptr,
		credentialTypePtr,
		C.int64_t(at.Unix()),
	))
}

// AuthorityFor returns a list of all credentials an issuer can currently issue
func (r *TrustedIssuerRegistry) AuthorityFor(issuer *Address) ([]string, error) {
	var collection *C.self_collection_string_buffer

	result := C.self_trusted_issuer_registry_authority_for(
		r.ptr,
		issuer.ptr,
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	collectionLen := int(C.self_collection_string_buffer_len(
		collection,
	))

	credentialTypes := make([]string, collectionLen)

	for i := 0; i < collectionLen; i++ {
		buf := C.self_collection_string_buffer_at(
			collection,
			C.size_t(i),
		)

		credentialTypes[i] = C.GoString(
			C.self_string_buffer_ptr(buf),
		)

		C.self_string_buffer_destroy(buf)
	}

	C.self_collection_string_buffer_destroy(
		collection,
	)

	return credentialTypes, nil
}
