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
)

type Relationship struct {
	ptr *C.self_pairwise_relationship
}

func newPairwiseRelationship(ptr *C.self_pairwise_relationship) *Relationship {
	r := &Relationship{
		ptr: ptr,
	}

	runtime.SetFinalizer(r, func(r *Relationship) {
		C.self_pairwise_relationship_destroy(
			r.ptr,
		)
	})

	return r
}

func (r *Relationship) Status() Status {
	return Status(C.self_pairwise_relationship_status(
		r.ptr,
	))
}

func (r *Relationship) With() *Identity {
	return newPairwiseIdentity(C.self_pairwise_relationship_with_identity(
		r.ptr,
	))
}

func (r *Relationship) As() *Identity {
	return newPairwiseIdentity(C.self_pairwise_relationship_as_identity(
		r.ptr,
	))
}
