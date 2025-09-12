package predicate

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

	"github.com/joinself/self-go-sdk/credential"
)

//go:linkname toVerifiableCredentialCollection github.com/joinself/self-go-sdk/credential.toVerifiableCredentialCollection
func toVerifiableCredentialCollection(credentials []*credential.VerifiableCredential) *C.self_collection_verifiable_credential

//go:linkname fromVerifiableCredentialCollection github.com/joinself/self-go-sdk/credential.fromVerifiableCredentialCollection
func fromVerifiableCredentialCollection(ptr *C.self_collection_verifiable_credential) []*credential.VerifiableCredential

type Tree struct {
	ptr *C.self_credential_predicate_tree
}

func newCredentialPredicateTree(ptr *C.self_credential_predicate_tree) *Tree {
	p := &Tree{
		ptr: ptr,
	}

	runtime.SetFinalizer(p, func(p *Tree) {
		C.self_credential_predicate_tree_destroy(
			p.ptr,
		)
	})

	return p
}

func credentialPredicateTreePtr(f *Tree) *C.self_credential_predicate_tree {
	return f.ptr
}

func NewTree(predicate *Predicate) *Tree {
	return newCredentialPredicateTree(
		C.self_credential_predicate_tree_init(
			predicate.ptr,
		),
	)
}

// FindOptimalMatch finds the optimal set of credentials that match the criteria in the predicate tree.
// Returns false if no valid match could be determined
func (p *Tree) FindOptimalMatch(credentials []*credential.VerifiableCredential) ([]*credential.VerifiableCredential, bool) {
	unfilteredCredentials := toVerifiableCredentialCollection(
		credentials,
	)

	filteredCredentials := C.self_credential_predicate_tree_find_optimal_match(
		p.ptr,
		unfilteredCredentials,
	)

	C.self_collection_verifiable_credential_destroy(unfilteredCredentials)

	if filteredCredentials == nil {
		return nil, false
	}

	candidates := fromVerifiableCredentialCollection(
		filteredCredentials,
	)

	C.self_collection_verifiable_credential_destroy(filteredCredentials)

	return candidates, true
}

// FindMissingPredicates finds missing predicates that are required but have not been met by
// the provided credentials. The report will contain required predicates, each containing a
// solultion of predicate options that will satisfy the requirements
func (p *Tree) FindMissingPredicates(credentials []*credential.VerifiableCredential) *Report {
	unfilteredCredentials := toVerifiableCredentialCollection(
		credentials,
	)

	defer C.self_collection_verifiable_credential_destroy(unfilteredCredentials)

	return newCredentialPredicateReport(
		C.self_credential_predicate_tree_find_missing_predicates(
			p.ptr,
			unfilteredCredentials,
		),
	)
}

// Graphviz render the predicate tree to graphviz dot format
func (p *Tree) Graphviz() string {
	dotBuf := C.self_credential_predicate_tree_graphviz(
		p.ptr,
	)

	dot := C.GoString(
		C.self_string_buffer_ptr(
			dotBuf,
		),
	)

	C.self_string_buffer_destroy(dotBuf)

	return dot
}
