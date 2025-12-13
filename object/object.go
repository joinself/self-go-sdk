package object

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

	"github.com/joinself/self-go-sdk/status"
)

type Object struct {
	ptr *C.self_object
}

func newObject(ptr *C.self_object) *Object {
	o := &Object{
		ptr: ptr,
	}

	runtime.AddCleanup(o, func(o *Object) {
		C.self_object_destroy(
			o.ptr,
		)
	}, o)

	return o
}

func objectPtr(o *Object) *C.self_object {
	return o.ptr
}

// New creates a new intended for sharing with others or storing to the accounts local storage
func New(mime string, data []byte) (*Object, error) {
	var object *C.self_object

	mimeType := C.CString(mime)
	dataBuf := C.CBytes(data)
	dataLen := C.size_t(len(data))

	result := C.self_object_create(
		&object,
		mimeType,
		(*C.uint8_t)(dataBuf),
		dataLen,
	)

	C.free(unsafe.Pointer(mimeType))
	C.free(dataBuf)

	if result > 0 {
		return nil, status.New(result)
	}

	return newObject(object), nil
}

// Id returns the hash of the encrypted data
func (o *Object) Id() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_object_id(
			o.ptr,
		)),
		32,
	)
}

// Hash returns the hash of the unencrypted data
func (o *Object) Hash() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_object_hash(
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
		44,
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

func fromObjectCollection(collection *C.self_collection_object) []*Object {
	collectionLen := int(C.self_collection_object_len(
		collection,
	))

	objects := make([]*Object, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_object_at(
			collection,
			C.size_t(i),
		)

		objects[i] = newObject(ptr)
	}

	return objects
}
