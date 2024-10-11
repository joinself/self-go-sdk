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
	"errors"
	"runtime"

	"github.com/joinself/self-go-sdk-next/credential"
	"github.com/joinself/self-go-sdk-next/keypair/signing"
	"github.com/joinself/self-go-sdk-next/object"
)

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk-next/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

type Introduction struct {
	ptr *C.self_message_content_introduction
}

func newIntroduction(ptr *C.self_message_content_introduction) *Introduction {
	c := &Introduction{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *Introduction) {
		C.self_message_content_introduction_destroy(
			c.ptr,
		)
	})

	return c
}

type IntroductionBuilder struct {
	ptr *C.self_message_content_introduction_builder
}

func newIntroductionBuilder(ptr *C.self_message_content_introduction_builder) *IntroductionBuilder {
	c := &IntroductionBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *IntroductionBuilder) {
		C.self_message_content_introduction_builder_destroy(
			c.ptr,
		)
	})

	return c
}

// DeocodeIntroduction decodes an introduction message
func DecodeIntroduction(msg *Message) (*Introduction, error) {
	content := C.self_message_message_content(msg.ptr)

	var introductionContent *C.self_message_content_introduction

	status := C.self_message_content_as_introduction(
		content,
		&introductionContent,
	)

	if status > 0 {
		return nil, errors.New("failed to decode introduction message")
	}

	return newIntroduction(introductionContent), nil
}

// DocumentAddress returns the document address of the sender
func (c *Introduction) DocumentAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_message_content_introduction_document_address(c.ptr))
}

// Presentations returns any presentations the sender wishes to share to assert their identity
func (c *Introduction) Presentations() []*credential.VerifiablePresentation {
	collection := C.self_message_content_introduction_presentations(c.ptr)
	presentations := fromVerifiablePresentationCollection(collection)
	C.self_collection_verifiable_presentation_destroy(collection)

	return presentations
}

// Assets returns the any assets that support the presentations shared by the sender
func (c *Introduction) Assets() []*object.Object {
	collection := C.self_message_content_introduction_assets(c.ptr)

	var attachments []*object.Object

	for i := 0; i < int(C.self_collection_object_len(collection)); i++ {
		attachments = append(attachments, newObject(
			C.self_collection_object_at(collection, C.size_t(i)),
		))
	}

	C.self_collection_object_destroy(collection)

	return attachments
}

// NewIntroduction constructs a new introduction message
func NewIntroduction() *IntroductionBuilder {
	return newIntroductionBuilder(C.self_message_content_introduction_builder_init())
}

// DocumentAddress sets the identity document address that you wish to be known as
func (b *IntroductionBuilder) DocumentAddress(address *signing.PublicKey) *IntroductionBuilder {
	C.self_message_content_introduction_builder_document_address(
		b.ptr,
		signingPublicKeyPtr(address),
	)

	return b
}

// Presentation adds a presentation to assist in sha
func (b *IntroductionBuilder) Presentation(presentation *credential.VerifiablePresentation) *IntroductionBuilder {
	C.self_message_content_introduction_builder_presentation(
		b.ptr,
		verifiablePresentationPtr(presentation),
	)

	return b
}

// Asset attaches an object to the introduction to support verification of presentation credentials
func (b *IntroductionBuilder) Asset(attachment *object.Object) *IntroductionBuilder {
	C.self_message_content_introduction_builder_asset(
		b.ptr,
		objectPtr(attachment),
	)

	return b
}

// Finish finalizes the introduction message and prepares it for sending
func (b *IntroductionBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	status := C.self_message_content_introduction_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if status > 0 {
		return nil, errors.New("failed to build introduction request")
	}

	return newContent(finishedContent), nil
}
