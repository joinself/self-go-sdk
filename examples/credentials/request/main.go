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
				case message.ContentTypeIntroduction:
					handleIntroduction(selfAccount, msg)
				case message.ContentTypeCredentialPresentationResponse:
					handleCredentialVerificationResponse(selfAccount, msg)
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

func handleCredentialVerificationResponse(selfAccount *account.Account, msg *event.Message) {
	credentialPresentationResponse, err := message.DecodeCredentialPresentationResponse(msg.Content())
	if err != nil {
		log.Fatal(err)
	}

	// a trust registry is a store of issuers and which credentials they can
	// issue. it is used to validate credentials, ensuring they were issued by
	// a trusted issuer.
	registry := credential.SandboxTrustedIssuerRegistry()
	responder := credential.AddressKey(msg.FromAddress())

	/*
		// we can also define new issuers and the credential types we trust them to issue, i.e.
		// allow them to share a profile credential containing their display name, etc
		registry.AddIssuer(responder)

		// grant our holder the authority to self issue a profile credential
		err = registry.GrantAuthority(responder, "ProfileCredential", time.Now(), nil)
		if err != nil {
			log.Fatal("failed to grant authority for a credential type", "error", err)
		}
	*/

	// validate all of the presentations and credentials. the returned credentials will be valid
	// for the provided holder address. the presentations and credentials will be validated, ensuring:
	// 1. they are well formatted
	// 2. their signatures are valid
	// 3. they have been issued by an authority for the specific credential type
	// 4. the issuer and holders keys are valid and have not been revoked
	// 5. the credentials have not been revoked (in an future release)
	// 6, any requirements have been met for specific credential types, i.e.
	//    (accompanying liveness credentials are provided)
	credentials, err := selfAccount.CredentialGraphValidFor(
		responder,
		registry,
		credentialPresentationResponse.Presentations(),
	)

	if err != nil {
		log.Fatal("failed to validate credential presentations", "error", err)
	}

	// iterate over the credentials that are valid for our holder
	for _, c := range credentials {
		claims, err := c.CredentialSubjectClaims()
		if err != nil {
			log.Fatalf("failed to parse credential claims, error: %s", err)
		}

		for k, v := range claims {
			log.Println(
				"credential value",
				"credentialType", c.CredentialType(),
				"field", k,
				"value", v,
			)
		}
	}

	os.Exit(0)
}

func handleDiscoveryResponse(selfAccount *account.Account, msg *event.Message) {
	_, err := message.DecodeDiscoveryResponse(msg.Content())
	if err != nil {
		log.Fatal(err)
	}

	sendCredentialPresentationRequest(selfAccount, msg)
}

func handleIntroduction(selfAccount *account.Account, msg *event.Message) {
	introduction, err := message.DecodeIntroduction(msg.Content())
	if err != nil {
		log.Fatal(err)
	}

	tokens, err := introduction.Tokens()
	if err != nil {
		log.Fatal(err)
	}

	for _, token := range tokens {
		err = selfAccount.TokenStore(msg.FromAddress(), msg.ToAddress(), msg.ToAddress(), token)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func sendCredentialPresentationRequest(selfAccount *account.Account, msg *event.Message) {
	content, err := message.NewCredentialPresentationRequest().
		Type([]string{"VerifiablePresentation", "CustomPresentation"}).
		Details(
			credential.CredentialTypeLiveness,
			[]*message.CredentialPresentationDetailParameter{
				message.NewCredentialPresentationDetailParameter(
					message.OperatorNotEquals,
					"sourceImageHash",
					"",
				),
			},
		).
		Details(
			credential.CredentialTypeEmail,
			[]*message.CredentialPresentationDetailParameter{
				message.NewCredentialPresentationDetailParameter(
					message.OperatorNotEquals,
					"emailAddress",
					"",
				),
			},
		).
		Finish()

	if err != nil {
		log.Fatal("failed to encode credential request message", "error", err)
	}

	// send the presentation request to the responding user
	err = selfAccount.MessageSend(msg.FromAddress(), content)
	if err != nil {
		log.Fatal("failed to send credential request message", "error", err)
	}
}
