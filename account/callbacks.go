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
typedef const self_dropped_event cself_dropped_event_t;
typedef const self_integrity_request cself_integrity_request_t;
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
extern void goOnDropped(void*, cself_dropped_event_t*);
extern void goOnLog(self_log_entry*);
extern void goOnResponse(void*, self_status);
extern struct self_platform_attestation* goOnIntegrity(void*, cself_integrity_request_t*);
extern struct self_platform_attestation* goOnIntegrityAdhoc(void*, cself_integrity_request_t*);

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

static void c_on_dropped(void *user_data, self_dropped_event *dropped) {
	goOnDropped(user_data, dropped);
}

static void c_on_log(self_log_entry *entry) {
	goOnLog(entry);
}

static struct self_platform_attestation* c_on_integrity(void *user_data, self_integrity_request *integrity) {
	return goOnIntegrity(user_data, integrity);
}

static struct self_platform_attestation* c_on_integrity_adhoc(void *user_data, self_integrity_request *integrity) {
	return goOnIntegrityAdhoc(user_data, integrity);
}

static self_account_callbacks *account_callbacks(bool integrity) {
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
	callbacks->on_dropped = c_on_dropped;

	if (integrity) {
	    callbacks->on_integrity = malloc(sizeof(*c_on_integrity));
        *callbacks->on_integrity = c_on_integrity;
	} else {
		callbacks->on_integrity = NULL;
	}

	return callbacks;
}

static self_account_config *account_config(
	self_account_target target,
	char *rpc_url,
	char *object_url,
	char *message_url,
	char *storage_path,
	uint8_t *storage_key_buf,
	size_t storage_key_len,
	uint32_t log_level
) {
	self_account_config *config = malloc(sizeof(self_account_config));

	config->target = target;
	config->rpc_endpoint = rpc_url;
	config->object_endpoint = object_url;
	config->messaging_endpoint = message_url;
	config->storage_path = storage_path;
	config->encryption_key_buf = storage_key_buf;
	config->encryption_key_len = storage_key_len;
	config->log_level = log_level;
	config->log_callback = c_on_log;

	return config;
}

static void account_config_destroy(self_account_config *config) {
	free(config);
}

static void account_callbacks_destroy(self_account_callbacks *callbacks) {
	if (callbacks == NULL) return;
	if (callbacks->on_integrity) {
		free(callbacks->on_integrity);
	}
	free(callbacks);
}

static void c_on_response(void *user_data, self_status response) {
	goOnResponse(user_data, response);
}

static self_status c_backup_key_restore(
	const char *rpc_endpoint,
	const char *object_endpoint,
	const struct self_verifiable_presentation *presentation,
	struct self_object *backup_image,
	struct self_object *restore_image,
	void *user_data,
	struct self_bytes_buffer **encryption_key,
	struct self_collection_verifiable_credential **credentials
) {
	self_on_integrity_cb callback_fn = c_on_integrity_adhoc;

	return self_backup_key_restore(
		rpc_endpoint,
		object_endpoint,
		presentation,
		backup_image,
		restore_image,
		&callback_fn,
		user_data,
		encryption_key,
		credentials
	);
};

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
	"runtime/cgo"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/platform"
	"github.com/joinself/self-go-sdk/status"
)

var responseOffset int64
var responseCallbacks sync.Map

type pairing struct {
	RequestID   []byte `json:"requestId"`
	PairingCode string `json:"pairingCode"`
}

//go:linkname newContent github.com/joinself/self-go-sdk/event.newContent
func newContent(m *C.self_message_content) *message.Content

//go:linkname newMessage github.com/joinself/self-go-sdk/event.newMessage
func newMessage(e *C.self_message) *event.Message

//go:linkname newCommit github.com/joinself/self-go-sdk/event.newCommit
func newCommit(e *C.self_commit) *event.Commit

//go:linkname newKeyPackage github.com/joinself/self-go-sdk/event.newKeyPackage
func newKeyPackage(e *C.self_key_package) *event.KeyPackage

//go:linkname newProposal github.com/joinself/self-go-sdk/event.newProposal
func newProposal(e *C.self_proposal) *event.Proposal

//go:linkname newReference github.com/joinself/self-go-sdk/event.newReference
func newReference(e *C.self_reference) *event.Reference

//go:linkname newWelcome github.com/joinself/self-go-sdk/event.newWelcome
func newWelcome(e *C.self_welcome) *event.Welcome

//go:linkname newDropped github.com/joinself/self-go-sdk/event.newDropped
func newDropped(e *C.self_dropped_event) *event.Dropped

func accountCallbacks(integrity bool) *C.self_account_callbacks {
	return C.account_callbacks(C.bool(integrity))
}

func accountConfig(
	target C.self_account_target,
	rpcURL *C.char,
	objectURL *C.char,
	messageURL *C.char,
	storagePath *C.char,
	storageKeyBuf *C.uint8_t,
	storageKeyLen C.size_t,
	logLevel C.uint32_t,
) *C.self_account_config {
	return C.account_config(target, rpcURL, objectURL, messageURL, storagePath, storageKeyBuf, storageKeyLen, logLevel)
}

func accountCallbacksDestroy(callbacks *C.self_account_callbacks) {
	C.account_callbacks_destroy(callbacks)
}

func accountConfigDestroy(config *C.self_account_config) {
	C.account_config_destroy(config)
}

//export goOnConnect
func goOnConnect(user_data unsafe.Pointer) {
	account := (*Account)(user_data)

	if atomic.LoadInt32(&account.status) == 0 {
		go func() {
			for {
				if account.config.SkipSetup {
					atomic.StoreInt32(&account.status, 1)
					break
				}

				pairingCode, unpaired, err := account.SDKPairingCode()
				if err != nil {
					time.Sleep(time.Millisecond * 100)
					continue
				}

				if !unpaired {
					atomic.StoreInt32(&account.status, 1)
					break
				}

				d := 13
				s := len(pairingCode) / d

				if account.config.LogLevel >= LogInfo {
					fmt.Printf(
						"%s BEGIN PAIRING CODE %s\n",
						strings.Repeat("=", (len(pairingCode)/d-20)/2),
						strings.Repeat("=", (len(pairingCode)/d-20)/2),
					)

					for i := 0; i < d; i++ {
						fmt.Println(pairingCode[s*i : s*i+s])
					}

					if len(pairingCode)%d > 0 {
						fmt.Println(pairingCode[len(pairingCode)-(len(pairingCode)%d):])
					}

					fmt.Printf(
						"%s END PAIRING CODE %s\n",
						strings.Repeat("=", (len(pairingCode)/d-18)/2),
						strings.Repeat("=", (len(pairingCode)/d-18)/2),
					)
				}

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

//export goOnDropped
func goOnDropped(user_data unsafe.Pointer, dropped *C.cself_dropped_event_t) {
	account := (*Account)(user_data)

	if account.callbacks.OnDropped != nil {
		account.callbacks.OnDropped(
			account,
			newDropped(dropped),
		)

		return
	}
}

//export goOnIntegrity
func goOnIntegrity(user_data unsafe.Pointer, integrity *C.cself_integrity_request_t) *C.self_platform_attestation {
	account := (*Account)(user_data)

	hashBuf := C.self_integrity_request_hash_buf(integrity)
	hashLen := C.self_integrity_request_hash_len(integrity)
	defer C.self_integrity_request_destroy(integrity)

	requestHash := C.GoBytes(
		unsafe.Pointer(hashBuf),
		C.int(hashLen),
	)

	return platformAttestationPtr(account.callbacks.onIntegrity(
		account,
		requestHash,
	))
}

//export goOnIntegrityAdhoc
func goOnIntegrityAdhoc(user_data unsafe.Pointer, integrity *C.cself_integrity_request_t) *C.self_platform_attestation {
	hashBuf := C.self_integrity_request_hash_buf(integrity)
	hashLen := C.self_integrity_request_hash_len(integrity)
	defer C.self_integrity_request_destroy(integrity)

	requestHash := C.GoBytes(
		unsafe.Pointer(hashBuf),
		C.int(hashLen),
	)

	handle := *(*cgo.Handle)(user_data)

	onIntegrity := handle.Value().(func(requestHash []byte) *platform.Attestation)

	return platformAttestationPtr(onIntegrity(
		requestHash,
	))
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

func backupKeyCreate(a *Account, presentation *credential.VerifiablePresentation, encryptionKey []byte) error {
	keyBuf := (*C.uint8_t)(C.CBytes(encryptionKey))
	keyLen := C.size_t(len(encryptionKey))

	result := C.self_account_backup_key_create(
		a.account,
		verifiablePresentationPtr(presentation),
		keyBuf,
		keyLen,
	)

	C.free(unsafe.Pointer(keyBuf))

	if result > 0 {
		return status.New(result)
	}

	return nil
}

func backupKeyRestore(target *Target, presentation *credential.VerifiablePresentation, backupImage, restoreImage *object.Object, onIntegrity func(requestHash []byte) *platform.Attestation) ([]byte, error) {
	var credentials *C.self_collection_verifiable_credential
	var keyBuf *C.self_bytes_buffer

	handle := cgo.NewHandle(onIntegrity)

	rpcURLBuf := C.CString(target.Rpc)
	objectURLBuf := C.CString(target.Object)

	result := C.c_backup_key_restore(
		rpcURLBuf,
		objectURLBuf,
		verifiablePresentationPtr(presentation),
		objectPtr(backupImage),
		objectPtr(restoreImage),
		unsafe.Pointer(&handle),
		&keyBuf,
		&credentials,
	)

	C.free(unsafe.Pointer(rpcURLBuf))
	C.free(unsafe.Pointer(objectURLBuf))
	handle.Delete()

	if result > 0 {
		return nil, status.New(uint32(result))
	}

	key := C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(keyBuf)),
		C.int(C.self_bytes_buffer_len(keyBuf)),
	)

	C.self_bytes_buffer_destroy(keyBuf)
	C.self_collection_verifiable_credential_destroy(credentials)

	return key, nil
}
