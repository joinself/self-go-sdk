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
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair/signing"
)

type Exchange struct {
	ptr *C.self_credential_exchange
}

func newExchange(ptr *C.self_credential_exchange) *Exchange {
	b := &Exchange{
		ptr: ptr,
	}

	runtime.SetFinalizer(b, func(b *Exchange) {
		C.self_credential_exchange_destroy(
			b.ptr,
		)
	})

	return b
}

func (e *Exchange) WithAddress() *signing.PublicKey {
	return newSigningPublicKey(C.self_credential_exchange_with_address(
		e.ptr,
	))
}

func (e *Exchange) CredentialHash() []byte {
	credentialHashPtr := C.self_credential_exchange_credential_hash(
		e.ptr,
	)

	credentialHash := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(credentialHashPtr)),
		C.int(C.self_bytes_buffer_len(credentialHashPtr)),
	)

	C.self_bytes_buffer_destroy(
		credentialHashPtr,
	)

	return credentialHash
}

func (e *Exchange) UnderLicense() *License {
	return nil
}

func (e *Exchange) SharedAt() time.Time {
	return time.Unix(int64(C.self_credential_exchange_shared_at(
		e.ptr,
	)), 0)
}

func fromCredentialExchangeCollection(collection *C.self_collection_credential_exchange) []*Exchange {
	collectionLen := int(C.self_collection_credential_exchange_len(
		collection,
	))

	exchanges := make([]*Exchange, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_credential_exchange_at(
			collection,
			C.size_t(i),
		)

		exchanges[i] = newExchange(ptr)
	}

	return exchanges
}
