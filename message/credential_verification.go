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
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/status"
)

//go:linkname newObject github.com/joinself/self-go-sdk/object.newObject
func newObject(ptr *C.self_object) *object.Object

//go:linkname objectPtr github.com/joinself/self-go-sdk/object.objectPtr
func objectPtr(ptr *object.Object) *C.self_object

type CredentialVerificationRequest struct {
	ptr *C.self_message_content_credential_verification_request
}

func newCredentialVerificationRequest(ptr *C.self_message_content_credential_verification_request) *CredentialVerificationRequest {
	c := &CredentialVerificationRequest{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialVerificationRequest) {
		C.self_message_content_credential_verification_request_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialVerificationResponse struct {
	ptr *C.self_message_content_credential_verification_response
}

func newCredentialVerificationResponse(ptr *C.self_message_content_credential_verification_response) *CredentialVerificationResponse {
	c := &CredentialVerificationResponse{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialVerificationResponse) {
		C.self_message_content_credential_verification_response_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialVerificationRequestBuilder struct {
	ptr *C.self_message_content_credential_verification_request_builder
}

func newCredentialVerificationRequestBuilder(ptr *C.self_message_content_credential_verification_request_builder) *CredentialVerificationRequestBuilder {
	c := &CredentialVerificationRequestBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialVerificationRequestBuilder) {
		C.self_message_content_credential_verification_request_builder_destroy(
			c.ptr,
		)
	})

	return c
}

type CredentialVerificationResponseBuilder struct {
	ptr *C.self_message_content_credential_verification_response_builder
}

func newCredentialVerificationResponseBuilder(ptr *C.self_message_content_credential_verification_response_builder) *CredentialVerificationResponseBuilder {
	c := &CredentialVerificationResponseBuilder{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialVerificationResponseBuilder) {
		C.self_message_content_credential_verification_response_builder_destroy(
			c.ptr,
		)
	})

	return c
}

// DecodeCredentialVerificationRequest decodes a message to a credential verification request
func DecodeCredentialVerificationRequest(content *Content) (*CredentialVerificationRequest, error) {
	contentPtr := contentPtr(content)

	var credentialVerificationRequestContent *C.self_message_content_credential_verification_request

	result := C.self_message_content_as_credential_verification_request(
		contentPtr,
		&credentialVerificationRequestContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCredentialVerificationRequest(credentialVerificationRequestContent), nil
}

// Type returns the type of credential that verification is being requested for
func (c *CredentialVerificationRequest) Type() []string {
	collection := C.self_message_content_credential_verification_request_credential_type(
		c.ptr,
	)

	credentialType := fromCredentialTypeCollection(
		collection,
	)

	C.self_collection_credential_type_destroy(
		collection,
	)

	return credentialType
}

// Proof returns associated verifiable credential proof to support the verification request
func (c *CredentialVerificationRequest) Proof() []*credential.VerifiablePresentation {
	collection := C.self_message_content_credential_verification_request_proof(
		c.ptr,
	)

	presentations := fromVerifiablePresentationCollection(
		collection,
	)

	C.self_collection_verifiable_presentation_destroy(
		collection,
	)

	return presentations
}

// Evidence returns associated data to be used as evidence to support the verification request
func (c *CredentialVerificationRequest) Evidence() []*CredentialVerificationEvidence {
	collection := C.self_message_content_credential_verification_request_evidence(
		c.ptr,
	)

	evidence := fromCredentialVerificationEvidenceCollection(
		collection,
	)

	C.self_collection_credential_verification_evidence_destroy(
		collection,
	)

	return evidence
}

// Parameters returns associated data to be used as parameters to support the verification request
func (c *CredentialVerificationRequest) Parameters() []*CredentialVerificationParameter {
	collection := C.self_message_content_credential_verification_request_parameters(
		c.ptr,
	)

	parameter := fromCredentialVerificationParameterCollection(
		collection,
	)

	C.self_collection_credential_verification_parameter_destroy(
		collection,
	)

	return parameter
}

// Type returns the time the request expires at
func (c *CredentialVerificationRequest) Expires() time.Time {
	return time.Unix(int64(C.self_message_content_credential_verification_request_expires(
		c.ptr,
	)), 0)
}

// NewCredentialVerificationRequest creates a new credential verification request
func NewCredentialVerificationRequest() *CredentialVerificationRequestBuilder {
	return newCredentialVerificationRequestBuilder(
		C.self_message_content_credential_verification_request_builder_init(),
	)
}

// Type sets the type of credential being requested
func (b *CredentialVerificationRequestBuilder) Type(credentialType []string) *CredentialVerificationRequestBuilder {
	collection := toCredentialTypeCollection(credentialType)

	C.self_message_content_credential_verification_request_builder_credential_type(
		b.ptr,
		collection,
	)

	C.self_collection_credential_type_destroy(
		collection,
	)

	return b
}

// Proof attaches proof to the credential verification request
func (b *CredentialVerificationRequestBuilder) Proof(proof *credential.VerifiablePresentation) *CredentialVerificationRequestBuilder {
	C.self_message_content_credential_verification_request_builder_proof(
		b.ptr,
		verifiablePresentationPtr(proof),
	)
	return b
}

// Evidence attaches evidence to the credential verification request
func (b *CredentialVerificationRequestBuilder) Evidence(evidenceType string, evidence *object.Object) *CredentialVerificationRequestBuilder {
	evidenceTypeC := C.CString(evidenceType)

	C.self_message_content_credential_verification_request_builder_evidence(
		b.ptr,
		evidenceTypeC,
		objectPtr(evidence),
	)

	C.free(unsafe.Pointer(evidenceTypeC))

	return b
}

// Evidence attaches evidence to the credential verification request
func (b *CredentialVerificationRequestBuilder) Parameter(key string, value any) *CredentialVerificationRequestBuilder {
	var paramValuePtr *C.self_message_content_parameter_value

	switch v := value.(type) {
	case []byte:
		valueBuf := C.CBytes(v)
		valueLen := len(v)

		paramValuePtr = C.self_message_content_parameter_value_bytes_create(
			(*C.uint8_t)(valueBuf),
			(C.size_t)(valueLen),
		)

		C.free(valueBuf)
	case string:
		valuePtr := C.CString(v)

		paramValuePtr = C.self_message_content_parameter_value_string_create(
			valuePtr,
		)

		C.free(unsafe.Pointer(valuePtr))
	case bool:
		paramValuePtr = C.self_message_content_parameter_value_bool_create(
			C.bool(v),
		)
	case int8:
		paramValuePtr = C.self_message_content_parameter_value_int8_create(
			C.int8_t(v),
		)
	case int16:
		paramValuePtr = C.self_message_content_parameter_value_int16_create(
			C.int16_t(v),
		)
	case int32:
		paramValuePtr = C.self_message_content_parameter_value_int32_create(
			C.int32_t(v),
		)
	case int64:
		paramValuePtr = C.self_message_content_parameter_value_int64_create(
			C.int64_t(v),
		)
	case uint8:
		paramValuePtr = C.self_message_content_parameter_value_uint8_create(
			C.uint8_t(v),
		)
	case uint16:
		paramValuePtr = C.self_message_content_parameter_value_uint16_create(
			C.uint16_t(v),
		)
	case uint32:
		paramValuePtr = C.self_message_content_parameter_value_uint32_create(
			C.uint32_t(v),
		)
	case uint64:
		paramValuePtr = C.self_message_content_parameter_value_uint64_create(
			C.uint64_t(v),
		)
	case float32:
		paramValuePtr = C.self_message_content_parameter_value_float32_create(
			C.float(v),
		)
	case float64:
		paramValuePtr = C.self_message_content_parameter_value_float64_create(
			C.double(v),
		)
	case [][]byte:
		collection := C.self_collection_bytes_buffer_init()

		for i := 0; i < len(v); i++ {
			valueBuf := C.CBytes(v[i])
			valueLen := len(v[i])

			C.self_collection_bytes_buffer_append(
				collection,
				(*C.uint8_t)(valueBuf),
				(C.size_t)(valueLen),
			)

			C.free(valueBuf)
		}

		paramValuePtr = C.self_message_content_parameter_value_array_bytes_create(
			collection,
		)

		C.self_collection_bytes_buffer_destroy(collection)
	case []string:
		collection := C.self_collection_string_buffer_init()

		for i := 0; i < len(v); i++ {
			valuePtr := C.CString(v[i])

			paramValuePtr = C.self_message_content_parameter_value_string_create(
				valuePtr,
			)

			C.self_collection_string_buffer_append(
				collection,
				valuePtr,
			)

			C.free(unsafe.Pointer(valuePtr))
		}

		paramValuePtr = C.self_message_content_parameter_value_array_string_create(
			collection,
		)

		C.self_collection_string_buffer_destroy(collection)
	default:
		return b
	}

	keyC := C.CString(key)

	C.self_message_content_credential_verification_request_builder_parameter(
		b.ptr,
		keyC,
		paramValuePtr,
	)

	C.free(unsafe.Pointer(keyC))

	return b
}

// Expires sets the time that the request expires at
func (b *CredentialVerificationRequestBuilder) Expires(expires time.Time) *CredentialVerificationRequestBuilder {
	C.self_message_content_credential_verification_request_builder_expires(
		b.ptr,
		C.int64_t(expires.Unix()),
	)
	return b
}

// Finish finalises the request and builds the content
func (b *CredentialVerificationRequestBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_credential_verification_request_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newContent(finishedContent), nil
}

// DecodeCredentialVerificationResponse decodes a message to a credential verification response
func DecodeCredentialVerificationResponse(content *Content) (*CredentialVerificationResponse, error) {
	contentPtr := contentPtr(content)

	var credentialVerificationResponseContent *C.self_message_content_credential_verification_response

	result := C.self_message_content_as_credential_verification_response(
		contentPtr,
		&credentialVerificationResponseContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCredentialVerificationResponse(credentialVerificationResponseContent), nil
}

// ResponseTo returns the id of the request that is being responded to
func (c *CredentialVerificationResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_message_content_credential_verification_response_response_to(
			c.ptr,
		)),
		20,
	)
}

// Status returns the status of the request
func (c *CredentialVerificationResponse) Status() ResponseStatus {
	return ResponseStatus(C.self_message_content_credential_verification_response_status(
		c.ptr,
	))
}

// Credentials returns verified credentials that have been asserted by the responder
func (c *CredentialVerificationResponse) Credentials() []*credential.VerifiableCredential {
	collection := C.self_message_content_credential_verification_response_verifiable_credentials(
		c.ptr,
	)

	credentials := fromVerifiableCredentialCollection(collection)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials
}

// NewCredentialVerificationResponse creates a new credential verification response
func NewCredentialVerificationResponse() *CredentialVerificationResponseBuilder {
	return newCredentialVerificationResponseBuilder(
		C.self_message_content_credential_verification_response_builder_init(),
	)
}

// ResponseTo sets the request id that is being responded to
func (b *CredentialVerificationResponseBuilder) ResponseTo(requestID []byte) *CredentialVerificationResponseBuilder {
	if len(requestID) != 20 {
		return b
	}

	requestIDBuf := C.CBytes(
		requestID,
	)

	C.self_message_content_credential_verification_response_builder_response_to(
		b.ptr,
		(*C.uint8_t)(requestIDBuf),
	)

	C.free(requestIDBuf)

	return b
}

// ResponseTo sets the request id that is being responded to
func (b *CredentialVerificationResponseBuilder) Status(status ResponseStatus) *CredentialVerificationResponseBuilder {
	C.self_message_content_credential_verification_response_builder_status(
		b.ptr,
		uint32(status),
	)

	return b
}

// VerifiableCredential attaches a verified credential to the response
func (b *CredentialVerificationResponseBuilder) VerifiableCredential(proof *credential.VerifiableCredential) *CredentialVerificationResponseBuilder {
	C.self_message_content_credential_verification_response_builder_verifiable_credential(
		b.ptr,
		verifiableCredentialPtr(proof),
	)
	return b
}

// Finish finalises the response and builds the content
func (b *CredentialVerificationResponseBuilder) Finish() (*Content, error) {
	var finishedContent *C.self_message_content

	result := C.self_message_content_credential_verification_response_builder_finish(
		b.ptr,
		&finishedContent,
	)

	if result > 0 {
		return nil, status.New(result)
	}
	return newContent(finishedContent), nil
}

type CredentialVerificationParameter struct {
	ptr *C.self_credential_verification_parameter
}

func newCredentialVerificationParameter(ptr *C.self_credential_verification_parameter) *CredentialVerificationParameter {
	c := &CredentialVerificationParameter{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialVerificationParameter) {
		C.self_credential_verification_parameter_destroy(
			c.ptr,
		)
	})

	return c
}

// Key returns the parameters key
func (c *CredentialVerificationParameter) Key() string {
	return C.GoString(
		C.self_credential_verification_parameter_parameter_key(
			c.ptr,
		),
	)
}

// Value returns the value of the parameter
func (c *CredentialVerificationParameter) Value() any {
	ptr := C.self_credential_verification_parameter_value(
		c.ptr,
	)

	defer func() {
		C.self_message_content_parameter_value_destroy(ptr)
	}()

	switch C.self_message_content_parameter_value_type_of(
		ptr,
	) {
	case C.PARAMETER_VALUE_BYTES:
		buf := C.self_message_content_parameter_value_as_bytes(
			ptr,
		)

		bytes := C.GoBytes(
			unsafe.Pointer(C.self_bytes_buffer_buf(
				buf,
			)),
			C.int(C.self_bytes_buffer_len(
				buf,
			)),
		)

		C.self_bytes_buffer_destroy(
			buf,
		)

		return bytes
	case C.PARAMETER_VALUE_STRING:
		buf := C.self_message_content_parameter_value_as_string(
			ptr,
		)

		str := C.GoString(
			C.self_string_buffer_ptr(
				buf,
			),
		)

		C.self_string_buffer_destroy(
			buf,
		)

		return str
	case C.PARAMETER_VALUE_INT8:
		return C.self_message_content_parameter_value_as_int8(
			ptr,
		)
	case C.PARAMETER_VALUE_INT16:
		return C.self_message_content_parameter_value_as_int16(
			ptr,
		)
	case C.PARAMETER_VALUE_INT32:
		return C.self_message_content_parameter_value_as_int32(
			ptr,
		)
	case C.PARAMETER_VALUE_INT64:
		return C.self_message_content_parameter_value_as_int64(
			ptr,
		)
	case C.PARAMETER_VALUE_UINT8:
		return C.self_message_content_parameter_value_as_uint8(
			ptr,
		)
	case C.PARAMETER_VALUE_UINT16:
		return C.self_message_content_parameter_value_as_uint16(
			ptr,
		)
	case C.PARAMETER_VALUE_UINT32:
		return C.self_message_content_parameter_value_as_uint32(
			ptr,
		)
	case C.PARAMETER_VALUE_UINT64:
		return C.self_message_content_parameter_value_as_uint64(
			ptr,
		)
	case C.PARAMETER_VALUE_FLOAT32:
		return C.self_message_content_parameter_value_as_float32(
			ptr,
		)
	case C.PARAMETER_VALUE_FLOAT64:
		return C.self_message_content_parameter_value_as_float64(
			ptr,
		)
	case C.PARAMETER_VALUE_ARRAY_BYTES:
		collection := C.self_message_content_parameter_value_as_array_bytes(
			ptr,
		)

		collectionLen := int(C.self_collection_bytes_buffer_len(collection))

		values := make([][]byte, collectionLen)

		for i := 0; i < collectionLen; i++ {
			buf := C.self_collection_bytes_buffer_at(collection, C.ulong(i))

			values[i] = C.GoBytes(
				unsafe.Pointer(C.self_bytes_buffer_buf(
					buf,
				)),
				C.int(C.self_bytes_buffer_len(
					buf,
				)),
			)

			C.self_bytes_buffer_destroy(
				buf,
			)
		}

		C.self_collection_bytes_buffer_destroy(
			collection,
		)

		return values
	case C.PARAMETER_VALUE_ARRAY_STRING:
		collection := C.self_message_content_parameter_value_as_array_string(
			ptr,
		)

		collectionLen := int(C.self_collection_string_buffer_len(collection))

		values := make([]string, collectionLen)

		for i := 0; i < collectionLen; i++ {
			buf := C.self_collection_string_buffer_at(collection, C.ulong(i))

			values[i] = C.GoString(C.self_string_buffer_ptr(
				buf,
			))

			C.self_string_buffer_destroy(
				buf,
			)
		}

		C.self_collection_string_buffer_destroy(
			collection,
		)

		return values
	default:
		return nil
	}
}

type CredentialVerificationEvidence struct {
	ptr *C.self_credential_verification_evidence
}

func newCredentialVerificationEvidence(ptr *C.self_credential_verification_evidence) *CredentialVerificationEvidence {
	c := &CredentialVerificationEvidence{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, func(c *CredentialVerificationEvidence) {
		C.self_credential_verification_evidence_destroy(
			c.ptr,
		)
	})

	return c
}

// EvidenceType returns the evidence type
func (c *CredentialVerificationEvidence) EvidenceType() string {
	return C.GoString(
		C.self_credential_verification_evidence_evidence_type(
			c.ptr,
		),
	)
}

// Object returns the object that makes up the content of the evidence
func (c *CredentialVerificationEvidence) Object() *object.Object {
	return newObject(C.self_credential_verification_evidence_object(
		c.ptr,
	))
}

func fromCredentialVerificationEvidenceCollection(collection *C.self_collection_credential_verification_evidence) []*CredentialVerificationEvidence {
	collectionLen := int(C.self_collection_credential_verification_evidence_len(
		collection,
	))

	details := make([]*CredentialVerificationEvidence, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_credential_verification_evidence_at(
			collection,
			C.size_t(i),
		)

		details[i] = newCredentialVerificationEvidence(ptr)
	}

	return details
}

func fromCredentialVerificationParameterCollection(collection *C.self_collection_credential_verification_parameter) []*CredentialVerificationParameter {
	collectionLen := int(C.self_collection_credential_verification_parameter_len(
		collection,
	))

	details := make([]*CredentialVerificationParameter, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_credential_verification_parameter_at(
			collection,
			C.size_t(i),
		)

		details[i] = newCredentialVerificationParameter(ptr)
	}

	return details
}
