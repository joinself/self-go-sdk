package credential

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

type Filter struct {
	ptr *C.self_credential_filter
}

func newCredentialFilter(ptr *C.self_credential_filter) *Filter {
	f := &Filter{
		ptr: ptr,
	}

	runtime.SetFinalizer(f, func(f *Filter) {
		C.self_credential_filter_destroy(
			f.ptr,
		)
	})

	return f
}

func credentialFilterPtr(f *Filter) *C.self_credential_filter {
	return f.ptr
}

// NewFilter creates a new filter that can be used to restrict a set of credentials
func NewFilter() *Filter {
	return newCredentialFilter(
		C.self_credential_filter_init(),
	)
}

// Equals checks if a credential field is equal to a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func (f *Filter) Equals(field, value string) *Filter {
	filterField := C.CString(field)
	filterValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(filterField))
		C.free(unsafe.Pointer(filterValue))
	}()

	C.self_credential_filter_equals(
		f.ptr,
		filterField,
		filterValue,
	)

	return f
}

// NotEquals checks if a credential field is not equal to a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func (f *Filter) NotEquals(field, value string) *Filter {
	filterField := C.CString(field)
	filterValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(filterField))
		C.free(unsafe.Pointer(filterValue))
	}()

	C.self_credential_filter_not_equals(
		f.ptr,
		filterField,
		filterValue,
	)

	return f
}

// GreaterThan checks if a credential field is greater than a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func (f *Filter) GreaterThan(field, value string) *Filter {
	filterField := C.CString(field)
	filterValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(filterField))
		C.free(unsafe.Pointer(filterValue))
	}()

	C.self_credential_filter_greater_than(
		f.ptr,
		filterField,
		filterValue,
	)

	return f
}

// GreaterThanOrEquals checks if a credential field is greater than or equal to a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func (f *Filter) GreaterThanOrEquals(field, value string) *Filter {
	filterField := C.CString(field)
	filterValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(filterField))
		C.free(unsafe.Pointer(filterValue))
	}()

	C.self_credential_filter_greater_than_or_equals(
		f.ptr,
		filterField,
		filterValue,
	)

	return f
}

// LessThan checks if a credential field is less than a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func (f *Filter) LessThan(field, value string) *Filter {
	filterField := C.CString(field)
	filterValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(filterField))
		C.free(unsafe.Pointer(filterValue))
	}()

	C.self_credential_filter_less_than(
		f.ptr,
		filterField,
		filterValue,
	)

	return f
}

// LessThanOrEquals checks if a credential field is less than or equal to a given value.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func (f *Filter) LessThanOrEquals(field, value string) *Filter {
	filterField := C.CString(field)
	filterValue := C.CString(value)

	defer func() {
		C.free(unsafe.Pointer(filterField))
		C.free(unsafe.Pointer(filterValue))
	}()

	C.self_credential_filter_less_than_or_equals(
		f.ptr,
		filterField,
		filterValue,
	)

	return f
}

// Empty checks if a credential field is empty.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func (f *Filter) Empty(field string) *Filter {
	filterField := C.CString(field)

	defer func() {
		C.free(unsafe.Pointer(filterField))
	}()

	C.self_credential_filter_empty(
		f.ptr,
		filterField,
	)

	return f
}

// NotEmpty checks if a credential field is not empty.
//
// The `field` parameter should be in the form of an RFC 6901 JSON Pointer
// See: https://datatracker.ietf.org/doc/html/rfc6901
func (f *Filter) NotEmpty(field string) *Filter {
	filterField := C.CString(field)

	defer func() {
		C.free(unsafe.Pointer(filterField))
	}()

	C.self_credential_filter_not_empty(
		f.ptr,
		filterField,
	)

	return f
}

func (f *Filter) Matches(credential *VerifiableCredential) bool {
	return bool(C.self_credential_filter_matches(
		f.ptr,
		credential.ptr,
	))
}

func FilterCredentials(credentials []*VerifiableCredential, filter *Filter) []*VerifiableCredential {
	return nil
}
