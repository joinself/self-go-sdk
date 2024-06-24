package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>

typedef const self_signing_public_key cself_signing_public_key_t;
typedef const self_message cself_message_t;
typedef const self_commit cself_commit_t;
typedef const self_key_package cself_key_package_t;
typedef const self_proposal cself_proposal_t;
typedef const self_welcome cself_welcome_t;
typedef const uint8_t cuint8_t;

extern void goOnConnect(void*);
extern void goOnDisconnect(void*, self_status);
extern void goOnMessage(void*, cself_message_t*);
extern void goOnCommit(void*, cself_commit_t*);
extern void goOnKeyPackage(void*, cself_key_package_t*);
extern void goOnProposal(void*, cself_proposal_t*);
extern void goOnWelcome(void*, cself_welcome_t*);

void c_on_connect(void* user_data) {
  goOnConnect(user_data);
}

void c_on_disconnect(void* user_data, self_status reason) {
  goOnDisconnect(user_data, reason);
}

void c_on_message(void *user_data, self_message *message) {
  goOnMessage(user_data, message);
}

void c_on_commit(void *user_data, self_commit *commit) {
  goOnCommit(user_data, commit);
}

void c_on_key_package(void *user_data, self_key_package *key_package) {
  goOnKeyPackage(user_data, key_package);
}

void c_on_proposal(void *user_data, self_proposal *proposal) {
  goOnProposal(user_data, proposal);
}

void c_on_welcome(void *user_data, self_welcome *welcome) {
	goOnWelcome(user_data, welcome);
}

self_account_callbacks *account_callbacks() {
	self_account_callbacks *callbacks = malloc(sizeof(self_account_callbacks));

	callbacks->on_connect = c_on_connect;
	callbacks->on_disconnect = c_on_disconnect;
	callbacks->on_message = c_on_message;
	callbacks->on_commit = c_on_commit;
	callbacks->on_key_package = c_on_key_package;
	callbacks->on_proposal = c_on_proposal;
	callbacks->on_welcome = c_on_welcome;

	return callbacks;
}
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
)

func accountCallbacks() *C.self_account_callbacks {
	return C.account_callbacks()
}

//export goOnConnect
func goOnConnect(user_data unsafe.Pointer) {
	(*Account)(user_data).callbacks.OnConnect()
}

//export goOnDisconnect
func goOnDisconnect(user_data unsafe.Pointer, reason C.self_status) {
	(*Account)(user_data).callbacks.OnDisconnect(nil)
}

//export goOnMessage
func goOnMessage(user_data unsafe.Pointer, msg *C.cself_message_t) {
	messageEvent := (*message.Message)(msg)

	runtime.SetFinalizer(messageEvent, func(msg *message.Message) {
		C.self_message_destroy(
			(*C.self_message)(msg),
		)
	})

	(*Account)(user_data).callbacks.OnMessage(
		(*Account)(user_data),
		messageEvent,
	)
}

//export goOnCommit
func goOnCommit(user_data unsafe.Pointer, commit *C.cself_commit_t) {
	commitEvent := (*message.Commit)(commit)

	runtime.SetFinalizer(commitEvent, func(commit *message.Commit) {
		C.self_commit_destroy(
			(*C.self_commit)(commit),
		)
	})
}

//export goOnKeyPackage
func goOnKeyPackage(user_data unsafe.Pointer, keyPackage *C.cself_key_package_t) {
	keyPackageEvent := (*message.KeyPackage)(keyPackage)

	runtime.SetFinalizer(keyPackageEvent, func(keyPackage *message.KeyPackage) {
		C.self_key_package_destroy(
			(*C.self_key_package)(keyPackage),
		)
	})

	err := (*Account)(user_data).ConnectionEstablish(
		(*signing.PublicKey)(C.self_key_package_to_address(keyPackage)),
		keyPackageEvent,
	)

	if err != nil {
		panic(err)
	}
}

//export goOnProposal
func goOnProposal(user_data unsafe.Pointer, proposal *C.cself_proposal_t) {
	proposalEvent := (*message.Proposal)(proposal)

	runtime.SetFinalizer(proposalEvent, func(proposal *message.Proposal) {
		C.self_proposal_destroy(
			(*C.self_proposal)(proposal),
		)
	})
}

//export goOnWelcome
func goOnWelcome(user_data unsafe.Pointer, welcome *C.cself_welcome_t) {
	welcomeEvent := (*message.Welcome)(welcome)

	runtime.SetFinalizer(welcomeEvent, func(welcome *message.Welcome) {
		C.self_welcome_destroy(
			(*C.self_welcome)(welcome),
		)
	})

	err := (*Account)(user_data).ConnectionAccept(
		(*signing.PublicKey)(C.self_welcome_to_address(welcome)),
		welcomeEvent,
	)

	if err != nil {
		panic(err)
	}
}
