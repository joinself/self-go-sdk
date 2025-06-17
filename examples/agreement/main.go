package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/object"
)

var requests sync.Map

func main() {
	cfg := &account.Config{
		StorageKey:  []byte("my secure random key"),
		StoragePath: "./storage",
		Environment: account.TargetSandbox,
		LogLevel:    account.LogWarn,
		Callbacks: account.Callbacks{
			OnWelcome: account.DefaultWelcomeAccept,
			OnMessage: func(selfAccount *account.Account, msg *event.Message) {
				switch event.ContentTypeOf(msg) {
				case message.ContentTypeDiscoveryResponse:
					handleDiscoveryResponse(msg)
				case message.ContentTypeCredentialVerificationResponse:
					handleCredentialVerificationResponse(msg)
				}
			},
		},
	}

	selfAccount, err := account.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	for {
		expiry := time.Now().Add(time.Minute * 5)

		keyPackage, err := selfAccount.ConnectionNegotiateOutOfBand(selfAccount.InboxDefault(), expiry)
		if err != nil {
			log.Fatal(err)
		}

		content, err := message.NewDiscoveryRequest().KeyPackage(keyPackage).Expires(expiry).Finish()
		if err != nil {
			log.Fatal(err)
		}

		discoveryCompleter := make(chan *event.Message, 1)
		requests.Store(hex.EncodeToString(content.ID()), discoveryCompleter)

		qrCode, err := event.NewAnonymousMessage(content).SetFlags(event.MessageFlagTargetSandbox).EncodeToQR(event.QREncodingUnicode)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(string(qrCode))

		discoveryResponse := <-discoveryCompleter

		// Read the terms.pdf file
		pdf := fpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, "Agreement")

		agreementBuf := bytes.NewBuffer(make([]byte, 1024))

		err = pdf.Output(agreementBuf)
		if err != nil {
			log.Fatal(err)
		}

		agreementTerms, err := object.New("application/pdf", agreementBuf.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		err = selfAccount.ObjectUpload(agreementTerms, false)
		if err != nil {
			log.Fatal(err)
		}

		// create the body of the credential agreement when validating the agreement, any validator
		// will be required to ensure they have been presented signed credentials for all parties
		claims := map[string]interface{}{
			"termsHash": hex.EncodeToString(agreementTerms.Hash()),
			"parties": []map[string]interface{}{
				{
					"type": "signatory",
					"id":   discoveryResponse.FromAddress().String(),
				},
				{
					"type": "signatory",
					"id":   selfAccount.InboxDefault().String(),
				},
			},
		}

		// create a credential to serve as our agreement the subject of our credential will be ourselves,
		// signifying our agreement to the terms. our counterparty will issue a credential in the same manner.
		unsignedAgreementCredential, err := credential.NewCredential().
			CredentialType([]string{"VerifiableCredential", "AgreementCredential"}).
			CredentialSubject(credential.AddressKey(selfAccount.InboxDefault())).
			CredentialSubjectClaims(claims).
			CredentialSubjectClaim("terms", hex.EncodeToString(agreementTerms.Id())).
			Issuer(credential.AddressKey(selfAccount.InboxDefault())).
			ValidFrom(time.Now()).
			SignWith(selfAccount.InboxDefault(), time.Now()).
			Finish()

		if err != nil {
			log.Fatal(err)
		}

		signedAgreementCredential, err := selfAccount.CredentialIssue(unsignedAgreementCredential)
		if err != nil {
			log.Fatal(err)
		}

		unsignedAgreementPresentation, err := credential.NewPresentation().
			PresentationType([]string{"VerifiablePresentation", "AgreementPresentation"}).
			Holder(credential.AddressKey(selfAccount.InboxDefault())).
			CredentialAdd(signedAgreementCredential).
			Finish()

		if err != nil {
			log.Fatal(err)
		}

		signedAgreementPresentation, err := selfAccount.PresentationIssue(unsignedAgreementPresentation)
		if err != nil {
			log.Fatal(err)
		}

		// create a new request and store a reference to it
		content, err = message.NewCredentialVerificationRequest().
			Type([]string{"VerifiableCredential", "AgreementCredential"}).
			Evidence("terms", agreementTerms).
			Proof(signedAgreementPresentation).
			Expires(time.Now().Add(time.Hour * 24)).
			Finish()

		if err != nil {
			log.Fatal(err)
		}

		verificationCompleter := make(chan *message.CredentialVerificationResponse, 1)
		requests.Store(hex.EncodeToString(content.ID()), verificationCompleter)

		err = selfAccount.MessageSend(discoveryResponse.FromAddress(), content)
		if err != nil {
			log.Fatal(err)
		}

		verificationResponse := <-verificationCompleter

		for _, c := range verificationResponse.Credentials() {
			err = c.Validate()
			if err != nil {
				log.Fatal("failed to validate credential", "error", err)
			}

			claims, err := c.CredentialSubjectClaims()
			if err != nil {
				log.Fatal(err)
			}

			parties, ok := claims["parties"].([]interface{})
			if !ok {
				log.Fatal("parties claim is not an array")
			}

			var isIssued, isSigner bool

			for _, subject := range parties {
				subjectDetails, ok := subject.(map[string]interface{})
				if !ok {
					log.Fatal("subject is not an object")
				}

				subjectType, ok := subjectDetails["type"].(string)
				if !ok || subjectType != "signatory" {
					continue
				}

				subjectID, ok := subjectDetails["id"].(string)
				if !ok {
					log.Fatal("subject id is not a string")
				}

				// check if the agreement issuer (alice) is provided in the agreement
				if subjectID == selfAccount.InboxDefault().String() {
					isIssued = true
				}

				// check if the responder is included as a signer in the agreement
				if subjectID == discoveryResponse.FromAddress().String() {
					isSigner = true
				}
			}

			if isIssued && isSigner {
				log.Println("Agreement is valid and signed by both parties")
				selfAccount.CredentialStore(c)
			} else {
				log.Fatal("Agreement is not valid or not signed by both parties")
			}
		}
		os.Exit(1)
	}
}

func handleDiscoveryResponse(msg *event.Message) {
	discoveryResponse, err := message.DecodeDiscoveryResponse(msg.Content())
	if err != nil {
		log.Fatal(err.Error())
	}

	completer, ok := requests.LoadAndDelete(hex.EncodeToString(discoveryResponse.ResponseTo()))
	if !ok {
		log.Fatal("received response to an unknown discovery request")
	}

	completer.(chan *event.Message) <- msg
}

func handleCredentialVerificationResponse(msg *event.Message) {
	credentialVerificationResponse, err := message.DecodeCredentialVerificationResponse(msg.Content())
	if err != nil {
		log.Fatal("failed to decode discovery response", "error", err)
	}

	completer, ok := requests.LoadAndDelete(hex.EncodeToString(credentialVerificationResponse.ResponseTo()))
	if !ok {
		log.Fatalf(
			"received response to unknown request. requestId: %s responseTo: %s",
			hex.EncodeToString(msg.ID()),
			hex.EncodeToString(credentialVerificationResponse.ResponseTo()),
		)
	}

	completer.(chan *message.CredentialVerificationResponse) <- credentialVerificationResponse
}
