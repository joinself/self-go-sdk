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
	"github.com/joinself/self-go-sdk/credential/predicate"
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
	/*
		account.SetLogFunc(func(level account.LogLevel, message string) {
			// disable logging
		})
	*/
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
		SkipSetup:   true,
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
				incomingMsg <- msg
			},
			OnKeyPackage: func(account *account.Account, keyPackage *event.KeyPackage) {
				_, err := account.ConnectionEstablish(
					keyPackage.ToAddress(),
					keyPackage.KeyPackage(),
				)
				if err != nil {
					panic(err)
				}
			},
			OnWelcome: func(account *account.Account, welcome *event.Welcome) {
				_, err := account.ConnectionAccept(
					welcome.ToAddress(),
					welcome.Welcome(),
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

	chatMessage, err := message.DecodeChat(messageFromAlice.Content())
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

	chatMessage, err = message.DecodeChat(messageFromBobby.Content())
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

	chatMessage, err = message.DecodeChat(messageFromBobby.Content())
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
	t.Skip("temporarily disabled")

	alice, _, _ := testAccount(t)

	data := make([]byte, 1024)
	rand.Read(data)

	encryptedObject, err := object.New(
		"application/octet-stream",
		data,
	)

	require.Nil(t, err)

	err = alice.ObjectUpload(
		encryptedObject,
		false,
	)

	require.Nil(t, err)

	err = alice.ObjectDownload(
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
		CredentialType(credential.TypePassport).
		CredentialSubject(credential.AddressAure(bobbyIdentifiers[0])).
		CredentialSubjectClaim("firstName", "bobby").
		Issuer(credential.AddressAure(aliceIdentifiers[0])).
		ValidFrom(time.Now()).
		SignWith(aliceKeys[0], time.Now()).
		Finish()

	require.Nil(t, err)

	passportVerifiableCredential, err := alice.CredentialIssue(passportCredential)
	require.Nil(t, err)
	assert.Equal(t, credential.TypePassport, passportVerifiableCredential.CredentialType()[0])

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
	verifiableCredentials, err := bobby.CredentialLookupByCredentialType(credential.TypePassport)
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

	aliceIntroduction, err := message.DecodeIntroduction(messageFromAlice.Content())
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

	bobbyIntroduction, err := message.DecodeIntroduction(messageFromBobby.Content())
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

	signingRequest, err := message.DecodeSigningRequest(messageFromAlice.Content())
	require.Nil(t, err)
	assert.True(t, signingRequest.RequiresLiveness())

	unsignedPayloads, err := signingRequest.Payloads()
	require.Nil(t, err)
	require.Len(t, unsignedPayloads, 1)

	unsignedPayload := unsignedPayloads[0]

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
		CredentialType(credential.TypeLiveness).
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
		Payload(unsignedPayload, []*signing.PublicKey{bobbyInvocation}).
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

	signingResponse, err := message.DecodeSigningResponse(messageFromBobby.Content())
	require.Nil(t, err)

	signedPayloads, err := signingResponse.Payloads()
	require.Nil(t, err)
	require.Len(t, signedPayloads, 1)

	signedPayload := signedPayloads[0]

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
	assert.Equal(t, credential.PresentationTypeLiveness, presentations[0].PresentationType()[0])
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

func TestAccountCredentialPresentationRequest(t *testing.T) {
	alice, aliceInbox, aliceWel := testAccount(t)
	bobby, bobbyInbox, _ := testAccount(t)

	aliceAddress, err := alice.InboxOpen()
	require.Nil(t, err)

	bobbyAddress, err := bobby.InboxOpen()
	require.Nil(t, err)
	err = alice.ConnectionNegotiate(
		aliceAddress,
		bobbyAddress,
		time.Now().Add(time.Hour),
	)

	require.Nil(t, err)

	now := time.Now()

	// wait for negotiation to finish
	<-aliceWel

	// alice issues a self signed credential
	aliceCredential, err := credential.NewCredential().
		CredentialType("ContactCredential").
		CredentialSubject(credential.AddressKey(aliceAddress)).
		CredentialSubjectClaims(
			map[string]interface{}{
				"contact": map[string]interface{}{
					"name":        "Alice",
					"phoneNumber": "+1 555-12345",
					"provider":    "Verison",
				},
			},
		).
		ValidFrom(now.Add(-time.Second)).
		ValidUntil(now.Add(time.Second)).
		Issuer(credential.AddressKey(aliceAddress)).
		SignWith(aliceAddress, now.Add(-time.Second)).
		Finish()

	require.Nil(t, err)

	aliceVerifiableCredential, err := alice.CredentialIssue(
		aliceCredential,
	)
	require.Nil(t, err)

	err = alice.CredentialStore(aliceVerifiableCredential)
	require.Nil(t, err)

	// request credential from alice
	// explicity request a credential that:
	// 1. is a ContactCredential
	// 2. has a phoneNumber with the country code +1
	// 3. has a provider that's either Verison, AT&T or Mint
	contactPredicate := predicate.Contains(
		credential.FieldType,
		"ContactCredential",
	).And(
		// match all fields in the credential's claims, as credentials with unrequested fields
		// will be omitted by the responder to avoid unintended disclosure of information that
		// was not explicitly requested
		predicate.NotEmpty(
			credential.FieldSubjectClaims,
		),
	).And(
		predicate.Contains(
			"/credentialSubject/contact/phoneNumber",
			"+1 ",
		),
	).And(
		predicate.OneOf(
			"/credentialSubject/contact/provider",
			[]string{"Verison", "AT&T", "Mint"},
		),
	)

	credentialPresentationRequest, err := message.NewCredentialPresentationRequest().
		PresentationType("ContactPresentation").
		Predicates(predicate.NewTree(contactPredicate)).
		Finish()

	require.Nil(t, err)

	err = bobby.MessageSend(aliceAddress, credentialPresentationRequest)
	require.Nil(t, err)

	messageFromBobby := <-aliceInbox

	requestFromBobby, err := message.DecodeCredentialPresentationRequest(
		messageFromBobby.Content(),
	)
	require.Nil(t, err)

	// inspect the requirements of the request and gather all credentials
	// that might match the criteria
	var candidates []*credential.VerifiableCredential

	predicates := requestFromBobby.
		Predicates()

	report := predicates.FindMissingPredicates(candidates)

	for _, requirement := range report.Requirements() {
		for _, option := range requirement.Options() {
			current := len(candidates)

			for _, predicator := range option {
				if predicator.Field() == credential.FieldType {
					// lookup credentials by their type.
					// NOTE predicates might not include a credential type, so in
					// some circumstances looking up all credentials may be better
					credentials, err := alice.CredentialLookupByCredentialType(
						predicator.Values()[0],
					)
					require.Nil(t, err)

					if len(credentials) > 0 {
						candidates = append(candidates, credentials...)
						break
					}
				}
			}

			if len(candidates) > current {
				// we've found credentials that might match the criteria
				// so we can skip trying other options that might satisfy
				// the predicates
				break
			}
		}
	}

	// restrict our set of credentials based on the requesters predicates
	// if we do not have a solution to the predicates with the credentials
	// we have provided, then we will need to either find more or request
	// verification of new credentials
	credentials, solution := predicates.FindOptimalMatch(
		candidates,
	)

	if !solution {
		// we dont have a solution with the credentials available, so respond with NOT_FOUND
		credentialPresentationResponse, err := message.NewCredentialPresentationResponse().
			ResponseTo(messageFromBobby.ID()).
			Status(message.ResponseStatusNotFound).
			Finish()

		require.Nil(t, err)

		err = alice.MessageSend(bobbyAddress, credentialPresentationResponse)
		require.Nil(t, err)
	}

	// create a presentation containing the matched credentials
	alicePresentation, err := credential.NewPresentation().
		PresentationType(requestFromBobby.PresentationType()...).
		CredentialAdd(credentials...).
		Holder(credential.AddressKey(aliceAddress)).
		Finish()

	require.Nil(t, err)

	aliceVerifiablePresentation, err := alice.PresentationIssue(alicePresentation)
	require.Nil(t, err)

	credentialPresentationResponse, err := message.NewCredentialPresentationResponse().
		ResponseTo(messageFromBobby.ID()).
		Status(message.ResponseStatusOk).
		VerifiablePresentation(aliceVerifiablePresentation).
		Finish()

	require.Nil(t, err)

	err = alice.MessageSend(bobbyAddress, credentialPresentationResponse)
	require.Nil(t, err)

	messageFromAlice := <-bobbyInbox

	responseFromAlice, err := message.DecodeCredentialPresentationResponse(
		messageFromAlice.Content(),
	)
	require.Nil(t, err)
	require.Equal(t, message.ResponseStatusOk, responseFromAlice.Status())

	presentations := responseFromAlice.Presentations()
	require.Len(t, presentations, 1)

	err = presentations[0].Validate()
	require.Nil(t, err)

	credentials = presentations[0].Credentials()
	require.Len(t, credentials, 1)

	assert.Equal(t, credential.AddressKey(aliceAddress).String(), credentials[0].Issuer().String())
	assert.Equal(t, credential.AddressKey(aliceAddress).String(), credentials[0].CredentialSubject().String())

	claims, err := credentials[0].CredentialSubjectClaims()
	require.Nil(t, err)

	contact, ok := claims["contact"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Alice", contact["name"])
	assert.Equal(t, "+1 555-12345", contact["phoneNumber"])
	assert.Equal(t, "Verison", contact["provider"])
}

func TestAccountSDKSetup(t *testing.T) {
	alice, aliceInbox, _ := testAccount(t)
	bobby, _, bobbyWel := testAccount(t)

	aliceAddress, err := alice.InboxOpen()
	require.Nil(t, err)

	// fmt.Println("alice:", aliceAddress)
	// fmt.Println("bobby:", bobbyAddress)

	// create an identity document for alice that serves as our
	// application identity document the sdk will be paired to
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

	// create a sdk pairing code for bobbys account
	code, unpaired, err := bobby.SDKPairingCode()
	require.Nil(t, err)
	assert.True(t, unpaired)

	anonymousMessage, err := event.AnonymousMessageDecodeFromString(code)
	require.Nil(t, err)

	discoveryRequest, err := message.DecodeDiscoveryRequest(anonymousMessage.Content())
	require.Nil(t, err)

	// negotiate a session with our sdk instance from the discovery request
	_, err = alice.ConnectionEstablish(
		aliceAddress,
		discoveryRequest.KeyPackage(),
	)

	require.Nil(t, err)

	// wait for negotiation to finish
	<-bobbyWel

	// send a response to bobby
	contentForBobby, err := message.NewDiscoveryResponse().
		ResponseTo(anonymousMessage.ID()).
		Status(message.ResponseStatusAccepted).
		Finish()

	require.Nil(t, err)

	err = alice.MessageSend(
		discoveryRequest.KeyPackage().FromAddress(),
		contentForBobby,
	)

	require.Nil(t, err)

	// wait for pairing request from sdk and ensure it's from an address that we've entered the code from
	messageFromBobby := wait(t, aliceInbox, time.Second)
	assert.Equal(t, discoveryRequest.KeyPackage().FromAddress().String(), messageFromBobby.FromAddress().String())

	pairingRequest, err := message.DecodeAccountPairingRequest(messageFromBobby.Content())
	require.Nil(t, err)

	// check the address matches the one we have been interacting with
	assert.Equal(t, discoveryRequest.KeyPackage().FromAddress().String(), pairingRequest.Address().String())

	aliceDocument, err := alice.IdentityResolve(aliceIdentifier)
	require.Nil(t, err)

	// create an operation to add the sdks key to the application identity document
	operation := aliceDocument.
		Create().
		GrantEmbedded(pairingRequest.Address(), identity.RoleVerification|identity.RoleAuthentication|identity.RoleMessaging).
		SignWith(aliceInvocation).
		Finish()

	err = alice.IdentitySign(operation)
	require.Nil(t, err)

	// respond to the pairing request
	// this response can also include credentials or presentations
	// that the sdk can use to identify itself to others
	contentForBobby, err = message.NewAccountPairingResponse().
		ResponseTo(messageFromBobby.ID()).
		Status(message.ResponseStatusCreated).
		DocumentAddress(aliceIdentifier).
		Operation(operation).
		Finish()

	require.Nil(t, err)

	err = alice.MessageSend(
		pairingRequest.Address(),
		contentForBobby,
	)

	require.Nil(t, err)

	start := time.Now()

	for {
		_, unpaired, err := bobby.SDKPairingCode()
		require.Nil(t, err)

		if !unpaired {
			break
		}

		time.Sleep(time.Millisecond * 100)

		if time.Since(start) > time.Second*2 {
			t.FailNow()
		}
	}

	time.Sleep(time.Second)

	aliceDocument, err = bobby.IdentityResolve(aliceIdentifier)
	require.Nil(t, err)
	assert.True(t, aliceDocument.HasRolesAt(pairingRequest.Address(), identity.RoleVerification|identity.RoleAuthentication|identity.RoleMessaging, time.Now()))
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
	assert.Len(t, inboxes, 1)
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

	chatMessage, err := message.DecodeChat(messageFromAlice.Content())
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

	chatMessage, err = message.DecodeChat(messageFromBobby.Content())
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

	chatMessage, err = message.DecodeChat(messageFromBobby.Content())
	require.Nil(t, err)
	assert.Equal(t, "hello again!", chatMessage.Message())
}
