package message

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
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/credential"
)

type CredentialPresentationRequest C.self_message_content_credential_presentation_request
type CredentialPresentationResponse C.self_message_content_credential_presentation_response
type CredentialPresentationRequestBuilder C.self_message_content_credential_presentation_request_builder
type CredentialPresentationResponseBuilder C.self_message_content_credential_presentation_response_builder

// DecodeCredentialPresentationRequest decodes a message to a credential presentation request
func DecodeCredentialPresentationRequest(msg *Message) (*CredentialPresentationRequest, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var credentialPresentationRequestContent *C.self_message_content_credential_presentation_request

	status := C.self_message_content_as_credential_presentation_request(
		content,
		&credentialPresentationRequestContent,
	)

	if status > 0 {
		return nil, errors.New("failed to decode credential presentation request message")
	}

	credentialPresentationRequest := (*CredentialPresentationRequest)(credentialPresentationRequestContent)

	runtime.SetFinalizer(credentialPresentationRequest, func(credentialPresentationRequest *CredentialPresentationRequest) {
		C.self_message_content_credential_presentation_request_destroy(
			(*C.self_message_content_credential_presentation_request)(credentialPresentationRequest),
		)
	})

	return credentialPresentationRequest, nil
}

// Type returns the type of credential that presentation is being requested for
func (c *CredentialPresentationRequest) Type() *credential.PresentationTypeCollection {
	collection := (*credential.PresentationTypeCollection)(C.self_message_content_credential_presentation_request_presentation_type(
		(*C.self_message_content_credential_presentation_request)(c),
	))

	runtime.SetFinalizer(collection, func(collection *credential.PresentationTypeCollection) {
		C.self_collection_presentation_type_destroy(
			(*C.self_collection_presentation_type)(collection),
		)
	})

	return collection
}

// Details returns details of the requested credential presentations
func (c *CredentialPresentationRequest) Details() *credential.CredentialPresentationDetailCollection {
	return (*credential.CredentialPresentationDetailCollection)(C.self_message_content_credential_presentation_request_details(
		(*C.self_message_content_credential_presentation_request)(c),
	))
}

// Type returns the time the request expires at
func (c *CredentialPresentationRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_credential_presentation_request_expires(
		(*C.self_message_content_credential_presentation_request)(c),
	)), 0)
}

// NewCredentialPresentationRequest creates a new credential presentation request
func NewCredentialPresentationRequest() *CredentialPresentationRequestBuilder {
	builder := (*CredentialPresentationRequestBuilder)(C.self_message_content_credential_presentation_request_builder_init())

	runtime.SetFinalizer(builder, func(builder *CredentialPresentationRequestBuilder) {
		C.self_message_content_credential_presentation_request_builder_destroy(
			(*C.self_message_content_credential_presentation_request_builder)(builder),
		)
	})

	return builder
}

// Type sets the type of presentation being requested
func (b *CredentialPresentationRequestBuilder) Type(presentationType *credential.PresentationTypeCollection) *CredentialPresentationRequestBuilder {
	C.self_message_content_credential_presentation_request_builder_presentation_type(
		(*C.self_message_content_credential_presentation_request_builder)(b),
		(*C.self_collection_presentation_type)(presentationType),
	)
	return b
}

// Details specifies the details of the credentials being requested for presentation
func (b *CredentialPresentationRequestBuilder) Details(credentialType *credential.CredentialTypeCollection, subject string) *CredentialPresentationRequestBuilder {
	subjectC := C.CString(subject)

	C.self_message_content_credential_presentation_request_builder_details(
		(*C.self_message_content_credential_presentation_request_builder)(b),
		(*C.self_collection_credential_type)(credentialType),
		subjectC,
	)

	C.free(unsafe.Pointer(subjectC))

	return b
}

// Expires sets the time that the request expires at
func (b *CredentialPresentationRequestBuilder) Expires(expires time.Time) *CredentialPresentationRequestBuilder {
	C.self_message_content_credential_presentation_request_builder_expires(
		(*C.self_message_content_credential_presentation_request_builder)(b),
		C.long(expires.Unix()),
	)
	return b
}

// Finish finalises the request and builds the content
func (b *CredentialPresentationRequestBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	status := C.self_message_content_credential_presentation_request_builder_finish(
		(*C.self_message_content_credential_presentation_request_builder)(b),
		&finishedContent,
	)

	if status > 0 {
		return nil, errors.New("failed to build credential verificaiton request")
	}

	content := (*Content)(finishedContent)

	runtime.SetFinalizer(content, func(content *Content) {
		C.self_message_content_destroy(
			(*C.self_message_content)(content),
		)
	})

	return content, nil
}

// DecodeCredentialPresentationResponse decodes a message to a credential presentation response
func DecodeCredentialPresentationResponse(msg *Message) (*CredentialPresentationResponse, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var credentialPresentationResponseContent *C.self_message_content_credential_presentation_response

	status := C.self_message_content_as_credential_presentation_response(
		content,
		&credentialPresentationResponseContent,
	)

	if status > 0 {
		return nil, errors.New("failed to decode credential presentation response message")
	}

	credentialPresentationResponse := (*CredentialPresentationResponse)(credentialPresentationResponseContent)

	runtime.SetFinalizer(credentialPresentationResponse, func(credentialPresentationResponse *CredentialPresentationResponse) {
		C.self_message_content_credential_presentation_response_destroy(
			(*C.self_message_content_credential_presentation_response)(credentialPresentationResponse),
		)
	})

	return credentialPresentationResponse, nil
}

// Status returns the status of the request
func (c *CredentialPresentationResponse) Status() int {
	return 0
}

// Credentials returns verified credentials that have been asserted by the responder
func (c *CredentialPresentationResponse) Credentials() *credential.VerifiableCredentialCollection {
	return (*credential.VerifiableCredentialCollection)(C.self_message_content_credential_presentation_response_verifiable_presentations(
		(*C.self_message_content_credential_presentation_response)(c),
	))
}

// NewCredentialPresentationResponse creates a new credential presentation response
func NewCredentialPresentationResponse() *CredentialPresentationResponseBuilder {
	builder := (*CredentialPresentationResponseBuilder)(C.self_message_content_credential_presentation_response_builder_init())

	runtime.SetFinalizer(builder, func(builder *CredentialPresentationResponseBuilder) {
		C.self_message_content_credential_presentation_response_builder_destroy(
			(*C.self_message_content_credential_presentation_response_builder)(builder),
		)
	})

	return builder
}

// VerifiableCredential attaches a verified presentation of credentails to the response
func (b *CredentialPresentationResponseBuilder) VerifiableCredential(presentation *credential.VerifiablePresentation) *CredentialPresentationResponseBuilder {
	C.self_message_content_credential_presentation_response_builder_verifiable_presentation(
		(*C.self_message_content_credential_presentation_response_builder)(b),
		(*C.self_verifiable_presentation)(presentation),
	)
	return b
}

// Finish finalises the response and builds the content
func (b *CredentialPresentationResponseBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	status := C.self_message_content_credential_presentation_response_builder_finish(
		(*C.self_message_content_credential_presentation_response_builder)(b),
		&finishedContent,
	)

	if status > 0 {
		return nil, errors.New("failed to build credential verificaiton response")
	}

	content := (*Content)(finishedContent)

	runtime.SetFinalizer(content, func(content *Content) {
		C.self_message_content_destroy(
			(*C.self_message_content)(content),
		)
	})

	return content, nil
}
