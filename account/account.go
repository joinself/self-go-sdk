package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration
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
	"time"
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

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk-next/keypair/signing.newSigningPublicKey
func newSigningPublicKey(ptr *C.self_signing_public_key) *signing.PublicKey

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk-next/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

//go:linkname newExchangePublicKey github.com/joinself/self-go-sdk-next/keypair/exchange.newExchangePublicKey
func newExchangePublicKey(ptr *C.self_exchange_public_key) *exchange.PublicKey

//go:linkname exchangePublicKeyPtr github.com/joinself/self-go-sdk-next/keypair/exchange.exchangePublicKeyPtr
func exchangePublicKeyPtr(p *exchange.PublicKey) *C.self_exchange_public_key

//go:linkname newIdentityDocument github.com/joinself/self-go-sdk-next/identity.newIdentityDocument
func newIdentityDocument(ptr *C.self_identity_document) *identity.Document

//go:linkname operationPtr github.com/joinself/self-go-sdk-next/identity.operationPtr
func operationPtr(o *identity.Operation) *C.self_identity_operation

//go:linkname credentialPtr github.com/joinself/self-go-sdk-next/credential.credentialPtr
func credentialPtr(c *credential.Credential) *C.self_credential

//go:linkname presentationPtr github.com/joinself/self-go-sdk-next/credential.presentationPtr
func presentationPtr(c *credential.Presentation) *C.self_presentation

//go:linkname newVerfiablePresentation github.com/joinself/self-go-sdk-next/credential.newVerfiablePresentation
func newVerfiablePresentation(ptr *C.self_verifiable_presentation) *credential.VerifiablePresentation

//go:linkname newVerifiableCredential github.com/joinself/self-go-sdk-next/credential.newVerifiableCredential
func newVerifiableCredential(ptr *C.self_verifiable_credential) *credential.VerifiableCredential

//go:linkname verifiableCredentialPtr github.com/joinself/self-go-sdk-next/credential.verifiableCredentialPtr
func verifiableCredentialPtr(v *credential.VerifiableCredential) *C.self_verifiable_credential

//go:linkname objectPtr github.com/joinself/self-go-sdk-next/object.objectPtr
func objectPtr(o *object.Object) *C.self_object

//go:linkname keyPackagePtr github.com/joinself/self-go-sdk-next/message.keyPackagePtr
func keyPackagePtr(k *message.KeyPackage) *C.self_key_package

//go:linkname welcomePtr github.com/joinself/self-go-sdk-next/message.welcomePtr
func welcomePtr(w *message.Welcome) *C.self_welcome

//go:linkname contentPtr github.com/joinself/self-go-sdk-next/message.contentPtr
func contentPtr(c *message.Content) *C.self_message_content

//go:linkname fromSigningPublicKeyCollection github.com/joinself/self-go-sdk-next/keypair/signing.fromSigningPublicKeyCollection
func fromSigningPublicKeyCollection(ptr *C.self_collection_signing_public_key) []*signing.PublicKey

//go:linkname fromVerifiableCredentialCollection github.com/joinself/self-go-sdk-next/credential.fromVerifiableCredentialCollection
func fromVerifiableCredentialCollection(ptr *C.self_collection_verifiable_credential) []*credential.VerifiableCredential

//go:linkname toCredentialTypeCollection github.com/joinself/self-go-sdk-next/credential.toCredentialTypeCollection
func toCredentialTypeCollection(c []string) *C.self_collection_credential_type

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

	rpcURLBuf := C.CString(cfg.rpcURL())
	objectURLBuf := C.CString(cfg.objectURL())
	messagingURLBuf := C.CString(cfg.messageURL())
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

// Init creates a new account, without any configuration
func Init() *Account {
	account := &Account{
		account: C.self_account_init(),
	}

	runtime.SetFinalizer(account, func(account *Account) {
		unpin(account)

		C.self_account_destroy(
			account.account,
		)
	})

	return account
}

// Configure configures an unconfigured account, fails if the account has already been configured
func (a *Account) Configure(cfg *Config) error {
	if a.callbacks != nil {
		return errors.New("account already configured")
	}

	a.callbacks = &cfg.Callbacks

	cfg.defaults()

	rpcURLBuf := C.CString(cfg.rpcURL())
	objectURLBuf := C.CString(cfg.objectURL())
	messagingURLBuf := C.CString(cfg.messageURL())
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
	pin(a)

	status := C.self_account_configure(
		a.account,
		rpcURLBuf,
		objectURLBuf,
		messagingURLBuf,
		storagePathBuf,
		storageKeyBuf,
		storageKeyLen,
		uint32(cfg.LogLevel),
		accountCallbacks(),
		unsafe.Pointer(a),
	)

	if status > 0 {
		return errors.New("configuring account failed")
	}

	return nil
}

// KeychainSigningCreate creates a new signing keypair
func (a *Account) KeychainSigningCreate() (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	status := C.self_account_keychain_signing_create(
		a.account,
		&address,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}

	return newSigningPublicKey(address), nil
}

// KeychainExchangeCreate creates a new exchange keypair
func (a *Account) KeychainExchangeCreate() (*exchange.PublicKey, error) {
	var address *C.self_exchange_public_key

	status := C.self_account_keychain_exchange_create(
		a.account,
		&address,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}
	return newExchangePublicKey(address), nil
}

// KeychainSigningAssociatedWith lists all keys associated with a identity that posess the specified set of roles
func (a *Account) KeychainSigningAssociatedWith(address *signing.PublicKey, roles identity.Role) ([]*signing.PublicKey, error) {
	var collection *C.self_collection_signing_public_key

	status := C.self_account_keychain_signing_associated_with(
		a.account,
		signingPublicKeyPtr(address),
		C.uint64_t(roles),
		&collection,
	)

	if status > 0 {
		return nil, errors.New("failed to create keypair")
	}

	keys := fromSigningPublicKeyCollection(
		collection,
	)

	C.self_collection_signing_public_key_destroy(
		collection,
	)

	return keys, nil
}

// IdentityList lists identities associated with or owned by the account
func (a *Account) IdentityList() ([]*signing.PublicKey, error) {
	var collection *C.self_collection_signing_public_key

	status := C.self_account_identity_list(
		a.account,
		&collection,
	)

	if status > 0 {
		return nil, errors.New("failed to list identities")
	}

	keys := fromSigningPublicKeyCollection(
		collection,
	)

	C.self_collection_signing_public_key_destroy(
		collection,
	)

	return keys, nil
}

// IdentityResolve resolves an identity document
func (a *Account) IdentityResolve(address *signing.PublicKey) (*identity.Document, error) {
	var document *C.self_identity_document

	status := C.self_account_identity_resolve(
		a.account,
		signingPublicKeyPtr(address),
		&document,
	)

	if status > 0 {
		return nil, errors.New("failed to resolve identity")
	}

	return newIdentityDocument(document), nil
}

// IdentityExecute executes an operation that creates or modifies a document
func (a *Account) IdentityExecute(operation *identity.Operation) error {
	status := C.self_account_identity_execute(
		a.account,
		operationPtr(operation),
	)

	if status > 0 {
		return errors.New("failed to execute operation")
	}

	return nil
}

// CredentialIssue signs and issues a verifiable credential
func (a *Account) CredentialIssue(unverifiedCredential *credential.Credential) (*credential.VerifiableCredential, error) {
	var verifiableCredential *C.self_verifiable_credential

	status := C.self_account_credential_issue(
		a.account,
		credentialPtr(unverifiedCredential),
		&verifiableCredential,
	)

	if status > 0 {
		return nil, errors.New("failed to issue credential")
	}

	return newVerifiableCredential(verifiableCredential), nil
}

// CredentialStore stores a verifiable credential
func (a *Account) CredentialStore(verifiedCredential *credential.VerifiableCredential) error {
	status := C.self_account_credential_store(
		a.account,
		verifiableCredentialPtr(verifiedCredential),
	)

	if status > 0 {
		return errors.New("failed to store credential")
	}

	return nil
}

// CredentialLookupByIssuer looks up credentials issued by a specific issuer
func (a *Account) CredentialLookupByIssuer(issuer *signing.PublicKey) ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	status := C.self_account_credential_lookup_by_issuer(
		a.account,
		signingPublicKeyPtr(issuer),
		&collection,
	)

	if status > 0 {
		return nil, errors.New("failed to store credential")
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialLookupByBearer looks up credentials held by a specific bearer
func (a *Account) CredentialLookupByBearer(bearer *signing.PublicKey) ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	status := C.self_account_credential_lookup_by_bearer(
		a.account,
		signingPublicKeyPtr(bearer),
		&collection,
	)

	if status > 0 {
		return nil, errors.New("failed to store credential")
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialLookupBySubject looks up credentials matching a specific credential type
func (a *Account) CredentialLookupByCredentialType(credentialType []string) ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	credentialTypeCollection := toCredentialTypeCollection(credentialType)

	status := C.self_account_credential_lookup_by_credential_type(
		a.account,
		credentialTypeCollection,
		&collection,
	)

	if status > 0 {
		return nil, errors.New("failed to store credential")
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_credential_type_destroy(
		credentialTypeCollection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// PresentationIssue signs and issues a verifiable presentation
func (a *Account) PresentationIssue(presentation *credential.Presentation) (*credential.VerifiablePresentation, error) {
	var verifiablePresentation *C.self_verifiable_presentation

	status := C.self_account_presentation_issue(
		a.account,
		presentationPtr(presentation),
		&verifiablePresentation,
	)

	if status > 0 {
		return nil, errors.New("failed to issue credential")
	}

	return newVerfiablePresentation(verifiablePresentation), nil
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

	return newSigningPublicKey(address), nil
}

// InboxClose closes an existing inbox permanently
func (a *Account) InboxClose(address *signing.PublicKey) error {
	status := C.self_account_inbox_close(
		a.account,
		signingPublicKeyPtr(address),
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
		signingPublicKeyPtr(asAddress),
		objectPtr(obj),
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
		signingPublicKeyPtr(asAddress),
		objectPtr(obj),
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
		signingPublicKeyPtr(asAddress),
		signingPublicKeyPtr(withAddress),
	)

	if status > 0 {
		return fmt.Errorf("failed negotiate connection, code: %d", status)
	}

	return nil
}

// ConnectionNegotiateOutOfBand negotiates a new encrypted group connection with an address. returns a
// key pacakge for use in an out of band message, like an anonymous message encoded to a QR code
func (a *Account) ConnectionNegotiateOutOfBand(asAddress *signing.PublicKey, expires time.Time) (*message.KeyPackage, error) {
	var keyPackage *C.self_key_package

	status := C.self_account_connection_negotiate_out_of_band(
		a.account,
		signingPublicKeyPtr(asAddress),
		C.int64_t(expires.Unix()),
		&keyPackage,
	)

	if status > 0 {
		return nil, fmt.Errorf("failed negotiate connection, code: %d", status)
	}

	return newKeyPackage(keyPackage), nil
}

// ConnectionEstablish establishes and sets up an encrypted connection with an address via a new group inbox
// using the key package the initiator sent to us, returns the address of the group
func (a *Account) ConnectionEstablish(asAddress *signing.PublicKey, keyPackage *message.KeyPackage) (*signing.PublicKey, error) {
	var groupAddress *C.self_signing_public_key

	status := C.self_account_connection_establish(
		a.account,
		signingPublicKeyPtr(asAddress),
		&groupAddress,
		keyPackagePtr(keyPackage),
	)

	if status > 0 {
		return nil, fmt.Errorf("failed establish connection, code: %d", status)
	}

	return newSigningPublicKey(groupAddress), nil
}

// ConnectionAccept accepts a welcome to a encrypted group, returns the address of the group
func (a *Account) ConnectionAccept(asAddress *signing.PublicKey, welcome *message.Welcome) (*signing.PublicKey, error) {
	var groupAddress *C.self_signing_public_key

	status := C.self_account_connection_accept(
		a.account,
		signingPublicKeyPtr(asAddress),
		&groupAddress,
		welcomePtr(welcome),
	)

	if status > 0 {
		return nil, fmt.Errorf("failed accept connection, code: %d", status)
	}

	return newSigningPublicKey(groupAddress), nil
}

// MessageSend sends a message to an address that we have established an encrypted group with
func (a *Account) MessageSend(toAddress *signing.PublicKey, content *message.Content) error {
	status := C.self_account_message_send(
		a.account,
		signingPublicKeyPtr(toAddress),
		contentPtr(content),
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
		a.account,
	)

	if status > 0 {
		return errors.New("failed to close account")
	}

	return nil
}
