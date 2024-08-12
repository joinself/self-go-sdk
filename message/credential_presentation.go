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

	"github.com/joinself/self-go-sdk-next/credential"
)

type CredentialPresentationRequest C.self_message_content_credential_presentation_request
type CredentialPresentationResponse C.self_message_content_credential_presentation_response
type CredentialPresentationRequestBuilder C.self_message_content_credential_presentation_request_builder
type CredentialPresentationResponseBuilder C.self_message_content_credential_presentation_response_builder

// DecodeCredentialPresentationRequest decodes a message to a credential presentation request
func DecodeCredentialPresentationRequest(msg *Message) (*CredentialPresentationRequest, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var credentialPresentationRequestContent *C.self_message_content_credential_presentation_request
	credentialPresentationRequestContentPtr := &credentialPresentationRequestContent

	status := C.self_message_content_as_credential_presentation_request(
		content,
		credentialPresentationRequestContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to decode credential presentation request message")
	}

	runtime.SetFinalizer(credentialPresentationRequestContentPtr, func(credentialPresentationRequest **C.self_message_content_credential_presentation_request) {
		C.self_message_content_credential_presentation_request_destroy(
			*credentialPresentationRequest,
		)
	})

	return (*CredentialPresentationRequest)(*credentialPresentationRequestContentPtr), nil
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
	finishedContentPtr := &finishedContent

	status := C.self_message_content_credential_presentation_request_builder_finish(
		(*C.self_message_content_credential_presentation_request_builder)(b),
		finishedContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to build credential verificaiton request")
	}

	runtime.SetFinalizer(finishedContentPtr, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(*finishedContentPtr), nil
}

// DecodeCredentialPresentationResponse decodes a message to a credential presentation response
func DecodeCredentialPresentationResponse(msg *Message) (*CredentialPresentationResponse, error) {
	content := C.self_message_message_content((*C.self_message)(msg))

	var credentialPresentationResponseContent *C.self_message_content_credential_presentation_response
	credentialPresentationResponseContentPtr := &credentialPresentationResponseContent

	status := C.self_message_content_as_credential_presentation_response(
		content,
		credentialPresentationResponseContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to decode credential presentation response message")
	}

	runtime.SetFinalizer(credentialPresentationResponseContentPtr, func(credentialPresentationResponse **C.self_message_content_credential_presentation_response) {
		C.self_message_content_credential_presentation_response_destroy(
			*credentialPresentationResponse,
		)
	})

	return (*CredentialPresentationResponse)(*credentialPresentationResponseContentPtr), nil
}

// ResponseTo returns the id of the request that is being responded to
func (c *CredentialPresentationResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_credential_presentation_response_response_to(
			(*C.self_message_content_credential_presentation_response)(c),
		)),
		20,
	)
}

// Status returns the status of the request
func (c *CredentialPresentationResponse) Status() ResponseStatus {
	return ResponseStatus(C.self_message_content_credential_presentation_response_status(
		(*C.self_message_content_credential_presentation_response)(c),
	))
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

// ResponseTo sets the request id that is being responded to
func (b *CredentialPresentationResponseBuilder) ResponseTo(requestID []byte) *CredentialPresentationResponseBuilder {
	if len(requestID) != 20 {
		return b
	}

	requestIDBuf := C.CBytes(
		requestID,
	)

	C.self_message_content_credential_presentation_response_builder_response_to(
		(*C.self_message_content_credential_presentation_response_builder)(b),
		(*C.uint8_t)(requestIDBuf),
	)

	C.free(requestIDBuf)

	return b
}

// ResponseTo sets the request id that is being responded to
func (b *CredentialPresentationResponseBuilder) Status(status ResponseStatus) *CredentialPresentationResponseBuilder {
	C.self_message_content_credential_presentation_response_builder_status(
		(*C.self_message_content_credential_presentation_response_builder)(b),
		uint32(status),
	)

	return b
}

// VerifiablePresentation attaches a verified presentation of credentails to the response
func (b *CredentialPresentationResponseBuilder) VerifiablePresentation(presentation *credential.VerifiablePresentation) *CredentialPresentationResponseBuilder {
	C.self_message_content_credential_presentation_response_builder_verifiable_presentation(
		(*C.self_message_content_credential_presentation_response_builder)(b),
		(*C.self_verifiable_presentation)(presentation),
	)
	return b
}

// Finish finalises the response and builds the content
func (b *CredentialPresentationResponseBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content
	finishedContentPtr := &finishedContent

	status := C.self_message_content_credential_presentation_response_builder_finish(
		(*C.self_message_content_credential_presentation_response_builder)(b),
		finishedContentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to build credential verificaiton response")
	}

	runtime.SetFinalizer(finishedContentPtr, func(content **C.self_message_content) {
		C.self_message_content_destroy(
			*content,
		)
	})

	return (*Content)(*finishedContentPtr), nil
}
