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

const (
	PresentationTypePassport                    = "PassportPresentation"
	PresentationTypeFacialComparsion            = "FacialComparisonPresentation"
	PresentationTypeLivenessAndFacialComparsion = "LivenessAndFacialComparisonPresentation"
	PresentationTypeBiometricAnchor             = "BiometricAnchorPresentation"
	PresentationTypeSharingAgreement            = "SharingAgreementPresentation"
	PresentationTypeProfile                     = "ProfilePresentation"
	PresentationTypeContactDetails              = "ContactDetailsPresentation"
	PresentationTypeApplicationPublisher        = "ApplicationPublisherPresentation"
)

type Presentation struct {
	ptr *C.self_presentation
}

func newPresentation(ptr *C.self_presentation) *Presentation {
	p := &Presentation{
		ptr: ptr,
	}

	runtime.SetFinalizer(p, func(b *Presentation) {
		C.self_presentation_destroy(
			p.ptr,
		)
	})

	return p
}

func presentationPtr(p *Presentation) *C.self_presentation {
	return p.ptr
}

type PresentationBuilder struct {
	ptr *C.self_presentation_builder
}

func newPresentationBuilder(ptr *C.self_presentation_builder) *PresentationBuilder {
	b := &PresentationBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(b, func(b *PresentationBuilder) {
		C.self_presentation_builder_destroy(
			b.ptr,
		)
	})

	return b
}

type VerifiablePresentation struct {
	ptr *C.self_verifiable_presentation
}

func newVerifiablePresentation(ptr *C.self_verifiable_presentation) *VerifiablePresentation {
	return &VerifiablePresentation{
		ptr: ptr,
	}
}

func verifiablePresentationPtr(v *VerifiablePresentation) *C.self_verifiable_presentation {
	return v.ptr
}

// NewPresentation creates a new presentation builder
func NewPresentation() *PresentationBuilder {
	return newPresentationBuilder(C.self_presentation_builder_init())
}

// PresentationType sets the type of presentation
func (b *PresentationBuilder) PresentationType(presentationType ...string) *PresentationBuilder {
	collection := toPresentationTypeCollection(presentationType)

	C.self_presentation_builder_presentation_type(
		b.ptr,
		collection,
	)

	C.self_collection_presentation_type_destroy(
		collection,
	)

	return b
}

// Holder sets the address of the credentials holder/bearer
func (b *PresentationBuilder) Holder(holderAddress *Address) *PresentationBuilder {
	C.self_presentation_builder_holder(
		b.ptr,
		holderAddress.ptr,
	)
	return b
}

// CredentialAdd adds a verifiable credential to the presentation
func (b *PresentationBuilder) CredentialAdd(credentials ...*VerifiableCredential) *PresentationBuilder {
	for i := range credentials {
		C.self_presentation_builder_credential_add(
			b.ptr,
			credentials[i].ptr,
		)
	}

	return b
}

// Finish generates and prepares the presentation for being signed by an account
func (b *PresentationBuilder) Finish() (*Presentation, error) {
	var presentation *C.self_presentation

	result := C.self_presentation_builder_finish(
		b.ptr,
		&presentation,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newPresentation(presentation), nil
}

// DecodeVerifiablePresentation decodes a verifiable credential from it's json form
func DecodeVerifiablePresentation(encodedPresentation []byte) (*VerifiablePresentation, error) {
	var verifiablePresentation *C.self_verifiable_presentation

	encodedBuf := C.CBytes(encodedPresentation)
	encodedLen := len(encodedPresentation)

	defer func() {
		C.free(encodedBuf)
	}()

	result := C.self_verifiable_presentation_decode(
		&verifiablePresentation,
		(*C.uint8_t)(encodedBuf),
		(C.size_t)(encodedLen),
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newVerifiablePresentation(verifiablePresentation), nil
}

// PresentationType returns the type of presentation
func (p *VerifiablePresentation) PresentationType() []string {
	collection := C.self_verifiable_presentation_type_of(
		p.ptr,
	)

	presentationType := fromPresentationTypeCollection(collection)

	C.self_collection_presentation_type_destroy(
		collection,
	)

	return presentationType
}

// Holder returns the subject of the credential's holder
func (p *VerifiablePresentation) Holder() *Address {
	return newAddress(C.self_verifiable_presentation_holder(
		p.ptr,
	))
}

// Credential returns the verifiable credentials contained in the presentation
func (p *VerifiablePresentation) Credentials() []*VerifiableCredential {
	collection := C.self_verifiable_presentation_credentials(
		p.ptr,
	)

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials
}

// Validate validates the contents of the presentation and it's signatures
func (p *VerifiablePresentation) Validate() error {
	result := C.self_verifiable_presentation_validate(
		p.ptr,
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// Encode returns a json encoded verifiable presentation
func (p *VerifiablePresentation) Encode() ([]byte, error) {
	var encodedPresentationBuffer *C.self_bytes_buffer
	encodedPresentationBufferPtr := &encodedPresentationBuffer

	result := C.self_verifiable_presentation_encode(
		p.ptr,
		encodedPresentationBufferPtr,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	encodedPresentation := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(*encodedPresentationBufferPtr)),
		C.int(C.self_bytes_buffer_len(*encodedPresentationBufferPtr)),
	)

	C.self_bytes_buffer_destroy(
		*encodedPresentationBufferPtr,
	)

	return encodedPresentation, nil
}

func toPresentationTypeCollection(presentationType []string) *C.self_collection_presentation_type {
	collection := C.self_collection_presentation_type_init()

	for i := 0; i < len(presentationType); i++ {
		typ := C.CString(presentationType[i])

		C.self_collection_presentation_type_append(
			collection,
			typ,
		)

		C.free(unsafe.Pointer(typ))
	}

	return collection
}

func fromPresentationTypeCollection(collection *C.self_collection_presentation_type) []string {
	collectionLen := int(C.self_collection_presentation_type_len(
		collection,
	))

	presentationType := make([]string, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_presentation_type_at(
			collection,
			C.size_t(i),
		)

		presentationType[i] = C.GoString(ptr)
	}

	return presentationType
}

func toVerifiablePresentationCollection(presentations []*VerifiablePresentation) *C.self_collection_verifiable_presentation {
	collection := C.self_collection_verifiable_presentation_init()

	for i := 0; i < len(presentations); i++ {

		C.self_collection_verifiable_presentation_append(
			collection,
			verifiablePresentationPtr(presentations[i]),
		)
	}

	return collection
}

func fromVerifiablePresentationCollection(collection *C.self_collection_verifiable_presentation) []*VerifiablePresentation {
	collectionLen := int(C.self_collection_verifiable_presentation_len(
		collection,
	))

	presentations := make([]*VerifiablePresentation, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_verifiable_presentation_at(
			collection,
			C.size_t(i),
		)

		presentations[i] = newVerifiablePresentation(ptr)
	}

	return presentations
}
