package predicate

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

type PredicatorType int

const (
	PredicatorTypeUnknown             PredicatorType = 1<<63 - 1
	PredicatorTypeEquals              PredicatorType = C.PREDICATOR_EQUALS
	PredicatorTypeNotEquals           PredicatorType = C.PREDICATOR_NOT_EQUALS
	PredicatorTypeGreaterThan         PredicatorType = C.PREDICATOR_GREATER_THAN
	PredicatorTypeGreaterThanOrEquals PredicatorType = C.PREDICATOR_GREATER_THAN_OR_EQUALS
	PredicatorTypeLessThan            PredicatorType = C.PREDICATOR_LESS_THAN
	PredicatorTypeLessThanOrEquals    PredicatorType = C.PREDICATOR_LESS_THAN_OR_EQUALS
	PredicatorTypeContains            PredicatorType = C.PREDICATOR_CONTAINS
	PredicatorTypeNotContains         PredicatorType = C.PREDICATOR_NOT_CONTAINS
	PredicatorTypeOneOf               PredicatorType = C.PREDICATOR_ONE_OF
	PredicatorTypeNotOneOf            PredicatorType = C.PREDICATOR_NOT_ONE_OF
	PredicatorTypeEmpty               PredicatorType = C.PREDICATOR_EMPTY
	PredicatorTypeNotEmpty            PredicatorType = C.PREDICATOR_NOT_EMPTY
)

func (t PredicatorType) String() string {
	switch t {
	case PredicatorTypeEquals:
		return "Equals"
	case PredicatorTypeNotEquals:
		return "NotEquials"
	case PredicatorTypeGreaterThan:
		return "GreaterThan"
	case PredicatorTypeGreaterThanOrEquals:
		return "GreaterThanOrEquals"
	case PredicatorTypeLessThan:
		return "LessThan"
	case PredicatorTypeLessThanOrEquals:
		return "LessThanOrEquals"
	case PredicatorTypeContains:
		return "Contains"
	case PredicatorTypeNotContains:
		return "NotContains"
	case PredicatorTypeOneOf:
		return "OneOf"
	case PredicatorTypeNotOneOf:
		return "NotOneOf"
	case PredicatorTypeEmpty:
		return "Empty"
	case PredicatorTypeNotEmpty:
		return "NotEmpty"
	default:
		return "Unknown"
	}
}

type Predicator struct {
	ptr *C.self_credential_predicator
}

func newCredentialPredicator(ptr *C.self_credential_predicator) *Predicator {
	p := &Predicator{
		ptr: ptr,
	}

	runtime.AddCleanup(p, func(p *Predicator) {
		C.self_credential_predicator_destroy(
			p.ptr,
		)
	}, p)

	return p
}

func (p *Predicator) Type() PredicatorType {
	return PredicatorType(
		C.self_credential_predicator_predicator_type(
			p.ptr,
		),
	)
}

func (p *Predicator) Field() string {
	fieldBuf := C.self_credential_predicator_field(
		p.ptr,
	)

	field := C.GoString(
		C.self_string_buffer_ptr(
			fieldBuf,
		),
	)

	C.self_string_buffer_destroy(fieldBuf)

	return field
}

func (p *Predicator) Values() []string {
	collection := C.self_credential_predicator_values(
		p.ptr,
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
}
