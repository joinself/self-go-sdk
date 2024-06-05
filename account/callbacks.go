package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>

typedef const self_signing_public_key cself_signing_public_key_t;
typedef const self_message cself_message_t;
typedef const uint8_t cuint8_t;

extern void goOnConnect(void*);
extern void goOnDisconnect(void*, self_status);
extern void goOnMessage(void*, cself_message_t*);
extern void goOnCommit(void*, cself_signing_public_key_t*, cself_signing_public_key_t*, cuint8_t*, size_t);
extern void goOnKeyPackage(void*, cself_signing_public_key_t*, cself_signing_public_key_t*, cuint8_t*, size_t);
extern void goOnProposal(void*, cself_signing_public_key_t*, cself_signing_public_key_t*, cuint8_t*, size_t);
extern void goOnWelcome(void*, cself_signing_public_key_t*, cself_signing_public_key_t*, cuint8_t*, size_t, cuint8_t*, size_t);

void c_on_connect(void* user_data) {
  goOnConnect(user_data);
}

void c_on_disconnect(void* user_data, self_status reason) {
  goOnDisconnect(user_data, reason);
}

void c_on_message(void *user_data, const self_message *message) {
  goOnMessage(user_data, message);
}

void c_on_commit(void *user_data, const self_signing_public_key *sender, const self_signing_public_key *recipient, const uint8_t *commit_buf, size_t commit_len) {
  goOnCommit(user_data, sender, recipient, commit_buf, commit_len);
}

void c_on_key_package(void *user_data, const self_signing_public_key *sender, const self_signing_public_key *recipient, const uint8_t *key_package_buf, size_t key_package_len) {
  goOnKeyPackage(user_data, sender, recipient, key_package_buf, key_package_len);
}

void c_on_proposal(void *user_data, const self_signing_public_key *sender, const self_signing_public_key *recipient, const uint8_t *proposal_buf, size_t proposal_len) {
  goOnProposal(user_data, sender, recipient, proposal_buf, proposal_len);
}

void c_on_welcome(void *user_data, const self_signing_public_key *sender, const self_signing_public_key *recipient, const uint8_t *welcome_buf, size_t welcome_len, const uint8_t *subscription_token_buf, size_t subscription_token_len) {
	goOnWelcome(user_data, sender, recipient, welcome_buf, welcome_len, subscription_token_buf, subscription_token_len);
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
	"fmt"
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
	(*Account)(user_data).callbacks.OnMessage(
		(*Account)(user_data),
		(*message.Message)(msg),
	)
}

//export goOnCommit
func goOnCommit(user_data unsafe.Pointer, fromAddress *C.cself_signing_public_key_t, toAddress *C.cself_signing_public_key_t, commitBuf *C.cuint8_t, commitLen C.size_t) {
}

//export goOnKeyPackage
func goOnKeyPackage(user_data unsafe.Pointer, fromAddress *C.cself_signing_public_key_t, toAddress *C.cself_signing_public_key_t, keyPackageBuf *C.cuint8_t, keyPackageLen C.size_t) {
	fmt.Println("got key package...")
	err := (*Account)(user_data).ConnectionEstablish(
		(*signing.PublicKey)(toAddress),
		(*signing.PublicKey)(fromAddress),
		C.GoBytes(
			unsafe.Pointer(keyPackageBuf),
			C.int(keyPackageLen),
		),
	)
	if err != nil {
		panic(err)
	}
}

//export goOnProposal
func goOnProposal(user_data unsafe.Pointer, fromAddress *C.cself_signing_public_key_t, toAddress *C.cself_signing_public_key_t, proposalBuf *C.cuint8_t, proposalLen C.size_t) {
}

//export goOnWelcome
func goOnWelcome(user_data unsafe.Pointer, fromAddress *C.cself_signing_public_key_t, toAddress *C.cself_signing_public_key_t, welcomeBuf *C.cuint8_t, welcomeLen C.size_t, notificationTokenBuf *C.cuint8_t, notificationTokenLen C.size_t) {
	fmt.Println("got welcome...")
	err := (*Account)(user_data).ConnectionAccept(
		(*signing.PublicKey)(toAddress),
		C.GoBytes(
			unsafe.Pointer(welcomeBuf),
			C.int(welcomeLen),
		),
		C.GoBytes(
			unsafe.Pointer(notificationTokenBuf),
			C.int(notificationTokenLen),
		),
	)
	if err != nil {
		panic(err)
	}
}
