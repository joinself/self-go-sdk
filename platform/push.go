package platform

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

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

type Push struct {
	ptr *C.self_platform_push
}

// NOTE these functions are all mobile specific so aren't exported
func newPlatformPush(ptr *C.self_platform_push) *Push {
	p := &Push{
		ptr: ptr,
	}

	runtime.AddCleanup(p, func(p *Push) {
		C.self_platform_push_destroy(
			p.ptr,
		)
	}, p)

	return p
}

func platformPushPtr(p *Push) *C.self_platform_push {
	return p.ptr
}

func platformPushFCM(applicationAddress *signing.PublicKey, credential string) *Push {
	credentialPtr := C.CString(credential)

	fcm := C.self_platform_push_fcm(
		signingPublicKeyPtr(applicationAddress),
		credentialPtr,
	)

	C.free(unsafe.Pointer(credentialPtr))

	return newPlatformPush(fcm)
}

func platformPushAPNS(applicationAddress *signing.PublicKey, credential string) *Push {
	credentialPtr := C.CString(credential)

	fcm := C.self_platform_push_apns(
		signingPublicKeyPtr(applicationAddress),
		credentialPtr,
	)

	C.free(unsafe.Pointer(credentialPtr))

	return newPlatformPush(fcm)
}
