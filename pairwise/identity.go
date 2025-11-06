package pairwise

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
	"unsafe"

	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname newAddress github.com/joinself/self-go-sdk/credential.newAddress
func newAddress(*C.self_credential_address) *credential.Address

type Identity struct {
	ptr *C.self_pairwise_identity
}

func newPairwiseIdentity(ptr *C.self_pairwise_identity) *Identity {
	i := &Identity{
		ptr: ptr,
	}

	runtime.SetFinalizer(i, func(i *Identity) {
		C.self_pairwise_identity_destroy(
			i.ptr,
		)
	})

	return i
}

func pairwiseIdentityPtr(r *Identity) *C.self_pairwise_identity {
	return r.ptr
}

func DecodeIdentity(encodedIdentity []byte) (*Identity, error) {
	var identity *C.self_pairwise_identity

	identityBuf := C.CBytes(encodedIdentity)
	identityLen := len(encodedIdentity)

	result := C.self_pairwise_identity_decode(
		(*C.uint8_t)(identityBuf),
		C.size_t(identityLen),
		&identity,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newPairwiseIdentity(identity), nil
}

func (i *Identity) DocumentAddress() *credential.Address {
	return newAddress(
		C.self_pairwise_identity_document_address(
			i.ptr,
		),
	)
}

func (i *Identity) BiometricAnchor() []byte {
	anchor := C.self_pairwise_identity_biometric_anchor_hash(
		i.ptr,
	)

	if anchor == nil {
		return nil
	}

	return C.GoBytes(
		unsafe.Pointer(anchor),
		20,
	)
}

func (i *Identity) Bytes() []byte {
	return nil
}
