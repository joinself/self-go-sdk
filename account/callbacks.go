package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>

typedef const self_signing_public_key cself_signing_public_key_t;
typedef const self_message cself_message_t;
typedef const self_commit cself_commit_t;
typedef const self_key_package cself_key_package_t;
typedef const self_proposal cself_proposal_t;
typedef const self_reference cself_reference_t;
typedef const self_welcome cself_welcome_t;
typedef const uint8_t cuint8_t;

extern void goOnConnect(void*);
extern void goOnDisconnect(void*, self_status);
extern void goOnAcknowledgement(void*, cself_reference_t*);
extern void goOnError(void*, cself_reference_t*, self_status);
extern void goOnMessage(void*, cself_message_t*);
extern void goOnCommit(void*, cself_commit_t*);
extern void goOnKeyPackage(void*, cself_key_package_t*);
extern void goOnProposal(void*, cself_proposal_t*);
extern void goOnWelcome(void*, cself_welcome_t*);
extern void goOnLog(self_log_entry*);
extern void goOnResponse(void*, self_status);

static void c_on_connect(void* user_data) {
  goOnConnect(user_data);
}

static void c_on_disconnect(void* user_data, self_status reason) {
  goOnDisconnect(user_data, reason);
}

static void c_on_acknowledgement(void* user_data, self_reference *reference) {
  goOnAcknowledgement(user_data, reference);
}

static void c_on_error(void* user_data, self_reference *reference, self_status reason) {
  goOnError(user_data, reference, reason);
}

static void c_on_message(void *user_data, self_message *message) {
  goOnMessage(user_data, message);
}

static void c_on_commit(void *user_data, self_commit *commit) {
  goOnCommit(user_data, commit);
}

static void c_on_key_package(void *user_data, self_key_package *key_package) {
  goOnKeyPackage(user_data, key_package);
}

static void c_on_proposal(void *user_data, self_proposal *proposal) {
  goOnProposal(user_data, proposal);
}

static void c_on_welcome(void *user_data, self_welcome *welcome) {
	goOnWelcome(user_data, welcome);
}

static void c_on_log(self_log_entry *entry) {
	goOnLog(entry);
}

static self_account_callbacks *account_callbacks() {
	self_account_callbacks *callbacks = malloc(sizeof(self_account_callbacks));

	callbacks->on_connect = c_on_connect;
	callbacks->on_disconnect = c_on_disconnect;
	callbacks->on_acknowledgement = c_on_acknowledgement;
	callbacks->on_error = c_on_error;
	callbacks->on_message = c_on_message;
	callbacks->on_commit = c_on_commit;
	callbacks->on_key_package = c_on_key_package;
	callbacks->on_proposal = c_on_proposal;
	callbacks->on_welcome = c_on_welcome;
	callbacks->on_log = c_on_log;
	callbacks->on_integrity = NULL;

	return callbacks;
}

static void c_on_response(void *user_data, self_status response) {
	goOnResponse(user_data, response);
}


//static void c_self_account_message_send_async(
//	struct self_account *account,
//    const struct self_signing_public_key *to_address,
//    const struct self_message_content *content,
//    void *user_data
//) {
//	self_on_response_cb callback_fn = c_on_response;

//	self_account_message_send_async(
//		account,
//		to_address,
//		content,
//		&callback_fn,
//		user_data
//	);
//};
*/
import "C"
import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk-next/event"
	"github.com/joinself/self-go-sdk-next/message"
	"github.com/joinself/self-go-sdk-next/status"
)

var responseOffset int64
var responseCallbacks sync.Map

//go:linkname newContent github.com/joinself/self-go-sdk-next/event.newContent
func newContent(m *C.self_message_content) *event.Content

//go:linkname newMessage github.com/joinself/self-go-sdk-next/event.newMessage
func newMessage(e *C.self_message) *event.Message

//go:linkname newCommit github.com/joinself/self-go-sdk-next/event.newCommit
func newCommit(e *C.self_commit) *event.Commit

//go:linkname newKeyPackage github.com/joinself/self-go-sdk-next/event.newKeyPackage
func newKeyPackage(e *C.self_key_package) *event.KeyPackage

//go:linkname newProposal github.com/joinself/self-go-sdk-next/event.newProposal
func newProposal(e *C.self_proposal) *event.Proposal

//go:linkname newReference github.com/joinself/self-go-sdk-next/event.newReference
func newReference(e *C.self_reference) *event.Reference

//go:linkname newWelcome github.com/joinself/self-go-sdk-next/event.newWelcome
func newWelcome(e *C.self_welcome) *event.Welcome

func accountCallbacks() *C.self_account_callbacks {
	return C.account_callbacks()
}

//export goOnConnect
func goOnConnect(user_data unsafe.Pointer) {
	account := (*Account)(user_data)

	if atomic.LoadInt32(&account.status) == 0 {
		go func() {
			for {
				inboxes, err := account.InboxList()
				if err != nil {
					fmt.Printf("[WARN] failed to list inboxes: %s\n", err.Error())
					time.Sleep(time.Second * 10)
					continue
				}

				if len(inboxes) > 0 {
					atomic.StoreInt32(&account.status, 1)
					break
				}

				inbox, err := account.InboxOpen()
				if err != nil {
					fmt.Printf("[WARN] failed to create inbox: %s\n", err.Error())
					time.Sleep(time.Second)
					continue
				}

				keyPackage, err := account.ConnectionNegotiateOutOfBand(inbox, time.Now().Add(time.Hour*8760))
				if err != nil {
					fmt.Printf("[WARN] failed to generate pairing key package: %s\n", err.Error())
					time.Sleep(time.Second * 10)
					continue
				}

				content, err := message.NewDiscoveryRequest().
					KeyPackage(keyPackage).
					Expires(time.Now().Add(time.Hour * 8760)).
					Finish()

				if err != nil {
					fmt.Printf("[WARN] failed to generate pairing discovery request: %s\n", err.Error())
					time.Sleep(time.Second * 10)
					continue
				}

				content.ID()

				pairingCode, err := event.NewAnonymousMessage(content).EncodeToString()
				if err != nil {
					fmt.Printf("[WARN] failed to enncode pairing discovery request: %s\n", err.Error())
					time.Sleep(time.Second * 10)
					continue
				}

				account.pairing = &pairingCode

				logger()(LogInfo, fmt.Sprintf("SDK Pairing Code: %s\n", pairingCode))

				atomic.StoreInt32(&account.status, 1)

				break
			}
		}()
	}

	if account.callbacks.OnConnect != nil {
		(*Account)(user_data).callbacks.OnConnect(account)
	}
}

//export goOnDisconnect
func goOnDisconnect(user_data unsafe.Pointer, reason C.self_status) {
	account := (*Account)(user_data)

	var err error
	if reason > 0 {
		err = status.New(uint32(reason))
	}

	if account.callbacks.OnDisconnect != nil {
		account.callbacks.OnDisconnect(account, err)
	}
}

//export goOnAcknowledgement
func goOnAcknowledgement(user_data unsafe.Pointer, reference *C.cself_reference_t) {
	account := (*Account)(user_data)

	if account.callbacks.OnAcknowledgement != nil {
		account.callbacks.OnAcknowledgement(account, newReference(reference))
	}
}

//export goOnError
func goOnError(user_data unsafe.Pointer, reference *C.cself_reference_t, reason C.self_status) {
	account := (*Account)(user_data)

	if account.callbacks.OnAcknowledgement != nil {
		account.callbacks.OnError(account, newReference(reference), fmt.Errorf("delivery failed, status: %d", reason))
	}
}

//export goOnMessage
func goOnMessage(user_data unsafe.Pointer, msg *C.cself_message_t) {
	account := (*Account)(user_data)

	account.callbacks.OnMessage(
		account,
		newMessage(msg),
	)
}

//export goOnCommit
func goOnCommit(user_data unsafe.Pointer, commit *C.cself_commit_t) {
	account := (*Account)(user_data)

	if account.callbacks.OnCommit != nil {
		account.callbacks.OnCommit(
			(*Account)(user_data),
			newCommit(commit),
		)

		return
	}
}

//export goOnKeyPackage
func goOnKeyPackage(user_data unsafe.Pointer, keyPackage *C.cself_key_package_t) {
	account := (*Account)(user_data)

	if account.callbacks.OnKeyPackage != nil {
		account.callbacks.OnKeyPackage(
			account,
			newKeyPackage(keyPackage),
		)

		return
	}
}

//export goOnProposal
func goOnProposal(user_data unsafe.Pointer, proposal *C.cself_proposal_t) {
	account := (*Account)(user_data)

	if account.callbacks.OnProposal != nil {
		account.callbacks.OnProposal(
			account,
			newProposal(proposal),
		)

		return
	}
}

//export goOnWelcome
func goOnWelcome(user_data unsafe.Pointer, welcome *C.cself_welcome_t) {
	account := (*Account)(user_data)

	if account.callbacks.OnWelcome != nil {
		account.callbacks.OnWelcome(
			account,
			newWelcome(welcome),
		)

		return
	}
}

//export goOnLog
func goOnLog(entry *C.self_log_entry) {
	level := C.self_log_entry_level(entry)
	message := C.GoString(C.self_log_entry_args(entry))

	logger()(
		LogLevel(level),
		message,
	)

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
			(callback).(func(error))(status.New(uint32(response)))
		} else {
			(callback).(func(error))(nil)
		}
	}
}

/*
func accountMessageSendAsync(account *Account, toAddress *signing.PublicKey, content *message.Content, callback func(err error)) {
	offset := atomic.AddInt64(&responseOffset, 1)
	responseCallbacks.Store(offset, callback)

	C.c_self_account_message_send_async(
		account.account,
		signingPublicKeyPtr(toAddress),
		contentPtr(content),
		unsafe.Pointer(&offset),
	)
}
*/
