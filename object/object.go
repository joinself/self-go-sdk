package object

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
	"unsafe"
)

type Object C.self_object

func Encrypted(mime string, data []byte) (*Object, error) {
	var object *C.self_object
	objectPtr := &object

	mimeType := C.CString(mime)
	dataBuf := C.CBytes(data)
	dataLen := C.ulong(len(data))

	status := C.self_object_create_encrypted(
		objectPtr,
		mimeType,
		(*C.uint8_t)(dataBuf),
		dataLen,
	)

	C.free(unsafe.Pointer(mimeType))
	C.free(dataBuf)

	if status > 0 {
		return nil, errors.New("object creation failed")
	}

	runtime.SetFinalizer(objectPtr, func(object **C.self_object) {
		C.self_object_destroy(
			*object,
		)
	})

	return (*Object)(*objectPtr), nil
}

func Unencrypted(mime string, data []byte) (*Object, error) {
	var object *C.self_object
	objectPtr := &object

	mimeType := C.CString(mime)
	dataBuf := C.CBytes(data)
	dataLen := C.ulong(len(data))

	status := C.self_object_create_unencrypted(
		objectPtr,
		mimeType,
		(*C.uint8_t)(dataBuf),
		dataLen,
	)

	C.free(unsafe.Pointer(mimeType))
	C.free(dataBuf)

	if status > 0 {
		return nil, errors.New("object creation failed")
	}

	runtime.SetFinalizer(object, func(object **C.self_object) {
		C.self_object_destroy(
			*object,
		)
	})

	return (*Object)(*objectPtr), nil
}

// Id returns the id hash of the encrypted data
func (o *Object) Id() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_object_id(
			(*C.self_object)(o),
		)),
		32,
	)
}

// MimeType returns the objects mime type
func (o *Object) MimeType() string {
	return C.GoString(
		C.self_object_mime(
			(*C.self_object)(o),
		),
	)
}

// Key returns the objects encryption key or nil if not present
func (o *Object) Key() []byte {
	key := C.self_object_key(
		(*C.self_object)(o),
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
			(*C.self_object)(o),
		)),
		C.int(C.self_object_data_len(
			(*C.self_object)(o),
		)),
	)
}
