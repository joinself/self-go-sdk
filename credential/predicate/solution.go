package predicate

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

type Solution struct {
	ptr *C.self_credential_predicate_solution
}

func newCredentialPredicateSolution(ptr *C.self_credential_predicate_solution) *Solution {
	s := &Solution{
		ptr: ptr,
	}

	runtime.AddCleanup(s, func(s *Solution) {
		C.self_credential_predicate_solution_destroy(
			s.ptr,
		)
	}, s)

	return s
}

// Options returns options and variations that will solve the required predicates
func (s *Solution) Options() [][]*Predicator {
	optionsLen := int(C.self_credential_predicate_solution_len(s.ptr))

	options := make([][]*Predicator, optionsLen)

	for i := 0; i < optionsLen; i++ {
		optionPtr := C.self_credential_predicate_solution_at(s.ptr, C.ulong(i))

		options[i] = fromCredentialPredicatorCollection(
			optionPtr,
		)

		C.self_collection_credential_predicator_destroy(optionPtr)
	}

	return options
}

func fromCredentialPredicatorCollection(collection *C.self_collection_credential_predicator) []*Predicator {
	collectionLen := int(C.self_collection_credential_predicator_len(
		collection,
	))

	predicates := make([]*Predicator, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_credential_predicator_at(
			collection,
			C.size_t(i),
		)

		predicates[i] = newCredentialPredicator(ptr)
	}

	return predicates
}
