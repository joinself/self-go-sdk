package object

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk -Wl,--allow-multiple-definition
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type Object struct {
	ptr *C.self_object
}

func newObject(ptr *C.self_object) *Object {
	o := &Object{
		ptr: ptr,
	}

	runtime.SetFinalizer(o, func(o *Object) {
		C.self_object_destroy(
			o.ptr,
		)
	})

	return o
}

func objectPtr(o *Object) *C.self_object {
	return o.ptr
}

// Encrypted creates a new encrypted object intended for sharing with others from some data
func Encrypted(mime string, data []byte) (*Object, error) {
	var object *C.self_object

	mimeType := C.CString(mime)
	dataBuf := C.CBytes(data)
	dataLen := C.size_t(len(data))

	status := C.self_object_create_encrypted(
		&object,
		mimeType,
		(*C.uint8_t)(dataBuf),
		dataLen,
	)

	C.free(unsafe.Pointer(mimeType))
	C.free(dataBuf)

	if status > 0 {
		return nil, errors.New("object creation failed")
	}

	return newObject(object), nil
}

// Unencrypted creates a new unencrypted object intended to be stored to the accounts local storage
func Unencrypted(mime string, data []byte) (*Object, error) {
	var object *C.self_object

	mimeType := C.CString(mime)
	dataBuf := C.CBytes(data)
	dataLen := C.size_t(len(data))

	status := C.self_object_create_unencrypted(
		&object,
		mimeType,
		(*C.uint8_t)(dataBuf),
		dataLen,
	)

	C.free(unsafe.Pointer(mimeType))
	C.free(dataBuf)

	if status > 0 {
		return nil, errors.New("object creation failed")
	}

	return newObject(object), nil
}

// Id returns the id hash of the encrypted data
func (o *Object) Id() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_object_id(
			o.ptr,
		)),
		32,
	)
}

// MimeType returns the objects mime type
func (o *Object) MimeType() string {
	return C.GoString(
		C.self_object_mime(
			o.ptr,
		),
	)
}

// Key returns the objects encryption key or nil if not present
func (o *Object) Key() []byte {
	key := C.self_object_key(
		o.ptr,
	)

	if key == nil {
		return nil
	}

	return C.GoBytes(
		unsafe.Pointer(key),
		32,
	)
}

// Data returns the objects data
func (o *Object) Data() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_object_data_buf(
			o.ptr,
		)),
		C.int(C.self_object_data_len(
			o.ptr,
		)),
	)
}
