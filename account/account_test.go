package account_test

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/identity"
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	account.SetLogFunc(func(level account.LogLevel, message string) {
		// disable logging
	})
}

func testAccount(t testing.TB) (*account.Account, chan *event.Message, chan *event.Welcome) {
	return testAccountWithPath(t, ":memory:")
}

func testAccountWithPath(t testing.TB, path string) (*account.Account, chan *event.Message, chan *event.Welcome) {
	incomingMsg := make(chan *event.Message, 1024)
	incomingWel := make(chan *event.Welcome, 1024)

	if path != ":memory:" {
		path = path + "/self.db"
	}

	signal := make(chan bool, 1)

	cfg := &account.Config{
		StorageKey:  make([]byte, 32),
		StoragePath: path,
		Environment: account.TargetSandbox,
		LogLevel:    account.LogError,
		Callbacks: account.Callbacks{
			OnConnect: func(account *account.Account) {
				signal <- true
			},
			OnDisconnect: func(account *account.Account, err error) {
				// require.Nil(t, err)
			},
			OnAcknowledgement: func(account *account.Account, reference *event.Reference) {
				// fmt.Println("acknowledged", hex.EncodeToString(reference.ID()))
			},
			OnError: func(account *account.Account, reference *event.Reference, err error) {
				// fmt.Println("errored", hex.EncodeToString(reference.ID()), err)
			},
			OnMessage: func(account *account.Account, msg *event.Message) {
				switch event.ContentType(msg) {
				case event.TypeChat:
					incomingMsg <- msg
				}
			},
			OnKeyPackage: func(account *account.Account, keyPackage *event.KeyPackage) {
				_, err := account.ConnectionEstablish(
					keyPackage.ToAddress(),
					keyPackage,
				)
				if err != nil {
					panic(err)
				}
			},
			OnWelcome: func(account *account.Account, welcome *event.Welcome) {
				_, err := account.ConnectionAccept(
					welcome.ToAddress(),
					welcome,
				)
				if err != nil {
					panic(err)
				}

				incomingWel <- welcome
			},
		},
	}

	acc, err := account.New(cfg)
	require.Nil(t, err)
	<-signal

	return acc, incomingMsg, incomingWel
}

func testRegisterIdentity(t testing.TB, account *account.Account) {
	identityKey, err := account.KeychainSigningCreate()
	require.Nil(t, err)
	invocationKey, err := account.KeychainSigningCreate()
	require.Nil(t, err)
	multiroleKey, err := account.KeychainSigningCreate()
	require.Nil(t, err)

	//exchangeKey, err := account.KeychainExchangeCreate()
	//require.Nil(t, err)

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

func wait(t testing.TB, ch chan *event.Message, timeout time.Duration) *event.Message {
	select {
	case <-time.After(timeout):
		require.Nil(t, errors.New("timeout"))
		return nil
	case m := <-ch:
		return m
	}
}

func TestAccountMessaging(t *testing.T) {
	alice, aliceInbox, aliceWel := testAccount(t)
	bobby, bobbyInbox, _ := testAccount(t)

	aliceAddress, err := alice.InboxOpen()
	require.Nil(t, err)

	bobbyAddress, err := bobby.InboxOpen()
	require.Nil(t, err)

	// fmt.Println("alice:", aliceAddress)
	// fmt.Println("bobby:", bobbyAddress)

	err = alice.ConnectionNegotiate(
		aliceAddress,
		bobbyAddress,
		time.Now().Add(time.Hour),
	)

	require.Nil(t, err)

	// wait for negotiation to finish
	<-aliceWel

	contentForBobby, err := message.NewChat().
		Message("hello").
		Finish()

	require.Nil(t, err)

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

	contentForAlice, err := message.NewChat().
		Message("hi!").
		Finish()

	require.Nil(t, err)

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

	contentForAlice, err = message.NewChat().
		Message("hello again!").
		Finish()

	require.Nil(t, err)

	start := time.Now()
	err = bobby.MessageSend(
		aliceAddress,
		contentForAlice,
	)

	require.Nil(t, err)

	messageFromBobby = wait(t, aliceInbox, time.Second)
	assert.Equal(t, bobbyAddress.String(), messageFromBobby.FromAddress().String())

	fmt.Println("sent and received in", time.Since(start))

	chatMessage, err = message.DecodeChat(messageFromBobby)
	require.Nil(t, err)
	assert.Equal(t, "hello again!", chatMessage.Message())

	aliceGroupWith, err := alice.GroupWith(bobbyAddress)
	require.Nil(t, err)

	bobbyGroupWith, err := bobby.GroupWith(aliceAddress)
	require.Nil(t, err)

	assert.True(t, aliceGroupWith.Matches(bobbyGroupWith))

	aliceMemberAs, err := alice.GroupMemberAs(aliceGroupWith)
	require.Nil(t, err)
	assert.True(t, aliceAddress.Matches(aliceMemberAs))

}

func TestAccountIdentity(t *testing.T) {
	alice, _, _ := testAccount(t)

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
	assert.Equal(t, len(keys), 1)

	invocationKey = keys[0]
	require.NotNil(t, invocationKey)

	document, err := alice.IdentityResolve(identityKey)
	require.Nil(t, err)

	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleVerification, time.Now()))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleAuthentication, time.Now()))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleMessaging, time.Now()))

	now := time.Now()

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

	// we have removed the authentication role, but this won't be reflected immediately upon
	// querying, as we recently queried for this identity.
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleVerification, now))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleMessaging, now))
	assert.True(t, document.HasRolesAt(multiroleKey, identity.RoleAuthentication, now))
}

func TestAccountObject(t *testing.T) {
	alice, _, _ := testAccount(t)

	asAddress, err := alice.KeychainSigningCreate()
	require.Nil(t, err)

	data := make([]byte, 1024)
	rand.Read(data)

	encryptedObject, err := object.New(
		"application/octet-stream",
		data,
	)

	require.Nil(t, err)

	err = alice.ObjectUpload(
		asAddress,
		encryptedObject,
		false,
	)

	require.Nil(t, err)

	err = alice.ObjectDownload(
		asAddress,
		encryptedObject,
	)

	require.Nil(t, err)
	assert.Equal(t, data, encryptedObject.Data())
}

func TestAccountCredentials(t *testing.T) {
	alice, _, _ := testAccount(t)
	bobby, _, _ := testAccount(t)

	testRegisterIdentity(t, alice)
	testRegisterIdentity(t, bobby)

	aliceIdentifiers, err := alice.IdentityList()
	require.Nil(t, err)
	require.Len(t, aliceIdentifiers, 1)

	aliceKeys, err := alice.KeychainSigningAssociatedWith(
		aliceIdentifiers[0],
		identity.RoleAssertion,
	)
	require.Nil(t, err)
	require.Len(t, aliceKeys, 1)

	bobbyIdentifiers, err := bobby.IdentityList()
	require.Nil(t, err)
	require.Len(t, bobbyIdentifiers, 1)

	bobbyKeys, err := bobby.KeychainSigningAssociatedWith(
		bobbyIdentifiers[0],
		identity.RoleAssertion,
	)
	require.Nil(t, err)
	require.Len(t, bobbyKeys, 1)

	// generate a credential for bobby, issued by alice
	passportCredential, err := credential.NewCredential().
		CredentialType(credential.CredentialTypePassport).
		CredentialSubject(credential.AddressAure(bobbyIdentifiers[0])).
		CredentialSubjectClaim("firstName", "bobby").
		Issuer(credential.AddressAure(aliceIdentifiers[0])).
		ValidFrom(time.Now()).
		SignWith(aliceKeys[0], time.Now()).
		Finish()

	require.Nil(t, err)

	passportVerifiableCredential, err := alice.CredentialIssue(passportCredential)
	require.Nil(t, err)
	assert.Equal(t, credential.CredentialTypePassport, passportVerifiableCredential.CredentialType())

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
	verifiableCredentials, err := bobby.CredentialLookupByCredentialType(credential.CredentialTypePassport)
	require.Nil(t, err)
	require.Len(t, verifiableCredentials, 1)

	passportPresentation, err := credential.NewPresentation().
		PresentationType(credential.PresentationTypePassport).
		Holder(
			credential.AddressAureWithKey(
				bobbyIdentifiers[0],
				bobbyKeys[0],
			),
		).
		CredentialAdd(verifiableCredentials[0]).
		Finish()

	require.Nil(t, err)

	passportVerifiablePresentation, err := bobby.PresentationIssue(passportPresentation)
	require.Nil(t, err)

	err = passportVerifiablePresentation.Validate()
	require.Nil(t, err)

	credentials := passportVerifiablePresentation.Credentials()
	assert.Len(t, credentials, 1)
}

func TestAccountPersistence(t *testing.T) {
	alicePath, bobbyPath := t.TempDir(), t.TempDir()

	alice, _, aliceWel := testAccountWithPath(t, alicePath)
	bobby, _, _ := testAccountWithPath(t, bobbyPath)

	aliceAddress, err := alice.InboxOpen()
	require.Nil(t, err)

	bobbyAddress, err := bobby.InboxOpen()
	require.Nil(t, err)

	// establish an encrypted connection
	err = alice.ConnectionNegotiate(
		aliceAddress,
		bobbyAddress,
		time.Now().Add(time.Hour),
	)

	require.Nil(t, err)

	// wait for negotiation to finish
	<-aliceWel

	// close down alices account
	err = alice.Close()
	require.Nil(t, err)

	time.Sleep(time.Second)

	// send alice a bunch of messages
	contentForAlice, err := message.NewChat().
		Message("hello").
		Finish()

	require.Nil(t, err)

	for i := 0; i < 100; i++ {
		err = bobby.MessageSend(
			aliceAddress,
			contentForAlice,
		)
		require.Nil(t, err)
	}

	// reopen alices account
	alice, aliceInbox, _ := testAccountWithPath(t, alicePath)

	// receive the messages from bobby
	for i := 0; i < 100; i++ {
		<-aliceInbox
	}

	inboxes, err := alice.InboxList()
	require.Nil(t, err)
	assert.Len(t, inboxes, 2)
}

func TestAccountDiscovery(t *testing.T) {
	alice, _, _ := testAccount(t)

	// alice opens an inbox
	address, err := alice.InboxOpen()
	require.Nil(t, err)

	keyPackage, err := alice.ConnectionNegotiateOutOfBand(
		address,
		time.Now().Add(time.Hour*99999),
	)

	require.Nil(t, err)

	// generate a discovery request
	content, err := message.NewDiscoveryRequest().
		Expires(time.Now().Add(time.Hour * 99999)).
		KeyPackage(keyPackage).
		Finish()

	require.Nil(t, err)

	anonymousMessage := event.NewAnonymousMessage(content)

	qrCode, err := anonymousMessage.EncodeToQR(event.QREncodingSVG)
	require.Nil(t, err)

	os.WriteFile("/tmp/qr.svg", qrCode, 0644)
}

func TestAccountSocketReconnect(t *testing.T) {
	t.Skip("manual test")
	testAccount(t)
	time.Sleep(time.Hour)
}

func TestFinalizers(t *testing.T) {
	t.Skip("manual test")
	alice, _, _ := testAccount(t)

	for i := 0; i < 100000; i++ {
		_, err := alice.IdentityList()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond)
		runtime.GC()
	}
}

func TestStoragePerformance(t *testing.T) {
	t.Skip("manual test")
	alice, aliceInbox, aliceWel := testAccount(t)
	bobby, bobbyInbox, _ := testAccount(t)

	aliceAddress, err := alice.InboxOpen()
	require.Nil(t, err)

	bobbyAddress, err := bobby.InboxOpen()
	require.Nil(t, err)

	// fmt.Println("alice:", aliceAddress)
	// fmt.Println("bobby:", bobbyAddress)

	err = alice.ConnectionNegotiate(
		aliceAddress,
		bobbyAddress,
		time.Now().Add(time.Hour),
	)

	require.Nil(t, err)

	// wait for negotiation to finish
	<-aliceWel

	contentForBobby, err := message.NewChat().
		Message("hello").
		Finish()

	require.Nil(t, err)

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

	contentForAlice, err := message.NewChat().
		Message("hi!").
		Finish()

	require.Nil(t, err)

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

	contentForAlice, err = message.NewChat().
		Message("hello again!").
		Finish()

	require.Nil(t, err)

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	for i := 0; i < 10000; i++ {
		err := bobby.MessageSend(
			aliceAddress,
			contentForAlice,
		)

		if err != nil {
			require.Nil(t, err)
		}
	}

	wg.Wait()

	messageFromBobby = wait(t, aliceInbox, time.Second)
	assert.Equal(t, bobbyAddress.String(), messageFromBobby.FromAddress().String())

	fmt.Println("sent and received in", time.Since(start))

	chatMessage, err = message.DecodeChat(messageFromBobby)
	require.Nil(t, err)
	assert.Equal(t, "hello again!", chatMessage.Message())
}
