package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/keypair/signing"
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
					handleDiscoveryResponse(selfAccount, msg)
				case message.ContentTypeCredentialVerificationResponse:
					handleCredentialVerificationResponse(selfAccount, msg)
				default:
					log.Printf("received unhandled event")
				}
			},
		},
	}

	selfAccount, err := account.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	expiry := time.Now().Add(time.Minute * 5)

	keyPackage, err := selfAccount.ConnectionNegotiateOutOfBand(selfAccount.InboxDefault(), expiry)
	if err != nil {
		log.Fatal(err)
	}

	content, err := message.NewDiscoveryRequest().KeyPackage(keyPackage).Expires(expiry).Finish()
	if err != nil {
		log.Fatal(err)
	}

	requests.Store(hex.EncodeToString(content.ID()), struct{}{})

	qrCode, err := event.NewAnonymousMessage(content).SetFlags(event.MessageFlagTargetSandbox).EncodeToQR(event.QREncodingUnicode)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(qrCode))

	runtime.Goexit()
}

func handleDiscoveryResponse(selfAccount *account.Account, msg *event.Message) {
	discoveryResponse, err := message.DecodeDiscoveryResponse(msg.Content())
	if err != nil {
		log.Fatal(err.Error())
	}

	_, ok := requests.LoadAndDelete(hex.EncodeToString(discoveryResponse.ResponseTo()))
	if !ok {
		log.Println("received response to an unknown discovery request")
		return
	}

	requestCredentialVerification(selfAccount, msg.FromAddress())
}

func handleCredentialVerificationResponse(selfAccount *account.Account, msg *event.Message) {
	credentialVerificationResponse, err := message.DecodeCredentialVerificationResponse(msg.Content())
	if err != nil {
		log.Fatal("failed to decode discovery response", "error", err)
	}

	for _, c := range credentialVerificationResponse.Credentials() {
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
			if subjectID == msg.FromAddress().String() {
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
}

func requestCredentialVerification(selfAccount *account.Account, responder *signing.PublicKey) {
	agreement, agreementTerms := createAgreementTerms(selfAccount, responder)

	// create a new request and store a reference to it
	content, err := message.NewCredentialVerificationRequest().
		Type([]string{"VerifiableCredential", "AgreementCredential"}).
		Evidence("terms", agreementTerms).
		Proof(agreement).
		Expires(time.Now().Add(time.Hour * 24)).
		Finish()

	if err != nil {
		log.Fatal(err)
	}

	err = selfAccount.MessageSend(responder, content)
	if err != nil {
		log.Fatal(err)
	}
}

func createAgreementTerms(selfAccount *account.Account, responder *signing.PublicKey) (*credential.VerifiablePresentation, *object.Object) {
	// Read the terms.pdf file
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Agreement")

	agreementBuf := bytes.NewBuffer(make([]byte, 1024))

	err := pdf.Output(agreementBuf)
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
				"id":   responder.String(),
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

	return signedAgreementPresentation, agreementTerms
}
