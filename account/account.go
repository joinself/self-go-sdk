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
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/joinself/self-go-sdk/identity"
	"github.com/joinself/self-go-sdk/keypair/exchange"
	"github.com/joinself/self-go-sdk/keypair/signing"
)

var pins = make(map[*Account]*runtime.Pinner)
var mu sync.Mutex

func pin(pointer *Account) {
	p := new(runtime.Pinner)
	p.Pin(pointer)
	p.Pin(pointer.callbacks)

	mu.Lock()
	pins[pointer] = p
	mu.Unlock()
}

func unpin(pointer *Account) {
	mu.Lock()
	pins[pointer].Unpin()
	mu.Unlock()
}

// Account a self account
type Account struct {
	account   *C.self_account
	callbacks *Callbacks
}

// New creates a new self account
func New(cfg *Config) (*Account, error) {
	account := &Account{
		account:   C.self_account_init(),
		callbacks: &cfg.Callbacks,
	}

	// rpcURLBuf := C.CString("http://127.0.0.1:3000/")
	// messagingURLBuf := C.CString("ws://127.0.0.1:3001/")
	rpcURLBuf := C.CString("http://127.0.0.1:8080/")
	messagingURLBuf := C.CString("ws://127.0.0.1:8088/")
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

	runtime.SetFinalizer(account, func(account *Account) {
		unpin(account)

		C.self_account_destroy(
			account.account,
		)
	})

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

// KeypairSigningCreate creates a new signing keypair
func (a *Account) KeypairSigningCreate() (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	status := C.self_account_keypair_signing_create(
		a.account,
		&address,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}

	return (*signing.PublicKey)(address), nil
}

// KeypairExchangeCreate creates a new exchange keypair
func (a *Account) KeypairExchangeCreate() (*exchange.PublicKey, error) {
	var address *C.self_exchange_public_key

	status := C.self_account_keypair_exchange_create(
		a.account,
		&address,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}

	return (*exchange.PublicKey)(address), nil
}

// KeypairSigningAssociatedWith lists all keys associated with a identity posessing the specified set of roles
func (a *Account) KeypairSigningAssociatedWith(address *signing.PublicKey, roles identity.Role) (*signing.PublicKeyCollection, error) {
	var collection *C.self_collection_signing_public_key

	status := C.self_account_keypair_signing_associated_with(
		a.account,
		(*C.self_signing_public_key)(address),
		C.ulong(roles),
		&collection,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}

	return (*signing.PublicKeyCollection)(collection), nil
}

// IdentityResolve resolves an identity document
func (a *Account) IdentityResolve(address *signing.PublicKey) (*identity.Document, error) {
	var document *C.self_identity_document

	status := C.self_account_identity_resolve(
		a.account,
		(*C.self_signing_public_key)(address),
		&document,
	)

	if status > 0 {
		return nil, errors.New("failed to resolve identity")
	}

	return (*identity.Document)(document), nil
}

// IdentityExecute executes an operation that creates or modifies a document
func (a *Account) IdentityExecute(operation *identity.Operation) error {
	status := C.self_account_identity_execute(
		a.account,
		(*C.self_identity_operation)(operation),
	)

	if status > 0 {
		return errors.New("failed to execute operation")
	}

	return nil
}

// InboxOpen opens a new inbox that can be used to send and receive messages
func (a *Account) InboxOpen() (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	status := C.self_account_inbox_open(
		a.account,
		&address,
	)

	if status > 0 {
		return nil, errors.New("failed to open inbox")
	}

	return (*signing.PublicKey)(address), nil
}

// InboxClose closes an existing inbox permanently
func (a *Account) InboxClose(address *signing.PublicKey) error {
	status := C.self_account_inbox_close(
		a.account,
		(*C.self_signing_public_key)(address),
	)

	if status > 0 {
		return errors.New("failed to close inbox")
	}

	return nil
}

// ConnectionNegotiate negotiates a new encrypted group connection with an address. sends a key
// package to the recipient, which they will use to invite us to an encrypted group
func (a *Account) ConnectionNegotiate(asAddress *signing.PublicKey, withAddress *signing.PublicKey) error {
	status := C.self_account_connection_negotiate(
		a.account,
		(*C.self_signing_public_key)(asAddress),
		(*C.self_signing_public_key)(withAddress),
	)

	if status > 0 {
		return fmt.Errorf("failed negotiate connection, code: %d", status)
	}

	return nil
}

// ConnectionEstablish establishes and sets up an encrypted connection with an address via a new group inbox
// using the key package the initiator sent to us
func (a *Account) ConnectionEstablish(asAddress *signing.PublicKey, withAddress *signing.PublicKey, keyPackage []byte) error {
	keyPackageBuf := C.CBytes(keyPackage)

	defer func() {
		C.free(keyPackageBuf)
	}()

	status := C.self_account_connection_establish(
		a.account,
		(*C.self_signing_public_key)(asAddress),
		(*C.self_signing_public_key)(withAddress),
		(*C.uint8_t)(keyPackageBuf),
		C.size_t(len(keyPackage)),
	)

	if status > 0 {
		return fmt.Errorf("failed establish connection, code: %d", status)
	}

	return nil
}

// ConnectionAccept accepts a welcome to a encrypted group
func (a *Account) ConnectionAccept(asAddress *signing.PublicKey, welcome, notificationToken []byte) error {
	welcomeBuf := C.CBytes(welcome)
	notificationTokenBuf := C.CBytes(notificationToken)

	defer func() {
		C.free(welcomeBuf)
		C.free(notificationTokenBuf)
	}()

	status := C.self_account_connection_accept(
		a.account,
		(*C.self_signing_public_key)(asAddress),
		(*C.uint8_t)(welcomeBuf),
		C.size_t(len(welcome)),
		(*C.uint8_t)(notificationTokenBuf),
		C.size_t(len(notificationToken)),
	)

	if status > 0 {
		return fmt.Errorf("failed accept connection, code: %d", status)
	}

	return nil
}

// MessageSend sends a message to an address that we have established an encrypted group with
func (a *Account) MessageSend(toAddress *signing.PublicKey, message []byte) error {
	messageBuf := C.CBytes(message)

	defer func() {
		C.free(messageBuf)
	}()

	status := C.self_account_message_send(
		a.account,
		(*C.self_signing_public_key)(toAddress),
		(*C.uint8_t)(messageBuf),
		C.size_t(len(message)),
	)

	if status > 0 {
		return fmt.Errorf("failed message send, code: %d", status)
	}

	return nil
}
