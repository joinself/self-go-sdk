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
)

var (
	TermSingleUse *Term = NewTerm(time.Duration(C.CREDENTIAL_TERM_SINGLE_USE) * time.Second)
	TermHour      *Term = NewTerm(time.Duration(C.CREDENTIAL_TERM_HOUR) * time.Second)
	TermDay       *Term = NewTerm(time.Duration(C.CREDENTIAL_TERM_DAY) * time.Second)
	TermWeek      *Term = NewTerm(time.Duration(C.CREDENTIAL_TERM_WEEK) * time.Second)
	TermMonth     *Term = NewTerm(time.Duration(C.CREDENTIAL_TERM_MONTH) * time.Second)
	TermYear      *Term = NewTerm(time.Duration(C.CREDENTIAL_TERM_YEAR) * time.Second)
)

type Term struct {
	ptr *C.self_credential_term
}

func newCredentialTerm(ptr *C.self_credential_term) *Term {
	t := &Term{
		ptr: ptr,
	}

	runtime.AddCleanup(t, func(t *Term) {
		C.self_credential_term_destroy(
			t.ptr,
		)
	}, t)

	return t
}

func NewTerm(duration time.Duration) *Term {
	return newCredentialTerm(
		C.self_credential_term_create(
			C.uint64_t(duration.Seconds()),
		),
	)
}

func credentialTermPtr(t *Term) *C.self_credential_term {
	return t.ptr
}

func (t *Term) Duration() time.Duration {
	return time.Duration(C.self_credential_term_duration(
		t.ptr,
	)) * time.Second
}
