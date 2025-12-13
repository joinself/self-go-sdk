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
	"unsafe"

	"github.com/joinself/self-go-sdk/object"
)

type ContentSummaryDescriptionType int

const (
	ContentSummaryDescriptionTypeUnknown        ContentSummaryDescriptionType = C.CONTENT_SUMMARY_DESCRIPTION_UNKNOWN
	ContentSummaryDescriptionTypeChatMessage    ContentSummaryDescriptionType = C.CONTENT_SUMMARY_DESCRIPTION_CHAT_MESSAGE
	ContentSummaryDescriptionTypeChatReference  ContentSummaryDescriptionType = C.CONTENT_SUMMARY_DESCRIPTION_CHAT_REFERENCE
	ContentSummaryDescriptionTypeChatAttachment ContentSummaryDescriptionType = C.CONTENT_SUMMARY_DESCRIPTION_CHAT_ATTACHMENT
	ContentSummaryDescriptionTypeCredential     ContentSummaryDescriptionType = C.CONTENT_SUMMARY_DESCRIPTION_CREDENTIAL
	ContentSummaryDescriptionTypePresentation   ContentSummaryDescriptionType = C.CONTENT_SUMMARY_DESCRIPTION_PRESENTATION
	ContentSummaryDescriptionTypeAsset          ContentSummaryDescriptionType = C.CONTENT_SUMMARY_DESCRIPTION_ASSET
	ContentSummaryDescriptionTypeSignature      ContentSummaryDescriptionType = C.CONTENT_SUMMARY_DESCRIPTION_SIGNATURE
)

type ContentSummary struct {
	ptr *C.self_message_content_summary
}

func newContentSummary(ptr *C.self_message_content_summary) *ContentSummary {
	c := &ContentSummary{
		ptr: ptr,
	}

	runtime.AddCleanup(c, func(c *ContentSummary) {
		C.self_message_content_summary_destroy(
			c.ptr,
		)
	}, c)

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

	runtime.AddCleanup(c, func(c *ContentSummaryDescription) {
		C.self_message_content_summary_description_destroy(
			c.ptr,
		)
	}, c)

	return c
}

func contentSummaryDescriptionType(contentSummaryDescription *C.self_message_content_summary_description) ContentSummaryDescriptionType {
	switch C.self_message_content_summary_description_type_of(contentSummaryDescription) {

	case C.CONTENT_SUMMARY_DESCRIPTION_CHAT_MESSAGE:
		return ContentSummaryDescriptionTypeChatMessage
	case C.CONTENT_SUMMARY_DESCRIPTION_CHAT_REFERENCE:
		return ContentSummaryDescriptionTypeChatReference
	case C.CONTENT_SUMMARY_DESCRIPTION_CHAT_ATTACHMENT:
		return ContentSummaryDescriptionTypeChatAttachment
	case C.CONTENT_SUMMARY_DESCRIPTION_CREDENTIAL:
		return ContentSummaryDescriptionTypeCredential
	case C.CONTENT_SUMMARY_DESCRIPTION_PRESENTATION:
		return ContentSummaryDescriptionTypePresentation
	case C.CONTENT_SUMMARY_DESCRIPTION_ASSET:
		return ContentSummaryDescriptionTypeAsset
	case C.CONTENT_SUMMARY_DESCRIPTION_SIGNATURE:
		return ContentSummaryDescriptionTypeSignature
	default:
		return ContentSummaryDescriptionTypeUnknown
	}
}

func (d *ContentSummaryDescription) DescriptionType() ContentSummaryDescriptionType {
	return contentSummaryDescriptionType(
		d.ptr,
	)
}

func (d *ContentSummaryDescription) ChatMessage() (string, bool) {
	chatMessagePtr := C.self_message_content_summary_description_as_chat_message(
		d.ptr,
	)

	if chatMessagePtr == nil {
		return "", false
	}

	chatMessage := C.GoString(
		C.self_string_buffer_ptr(chatMessagePtr),
	)

	C.self_string_buffer_destroy(chatMessagePtr)

	return chatMessage, true
}

func (d *ContentSummaryDescription) ChatAttachment() (*object.Object, bool) {
	chatAttachmentPtr := C.self_message_content_summary_description_as_chat_attachment(
		d.ptr,
	)

	if chatAttachmentPtr == nil {
		return nil, false
	}

	return newObject(chatAttachmentPtr), true
}

func (d *ContentSummaryDescription) Credential() ([]string, bool) {
	collection := C.self_message_content_summary_description_as_credential(
		d.ptr,
	)

	if collection == nil {
		return nil, false
	}

	credentials := fromCredentialTypeCollection(collection)

	C.self_collection_credential_type_destroy(
		collection,
	)

	return credentials, true
}

func (d *ContentSummaryDescription) Presentation() ([]string, bool) {
	collection := C.self_message_content_summary_description_as_presentation(
		d.ptr,
	)

	if collection == nil {
		return nil, false
	}

	presentations := fromPresentationTypeCollection(collection)

	C.self_collection_presentation_type_destroy(
		collection,
	)

	return presentations, true
}

func (d *ContentSummaryDescription) Asset() (*object.Object, bool) {
	assetPtr := C.self_message_content_summary_description_as_asset(
		d.ptr,
	)

	if assetPtr == nil {
		return nil, false
	}

	return newObject(assetPtr), true
}

func (d *ContentSummaryDescription) Signature() (SigningPayloadType, bool) {
	switch C.self_message_content_summary_description_as_signature(d.ptr) {
	case C.SIGNING_PAYLOAD_IDENTITY_DOCUMENT_OPERATION:
		return SigningPayloadIdentityDocumentOperation, true
	default:
		return SigningPayloadUnknown, false
	}
}
