package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

var apins = make(map[*Account]*runtime.Pinner)
var cpins = make(map[*Callbacks]*runtime.Pinner)
var mu sync.Mutex

func pin(pointer *Account) {
	ap := new(runtime.Pinner)
	cp := new(runtime.Pinner)
	ap.Pin(pointer)
	cp.Pin(pointer.callbacks)

	mu.Lock()
	apins[pointer] = ap
	cpins[pointer.callbacks] = cp
	mu.Unlock()
}

func unpin(pointer *Account) {
	mu.Lock()
	apins[pointer].Unpin()
	cpins[pointer.callbacks].Unpin()
	mu.Unlock()
}

type Account struct {
	account   *C.self_account
	callbacks *Callbacks
}

func New(cfg *Config) (*Account, error) {
	account := &Account{
		account:   C.self_account_init(),
		callbacks: &cfg.Callbacks,
	}

	rpcURLBuf := C.CString("http://127.0.0.1:3000/")
	messagingURLBuf := C.CString("ws://127.0.0.1:3001/")
	storagePathBuf := C.CString(cfg.StoragePath)
	storageKeyBuf := (*C.uint8_t)(C.CBytes(cfg.StorageKey))
	storageKeyLen := C.size_t(len(cfg.StorageKey))

	defer func() {
		C.free(unsafe.Pointer(rpcURLBuf))
		C.free(unsafe.Pointer(messagingURLBuf))
		C.free(unsafe.Pointer(storagePathBuf))
		C.free(unsafe.Pointer(storageKeyBuf))
	}()

	// pin our account and callback pointers
	// so we can pass them as user-data to C
	pin(account)

	s := C.self_account_configure(
		account.account,
		rpcURLBuf,
		messagingURLBuf,
		storagePathBuf,
		storageKeyBuf,
		storageKeyLen,
		accountCallbacks(),
		unsafe.Pointer(account),
	)

	if s > 0 {
		return nil, errors.New("configuring account failed")
	}

	return account, nil
}

func (a *Account) InboxOpen() (*PublicKey, error) {
	var address *C.self_signing_public_key

	status := C.self_account_inbox_open(
		a.account,
		&address,
	)

	if status > 0 {
		return nil, errors.New("failed to open inbox")
	}

	return &PublicKey{public: address}, nil
}

func (a *Account) ConnectionNegotiate(asAddress *PublicKey, withAddress *PublicKey) error {
	C.self_account_connection_negotiate(
		a.account,
		asAddress.public,
		withAddress.public,
	)

	return nil
}

func (a *Account) ConnectionEstablish(asAddress *PublicKey, withAddress *PublicKey, keyPackage []byte) error {
	keyPackageBuf := C.CBytes(keyPackage)

	defer func() {
		C.free(keyPackageBuf)
	}()

	C.self_account_connection_establish(
		a.account,
		asAddress.public,
		withAddress.public,
		(*C.uint8_t)(keyPackageBuf),
		C.size_t(len(keyPackage)),
	)

	return nil
}

func (a *Account) ConnectionAccept(asAddress *PublicKey, welcome, notificationToken []byte) error {
	welcomeBuf := C.CBytes(welcome)
	notificationTokenBuf := C.CBytes(notificationToken)

	defer func() {
		C.free(welcomeBuf)
		C.free(notificationTokenBuf)
	}()

	C.self_account_connection_accept(
		a.account,
		asAddress.public,
		(*C.uint8_t)(welcomeBuf),
		C.size_t(len(welcome)),
		(*C.uint8_t)(notificationTokenBuf),
		C.size_t(len(notificationToken)),
	)

	return nil
}

func (a *Account) MessageSend(toAddress *PublicKey, message []byte) error {
	messageBuf := C.CBytes(message)

	defer func() {
		C.free(messageBuf)
	}()

	s := C.self_account_message_send(
		a.account,
		toAddress.public,
		(*C.uint8_t)(messageBuf),
		C.size_t(len(message)),
	)

	if s > 0 {
		return errors.New("message send failed")
	}

	return nil
}
