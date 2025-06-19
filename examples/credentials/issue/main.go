package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
)

func main() {
	cfg := &account.Config{
		StorageKey:  []byte("my secure random key"),
		StoragePath: "./storage",
		Environment: account.TargetSandbox,
		Callbacks: account.Callbacks{
			OnWelcome: account.DefaultWelcomeAccept,
			OnMessage: func(selfAccount *account.Account, msg *event.Message) {
				switch event.ContentTypeOf(msg) {
				case message.ContentTypeDiscoveryResponse:
					handleDiscoveryResponse(selfAccount, msg)
				default:
					log.Printf("received unhandled event")
				}
			},
		},
	}

	selfAccount, err := account.New(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	expires := time.Now().Add(time.Minute * 5)

	keyPackage, err := selfAccount.ConnectionNegotiateOutOfBand(selfAccount.InboxDefault(), expires)
	if err != nil {
		log.Fatal(err.Error())
	}

	content, err := message.NewDiscoveryRequest().KeyPackage(keyPackage).Expires(expires).Finish()
	if err != nil {
		log.Fatal(err.Error())
	}

	qrCode, err := event.NewAnonymousMessage(content).SetFlags(event.MessageFlagTargetSandbox).EncodeToQR(event.QREncodingUnicode)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(qrCode))

	runtime.Goexit()
}

func handleDiscoveryResponse(selfAccount *account.Account, msg *event.Message) {
	// specify the type of credential. It's possible for a credential to have
	// more than one type as credentials can express different claims, i.e
	// a credential that holds both contact details and a passport.
	credentialType := []string{"VerifiableCredential", "CustomerCredential"}

	// the subject address the credential will be issued for. as we don't
	// have an document address for our responder (that would be shared in
	// an introduction message), we can use a key method indicating it is
	// referencing the responders messaging address.
	subjectAddress := credential.AddressKey(msg.FromAddress())

	// the address that will be asserting and issuing the claims about the
	// subject. if our sdk has been paired to an application, we may use
	// the applications address as `credential.AddressAure(applicationAddress)`
	issuerAddress := credential.AddressKey(selfAccount.InboxDefault())

	// create a new customer credential for our responder.
	customerCredential, err := credential.NewCredential().
		CredentialType(credentialType).
		CredentialSubject(subjectAddress).
		CredentialSubjectClaims(map[string]any{
			"customer": map[string]any{
				"customerOf": issuerAddress,
			},
		}).
		Issuer(issuerAddress).
		ValidFrom(time.Now()).
		SignWith(selfAccount.InboxDefault(), time.Now()).
		Finish()

	if err != nil {
		log.Fatal(err)
	}

	// sign and issue the verifiable credential with our account
	customerVerifiableCredential, err := selfAccount.CredentialIssue(customerCredential)
	if err != nil {
		log.Fatal(err)
	}

	// create a new credential message to send an unsolicited credential to
	// a given address.
	content, err := message.NewCredential().VerifiableCredential(customerVerifiableCredential).Finish()
	if err != nil {
		log.Fatal(err)
	}

	// send the credential message to the responding user
	err = selfAccount.MessageSend(msg.FromAddress(), content)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)
	os.Exit(0)
}
