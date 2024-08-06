package credential

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

var (
	PresentationTypePassport             = newPresentationTypeCollection().Append("VerifiablePresentation").Append("PassportPresentation")
	PresentationTypeLiveness             = newPresentationTypeCollection().Append("VerifiablePresentation").Append("LivenessPresentation")
	PresentationTypeProfileImage         = newPresentationTypeCollection().Append("VerifiablePresentation").Append("ProfileImagePresentation")
	PresentationTypeApplicationPublisher = newPresentationTypeCollection().Append("VerifiablePresentation").Append("ApplicationPublisherPresentation")
)

type Presentation C.self_presentation
type PresentationBuilder C.self_presentation_builder
type VerifiablePresentation C.self_verifiable_presentation
type PresentationTypeCollection C.self_collection_presentation_type

// NewPresentation creates a new presentation builder
func NewPresentation() *PresentationBuilder {
	builder := (*PresentationBuilder)(C.self_presentation_builder_init())

	runtime.SetFinalizer(builder, func(builder *PresentationBuilder) {
		C.self_presentation_builder_destroy(
			(*C.self_presentation_builder)(builder),
		)
	})

	return builder
}

// PresentationType sets the type of presentation
func (b *PresentationBuilder) Presentationtype(presentationType *PresentationTypeCollection) *PresentationBuilder {
	C.self_presentation_builder_presentation_type(
		(*C.self_presentation_builder)(b),
		(*C.self_collection_presentation_type)(presentationType),
	)
	return b
}

// Holder sets the address of the credentials holder/bearer
func (b *PresentationBuilder) Holder(holderAddress *Address) *PresentationBuilder {
	C.self_presentation_builder_holder(
		(*C.self_presentation_builder)(b),
		(*C.self_credential_address)(holderAddress),
	)
	return b
}

// CredentialAdd adds a verifiable credential to the presentation
func (b *PresentationBuilder) CredentialAdd(credential *VerifiableCredential) *PresentationBuilder {
	C.self_presentation_builder_credential_add(
		(*C.self_presentation_builder)(b),
		(*C.self_verifiable_credential)(credential),
	)

	return b
}

// Finish generates and prepares the presentation for being signed by an account
func (b *PresentationBuilder) Finish() (*Presentation, error) {
	var presentation *C.self_presentation
	presentationPtr := &presentation

	status := C.self_presentation_builder_finish(
		(*C.self_presentation_builder)(b),
		presentationPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to create presentation")
	}

	runtime.SetFinalizer(presentationPtr, func(presentation **C.self_presentation) {
		C.self_presentation_destroy(
			*presentation,
		)
	})

	return (*Presentation)(*presentationPtr), nil
}

// PresentationType returns the type of presentation
func (p *VerifiablePresentation) PresentationType() *PresentationTypeCollection {
	collection := (*PresentationTypeCollection)(C.self_verifiable_presentation_presentation_type(
		(*C.self_verifiable_presentation)(p),
	))

	runtime.SetFinalizer(collection, func(collection *PresentationTypeCollection) {
		C.self_collection_presentation_type_destroy(
			(*C.self_collection_presentation_type)(collection),
		)
	})

	return collection
}

// Holder returns the subject of the credential's holder
func (p *VerifiablePresentation) Holder() *Address {
	holder := (*Address)(C.self_verifiable_presentation_holder(
		(*C.self_verifiable_presentation)(p),
	))

	runtime.SetFinalizer(holder, func(address *Address) {
		C.self_credential_address_destroy(
			(*C.self_credential_address)(address),
		)
	})

	return holder
}

// Credential returns the verifiable credentials contained in the presentation
func (p *VerifiablePresentation) Credentials() *VerifiableCredentialCollection {
	collection := C.self_verifiable_presentation_credentials(
		(*C.self_verifiable_presentation)(p),
	)

	c := (*VerifiableCredentialCollection)(collection)

	runtime.SetFinalizer(c, func(collection *VerifiableCredentialCollection) {
		C.self_collection_verifiable_credential_destroy(
			(*C.self_collection_verifiable_credential)(collection),
		)
	})

	return c
}

// Validate validates the contents of the presentation and it's signatures
func (p *VerifiablePresentation) Validate() error {
	status := C.self_verifiable_presentation_validate(
		(*C.self_verifiable_presentation)(p),
	)

	if status > 0 {
		return errors.New("presentation invalid")
	}

	return nil
}

func NewPresentationTypeCollection() *PresentationTypeCollection {
	collection := (*PresentationTypeCollection)(C.self_collection_presentation_type_init())

	runtime.SetFinalizer(collection, func(collection *PresentationTypeCollection) {
		C.self_collection_presentation_type_destroy(
			(*C.self_collection_presentation_type)(collection),
		)
	})

	return collection
}

func newPresentationTypeCollection() *PresentationTypeCollection {
	return (*PresentationTypeCollection)(C.self_collection_presentation_type_init())
}

func (c *PresentationTypeCollection) Length() int {
	return int(C.self_collection_presentation_type_len(
		(*C.self_collection_presentation_type)(c),
	))
}

func (c *PresentationTypeCollection) Get(index int) string {
	return C.GoString(C.self_collection_presentation_type_at(
		(*C.self_collection_presentation_type)(c),
		C.ulong(index),
	))
}

func (c *PresentationTypeCollection) Append(element string) *PresentationTypeCollection {
	elementC := C.CString(element)

	C.self_collection_presentation_type_append(
		(*C.self_collection_presentation_type)(c),
		elementC,
	)

	C.free(unsafe.Pointer(elementC))

	return c
}
