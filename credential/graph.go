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
	"unsafe"

	"github.com/joinself/self-go-sdk/status"
)

//go:linkname pairwiseIdentityInterfacePtr github.com/joinself/self-go-sdk/pairwise.pairwiseIdentityInterfacePtr
func pairwiseIdentityInterfacePtr(r interface{}) (*C.self_pairwise_identity, bool)

type IdentityRecord interface {
	DocumentAddress() *Address
}

// Graph a graph of verifiable credentials
type Graph struct {
	ptr *C.self_credential_graph
}

func newCredentialGraph(ptr *C.self_credential_graph) *Graph {
	g := &Graph{
		ptr: ptr,
	}

	runtime.SetFinalizer(g, func(g *Graph) {
		C.self_credential_graph_destroy(
			g.ptr,
		)
	})

	return g
}

func credentialGraphPtr(c *Graph) *C.self_credential_graph {
	return c.ptr
}

// ValidCredentialsFor gets all valid credentials for a given holder
func (g *Graph) ValidCredentialsFor(holder *Address) ([]*VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	result := C.self_credential_graph_valid_credentials_for(
		g.ptr,
		credentialAddressPtr(holder),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromVerifiableCredentialCollection(collection)

	C.self_collection_verifiable_credential_destroy(collection)

	return credentials, nil
}

// ValidDocumentFor returns true if a holder and it's identity document have met the minimum requirements based on credentials it holds
func (g *Graph) ValidDocumentFor(holder *Address) bool {
	return bool(C.self_credential_graph_valid_document_for(
		g.ptr,
		credentialAddressPtr(holder),
	))
}

// ValidAuthenticationFor returns true if the graph contains a valid LivenessAndFacialComparison credential held by the identity,
// optionally linked to a provided authentication challenge.
func (g *Graph) ValidAuthenticationFor(identity IdentityRecord, challenge []byte) bool {
	var challengeBuf *C.uint8_t

	if len(challenge) == 32 {
		challengeBuf := (*C.uint8_t)(C.CBytes(challenge))
		defer C.free(unsafe.Pointer(challengeBuf))
	}

	pairwiseIdentityPtr, ok := pairwiseIdentityInterfacePtr(identity)
	if !ok {
		return false
	}

	return bool(C.self_credential_graph_valid_authentication_for(
		g.ptr,
		pairwiseIdentityPtr,
		challengeBuf,
	))
}

// BiometricAnchorHashFor retireves the biometric anchor hash for a given holder, if it exists
func (g *Graph) BiometricAnchorHashFor(holder *Address) ([]byte, bool) {
	anchorHash := C.self_credential_graph_biometric_anchor_hash_for(
		g.ptr,
		credentialAddressPtr(holder),
	)

	if anchorHash == nil {
		return nil, false
	}

	hash := C.GoBytes(unsafe.Pointer(C.self_bytes_buffer_buf(anchorHash)), 32)
	C.self_bytes_buffer_destroy(anchorHash)

	return hash, true
}
