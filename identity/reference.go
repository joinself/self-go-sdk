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

	"github.com/joinself/self-go-sdk/keypair"
	"github.com/joinself/self-go-sdk/keypair/signing"
)

type Reference struct {
	ptr *C.self_identity_operation_description_reference
}

func newReferenceDescription(ptr *C.self_identity_operation_description_reference) *Reference {
	r := &Reference{
		ptr: ptr,
	}

	runtime.AddCleanup(r, func(r *Reference) {
		C.self_identity_operation_description_reference_destroy(
			r.ptr,
		)
	}, r)

	return r
}

func referenceDescriptionPtr(r *Reference) *C.self_identity_operation_description_reference {
	return r.ptr
}

// Address returns the address of the reference key
func (r *Reference) Address() keypair.PublicKey {
	switch C.self_identity_operation_description_reference_address_type(r.ptr) {
	case C.KEYPAIR_SIGNING:
		return newSigningPublicKey(
			C.self_identity_operation_description_reference_address_as_signing(
				r.ptr,
			),
		)
	case C.KEYPAIR_EXCHANGE:
		return newExchangePublicKey(
			C.self_identity_operation_description_reference_address_as_exchange(
				r.ptr,
			),
		)
	default:
		return nil
	}
}

// Controller returns the controller of the key or nil if not specified
func (r *Reference) Controller() *signing.PublicKey {
	controller := C.self_identity_operation_description_reference_controller(
		r.ptr,
	)

	if controller == nil {
		return nil
	}

	return newSigningPublicKey(controller)
}

// Method returns the controller of the controllers did if a controller is specified
func (r *Reference) Method() Method {
	return Method(C.self_identity_operation_description_reference_method(
		r.ptr,
	))
}
