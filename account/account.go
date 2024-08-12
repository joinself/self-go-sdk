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

	"github.com/joinself/self-go-sdk-next/credential"
	"github.com/joinself/self-go-sdk-next/identity"
	"github.com/joinself/self-go-sdk-next/keypair/exchange"
	"github.com/joinself/self-go-sdk-next/keypair/signing"
	"github.com/joinself/self-go-sdk-next/message"
	"github.com/joinself/self-go-sdk-next/object"
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

	cfg.defaults()

	rpcURLBuf := C.CString("http://rpc.next.sandbox.joinself.com/")
	objectURLBuf := C.CString("http://object.next.sandbox.joinself.com/")
	messagingURLBuf := C.CString("ws://message.next.sandbox.joinself.com/")
	storagePathBuf := C.CString(cfg.StoragePath)
	storageKeyBuf := (*C.uint8_t)(C.CBytes(cfg.StorageKey))
	storageKeyLen := C.size_t(len(cfg.StorageKey))

	defer func() {
		C.free(unsafe.Pointer(rpcURLBuf))
		C.free(unsafe.Pointer(objectURLBuf))
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

	status := C.self_account_configure(
		account.account,
		rpcURLBuf,
		objectURLBuf,
		messagingURLBuf,
		storagePathBuf,
		storageKeyBuf,
		storageKeyLen,
		uint32(cfg.LogLevel),
		accountCallbacks(),
		unsafe.Pointer(account),
	)

	if status > 0 {
		return nil, errors.New("configuring account failed")
	}

	return account, nil
}

// KeychainSigningCreate creates a new signing keypair
func (a *Account) KeychainSigningCreate() (*signing.PublicKey, error) {
	var address *C.self_signing_public_key
	addressPtr := &address

	status := C.self_account_keychain_signing_create(
		a.account,
		addressPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}

	runtime.SetFinalizer(addressPtr, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(*addressPtr), nil
}

// KeychainExchangeCreate creates a new exchange keypair
func (a *Account) KeychainExchangeCreate() (*exchange.PublicKey, error) {
	var address *C.self_exchange_public_key
	addressPtr := &address

	status := C.self_account_keychain_exchange_create(
		a.account,
		addressPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}

	runtime.SetFinalizer(addressPtr, func(address **C.self_exchange_public_key) {
		C.self_exchange_public_key_destroy(
			*address,
		)
	})

	return (*exchange.PublicKey)(*addressPtr), nil
}

// KeychainSigningAssociatedWith lists all keys associated with a identity that posess the specified set of roles
func (a *Account) KeychainSigningAssociatedWith(address *signing.PublicKey, roles identity.Role) (*signing.PublicKeyCollection, error) {
	var collection *C.self_collection_signing_public_key
	collectionPtr := &collection

	status := C.self_account_keychain_signing_associated_with(
		a.account,
		(*C.self_signing_public_key)(address),
		C.ulong(roles),
		collectionPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}

	c := (*signing.PublicKeyCollection)(collection)

	runtime.SetFinalizer(c, func(collection *signing.PublicKeyCollection) {
		C.self_collection_signing_public_key_destroy(
			(*C.self_collection_signing_public_key)(collection),
		)
	})

	return c, nil
}

// IdentityList lists identities associated with or owned by the account
func (a *Account) IdentityList() (*signing.PublicKeyCollection, error) {
	var collection *C.self_collection_signing_public_key
	collectionPtr := &collection

	status := C.self_account_identity_list(
		a.account,
		collectionPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to list identities")
	}

	runtime.SetFinalizer(collectionPtr, func(collection **C.self_collection_signing_public_key) {
		C.self_collection_signing_public_key_destroy(
			*collection,
		)
	})

	return (*signing.PublicKeyCollection)(*collectionPtr), nil
}

// IdentityResolve resolves an identity document
func (a *Account) IdentityResolve(address *signing.PublicKey) (*identity.Document, error) {
	var document *C.self_identity_document
	var documentPtr = &document

	status := C.self_account_identity_resolve(
		a.account,
		(*C.self_signing_public_key)(address),
		documentPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to resolve identity")
	}

	runtime.SetFinalizer(documentPtr, func(document **C.self_identity_document) {
		C.self_identity_document_destroy(
			*document,
		)
	})

	return (*identity.Document)(*documentPtr), nil
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

// CredentialIssue signs and issues a verifiable credential
func (a *Account) CredentialIssue(unverifiedCredential *credential.Credential) (*credential.VerifiableCredential, error) {
	var verifiableCredential *C.self_verifiable_credential
	verifiableCredentialPtr := &verifiableCredential

	status := C.self_account_credential_issue(
		(*C.self_account)(a.account),
		(*C.self_credential)(unverifiedCredential),
		verifiableCredentialPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to issue credential")
	}

	runtime.SetFinalizer(verifiableCredentialPtr, func(verifiableCredential **C.self_verifiable_credential) {
		C.self_verifiable_credential_destroy(
			*verifiableCredential,
		)
	})

	return (*credential.VerifiableCredential)(*verifiableCredentialPtr), nil
}

// CredentialStore stores a verifiable credential
func (a *Account) CredentialStore(verifiedCredential *credential.VerifiableCredential) error {
	status := C.self_account_credential_store(
		(*C.self_account)(a.account),
		(*C.self_verifiable_credential)(verifiedCredential),
	)

	if status > 0 {
		return errors.New("failed to store credential")
	}

	return nil
}

// CredentialLookupByIssuer looks up credentials issued by a specific issuer
func (a *Account) CredentialLookupByIssuer(issuer *signing.PublicKey) (*credential.VerifiableCredentialCollection, error) {
	var collection *C.self_collection_verifiable_credential
	collectionPtr := &collection

	status := C.self_account_credential_lookup_by_issuer(
		(*C.self_account)(a.account),
		(*C.self_signing_public_key)(issuer),
		collectionPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to store credential")
	}

	runtime.SetFinalizer(collectionPtr, func(collection **C.self_collection_verifiable_credential) {
		C.self_collection_verifiable_credential_destroy(
			*collection,
		)
	})

	return (*credential.VerifiableCredentialCollection)(*collectionPtr), nil
}

// CredentialLookupByBearer looks up credentials held by a specific bearer
func (a *Account) CredentialLookupByBearer(bearer *signing.PublicKey) (*credential.VerifiableCredentialCollection, error) {
	var collection *C.self_collection_verifiable_credential
	collectionPtr := &collection

	status := C.self_account_credential_lookup_by_bearer(
		(*C.self_account)(a.account),
		(*C.self_signing_public_key)(bearer),
		collectionPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to store credential")
	}

	runtime.SetFinalizer(collectionPtr, func(collection **C.self_collection_verifiable_credential) {
		C.self_collection_verifiable_credential_destroy(
			*collection,
		)
	})

	return (*credential.VerifiableCredentialCollection)(*collectionPtr), nil
}

// CredentialLookupBySubject looks up credentials matching a specific credential type
func (a *Account) CredentialLookupByCredentialType(credentialType *credential.CredentialTypeCollection) (*credential.VerifiableCredentialCollection, error) {
	var collection *C.self_collection_verifiable_credential
	collectionPtr := &collection

	status := C.self_account_credential_lookup_by_credential_type(
		(*C.self_account)(a.account),
		(*C.self_collection_credential_type)(credentialType),
		collectionPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to store credential")
	}

	runtime.SetFinalizer(collectionPtr, func(collection **C.self_collection_verifiable_credential) {
		C.self_collection_verifiable_credential_destroy(
			*collection,
		)
	})

	return (*credential.VerifiableCredentialCollection)(*collectionPtr), nil
}

// PresentationIssue signs and issues a verifiable presentation
func (a *Account) PresentationIssue(presentation *credential.Presentation) (*credential.VerifiablePresentation, error) {
	var verifiablePresentation *C.self_verifiable_presentation
	verifiablePresentationPtr := &verifiablePresentation

	status := C.self_account_presentation_issue(
		(*C.self_account)(a.account),
		(*C.self_presentation)(presentation),
		verifiablePresentationPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to issue credential")
	}

	runtime.SetFinalizer(verifiablePresentationPtr, func(verifiablePresentation **C.self_verifiable_presentation) {
		C.self_verifiable_presentation_destroy(
			*verifiablePresentation,
		)
	})

	return (*credential.VerifiablePresentation)(*verifiablePresentationPtr), nil
}

// InboxOpen opens a new inbox that can be used to send and receive messages
func (a *Account) InboxOpen() (*signing.PublicKey, error) {
	var address *C.self_signing_public_key
	addressPtr := &address

	status := C.self_account_inbox_open(
		a.account,
		addressPtr,
	)

	if status > 0 {
		return nil, errors.New("failed to open inbox")
	}

	runtime.SetFinalizer(addressPtr, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(*addressPtr), nil
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

// ObjectUpload uploads an encrypted object, optionally storing it our to local storage
func (a *Account) ObjectUpload(asAddress *signing.PublicKey, obj *object.Object, persistLocally bool) error {
	status := C.self_account_object_upload(
		a.account,
		(*C.self_signing_public_key)(asAddress),
		(*C.self_object)(obj),
		C.bool(persistLocally),
	)

	if status > 0 {
		return fmt.Errorf("failed object upload, code: %d", status)
	}

	return nil
}

// ObjectDownload downloads and decrypts an object
func (a *Account) ObjectDownload(asAddress *signing.PublicKey, obj *object.Object) error {
	status := C.self_account_object_download(
		a.account,
		(*C.self_signing_public_key)(asAddress),
		(*C.self_object)(obj),
	)

	if status > 0 {
		return fmt.Errorf("failed object download, code: %d", status)
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
// using the key package the initiator sent to us, returns the address of the group
func (a *Account) ConnectionEstablish(asAddress *signing.PublicKey, keyPackage *message.KeyPackage) (*signing.PublicKey, error) {
	var groupAddress *C.self_signing_public_key
	groupAddressPtr := &groupAddress

	status := C.self_account_connection_establish(
		a.account,
		(*C.self_signing_public_key)(asAddress),
		groupAddressPtr,
		(*C.self_key_package)(keyPackage),
	)

	if status > 0 {
		return nil, fmt.Errorf("failed establish connection, code: %d", status)
	}

	runtime.SetFinalizer(groupAddressPtr, func(address **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*address,
		)
	})

	return (*signing.PublicKey)(*groupAddressPtr), nil
}

// ConnectionAccept accepts a welcome to a encrypted group, returns the address of the group
func (a *Account) ConnectionAccept(asAddress *signing.PublicKey, welcome *message.Welcome) (*signing.PublicKey, error) {
	var groupAddress *C.self_signing_public_key
	groupAddressPtr := &groupAddress

	status := C.self_account_connection_accept(
		a.account,
		(*C.self_signing_public_key)(asAddress),
		groupAddressPtr,
		(*C.self_welcome)(welcome),
	)

	if status > 0 {
		return nil, fmt.Errorf("failed accept connection, code: %d", status)
	}

	runtime.SetFinalizer(groupAddressPtr, func(groupAddress **C.self_signing_public_key) {
		C.self_signing_public_key_destroy(
			*groupAddress,
		)
	})

	return (*signing.PublicKey)(*groupAddressPtr), nil
}

// MessageSend sends a message to an address that we have established an encrypted group with
func (a *Account) MessageSend(toAddress *signing.PublicKey, content *message.Content) error {
	status := C.self_account_message_send(
		a.account,
		(*C.self_signing_public_key)(toAddress),
		(*C.self_message_content)(content),
	)

	if status > 0 {
		return fmt.Errorf("failed message send, code: %d", status)
	}

	return nil
}

// MessageSend sends a message to an address that we have established an encrypted group with
func (a *Account) MessageSendAsync(toAddress *signing.PublicKey, content *message.Content, callback func(err error)) {
	accountMessageSendAsync(
		a,
		toAddress,
		content,
		callback,
	)
}

// Close shuts down the account
func (a *Account) Close() error {
	status := C.self_account_destroy(
		(*C.self_account)(a.account),
	)

	if status > 0 {
		return errors.New("failed to close account")
	}

	return nil
}
