package account_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAccount(t testing.TB) (*account.Account, chan *account.Message) {
	incoming := make(chan *account.Message, 1024)

	cfg := &account.Config{
		StorageKey:  make([]byte, 32),
		StoragePath: t.TempDir() + "self.db",
		Callbacks: account.Callbacks{
			OnConnect: func() {},
			OnDisconnect: func(err error) {
				require.Nil(t, err)
			},
			OnMessage: func(account *account.Account, message *account.Message) {
				fmt.Println(
					"to:", message.ToAddress(),
					"from:", message.FromAddress(),
					"message:", string(message.Message()),
				)

				incoming <- message
			},
		},
	}

	acc, err := account.New(cfg)
	require.Nil(t, err)

	return acc, incoming
}

func testRegisterIdentity(t testing.TB, account *account.Account) {
	identityKey, err := account.KeychainSigningCreate()
	require.Nil(t, err)
	invocationKey, err := account.KeychainSigningCreate()
	require.Nil(t, err)
	multiroleKey, err := account.KeychainSigningCreate()
	require.Nil(t, err)

	operation := identity.NewOperation().
		Identifier(identityKey).
		Sequence(0).
		Timestamp(time.Now()).
		GrantEmbedded(invocationKey, identity.RoleInvocation).
		GrantEmbedded(multiroleKey, identity.RoleVerification|identity.RoleAuthentication|identity.RoleAssertion|identity.RoleMessaging).
		SignWith(identityKey).
		SignWith(invocationKey).
		SignWith(multiroleKey).
		Finish()

	err = account.IdentityExecute(operation)
	require.Nil(t, err)
}

func wait(t testing.TB, ch chan *account.Message, timeout time.Duration) *account.Message {
	select {
	case <-time.After(timeout):
		require.Nil(t, errors.New("timeout"))
		return nil
	case m := <-ch:
		return m
	}
}

func TestAccountMessaging(t *testing.T) {
	alice, aliceInbox := testAccount(t)
	bobby, bobbyInbox := testAccount(t)

	aliceAddress, err := alice.InboxOpen()
	require.Nil(t, err)

	bobbyAddress, err := bobby.InboxOpen()
	require.Nil(t, err)

	fmt.Println("alice:", aliceAddress)
	fmt.Println("bobby:", bobbyAddress)

	err = alice.ConnectionNegotiate(
		aliceAddress,
		bobbyAddress,
	)

	require.Nil(t, err)

	// wait for negotiation to finish
	time.Sleep(time.Millisecond * 2000)

	fmt.Println("message send")
	// send a message from alice
	err = alice.MessageSend(
		bobbyAddress,
		[]byte("hello"),
	)

	require.Nil(t, err)

	fmt.Println("message send ok")

	message := wait(t, bobbyInbox, time.Second)
	assert.Equal(t, aliceAddress.String(), message.FromAddress().String())
	assert.Equal(t, []byte("hello"), message.Message())

	fmt.Println("message send")
	// send a response from bobby
	err = bobby.MessageSend(
		aliceAddress,
		[]byte("hi!"),
	)

	require.Nil(t, err)

	message = wait(t, aliceInbox, time.Second)
	assert.Equal(t, bobbyAddress.String(), message.FromAddress().String())
	assert.Equal(t, []byte("hi!"), message.Message())

	identityKey, err := alice.KeychainSigningCreate()
	require.Nil(t, err)
	invocationKey, err := alice.KeychainSigningCreate()
	require.Nil(t, err)
	multiroleKey, err := alice.KeychainSigningCreate()
	require.Nil(t, err)

	document := identity.NewDocument()
	operation := document.
		Create().
		Identifier(identityKey).
		GrantEmbedded(invocationKey, identity.RoleInvocation).
		GrantEmbedded(multiroleKey, identity.RoleAuthentication|identity.RoleMessaging).
		SignWith(identityKey).
		SignWith(invocationKey).
		SignWith(multiroleKey).
		Finish()

	err = alice.IdentityExecute(operation)
	require.Nil(t, err)
}

func TestAccountIdentity(t *testing.T) {
	alice, _ := testAccount(t)

	identityKey, err := alice.KeychainSigningCreate()
	require.Nil(t, err)
	invocationKey, err := alice.KeychainSigningCreate()
	require.Nil(t, err)
	multiroleKey, err := alice.KeychainSigningCreate()
	require.Nil(t, err)

	operation := identity.NewOperation().
		Identifier(identityKey).
		Sequence(0).
		Timestamp(time.Now()).
		GrantEmbedded(invocationKey, identity.RoleInvocation).
		GrantEmbedded(multiroleKey, identity.RoleVerification|identity.RoleAuthentication|identity.RoleMessaging).
		SignWith(identityKey).
		SignWith(invocationKey).
		SignWith(multiroleKey).
		Finish()

	err = alice.IdentityExecute(operation)
	require.Nil(t, err)

	keys, err := alice.KeychainSigningAssociatedWith(identityKey, identity.RoleInvocation)
	require.Nil(t, err)
	assert.Equal(t, keys.Length(), 1)

	invocationKey = keys.Get(0)
	require.NotNil(t, invocationKey)

	document, err := alice.IdentityResolve(identityKey)
	require.Nil(t, err)

	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleVerification, time.Now()))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleAuthentication, time.Now()))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleMessaging, time.Now()))

	operation = document.
		Create().
		Timestamp(time.Now().Add(time.Second)).
		Modify(multiroleKey, identity.RoleVerification|identity.RoleMessaging).
		SignWith(invocationKey).
		Finish()

	err = alice.IdentityExecute(operation)
	require.Nil(t, err)

	document, err = alice.IdentityResolve(identityKey)
	require.Nil(t, err)

	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleVerification, time.Now()))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleMessaging, time.Now()))
	assert.False(t, document.HasRolesAt(multiroleKey, identity.RoleAuthentication, time.Now()))
}

func TestAccountCredentials(t *testing.T) {
	alice, _ := testAccount(t)
	bobby, _ := testAccount(t)

	testRegisterIdentity(t, alice)
	testRegisterIdentity(t, bobby)

	aliceIdentifiers, err := alice.IdentityList()
	require.Nil(t, err)
	require.Equal(t, aliceIdentifiers.Length(), 1)

	aliceKeys, err := alice.KeychainSigningAssociatedWith(
		aliceIdentifiers.Get(0),
		identity.RoleAssertion,
	)
	require.Nil(t, err)
	require.Equal(t, aliceKeys.Length(), 1)

	bobbyIdentifiers, err := bobby.IdentityList()
	require.Nil(t, err)
	require.Equal(t, bobbyIdentifiers.Length(), 1)

	passportCredential, err := credential.NewCredential().
		CredentialType(credential.CredentialPassport).
		CredentialSubject(credential.AddressAure(bobbyIdentifiers.Get(0))).
		CredentialSubjectClaim("firstName", "bobby").
		Issuer(credential.AddressAure(aliceIdentifiers.Get(0))).
		ValidFrom(time.Now()).
		SignWith(aliceKeys.Get(0), time.Now()).
		Finish()

	require.Nil(t, err)

	passportVerifiableCredential, err := alice.CredentialIssue(passportCredential)
	require.Nil(t, err)
	assert.Equal(t, credential.CredentialPassport, passportVerifiableCredential.CredentialType())

	firstName, ok := passportVerifiableCredential.CredentialSubjectClaim("firstName")
	require.True(t, ok)
	assert.Equal(t, "bobby", firstName)

	_, ok = passportVerifiableCredential.CredentialSubjectClaim("lastName")
	require.False(t, ok)

	err = passportVerifiableCredential.Validate()
	require.Nil(t, err)

	err = bobby.CredentialStore(passportVerifiableCredential)
	require.Nil(t, err)
}
