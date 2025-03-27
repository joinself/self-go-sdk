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

	"github.com/joinself/self-go-sdk/keypair/signing"
)

// NOTE : types and functions in this file are mobile specific
// so are not exported publicly
type notification struct {
	ptr *C.self_notification
}

func newNotifiation(ptr *C.self_notification) *notification {
	e := &notification{
		ptr: ptr,
	}

	runtime.SetFinalizer(e, func(e *notification) {
		C.self_notification_destroy(
			e.ptr,
		)
	})

	return e
}

// ID returns the id of the notifications content
func notificationID(n *notification) []byte {
	return C.GoBytes(
		unsafe.Pointer(C.self_notification_id(n.ptr)),
		20,
	)
}

// FromAddress returns the address the event was sent by
func notificationFromAddress(n *notification) *signing.PublicKey {
	return newSigningPublicKey(
		C.self_notification_from_address(n.ptr),
	)
}

// ToAddress returns the address the event was addressed to
func notificationToAddress(n *notification) *signing.PublicKey {
	return newSigningPublicKey(
		C.self_notification_to_address(n.ptr),
	)
}

// Content returns the notifications content
func notificationContentSummary(n *notification) *ContentSummary {
	return newContentSummary(
		C.self_notification_content_summary(n.ptr),
	)
}
