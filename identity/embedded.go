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

type Embedded struct {
	ptr *C.self_identity_operation_description_embedded
}

func newEmbeddedDescription(ptr *C.self_identity_operation_description_embedded) *Embedded {
	e := &Embedded{
		ptr: ptr,
	}

	runtime.AddCleanup(e, func(ptr *C.self_identity_operation_description_embedded) {
		C.self_identity_operation_description_embedded_destroy(
			ptr,
		)
	}, e.ptr)

	return e
}

func embeddedDescriptionPtr(e *Embedded) *C.self_identity_operation_description_embedded {
	return e.ptr
}

// Address returns the address of the embedded key
func (e *Embedded) Address() keypair.PublicKey {
	switch C.self_identity_operation_description_embedded_address_type(e.ptr) {
	case C.KEYPAIR_SIGNING:
		return newSigningPublicKey(
			C.self_identity_operation_description_embedded_address_as_signing(
				e.ptr,
			),
		)
	case C.KEYPAIR_EXCHANGE:
		return newExchangePublicKey(
			C.self_identity_operation_description_embedded_address_as_exchange(
				e.ptr,
			),
		)
	default:
		return nil
	}
}

// Controller returns the controller of the key or nil if not specified
func (e *Embedded) Controller() *signing.PublicKey {
	controller := C.self_identity_operation_description_embedded_controller(
		e.ptr,
	)

	if controller == nil {
		return nil
	}

	return newSigningPublicKey(controller)
}
