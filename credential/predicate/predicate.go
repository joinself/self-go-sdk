package predicate

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

type Predicate struct {
	ptr *C.self_credential_predicate
}

func newCredentialPredicate(ptr *C.self_credential_predicate) *Predicate {
	p := &Predicate{
		ptr: ptr,
	}

	runtime.AddCleanup(p, func(ptr *C.self_credential_predicate) {
		C.self_credential_predicate_destroy(
			ptr,
		)
	}, p.ptr)

	return p
}

// Equals checks if a credential field is equal to a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func Equals(field, value string) *Predicate {
	predicateField := C.CString(field)
	predicateValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.free(unsafe.Pointer(predicateValue))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_equals(

			predicateField,
			predicateValue,
		),
	)
}

// NotEquals checks if a credential field is not equal to a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func NotEquals(field, value string) *Predicate {
	predicateField := C.CString(field)
	predicateValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.free(unsafe.Pointer(predicateValue))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_not_equals(

			predicateField,
			predicateValue,
		),
	)
}

// GreaterThan checks if a credential field is greater than a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func GreaterThan(field, value string) *Predicate {
	predicateField := C.CString(field)
	predicateValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.free(unsafe.Pointer(predicateValue))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_greater_than(
			predicateField,
			predicateValue,
		),
	)
}

// GreaterThanOrEquals checks if a credential field is greater than or equal to a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func GreaterThanOrEquals(field, value string) *Predicate {
	predicateField := C.CString(field)
	predicateValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.free(unsafe.Pointer(predicateValue))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_greater_than_or_equals(
			predicateField,
			predicateValue,
		),
	)
}

// LessThan checks if a credential field is less than a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func LessThan(field, value string) *Predicate {
	predicateField := C.CString(field)
	predicateValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.free(unsafe.Pointer(predicateValue))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_less_than(
			predicateField,
			predicateValue,
		),
	)
}

// LessThanOrEquals checks if a credential field is less than or equal to a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func LessThanOrEquals(field, value string) *Predicate {
	predicateField := C.CString(field)
	predicateValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.free(unsafe.Pointer(predicateValue))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_less_than_or_equals(
			predicateField,
			predicateValue,
		),
	)
}

// Contains checks if a credential field contains a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func Contains(field, value string) *Predicate {
	predicateField := C.CString(field)
	predicateValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.free(unsafe.Pointer(predicateValue))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_contains(
			predicateField,
			predicateValue,
		),
	)
}

// NotContains checks if a credential field does not contain a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func NotContains(field, value string) *Predicate {
	predicateField := C.CString(field)
	predicateValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.free(unsafe.Pointer(predicateValue))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_not_contains(
			predicateField,
			predicateValue,
		),
	)
}

// OneOf checks if a credential field is one of a given set of values.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func OneOf(field string, values []string) *Predicate {
	predicateField := C.CString(field)

	collection := C.self_collection_string_buffer_init()

	for i := 0; i < len(values); i++ {
		valuePtr := C.CString(values[i])

		C.self_collection_string_buffer_append(
			collection,
			valuePtr,
		)

		C.free(unsafe.Pointer(valuePtr))
	}

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.self_collection_string_buffer_destroy(collection)
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_one_of(
			predicateField,
			collection,
		),
	)
}

// NotOneOfchecks if a credential field is not one of a given set of values.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func NotOneOf(field string, values []string) *Predicate {
	predicateField := C.CString(field)

	collection := C.self_collection_string_buffer_init()

	for i := 0; i < len(values); i++ {
		valuePtr := C.CString(values[i])

		C.self_collection_string_buffer_append(
			collection,
			valuePtr,
		)

		C.free(unsafe.Pointer(valuePtr))
	}

	defer func() {
		C.free(unsafe.Pointer(predicateField))
		C.self_collection_string_buffer_destroy(collection)
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_not_one_of(
			predicateField,
			collection,
		),
	)
}

// Empty checks if a credential field is empty.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func Empty(field string) *Predicate {
	predicateField := C.CString(field)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_empty(
			predicateField,
		),
	)
}

// NotEmpty checks if a credential field is not empty.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func NotEmpty(field string) *Predicate {
	predicateField := C.CString(field)

	defer func() {
		C.free(unsafe.Pointer(predicateField))
	}()

	return newCredentialPredicate(
		C.self_credential_predicate_not_empty(

			predicateField,
		),
	)
}

// And joins two predicates together, requiring both to be true
func (p *Predicate) And(other *Predicate) *Predicate {
	return newCredentialPredicate(
		C.self_credential_predicate_and(
			p.ptr,
			other.ptr,
		),
	)
}

// Or joins two predicates together, requiring either to be true
func (p *Predicate) Or(other *Predicate) *Predicate {
	return newCredentialPredicate(
		C.self_credential_predicate_or(
			p.ptr,
			other.ptr,
		),
	)
}
