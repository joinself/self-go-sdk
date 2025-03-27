package event

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
)

type ContentSummary struct {
	ptr *C.self_message_content_summary
}

func newContentSummary(ptr *C.self_message_content_summary) *ContentSummary {
	c := &ContentSummary{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *ContentSummary) {
		C.self_message_content_summary_destroy(
			c.ptr,
		)
	})

	return c
}

func contentSummaryPtr(c *ContentSummary) *C.self_message_content_summary {
	return c.ptr
}

func contentSummaryType(contentSummary *C.self_message_content_summary) ContentType {
	switch C.self_message_content_summary_type_of(contentSummary) {
	case C.CONTENT_CUSTOM:
		return ContentTypeCustom
	case C.CONTENT_CHAT:
		return ContentTypeChat
	case C.CONTENT_RECEIPT:
		return ContentTypeReceipt
	case C.CONTENT_CREDENTIAL:
		return ContentTypeCredential
	case C.CONTENT_INTRODUCTION:
		return ContentTypeIntroduction
	case C.CONTENT_DISCOVERY_REQUEST:
		return ContentTypeDiscoveryRequest
	case C.CONTENT_DISCOVERY_RESPONSE:
		return ContentTypeDiscoveryResponse
	case C.CONTENT_SIGNING_REQUEST:
		return ContentTypeSigningRequest
	case C.CONTENT_SIGNING_RESPONSE:
		return ContentTypeSigningResponse
	case C.CONTENT_ACCOUNT_PAIRING_REQUEST:
		return ContentTypeAccountPairingRequest
	case C.CONTENT_ACCOUNT_PAIRING_RESPONSE:
		return ContentTypeAccountPairingResponse
	case C.CONTENT_CREDENTIAL_VERIFICATION_REQUEST:
		return ContentTypeCredentialVerificationRequest
	case C.CONTENT_CREDENTIAL_VERIFICATION_RESPONSE:
		return ContentTypeCredentialVerificationResponse
	case C.CONTENT_CREDENTIAL_PRESENTATION_REQUEST:
		return ContentTypeCredentialPresentationRequest
	case C.CONTENT_CREDENTIAL_PRESENTATION_RESPONSE:
		return ContentTypeCredentialPresentationResponse
	default:
		return ContentTypeUnknown
	}
}

// ID returns the unique id of the message
func (c *ContentSummary) ID() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_summary_id(c.ptr)),
		20,
	)
}

// ContentType get the content type
func (c *ContentSummary) ContentType() ContentType {
	return contentSummaryType(
		c.ptr,
	)
}

// Descriptions get the descriptions
func (c *ContentSummary) Descriptions() []*ContentSummaryDescription {
	collection := C.self_message_content_summary_descriptions(c.ptr)

	var descriptions []*ContentSummaryDescription

	for i := 0; i < int(C.self_collection_message_content_summary_description_len(collection)); i++ {
		descriptions = append(descriptions, newContentSummaryDescription(
			C.self_collection_message_content_summary_description_at(collection, C.size_t(i)),
		))
	}

	C.self_collection_message_content_summary_description_destroy(collection)

	return descriptions
}

type ContentSummaryDescription struct {
	ptr *C.self_message_content_summary_description
}

func newContentSummaryDescription(ptr *C.self_message_content_summary_description) *ContentSummaryDescription {
	c := &ContentSummaryDescription{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *ContentSummaryDescription) {
		C.self_message_content_summary_description_destroy(
			c.ptr,
		)
	})

	return c
}
