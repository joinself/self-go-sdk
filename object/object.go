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
	var objectPtr *C.self_object

	mimeType := C.CString(mime)
	dataBuf := C.CBytes(data)
	dataLen := C.ulong(len(data))

	status := C.self_object_create_encrypted(
		&objectPtr,
		mimeType,
		(*C.uint8_t)(dataBuf),
		dataLen,
	)

	C.free(unsafe.Pointer(mimeType))
	C.free(dataBuf)

	if status > 0 {
		return nil, errors.New("object creation failed")
	}

	object := (*Object)(objectPtr)

	runtime.SetFinalizer(object, func(object *Object) {
		C.self_object_destroy(
			(*C.self_object)(object),
		)
	})

	return object, nil
}

func Unencrypted(mime string, data []byte) (*Object, error) {
	var objectPtr *C.self_object

	mimeType := C.CString(mime)
	dataBuf := C.CBytes(data)
	dataLen := C.ulong(len(data))

	status := C.self_object_create_unencrypted(
		&objectPtr,
		mimeType,
		(*C.uint8_t)(dataBuf),
		dataLen,
	)

	C.free(unsafe.Pointer(mimeType))
	C.free(dataBuf)

	if status > 0 {
		return nil, errors.New("object creation failed")
	}

	object := (*Object)(objectPtr)

	runtime.SetFinalizer(object, func(object *Object) {
		C.self_object_destroy(
			(*C.self_object)(object),
		)
	})

	return object, nil
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
