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
extern void goOnLog(self_log_entry*);
extern void goOnResponse(void*, self_status);

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

void c_on_log(self_log_entry *entry) {
	goOnLog(entry);
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
	callbacks->on_log = c_on_log;

	return callbacks;
}

void c_on_response(void *user_data, self_status response) {
	goOnResponse(user_data, response);
}

void c_self_account_message_send_async(
	struct self_account *account,
    const struct self_signing_public_key *to_address,
    const struct self_message_content *content,
    void *user_data
) {
	self_on_response_cb callback_fn = c_on_response;

	self_account_message_send_async(
		account,
		to_address,
		content,
		&callback_fn,
		user_data
	);
};
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/joinself/self-go-sdk-next/keypair/signing"
	"github.com/joinself/self-go-sdk-next/message"
)

var responseOffset int64
var responseCallbacks sync.Map

func accountCallbacks() *C.self_account_callbacks {
	return C.account_callbacks()
}

//export goOnConnect
func goOnConnect(user_data unsafe.Pointer) {
	(*Account)(user_data).callbacks.OnConnect()
}

//export goOnDisconnect
func goOnDisconnect(user_data unsafe.Pointer, reason C.self_status) {
	// TODO handle reason
	if reason > 0 {
		fmt.Println(reason)
	}
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

	account := (*Account)(user_data)

	if account.callbacks.OnCommit != nil {
		account.callbacks.OnCommit(
			(*Account)(user_data),
			commitEvent,
		)

		return
	}
}

//export goOnKeyPackage
func goOnKeyPackage(user_data unsafe.Pointer, keyPackage *C.cself_key_package_t) {
	keyPackageEvent := (*message.KeyPackage)(keyPackage)

	runtime.SetFinalizer(keyPackageEvent, func(keyPackage *message.KeyPackage) {
		C.self_key_package_destroy(
			(*C.self_key_package)(keyPackage),
		)
	})

	account := (*Account)(user_data)

	if account.callbacks.OnKeyPackage != nil {
		account.callbacks.OnKeyPackage(
			(*Account)(user_data),
			keyPackageEvent,
		)

		return
	}

	_, err := account.ConnectionEstablish(
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

	account := (*Account)(user_data)

	if account.callbacks.OnProposal != nil {
		account.callbacks.OnProposal(
			(*Account)(user_data),
			proposalEvent,
		)

		return
	}
}

//export goOnWelcome
func goOnWelcome(user_data unsafe.Pointer, welcome *C.cself_welcome_t) {
	welcomeEvent := (*message.Welcome)(welcome)

	runtime.SetFinalizer(welcomeEvent, func(welcome *message.Welcome) {
		C.self_welcome_destroy(
			(*C.self_welcome)(welcome),
		)
	})

	account := (*Account)(user_data)

	if account.callbacks.OnWelcome != nil {
		account.callbacks.OnWelcome(
			(*Account)(user_data),
			welcomeEvent,
		)

		return
	}

	_, err := account.ConnectionAccept(
		(*signing.PublicKey)(C.self_welcome_to_address(welcome)),
		welcomeEvent,
	)

	if err != nil {
		panic(err)
	}
}

//export goOnLog
func goOnLog(entry *C.self_log_entry) {
	level := C.self_log_entry_level(entry)
	message := C.GoString(C.self_log_entry_args(entry))

	switch level {
	case C.LOG_ERROR:
		fmt.Printf("[ERROR] %s\n", message)
	case C.LOG_WARN:
		fmt.Printf("[WARN] %s\n", message)
	case C.LOG_INFO:
		fmt.Printf("[INFO] %s\n", message)
	case C.LOG_DEBUG:
		fmt.Printf("[DEBUG] %s\n", message)
	case C.LOG_TRACE:
		fmt.Printf("[TRACE] %s\n", message)
	}

	C.self_log_entry_destroy(entry)
}

//export goOnResponse
func goOnResponse(user_data unsafe.Pointer, response C.self_status) {
	offset := (*int64)(user_data)

	callback, ok := responseCallbacks.LoadAndDelete(*offset)
	if !ok {
		return
	}

	if callback != nil {
		if int(response) > 0 {
			(callback).(func(error))(errors.New("request failed"))
		} else {
			(callback).(func(error))(nil)
		}
	}
}

func accountMessageSendAsync(account *Account, toAddress *signing.PublicKey, content *message.Content, callback func(err error)) {
	offset := atomic.AddInt64(&responseOffset, 1)
	responseCallbacks.Store(offset, callback)

	C.c_self_account_message_send_async(
		account.account,
		(*C.self_signing_public_key)(toAddress),
		(*C.self_message_content)(content),
		unsafe.Pointer(&offset),
	)
}
