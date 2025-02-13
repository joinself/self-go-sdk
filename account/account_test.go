package account_test

import (
	"crypto/rand"
	"encoding/hex"
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
	"github.com/joinself/self-go-sdk/keypair"
	"github.com/joinself/self-go-sdk/keypair/signing"
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
		Environment: account.TargetProductionSandbox,
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
				incomingMsg <- msg
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

	bobbyMemberAs, err := alice.GroupMemberAs(bobbyGroupWith)
	require.Nil(t, err)
	assert.True(t, aliceAddress.Matches(bobbyMemberAs))
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

func TestAccountMessageSigning(t *testing.T) {
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

	// create an identity document for alice
	aliceIdentifier, err := alice.KeychainSigningCreate()
	require.Nil(t, err)

	aliceInvocation, err := alice.KeychainSigningCreate()
	require.Nil(t, err)

	aliceIdentityOperation := identity.NewOperation().
		Identifier(aliceIdentifier).
		Sequence(0).
		Timestamp(time.Now()).
		GrantEmbedded(aliceInvocation, identity.RoleInvocation).
		GrantEmbedded(aliceAddress, identity.RoleVerification|identity.RoleMessaging|identity.RoleAuthentication).
		SignWith(aliceIdentifier).
		SignWith(aliceInvocation).
		SignWith(aliceAddress).
		Finish()

	err = alice.IdentityExecute(aliceIdentityOperation)
	require.Nil(t, err)

	// create an identity document for bobby
	bobbyIdentifier, err := bobby.KeychainSigningCreate()
	require.Nil(t, err)

	bobbyInvocation, err := bobby.KeychainSigningCreate()
	require.Nil(t, err)

	bobbyAssertion, err := bobby.KeychainSigningCreate()
	require.Nil(t, err)

	bobbyIdentityOperation := identity.NewOperation().
		Identifier(bobbyIdentifier).
		Sequence(0).
		Timestamp(time.Now()).
		GrantEmbedded(bobbyInvocation, identity.RoleInvocation).
		GrantEmbedded(bobbyAssertion, identity.RoleAssertion).
		GrantEmbedded(bobbyAddress, identity.RoleVerification|identity.RoleMessaging|identity.RoleAuthentication).
		SignWith(bobbyIdentifier).
		SignWith(bobbyInvocation).
		SignWith(bobbyAssertion).
		SignWith(bobbyAddress).
		Finish()

	err = bobby.IdentityExecute(bobbyIdentityOperation)
	require.Nil(t, err)

	// exchange introductions
	contentForBobby, err := message.NewIntroduction().
		DocumentAddress(aliceIdentifier).
		Finish()

	require.Nil(t, err)

	err = alice.MessageSend(
		bobbyAddress,
		contentForBobby,
	)

	require.Nil(t, err)

	messageFromAlice := wait(t, bobbyInbox, time.Second)
	assert.Equal(t, aliceAddress.String(), messageFromAlice.FromAddress().String())

	aliceIntroduction, err := message.DecodeIntroduction(messageFromAlice)
	require.Nil(t, err)
	assert.True(t, aliceIntroduction.DocumentAddress().Matches(aliceIdentifier))

	contentForAlice, err := message.NewIntroduction().
		DocumentAddress(bobbyIdentifier).
		Finish()

	require.Nil(t, err)

	err = bobby.MessageSend(
		aliceAddress,
		contentForAlice,
	)

	require.Nil(t, err)

	messageFromBobby := wait(t, aliceInbox, time.Second)
	assert.Equal(t, bobbyAddress.String(), messageFromBobby.FromAddress().String())

	bobbyIntroduction, err := message.DecodeIntroduction(messageFromBobby)
	require.Nil(t, err)
	assert.True(t, bobbyIntroduction.DocumentAddress().Matches(bobbyIdentifier))

	// select which of bobbys keys we want to use
	bobbyIdentityDocument, err := alice.IdentityResolve(bobbyIntroduction.DocumentAddress())
	require.Nil(t, err)

	bobbyInvocationKeys := bobbyIdentityDocument.SigningKeysWithRoles(identity.RoleInvocation)
	require.Len(t, bobbyInvocationKeys, 1)
	assert.True(t, bobbyInvocation.Matches(bobbyInvocationKeys[0]))

	// create an identity document and keys that management can be shared
	sharedIdentifier, err := alice.KeychainSigningCreate()
	require.Nil(t, err)

	sharedIdentityOperation := identity.NewOperation().
		Identifier(sharedIdentifier).
		Sequence(0).
		Timestamp(time.Now()).
		GrantReferenced(
			identity.MethodAure,
			aliceIdentifier,
			aliceInvocation,
			identity.RoleInvocation,
		).
		GrantReferenced(
			identity.MethodAure,
			bobbyIdentifier,
			bobbyInvocation,
			identity.RoleInvocation,
		).
		SignWith(sharedIdentifier).
		SignWith(aliceInvocation).
		Finish()

	err = alice.IdentitySign(sharedIdentityOperation)
	require.Nil(t, err)

	contentForBobby, err = message.NewSigningRequest().
		Payload(message.NewSigningIdentityDocumentOperation(
			sharedIdentifier,
			sharedIdentityOperation,
		)).
		RequireLiveness().
		Finish()

	require.Nil(t, err)

	// send siging request to bobby
	err = alice.MessageSend(
		bobbyAddress,
		contentForBobby,
	)

	require.Nil(t, err)

	messageFromAlice = wait(t, bobbyInbox, time.Second)
	assert.Equal(t, aliceAddress.String(), messageFromAlice.FromAddress().String())

	signingRequest, err := message.DecodeSigningRequest(messageFromAlice)
	require.Nil(t, err)
	assert.True(t, signingRequest.RequiresLiveness())

	unsignedPayload := signingRequest.Payload()
	unsignedIdentityDocumentOperation, err := unsignedPayload.AsIdentityDocumentOperation()
	require.Nil(t, err)

	// validate the operation and check what keys it's adding
	unsignedOperation := unsignedIdentityDocumentOperation.Operation()
	actions := unsignedOperation.Actions()

	require.Len(t, actions, 2)
	assert.Equal(t, actions[0].Action(), identity.ActionGrant)
	assert.Equal(t, actions[0].Roles(), identity.RoleInvocation)
	assert.Equal(t, actions[0].Description(), identity.DescriptionReference)
	assert.Equal(t, actions[1].Action(), identity.ActionGrant)
	assert.Equal(t, actions[1].Roles(), identity.RoleInvocation)
	assert.Equal(t, actions[1].Description(), identity.DescriptionReference)

	reference := actions[0].Reference()
	assert.Equal(t, keypair.KeyTypeSigning, reference.Address().Type())
	assert.Equal(t, identity.MethodAure, reference.Method())
	assert.True(t, reference.Address().(*signing.PublicKey).Matches(aliceInvocation))
	assert.True(t, reference.Controller().Matches(aliceIdentifier))

	reference = actions[1].Reference()
	assert.Equal(t, keypair.KeyTypeSigning, reference.Address().Type())
	assert.Equal(t, identity.MethodAure, reference.Method())
	assert.True(t, reference.Address().(*signing.PublicKey).Matches(bobbyInvocation))
	assert.True(t, reference.Controller().Matches(bobbyIdentifier))

	// check alice has signed the operation (optional)
	assert.True(t, unsignedOperation.SignedBy(aliceInvocation))

	// request a liveness check using the hash of the operation (mocked here as a self signed credential)
	unverifiedCredential, err := credential.NewCredential().
		CredentialType(credential.CredentialTypeLiveness).
		CredentialSubject(credential.AddressAure(
			bobbyIdentifier,
		)).
		CredentialSubjectClaim(
			"requestHash",
			hex.EncodeToString(unsignedOperation.Hash()),
		).
		ValidFrom(time.Now()).
		Issuer(credential.AddressAure(
			bobbyIdentifier,
		)).
		SignWith(bobbyAssertion, time.Now()).
		Finish()

	require.Nil(t, err)

	verifiedCredential, err := bobby.CredentialIssue(unverifiedCredential)
	require.Nil(t, err)

	unverifiedPresentation, err := credential.NewPresentation().
		PresentationType(credential.PresentationTypeLiveness).
		CredentialAdd(verifiedCredential).
		Holder(credential.AddressAureWithKey(
			sharedIdentifier,
			bobbyInvocation,
		)).
		Finish()

	require.Nil(t, err)

	verifiedPresentation, err := bobby.PresentationIssue(unverifiedPresentation)
	require.Nil(t, err)

	// send the response and indicate which key is to be used to sign the operation
	contentForAlice, err = message.NewSigningResponse().
		ResponseTo(messageFromAlice.ID()).
		Status(message.ResponseStatusAccepted).
		Payload(unsignedPayload).
		SignWith(bobbyInvocation).
		Presentation(verifiedPresentation).
		Finish()

	require.Nil(t, err)

	// send a response from bobby
	err = bobby.MessageSend(
		aliceAddress,
		contentForAlice,
	)

	require.Nil(t, err)

	messageFromBobby = wait(t, aliceInbox, time.Second)
	assert.Equal(t, bobbyAddress.String(), messageFromBobby.FromAddress().String())

	signingResponse, err := message.DecodeSigningResponse(messageFromBobby)
	require.Nil(t, err)

	signedPayload := signingResponse.Payload()

	// IMPORTANT - verify we have signed this operation to ensure it's the same
	// one we sent in the original request BEFORE we execute it
	// we also must ensure that the response contains a liveness presentation
	// linked to the operation hash that was signed

	require.Equal(t, message.SigningPayloadIdentityDocumentOperation, signedPayload.PayloadType())

	signedIdentityDocumentOperation, err := signedPayload.AsIdentityDocumentOperation()
	require.Nil(t, err)

	assert.True(t, sharedIdentifier.Matches(signedIdentityDocumentOperation.DocumentAddress()))

	signedOperation := signedIdentityDocumentOperation.Operation()
	require.True(t, signedOperation.SignedBy(aliceInvocation))
	require.True(t, signedOperation.SignedBy(bobbyInvocation))

	presentations := signingResponse.Presentations()
	require.Len(t, presentations, 1)
	assert.Nil(t, presentations[0].Validate())
	assert.Equal(t, credential.PresentationTypeLiveness, presentations[0].PresentationType())
	assert.True(t, presentations[0].Holder().Address().Matches(sharedIdentifier))

	credentials := presentations[0].Credentials()
	require.Len(t, credentials, 1)
	assert.Nil(t, credentials[0].Validate())
	assert.True(t, credentials[0].CredentialSubject().Address().Matches(bobbyIdentifier))

	claimRequestHash, ok := credentials[0].CredentialSubjectClaim("requestHash")
	require.True(t, ok)
	assert.Equal(t, hex.EncodeToString(unsignedOperation.Hash()), claimRequestHash)
	assert.Equal(t, hex.EncodeToString(signedOperation.Hash()), claimRequestHash)

	// store the presentation to be used when asserting the liveness requirements for the operation has been met
	err = alice.PresentationStore(presentations[0])
	require.Nil(t, err)

	// execute the completed operation
	err = alice.IdentityExecute(signedOperation)
	require.Nil(t, err)
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
