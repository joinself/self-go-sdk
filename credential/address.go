package credential

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"

	"github.com/joinself/self-go-sdk/keypair/signing"
)

type Address = C.self_credential_address

func AddressAure(address *signing.PublicKey) *Address {
	a := (*Address)(C.self_credential_address_aure(
		(*C.self_signing_public_key)(address),
	))

	runtime.SetFinalizer(a, func(a *Address) {
		C.self_credential_address_destroy(
			(*C.self_credential_address)(a),
		)
	})

	return a
}

func AddressAureWithKey(address, key *signing.PublicKey) *Address {
	a := (*Address)(C.self_credential_address_aure_with_key(
		(*C.self_signing_public_key)(address),
		(*C.self_signing_public_key)(key),
	))

	runtime.SetFinalizer(a, func(a *Address) {
		C.self_credential_address_destroy(
			(*C.self_credential_address)(a),
		)
	})

	return a
}

func AddressKey(address *signing.PublicKey) *Address {
	a := (*Address)(C.self_credential_address_key(
		(*C.self_signing_public_key)(address),
	))

	runtime.SetFinalizer(a, func(a *Address) {
		C.self_credential_address_destroy(
			(*C.self_credential_address)(a),
		)
	})

	return a
}
