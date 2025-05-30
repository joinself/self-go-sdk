package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/crypto"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/identity"
	"github.com/joinself/self-go-sdk/keypair/exchange"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/object"
	"github.com/joinself/self-go-sdk/platform"
	"github.com/joinself/self-go-sdk/status"
	"github.com/joinself/self-go-sdk/token"
)

var pins = make(map[*Account]*runtime.Pinner)
var mu sync.Mutex

func pin(pointer *Account) unsafe.Pointer {
	p := new(runtime.Pinner)
	p.Pin(pointer)
	p.Pin(pointer.callbacks)

	mu.Lock()
	pins[pointer] = p
	mu.Unlock()

	return unsafe.Pointer(pointer)
}

func unpin(pointer *Account) {
	mu.Lock()
	pins[pointer].Unpin()
	mu.Unlock()
}

//go:linkname newSigningPublicKey github.com/joinself/self-go-sdk/keypair/signing.newSigningPublicKey
func newSigningPublicKey(ptr *C.self_signing_public_key) *signing.PublicKey

//go:linkname signingPublicKeyPtr github.com/joinself/self-go-sdk/keypair/signing.signingPublicKeyPtr
func signingPublicKeyPtr(p *signing.PublicKey) *C.self_signing_public_key

//go:linkname newExchangePublicKey github.com/joinself/self-go-sdk/keypair/exchange.newExchangePublicKey
func newExchangePublicKey(ptr *C.self_exchange_public_key) *exchange.PublicKey

//go:linkname exchangePublicKeyPtr github.com/joinself/self-go-sdk/keypair/exchange.exchangePublicKeyPtr
func exchangePublicKeyPtr(p *exchange.PublicKey) *C.self_exchange_public_key

//go:linkname newIdentityDocument github.com/joinself/self-go-sdk/identity.newIdentityDocument
func newIdentityDocument(ptr *C.self_identity_document) *identity.Document

//go:linkname operationPtr github.com/joinself/self-go-sdk/identity.operationPtr
func operationPtr(o *identity.Operation) *C.self_identity_operation

//go:linkname credentialPtr github.com/joinself/self-go-sdk/credential.credentialPtr
func credentialPtr(c *credential.Credential) *C.self_credential

//go:linkname presentationPtr github.com/joinself/self-go-sdk/credential.presentationPtr
func presentationPtr(c *credential.Presentation) *C.self_presentation

//go:linkname newVerifiablePresentation github.com/joinself/self-go-sdk/credential.newVerifiablePresentation
func newVerifiablePresentation(ptr *C.self_verifiable_presentation) *credential.VerifiablePresentation

//go:linkname newVerifiableCredential github.com/joinself/self-go-sdk/credential.newVerifiableCredential
func newVerifiableCredential(ptr *C.self_verifiable_credential) *credential.VerifiableCredential

//go:linkname verifiableCredentialPtr github.com/joinself/self-go-sdk/credential.verifiableCredentialPtr
func verifiableCredentialPtr(v *credential.VerifiableCredential) *C.self_verifiable_credential

//go:linkname verifiablePresentationPtr github.com/joinself/self-go-sdk/credential.verifiablePresentationPtr
func verifiablePresentationPtr(v *credential.VerifiablePresentation) *C.self_verifiable_presentation

//go:linkname toVerifiablePresentationCollection github.com/joinself/self-go-sdk/credential.toVerifiablePresentationCollection
func toVerifiablePresentationCollection(presentations []*credential.VerifiablePresentation) *C.self_collection_verifiable_presentation

//go:linkname trustedIssuerRegistryPtr github.com/joinself/self-go-sdk/credential.trustedIssuerRegistryPtr
func trustedIssuerRegistryPtr(r *credential.TrustedIssuerRegistry) *C.self_trusted_issuer_registry

//go:linkname credentialAddressPtr github.com/joinself/self-go-sdk/credential.credentialAddressPtr
func credentialAddressPtr(a *credential.Address) *C.self_credential_address

//go:linkname newObject github.com/joinself/self-go-sdk/object.newObject
func newObject(ptr *C.self_object) *object.Object

//go:linkname objectPtr github.com/joinself/self-go-sdk/object.objectPtr
func objectPtr(o *object.Object) *C.self_object

//go:linkname newToken github.com/joinself/self-go-sdk/token.newToken
func newToken(ptr *C.self_token) *token.Token

//go:linkname tokenPtr github.com/joinself/self-go-sdk/token.tokenPtr
func tokenPtr(t *token.Token) *C.self_token

//go:linkname platformPushPtr github.com/joinself/self-go-sdk/platform.platformPushPtr
func platformPushPtr(t *platform.Push) *C.self_platform_push

//go:linkname keyPackagePtr github.com/joinself/self-go-sdk/event.keyPackagePtr
func keyPackagePtr(k *event.KeyPackage) *C.self_key_package

//go:linkname welcomePtr github.com/joinself/self-go-sdk/event.welcomePtr
func welcomePtr(w *event.Welcome) *C.self_welcome

//go:linkname newCryptoKeyPackage github.com/joinself/self-go-sdk/crypto.newKeyPackage
func newCryptoKeyPackage(e *C.self_key_package, owned bool) *crypto.KeyPackage

//go:linkname cryptoKeyPackagePtr github.com/joinself/self-go-sdk/event.keyPackagePtr
func cryptoKeyPackagePtr(w *crypto.KeyPackage) *C.self_key_package

//go:linkname newCryptoWelcome github.com/joinself/self-go-sdk/crypto.newWelcome
func newCryptoWelcome(e *C.self_welcome, owned bool) *crypto.Welcome

//go:linkname cryptoWelcomePtr github.com/joinself/self-go-sdk/event.welcomePtr
func cryptoWelcomePtr(w *crypto.Welcome) *C.self_welcome

//go:linkname contentPtr github.com/joinself/self-go-sdk/message.contentPtr
func contentPtr(c *message.Content) *C.self_message_content

//go:linkname contentSummaryPtr github.com/joinself/self-go-sdk/message.contentSummaryPtr
func contentSummaryPtr(c *message.ContentSummary) *C.self_message_content_summary

//go:linkname platformAttestationPtr github.com/joinself/self-go-sdk/platform.platformAttestationPtr
func platformAttestationPtr(a *platform.Attestation) *C.self_platform_attestation

//go:linkname fromSigningPublicKeyCollection github.com/joinself/self-go-sdk/keypair/signing.fromSigningPublicKeyCollection
func fromSigningPublicKeyCollection(ptr *C.self_collection_signing_public_key) []*signing.PublicKey

//go:linkname fromCredentialExchangeCollection github.com/joinself/self-go-sdk/credential.fromCredentialExchangeCollection
func fromCredentialExchangeCollection(ptr *C.self_collection_credential_exchange) []*credential.Exchange

//go:linkname fromVerifiableCredentialCollection github.com/joinself/self-go-sdk/credential.fromVerifiableCredentialCollection
func fromVerifiableCredentialCollection(ptr *C.self_collection_verifiable_credential) []*credential.VerifiableCredential

//go:linkname fromVerifiablePresentationCollection github.com/joinself/self-go-sdk/credential.fromVerifiablePresentationCollection
func fromVerifiablePresentationCollection(ptr *C.self_collection_verifiable_presentation) []*credential.VerifiablePresentation

//go:linkname toCredentialTypeCollection github.com/joinself/self-go-sdk/credential.toCredentialTypeCollection
func toCredentialTypeCollection(c []string) *C.self_collection_credential_type

//go:linkname toPresentationTypeCollection github.com/joinself/self-go-sdk/credential.toPresentationTypeCollection
func toPresentationTypeCollection(c []string) *C.self_collection_presentation_type

// Account a self account
type Account struct {
	account   *C.self_account
	callbacks *Callbacks
	config    *Config
	status    int32
}

// New creates a new self account
func New(cfg *Config) (*Account, error) {
	account := &Account{
		account:   C.self_account_init(),
		callbacks: &cfg.Callbacks,
		config:    cfg,
	}

	cfg.defaults()

	rpcURLBuf := C.CString(cfg.Environment.Rpc)
	objectURLBuf := C.CString(cfg.Environment.Object)
	messagingURLBuf := C.CString(cfg.Environment.Message)
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
	pinnedAccount := pin(account)

	runtime.SetFinalizer(account, func(account *Account) {
		unpin(account)

		C.self_account_destroy(
			account.account,
		)
	})

	result := C.self_account_configure(
		account.account,
		rpcURLBuf,
		objectURLBuf,
		messagingURLBuf,
		storagePathBuf,
		storageKeyBuf,
		storageKeyLen,
		uint32(cfg.LogLevel),
		accountCallbacks(
			cfg.Callbacks.onIntegrity != nil,
		),
		pinnedAccount,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	if !cfg.SkipReady {
		for atomic.LoadInt32(&account.status) == 0 {
			time.Sleep(time.Millisecond)
		}
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
	a.config = cfg

	cfg.defaults()

	rpcURLBuf := C.CString(cfg.Environment.Rpc)
	objectURLBuf := C.CString(cfg.Environment.Object)
	messagingURLBuf := C.CString(cfg.Environment.Message)
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
	pinnedAccount := pin(a)

	result := C.self_account_configure(
		a.account,
		rpcURLBuf,
		objectURLBuf,
		messagingURLBuf,
		storagePathBuf,
		storageKeyBuf,
		storageKeyLen,
		uint32(cfg.LogLevel),
		accountCallbacks(
			cfg.Callbacks.onIntegrity != nil,
		),
		pinnedAccount,
	)

	if result > 0 {
		return status.New(result)
	}

	if !cfg.SkipReady {
		for atomic.LoadInt32(&a.status) == 0 {
			time.Sleep(time.Millisecond)
		}
	}

	return nil
}

// SDKPairingCode returns the sdk pairing code used in linking an sdk instance to
// an application identity. If the sdk has already been linked, or the pairing
// code is not yet available, this will return false
func (a *Account) SDKPairingCode() (string, bool, error) {
	var buffer *C.self_string_buffer

	result := C.self_account_sdk_pairing_code(
		a.account,
		&buffer,
	)

	if result > 0 {
		return "", false, status.New(result)
	}

	if buffer == nil {
		return "", false, nil
	}

	code := C.GoString(
		C.self_string_buffer_ptr(buffer),
	)

	C.self_string_buffer_destroy(buffer)

	return code, true, nil
}

// KeychainSigningCreate creates a new signing keypair
func (a *Account) KeychainSigningCreate() (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	result := C.self_account_keychain_signing_create(
		a.account,
		&address,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningPublicKey(address), nil
}

// KeychainSigningImport imports an existing ed25519 signing keypair
func (a *Account) KeychainSigningImport(seed []byte) (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	seedBuf := C.CBytes(seed)

	result := C.self_account_keychain_signing_import(
		a.account,
		(*C.uint8_t)(seedBuf),
		&address,
	)

	C.free(seedBuf)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningPublicKey(address), nil
}

// KeychainExchangeCreate creates a new exchange keypair
func (a *Account) KeychainExchangeCreate() (*exchange.PublicKey, error) {
	var address *C.self_exchange_public_key

	result := C.self_account_keychain_exchange_create(
		a.account,
		&address,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newExchangePublicKey(address), nil
}

// KeychainSigningAssociatedWith lists all keys associated with a identity that posess the specified set of roles
func (a *Account) KeychainSigningAssociatedWith(address *signing.PublicKey, roles identity.Role) ([]*signing.PublicKey, error) {
	var collection *C.self_collection_signing_public_key

	result := C.self_account_keychain_signing_associated_with(
		a.account,
		signingPublicKeyPtr(address),
		C.uint64_t(roles),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	keys := fromSigningPublicKeyCollection(
		collection,
	)

	C.self_collection_signing_public_key_destroy(
		collection,
	)

	return keys, nil
}

// KeychainSigningAssociatedTo lists all document addresses a key is associated to
func (a *Account) KeychainSigningAssociatedTo(address *signing.PublicKey) ([]*signing.PublicKey, error) {
	var collection *C.self_collection_signing_public_key

	result := C.self_account_keychain_signing_associated_to(
		a.account,
		signingPublicKeyPtr(address),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
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

	result := C.self_account_identity_list(
		a.account,
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
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

	result := C.self_account_identity_resolve(
		a.account,
		signingPublicKeyPtr(address),
		&document,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newIdentityDocument(document), nil
}

// IdentityExecute executes an operation that creates or modifies a document
func (a *Account) IdentityExecute(operation *identity.Operation) error {
	result := C.self_account_identity_execute(
		a.account,
		operationPtr(operation),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// IdentitySign signs an operation that can later be executed
func (a *Account) IdentitySign(operation *identity.Operation) error {
	result := C.self_account_identity_sign(
		a.account,
		operationPtr(operation),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// CredentialIssue signs and issues a verifiable credential
func (a *Account) CredentialIssue(unverifiedCredential *credential.Credential) (*credential.VerifiableCredential, error) {
	var verifiableCredential *C.self_verifiable_credential

	result := C.self_account_credential_issue(
		a.account,
		credentialPtr(unverifiedCredential),
		&verifiableCredential,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newVerifiableCredential(verifiableCredential), nil
}

// CredentialGraphValidFor validates and filters credentials from a given collection of verifiable presentations, validating credentials, presentations and ensuring that an issuer, subject or holder keys have not been revoked and were valid at the time of use.
func (a *Account) CredentialGraphValidFor(address *credential.Address, registry *credential.TrustedIssuerRegistry, presentations []*credential.VerifiablePresentation) ([]*credential.VerifiableCredential, error) {
	var verifiableCredentials *C.self_collection_verifiable_credential

	verifiablePresentations := toVerifiablePresentationCollection(presentations)

	result := C.self_account_credential_graph_valid_for(
		a.account,
		credentialAddressPtr(address),
		trustedIssuerRegistryPtr(registry),
		verifiablePresentations,
		&verifiableCredentials,
	)

	C.self_collection_verifiable_presentation_destroy(
		verifiablePresentations,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromVerifiableCredentialCollection(
		verifiableCredentials,
	)

	C.self_collection_verifiable_credential_destroy(
		verifiableCredentials,
	)

	return credentials, nil
}

// CredentialStore stores a verifiable credential
func (a *Account) CredentialStore(verifiedCredential *credential.VerifiableCredential) error {
	result := C.self_account_credential_store(
		a.account,
		verifiableCredentialPtr(verifiedCredential),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// CredentialLookup looks up all credentials stored to the account
func (a *Account) CredentialLookup() ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	result := C.self_account_credential_lookup(
		a.account,
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialLookupByIssuer looks up credentials issued by a specific issuer
func (a *Account) CredentialLookupByIssuer(issuer *signing.PublicKey) ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	result := C.self_account_credential_lookup_by_issuer(
		a.account,
		signingPublicKeyPtr(issuer),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
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

	result := C.self_account_credential_lookup_by_bearer(
		a.account,
		signingPublicKeyPtr(bearer),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialLookupByCredentialType looks up credentials matching a specific credential type
func (a *Account) CredentialLookupByCredentialType(credentialType []string) ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	credentialTypeCollection := toCredentialTypeCollection(credentialType)

	result := C.self_account_credential_lookup_by_credential_type(
		a.account,
		credentialTypeCollection,
		&collection,
	)

	C.self_collection_credential_type_destroy(
		credentialTypeCollection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialLookupByCredentialHash looks up credentials held by it's hash
func (a *Account) CredentialLookupByCredentialHash(credentialHash []byte) ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	hashPtr := C.CBytes(credentialHash)
	hashLen := len(credentialHash)

	result := C.self_account_credential_lookup_by_credential_hash(
		a.account,
		(*C.uint8_t)(hashPtr),
		C.size_t(hashLen),
		&collection,
	)

	C.free(hashPtr)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialSharedWithAddress returns all credentials shared with a given address of a given credential type
func (a *Account) CredentialSharedWithAddress(withAddress *signing.PublicKey) ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	result := C.self_account_credential_shared_with_address(
		a.account,
		signingPublicKeyPtr(withAddress),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialSharedWithAddress returns all credentials shared with a given address
func (a *Account) CredentialSharedWithAddressByCredentialType(withAddress *signing.PublicKey, credentialType []string) ([]*credential.VerifiableCredential, error) {
	var collection *C.self_collection_verifiable_credential

	credentialTypeCollection := toCredentialTypeCollection(credentialType)

	result := C.self_account_credential_shared_with_address_by_credential_type(
		a.account,
		signingPublicKeyPtr(withAddress),
		credentialTypeCollection,
		&collection,
	)

	C.self_collection_credential_type_destroy(
		credentialTypeCollection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromVerifiableCredentialCollection(
		collection,
	)

	C.self_collection_verifiable_credential_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialExchangeTrack tracks an credential exchange with a given addresss
func (a *Account) CredentialExchangeTrack(withAddress *signing.PublicKey, credential *credential.VerifiableCredential, underLicense *credential.License) error {
	result := C.self_account_credential_exchange_track(
		a.account,
		signingPublicKeyPtr(withAddress),
		verifiableCredentialPtr(credential),
		nil,
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// CredentialExchangeLog returns a log of credentials shared
func (a *Account) CredentialExchangeLog() ([]*credential.Exchange, error) {
	var collection *C.self_collection_credential_exchange

	result := C.self_account_credential_exchange_log(
		a.account,
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromCredentialExchangeCollection(
		collection,
	)

	C.self_collection_credential_exchange_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialExchangeLogWithAddress returns a log of credentials shared with an address
func (a *Account) CredentialExchangeLogWithAddress(withAddress *signing.PublicKey) ([]*credential.Exchange, error) {
	var collection *C.self_collection_credential_exchange

	result := C.self_account_credential_exchange_log_with_address(
		a.account,
		signingPublicKeyPtr(withAddress),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromCredentialExchangeCollection(
		collection,
	)

	C.self_collection_credential_exchange_destroy(
		collection,
	)

	return credentials, nil
}

// CredentialExchangeLogCredential returns a log of every exchange of a given credential
func (a *Account) CredentialExchangeLogCredential(verifiableCredential *credential.VerifiableCredential) ([]*credential.Exchange, error) {
	var collection *C.self_collection_credential_exchange

	result := C.self_account_credential_exchange_log_credential(
		a.account,
		verifiableCredentialPtr(verifiableCredential),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	credentials := fromCredentialExchangeCollection(
		collection,
	)

	C.self_collection_credential_exchange_destroy(
		collection,
	)

	return credentials, nil
}

// PresentationIssue signs and issues a verifiable presentation
func (a *Account) PresentationIssue(presentation *credential.Presentation) (*credential.VerifiablePresentation, error) {
	var verifiablePresentation *C.self_verifiable_presentation

	result := C.self_account_presentation_issue(
		a.account,
		presentationPtr(presentation),
		&verifiablePresentation,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newVerifiablePresentation(verifiablePresentation), nil
}

// PresentationSign signs a verifiable presentation
func (a *Account) PresentationSign(verifiedPresentation *credential.VerifiablePresentation) error {
	result := C.self_account_presentation_sign(
		a.account,
		verifiablePresentationPtr(verifiedPresentation),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// PresentationStore stores a verifiable presentation
func (a *Account) PresentationStore(verifiedPresentation *credential.VerifiablePresentation) error {
	result := C.self_account_presentation_store(
		a.account,
		verifiablePresentationPtr(verifiedPresentation),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// PresentationLookupByHolder looks up presentations inteded for a specific holder
func (a *Account) PresentationLookupByHolder(holder *signing.PublicKey) ([]*credential.VerifiablePresentation, error) {
	var collection *C.self_collection_verifiable_presentation

	result := C.self_account_presentation_lookup_by_holder(
		a.account,
		signingPublicKeyPtr(holder),
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	presentations := fromVerifiablePresentationCollection(
		collection,
	)

	C.self_collection_verifiable_presentation_destroy(
		collection,
	)

	return presentations, nil
}

// PresentationLookupByPresentationType looks up presentations matching a specific presentation type
func (a *Account) PresentationLookupByPresentationType(presentationType []string) ([]*credential.VerifiablePresentation, error) {
	var collection *C.self_collection_verifiable_presentation

	presentationTypeCollection := toPresentationTypeCollection(presentationType)

	result := C.self_account_presentation_lookup_by_presentation_type(
		a.account,
		presentationTypeCollection,
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	presentations := fromVerifiablePresentationCollection(
		collection,
	)

	C.self_collection_presentation_type_destroy(
		presentationTypeCollection,
	)

	C.self_collection_verifiable_presentation_destroy(
		collection,
	)

	return presentations, nil
}

// TokenStore stores a token
func (a *Account) TokenStore(fromAddress, toAddress, forAddress *signing.PublicKey, token *token.Token) error {
	result := C.self_account_token_store(
		a.account,
		signingPublicKeyPtr(fromAddress),
		signingPublicKeyPtr(toAddress),
		signingPublicKeyPtr(forAddress),
		tokenPtr(token),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// InboxOpen opens a new inbox that can be used to send and receive messages
func (a *Account) InboxOpen() (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	result := C.self_account_inbox_open(
		a.account,
		0,
		&address,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningPublicKey(address), nil
}

// InboxOpenWithExpiry opens a new inbox that can be used to send and receive messages that expires after a given time period
func (a *Account) InboxOpenWithExpiry(expires time.Time) (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	result := C.self_account_inbox_open(
		a.account,
		C.int64_t(expires.Unix()),
		&address,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningPublicKey(address), nil
}

// InboxClose closes an existing inbox permanently
func (a *Account) InboxClose(address *signing.PublicKey) error {
	result := C.self_account_inbox_close(
		a.account,
		signingPublicKeyPtr(address),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// InboxList lists all inboxes
func (a *Account) InboxList() ([]*signing.PublicKey, error) {
	var collection *C.self_collection_signing_public_key

	result := C.self_account_inbox_list(
		a.account,
		&collection,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	inboxes := fromSigningPublicKeyCollection(
		collection,
	)

	C.self_collection_signing_public_key_destroy(
		collection,
	)

	return inboxes, nil
}

// GroupWith returns the address of the encrypted group that has been
// negotiated with another address.
// If there is no existing group, this will returnn nil
func (a *Account) GroupWith(withAddress *signing.PublicKey) (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	result := C.self_account_group_with(
		a.account,
		signingPublicKeyPtr(withAddress),
		&address,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningPublicKey(address), nil
}

// GroupMemberAs returns the address used to interact with a given group
// If there is no existing group, this will returnn nil
func (a *Account) GroupMemberAs(groupAddress *signing.PublicKey) (*signing.PublicKey, error) {
	var address *C.self_signing_public_key

	result := C.self_account_group_member_as(
		a.account,
		signingPublicKeyPtr(groupAddress),
		&address,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningPublicKey(address), nil
}

// ValueKeys returns all keys for key value pairs stored on the account
// an optional param can be passed to filter keys with a given prefix
func (a *Account) ValueKeys(prefix ...string) ([]string, error) {
	var collection *C.self_collection_value_key

	var pfx *C.char

	if len(prefix) > 0 {
		pfx = C.CString(prefix[0])
	}

	result := C.self_account_value_keys(
		a.account,
		pfx,
		&collection,
	)

	if pfx != nil {
		C.free(unsafe.Pointer(pfx))
	}

	if result > 0 {
		return nil, status.New(result)
	}

	collectionLen := int(C.self_collection_value_key_len(
		collection,
	))

	keys := make([]string, collectionLen)

	for i := 0; i < collectionLen; i++ {
		ptr := C.self_collection_value_key_at(
			collection,
			C.size_t(i),
		)

		keys[i] = C.GoString(ptr)
	}

	C.self_collection_value_key_destroy(collection)

	return keys, nil
}

// ValueLookup looks up a value by it's key
func (a *Account) ValueLookup(key string) ([]byte, error) {
	var value *C.self_bytes_buffer

	keyPtr := C.CString(key)

	result := C.self_account_value_lookup(
		a.account,
		keyPtr,
		&value,
	)

	C.free(unsafe.Pointer(keyPtr))

	if result > 0 {
		return nil, status.New(result)
	}

	defer C.self_bytes_buffer_destroy(
		value,
	)

	return C.GoBytes(
		unsafe.Pointer(C.self_bytes_buffer_buf(value)),
		C.int(C.self_bytes_buffer_len(value)),
	), nil
}

// ValueStore stores a value to the accounts storage
func (a *Account) ValueStore(key string, value []byte) error {
	keyPtr := C.CString(key)
	valueBuf := C.CBytes(value)
	valueLen := len(value)

	result := C.self_account_value_store(
		a.account,
		keyPtr,
		(*C.uint8_t)(valueBuf),
		C.size_t(valueLen),
	)

	C.free(unsafe.Pointer(keyPtr))
	C.free(valueBuf)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// ValueStoreWithExpiry stores a value to the accounts storage with an expiry
func (a *Account) ValueStoreWithExpiry(key string, value []byte, expires time.Time) error {
	keyPtr := C.CString(key)
	valueBuf := C.CBytes(value)
	valueLen := len(value)

	result := C.self_account_value_store_with_expiry(
		a.account,
		keyPtr,
		(*C.uint8_t)(valueBuf),
		C.size_t(valueLen),
		C.int64_t(expires.Unix()),
	)

	C.free(unsafe.Pointer(keyPtr))
	C.free(valueBuf)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// ValueRemove removes a value by it's key
func (a *Account) ValueRemove(key string) error {
	keyPtr := C.CString(key)

	result := C.self_account_value_remove(
		a.account,
		keyPtr,
	)

	C.free(unsafe.Pointer(keyPtr))

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// ObjectUpload uploads an encrypted object, optionally storing it our to local storage
func (a *Account) ObjectUpload(obj *object.Object, persistLocally bool) error {
	result := C.self_account_object_upload(
		a.account,
		objectPtr(obj),
		C.bool(persistLocally),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// ObjectDownload downloads and decrypts an object
func (a *Account) ObjectDownload(obj *object.Object) error {
	result := C.self_account_object_download(
		a.account,
		objectPtr(obj),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// ObjectStore stores an object to local storage
func (a *Account) ObjectStore(obj *object.Object) error {
	result := C.self_account_object_store(
		a.account,
		objectPtr(obj),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// ObjectRetrieve downloads and decrypts an object
func (a *Account) ObjectRetrieve(hash []byte) (*object.Object, error) {
	var objPtr *C.self_object

	hashPtr := C.CBytes(hash)

	result := C.self_account_object_retrieve(
		a.account,
		(*C.uint8_t)(hashPtr),
		&objPtr,
	)

	C.free(hashPtr)

	if result > 0 {
		return nil, status.New(result)
	}

	return newObject(objPtr), nil
}

// ConnectionNegotiate negotiates a new encrypted group connection with an address. sends a key
// package to the recipient, which they will use to invite us to an encrypted group
func (a *Account) ConnectionNegotiate(asAddress *signing.PublicKey, withAddress *signing.PublicKey, expires time.Time) error {
	result := C.self_account_connection_negotiate(
		a.account,
		signingPublicKeyPtr(asAddress),
		signingPublicKeyPtr(withAddress),
		C.int64_t(expires.Unix()),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// ConnectionNegotiateOutOfBand negotiates a new encrypted group connection with an address. returns a
// key pacakge for use in an out of band message, like an anonymous message encoded to a QR code
func (a *Account) ConnectionNegotiateOutOfBand(asAddress *signing.PublicKey, expires time.Time) (*crypto.KeyPackage, error) {
	var keyPackage *C.self_key_package

	result := C.self_account_connection_negotiate_out_of_band(
		a.account,
		signingPublicKeyPtr(asAddress),
		C.int64_t(expires.Unix()),
		&keyPackage,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newCryptoKeyPackage(keyPackage, true), nil
}

// ConnectionEstablish establishes and sets up an encrypted connection with an address via a new group inbox
// using the key package the initiator sent to us, returns the address of the group
func (a *Account) ConnectionEstablish(asAddress *signing.PublicKey, keyPackage *crypto.KeyPackage) (*signing.PublicKey, error) {
	var groupAddress *C.self_signing_public_key

	result := C.self_account_connection_establish(
		a.account,
		signingPublicKeyPtr(asAddress),
		cryptoKeyPackagePtr(keyPackage),
		&groupAddress,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningPublicKey(groupAddress), nil
}

// ConnectionAccept accepts a welcome to a encrypted group, returns the address of the group
func (a *Account) ConnectionAccept(asAddress *signing.PublicKey, welcome *crypto.Welcome) (*signing.PublicKey, error) {
	var groupAddress *C.self_signing_public_key

	result := C.self_account_connection_accept(
		a.account,
		signingPublicKeyPtr(asAddress),
		cryptoWelcomePtr(welcome),
		&groupAddress,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newSigningPublicKey(groupAddress), nil
}

// MessageSend sends a message to an address that we have established an encrypted group with
// the OnAcknowledgement and OnError callback will be invoked upon receiving the servers response,
// referencing the id of the messages content
func (a *Account) MessageSend(toAddress *signing.PublicKey, content *message.Content) error {
	result := C.self_account_message_send(
		a.account,
		signingPublicKeyPtr(toAddress),
		contentPtr(content),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// NotificationSend sends a push notification
func (a *Account) NotificationSend(toAddress *signing.PublicKey, summary *message.ContentSummary) error {
	result := C.self_account_notification_send(
		a.account,
		signingPublicKeyPtr(toAddress),
		contentSummaryPtr(summary),
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// Close shuts down the account
func (a *Account) Close() error {
	result := C.self_account_destroy(
		a.account,
	)

	if result > 0 {
		return status.New(result)
	}

	return nil
}

// issues a push token. mobile specific so not exported
func tokenIssuePush(a *Account, forAddress *signing.PublicKey, providerAddress *exchange.PublicKey, pushCredential *platform.Push, delegatable bool) (*token.Token, error) {
	var token *C.self_token

	result := C.self_account_token_issue_push(
		a.account,
		signingPublicKeyPtr(forAddress),
		exchangePublicKeyPtr(providerAddress),
		platformPushPtr(pushCredential),
		C.bool(delegatable),
		&token,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newToken(token), nil
}

// registers and requests an identity token. mobile specific so not exported
func sdkRegister(a *Account, forAddress *signing.PublicKey) (*token.Token, error) {
	var token *C.self_token
	var address *C.self_signing_public_key

	if forAddress != nil {
		address = signingPublicKeyPtr(forAddress)
	}

	result := C.self_account_sdk_register(
		a.account,
		address,
		&token,
	)

	if result > 0 {
		return nil, status.New(result)
	}

	return newToken(token), nil
}
