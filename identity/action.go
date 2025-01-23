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
)

type ActionSummary struct {
	ptr *C.self_identity_operation_action
}

func newOperationAction(ptr *C.self_identity_operation_action) *ActionSummary {
	a := &ActionSummary{
		ptr: ptr,
	}

	runtime.SetFinalizer(a, func(a *ActionSummary) {
		C.self_identity_operation_action_destroy(
			a.ptr,
		)
	})

	return a
}

func operationActionPtr(a *ActionSummary) *C.self_identity_operation_action {
	return a.ptr
}

// Action returns the action being performed
func (a *ActionSummary) Action() Action {
	return Action(C.self_identity_operation_action_action_type(
		a.ptr,
	))
}

// Address returns the address
func (a *ActionSummary) Description() Description {
	return Description(C.self_identity_operation_action_description_type(
		a.ptr,
	))
}

// Embedded returns the embedded key description, or nil if there is no embedded key description
func (a *ActionSummary) Embedded() *Embedded {
	embedded := C.self_identity_operation_action_description_as_embedded(
		a.ptr,
	)

	if embedded == nil {
		return nil
	}

	return newEmbeddedDescription(embedded)
}

// Reference returns the reference key description, or nil if there is no reference key description
func (a *ActionSummary) Reference() *Reference {
	reference := C.self_identity_operation_action_description_as_reference(
		a.ptr,
	)

	if reference == nil {
		return nil
	}

	return newReferenceDescription(reference)
}

// From returns the time the action will take effect from
func (a *ActionSummary) From() time.Time {
	return time.Unix(int64(C.self_identity_operation_action_from(
		a.ptr,
	)), 0)
}

// Roles returns the roles the action will put into effect, if applicable (grant, modify)
func (a *ActionSummary) Roles() Role {
	return Role(C.self_identity_operation_action_roles(
		a.ptr,
	))
}
