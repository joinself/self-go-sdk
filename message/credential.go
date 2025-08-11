package message

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
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/status"
)

type Credential struct {
	ptr *C.self_message_content_credential
}

func newCredential(ptr *C.self_message_content_credential) *Credential {
	c := &Credential{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *Credential) {
		C.self_message_content_credential_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialBuilder struct {
	ptr *C.self_message_content_credential_builder
}

func newCredentialBuilder(ptr *C.self_message_content_credential_builder) *CredentialBuilder {
	c := &CredentialBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialBuilder) {
		C.self_message_content_credential_builder_destroy(
			c.ptr,
		)
	})

	return c
}

// DecodeCredential decodes a message to a credential
func DecodeCredential(content *Content) (*Credential, error) {
	contentPtr := contentPtr(content)

	var credentialContent *C.self_message_content_credential

	result := C.self_message_content_as_credential(
		contentPtr,
		&credentialContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCredential(credentialContent), nil
}

// Presentations returns veriable presentations that have been asserted by the sender
func (c *Credential) Presentations() []*credential.VerifiablePresentation {
	collection := C.self_message_content_credential_verifiable_presentations(
		c.ptr,
	)

	credentials := fromVerifiablePresentationCollection(
		collection,
	)

	C.self_collection_verifiable_presentation_destroy(
		collection,
	)

	return credentials
}

// Credentials returns veriable credentials that have been asserted by the sender
func (c *Credential) Credentials() []*credential.VerifiableCredential {
	collection := C.self_message_content_credential_verifiable_credentials(
		c.ptr,
	)

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials
}

// Assets returns assets that are used in support of attested verifiable credenntials
func (c *Credential) Assets() []*object.Object {
	collection := C.self_message_content_credential_assets(
		c.ptr,
	)

	credentials := fromObjectCollection(
		collection,
	)

	C.self_collection_object_destroy(
		collection,
	)

	return credentials
}

// NewCredential creates a new credential
func NewCredential() *CredentialBuilder {
	return newCredentialBuilder(
		C.self_message_content_credential_builder_init(),
	)
}

// VerifiablePresentation attaches a verified presentation of credentails to the message
func (b *CredentialBuilder) VerifiablePresentation(presentations ...*credential.VerifiablePresentation) *CredentialBuilder {
	for i := range presentations {
		C.self_message_content_credential_builder_verifiable_presentation(
			b.ptr,
			verifiablePresentationPtr(presentations[i]),
		)
	}
	return b
}

// VerifiableCredential attaches a verified credential to the message
func (b *CredentialBuilder) VerifiableCredential(credentials ...*credential.VerifiableCredential) *CredentialBuilder {
	for i := range credentials {
		C.self_message_content_credential_builder_verifiable_credential(
			b.ptr,
			verifiableCredentialPtr(credentials[i]),
		)
	}
	return b
}

// Asset attaches an asset used to support a verifiable credential to the message
func (b *CredentialBuilder) Asset(assets ...*object.Object) *CredentialBuilder {
	for i := range assets {
		C.self_message_content_credential_builder_asset(
			b.ptr,
			objectPtr(assets[i]),
		)
	}
	return b
}

// Finish finalises the response and builds the content
func (b *CredentialBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_credential_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}
