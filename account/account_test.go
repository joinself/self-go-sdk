package account_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/identity"
	"github.com/joinself/self-go-sdk/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAccount(t testing.TB) (*account.Account, chan *message.Message) {
	incoming := make(chan *message.Message, 1024)

	cfg := &account.Config{
		StorageKey:  make([]byte, 32),
		StoragePath: t.TempDir() + "self.db",
		Callbacks: account.Callbacks{
			OnConnect: func() {},
			OnDisconnect: func(err error) {
				require.Nil(t, err)
			},
			OnMessage: func(account *account.Account, msg *message.Message) {
				switch message.ContentType(msg) {
				case message.TypeChat:
					content, err := message.DecodeChat(msg)
					require.Nil(t, err)

					fmt.Println(
						"to:", msg.ToAddress(),
						"from:", msg.FromAddress(),
						"message:", string(content.Message()),
					)

					incoming <- msg
				}

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

func wait(t testing.TB, ch chan *message.Message, timeout time.Duration) *message.Message {
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

	contentForBobby := message.NewChat().
		Message("hello").
		Finish()

	// send a message from alice
	err = alice.MessageSend(
		bobbyAddress,
		contentForBobby,
	)

	require.Nil(t, err)

	messageFromAlice := wait(t, bobbyInbox, time.Second)
	assert.Equal(t, aliceAddress.String(), messageFromAlice.FromAddress().String())

	chatMessage, err := message.DecodeChat(messageFromAlice)
	require.Nil(t, err)
	assert.Equal(t, "hello", chatMessage.Message())

	contentForAlice := message.NewChat().
		Message("hi!").
		Finish()

	// send a response from bobby
	err = bobby.MessageSend(
		aliceAddress,
		contentForAlice,
	)

	require.Nil(t, err)

	messageFromBobby := wait(t, aliceInbox, time.Second)
	assert.Equal(t, bobbyAddress.String(), messageFromBobby.FromAddress().String())

	chatMessage, err = message.DecodeChat(messageFromBobby)
	require.Nil(t, err)
	assert.Equal(t, "hi!", chatMessage.Message())

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
		GrantEmbedded(multiroleKey, identity.RoleVerification|identity.RoleAuthentication|identity.RoleMessaging).
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

	now := time.Now()

	// we have removed the authentication role, but this won't be reflected immediately upon
	// querying, as we recently queried for this identity.
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleVerification, now))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleMessaging, now))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleAuthentication, now))
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

	bobbyKeys, err := bobby.KeychainSigningAssociatedWith(
		bobbyIdentifiers.Get(0),
		identity.RoleAssertion,
	)
	require.Nil(t, err)
	require.Equal(t, bobbyKeys.Length(), 1)

	// generate a credential for bobby, issued by alice
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

	// store the credential on bobbys account
	err = bobby.CredentialStore(passportVerifiableCredential)
	require.Nil(t, err)

	// retrieve the credential from bobbys account
	verifiableCredentials, err := bobby.CredentialLookupByCredentialType(credential.CredentialPassport)
	require.Nil(t, err)
	require.Equal(t, 1, verifiableCredentials.Length())

	passportPresentation, err := credential.NewPresentation().
		Presentationtype(credential.PresentationPassport).
		Holder(
			credential.AddressAureWithKey(
				bobbyIdentifiers.Get(0),
				bobbyKeys.Get(0),
			),
		).
		CredentialAdd(verifiableCredentials.Get(0)).
		Finish()

	require.Nil(t, err)

	passportVerifiablePresentation, err := bobby.PresentationIssue(passportPresentation)
	require.Nil(t, err)

	err = passportVerifiablePresentation.Validate()
	require.Nil(t, err)

	credentials := passportVerifiablePresentation.Credentials()
	assert.Equal(t, 1, credentials.Length())
}
