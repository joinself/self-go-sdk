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
)

const (
	PresentationPassport             = PresentationType(C.PRESENTATION_PASSPORT)
	PresentationLiveness             = PresentationType(C.PRESENTATION_LIVENESS)
	PresentationProfileImage         = PresentationType(C.PRESENTATION_PROFILE_IMAGE)
	PresentationApplicationPublisher = PresentationType(C.PRESENTATION_APPLICATION_PUBLISHER)
)

type PresentationType = C.self_presentation_type
type Presentation C.self_presentation
type PresentationBuilder C.self_presentation_builder
type VerifiablePresentation C.self_verifiable_presentation

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
func (b *PresentationBuilder) Presentationtype(presentationType PresentationType) *PresentationBuilder {
	C.self_presentation_builder_presentation_type(
		(*C.self_presentation_builder)(b),
		uint32(presentationType),
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
	var presentationPtr *C.self_presentation

	status := C.self_presentation_builder_finish(
		(*C.self_presentation_builder)(b),
		&presentationPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to create presentation")
	}

	presentation := (*Presentation)(presentationPtr)

	runtime.SetFinalizer(presentation, func(presentation *Presentation) {
		C.self_presentation_destroy(
			(*C.self_presentation)(presentation),
		)
	})

	return presentation, nil
}

// PresentationType returns the type of presentation
func (p *VerifiablePresentation) PresentationType() PresentationType {
	return PresentationType(C.self_verifiable_presentation_presentation_type(
		(*C.self_verifiable_presentation)(p),
	))
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
