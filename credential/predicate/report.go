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

type Report struct {
	ptr *C.self_credential_predicate_report
}

func newCredentialPredicateReport(ptr *C.self_credential_predicate_report) *Report {
	r := &Report{
		ptr: ptr,
	}

	runtime.AddCleanup(r, func(ptr *C.self_credential_predicate_report) {
		C.self_credential_predicate_report_destroy(
			ptr,
		)
	}, r.ptr)

	return r
}

// Requirements returns a list of solutions to required predicates that have not been satisfied
func (r *Report) Requirements() []*Solution {
	solutionsLen := int(C.self_credential_predicate_report_requirements_len(r.ptr))

	solutions := make([]*Solution, solutionsLen)

	for i := 0; i < solutionsLen; i++ {
		solutionPtr := C.self_credential_predicate_report_requirements_at(r.ptr, C.ulong(i))

		solutions[i] = newCredentialPredicateSolution(
			solutionPtr,
		)
	}

	return solutions
}
