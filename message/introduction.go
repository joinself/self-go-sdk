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
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/pairwise"
	"github.com/joinself/self-go-sdk/status"
	"github.com/joinself/self-go-sdk/token"
)

//go:linkname newAddress github.com/joinself/self-go-sdk/credential.newAddress
func newAddress(*C.self_credential_address) *credential.Address

//go:linkname credentialAddressPtr github.com/joinself/self-go-sdk/credential.credentialAddressPtr
func credentialAddressPtr(ptr *credential.Address) *C.self_credential_address

//go:linkname newToken github.com/joinself/self-go-sdk/token.newToken
func newToken(ptr *C.self_token) *token.Token

//go:linkname tokenPtr github.com/joinself/self-go-sdk/token.tokenPtr
func tokenPtr(ptr *token.Token) *C.self_token

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

//go:linkname newPairwiseIntroduction github.com/joinself/self-go-sdk/pairwise.newPairwiseIntroduction
func newPairwiseIntroduction(ptr *C.self_pairwise_introduction) *pairwise.Introduction

type Introduction struct {
	ptr *C.self_message_content_introduction
}

func newIntroduction(ptr *C.self_message_content_introduction) *Introduction {
	c := &Introduction{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(c *Introduction) {
		C.self_message_content_introduction_destroy(
			c.ptr,
		)
	}, c)

	return c
}

type IntroductionBuilder struct {
	ptr *C.self_message_content_introduction_builder
}

func newIntroductionBuilder(ptr *C.self_message_content_introduction_builder) *IntroductionBuilder {
	c := &IntroductionBuilder{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(c *IntroductionBuilder) {
		C.self_message_content_introduction_builder_destroy(
			c.ptr,
		)
	}, c)

	return c
}

// DeocodeIntroduction decodes an introduction message
func DecodeIntroduction(content *Content) (*Introduction, error) {
	contentPtr := contentPtr(content)

	var introductionContent *C.self_message_content_introduction

	result := C.self_message_content_as_introduction(
		contentPtr,
		&introductionContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newIntroduction(introductionContent), nil
}

// Introduction returns the pairwise introduction from the sender
func (c *Introduction) Introduction() (*pairwise.Introduction, error) {
	var introduction *C.self_pairwise_introduction

	result := C.self_message_content_introduction_introduction(
		c.ptr,
		&introduction,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newPairwiseIntroduction(introduction), nil
}

// DocumentAddress returns the document address of the sender
func (c *Introduction) DocumentAddress() *credential.Address {
	return newAddress(C.self_message_content_introduction_document_address(c.ptr))
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

// Tokens returns the any tokens that the sender has issued
func (c *Introduction) Tokens() ([]*token.Token, error) {
	var collection *C.self_collection_token

	result := C.self_message_content_introduction_tokens(
		c.ptr,
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	var tokens []*token.Token

	for i := 0; i < int(C.self_collection_token_len(collection)); i++ {
		tokens = append(tokens, newToken(
			C.self_collection_token_at(collection, C.size_t(i)),
		))
	}

	C.self_collection_token_destroy(collection)

	return tokens, nil
}

// NewIntroduction constructs a new introduction message
func NewIntroduction() *IntroductionBuilder {
	return newIntroductionBuilder(C.self_message_content_introduction_builder_init())
}

// DocumentAddress sets the identity document address that you wish to be known as
func (b *IntroductionBuilder) DocumentAddress(address *credential.Address) *IntroductionBuilder {
	C.self_message_content_introduction_builder_document_address(
		b.ptr,
		credentialAddressPtr(address),
	)

	return b
}

// Presentation adds a presentation to assist in sha
func (b *IntroductionBuilder) Presentation(presentations ...*credential.VerifiablePresentation) *IntroductionBuilder {
	for i := range presentations {
		C.self_message_content_introduction_builder_presentation(
			b.ptr,
			verifiablePresentationPtr(presentations[i]),
		)
	}

	return b
}

// Asset attaches an object to the introduction to support verification of presentation credentials
func (b *IntroductionBuilder) Asset(attachments ...*object.Object) *IntroductionBuilder {
	for i := range attachments {
		C.self_message_content_introduction_builder_asset(
			b.ptr,
			objectPtr(attachments[i]),
		)
	}

	return b
}

// Token attaches a token to the introduction that can be used by the recipient(s)
func (b *IntroductionBuilder) Token(tokens ...*token.Token) *IntroductionBuilder {
	for i := range tokens {
		C.self_message_content_introduction_builder_token(
			b.ptr,
			tokenPtr(tokens[i]),
		)
	}

	return b
}

// Finish finalizes the introduction message and prepares it for sending
func (b *IntroductionBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_introduction_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}
