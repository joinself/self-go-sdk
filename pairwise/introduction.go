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

	"github.com/joinself/self-go-sdk/credential"
)

//go:linkname fromVerifiablePresentationCollection github.com/joinself/self-go-sdk/credential.fromVerifiablePresentationCollection
func fromVerifiablePresentationCollection(ptr *C.self_collection_verifiable_presentation) []*credential.VerifiablePresentation

type Introduction struct {
	ptr *C.self_pairwise_introduction
}

func newPairwiseIntroduction(ptr *C.self_pairwise_introduction) *Introduction {
	i := &Introduction{
		ptr: ptr,
	}

	runtime.AddCleanup(i, func(i *Introduction) {
		C.self_pairwise_introduction_destroy(
			i.ptr,
		)
	}, i)

	return i
}

func pairwiseIntroductionPtr(r *Introduction) *C.self_pairwise_introduction {
	return r.ptr
}

// DocumentAddress returns the document address shared by the sender
func (i *Introduction) DocumentAddress() *credential.Address {
	return newAddress(
		C.self_pairwise_introduction_document_address(
			i.ptr,
		),
	)
}

// Presentations returns the verifiable presentations shared by the sender
func (i *Introduction) Presentations() []*credential.VerifiablePresentation {
	collection := C.self_pairwise_introduction_presentations(
		i.ptr,
	)

	presentations := fromVerifiablePresentationCollection(
		collection,
	)

	C.self_collection_verifiable_presentation_destroy(
		collection,
	)

	return presentations
}
