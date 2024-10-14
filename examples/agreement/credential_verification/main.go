package main

import (
	"bytes"
	"encoding/hex"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-pdf/fpdf"
	"github.com/joinself/self-go-sdk-next/account"
	"github.com/joinself/self-go-sdk-next/credential"
	"github.com/joinself/self-go-sdk-next/examples/common"
	"github.com/joinself/self-go-sdk-next/message"
	"github.com/joinself/self-go-sdk-next/object"
)

var requests sync.Map

func main() {
	discoveryManager := common.NewDiscoveryManager()
	// Load the default account config
	cfg := common.NewExamplesDefaultConfig()

	// Override the default OnMessagecallback to handle discovery responses and chat messages
	cfg.Callbacks.OnMessage = func(selfAccount *account.Account, msg *message.Message) {
		switch message.ContentType(msg) {
		case message.TypeDiscoveryResponse:
			discoveryManager.HandleDiscoveryResponse(selfAccount, msg)
		case message.TypeCredentialVerificationResponse:
			log.Info(
				"received response to credential verification request",
				"from", msg.FromAddress().String(),
				"requestId", hex.EncodeToString(msg.ID()),
			)

			credentialVerificationResponse, err := message.DecodeCredentialVerificationResponse(msg)
			if err != nil {
				log.Warn("failed to decode discovery response", "error", err)
				return
			}

			completer, ok := requests.LoadAndDelete(hex.EncodeToString(credentialVerificationResponse.ResponseTo()))
			if !ok {
				log.Warn(
					"received response to unknown request",
					"requestId", hex.EncodeToString(msg.ID()),
					"responseTo", hex.EncodeToString(credentialVerificationResponse.ResponseTo()),
				)
				return
			}

			completer.(chan *message.CredentialVerificationResponse) <- credentialVerificationResponse
		}
	}

	log.Info("initializing self account")

	// initialize and load the account
	selfAccount, err := account.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize account", "error", err)
	}

	// TODO : this will look slightly different in production.
	// right now, we can just open an inbox to send and receive
	// messages from it. In the future we will hide some of this
	// and do proper linking with the application identity.
	// NB: this does not need to happen every time we start the SDK,
	// only once!
	inboxAddress, err := selfAccount.InboxOpen()
	if err != nil {
		log.Fatal("failed to open account inbox", "error", err)
	}

	log.Info("initialized account success")

	for {
		// generate a discovery request and display a QR code that the user can scan
		completer, err := discoveryManager.GenerateAndDisplayQRCode(selfAccount, inboxAddress)
		if err != nil {
			log.Fatal("failed to generate and display QR code", "error", err)
		}
		// wait for a discovery response from the user
		discoveryResponse := <-completer
		log.Info(
			"received response to discovery request",
			"requestId", hex.EncodeToString(discoveryResponse.ID()),
		)

		responderAddress := discoveryResponse.FromAddress()

		// Read the terms.pdf file
		pdf := fpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, "Agreement")

		agreementBuf := bytes.NewBuffer(make([]byte, 1024))

		err = pdf.Output(agreementBuf)
		if err != nil {
			log.Fatal("failed to build agreement pdf", "error", err.Error())
		}

		agreementTerms, err := object.New(
			"application/pdf",
			agreementBuf.Bytes(),
		)

		if err != nil {
			log.Error("failed to encrypt object", "error", err)
		}

		err = selfAccount.ObjectUpload(
			inboxAddress,
			agreementTerms,
			false,
		)

		if err != nil {
			log.Error("failed to upload object", "error", err)
		}

		// create the body of the credential agreement
		claims := make(map[string]interface{})

		// the id of our unencrypted object is the SHA3 hash of our terms
		claims["termsHash"] = hex.EncodeToString(agreementTerms.Id())

		// TODO we can support different signatory types [signatory, witness, promisee, promisor, etc]
		// OR we leave that up to the content of the terms to define?
		aliceParty := map[string]string{
			"type": "signatory",
			"id":   inboxAddress.String(),
		}

		bobbyParty := map[string]string{
			"type": "signatory",
			"id":   responderAddress.String(),
		}

		// when validating the agreement, any validator will be
		// required to ensure they have been presented signed
		// credentials for all parties
		claims["parties"] = []map[string]string{aliceParty, bobbyParty}

		// create a credential to serve as our agreement
		// the subject of our credential will be ourselves,
		// signifying our agreement to the terms.
		// our counterparty will issue a credential in the
		// same manner.
		unsignedAgreementCredential, err := credential.NewCredential().
			CredentialType([]string{"VerifiableCredential", "AgreementCredential"}).
			CredentialSubject(credential.AddressKey(responderAddress)).
			CredentialSubjectClaims(claims).
			// CredentialSubjectClaim("terms", hex.EncodeToString(agreementTerms.Id())).
			Issuer(credential.AddressKey(inboxAddress)).
			ValidFrom(time.Now()).
			SignWith(inboxAddress, time.Now()).
			Finish()

		if err != nil {
			log.Error("failed to create credential", "error", err)
		}

		signedAgreementCredential, err := selfAccount.CredentialIssue(unsignedAgreementCredential)
		if err != nil {
			log.Error("failed to issue credential", "error", err)
		}

		// create a new request and store a reference to it
		content, err := message.NewCredentialVerificationRequest().
			Type([]string{"VerifiableCredential", "AgreementCredential"}).
			Evidence("terms", agreementTerms).
			Proof(signedAgreementCredential).
			Expires(time.Now().Add(time.Hour * 24)).
			Finish()

		if err != nil {
			log.Fatal("failed to encode credential request message", "error", err)
		}

		verificationCompleter := make(chan *message.CredentialVerificationResponse, 1)

		requests.Store(
			hex.EncodeToString(content.ID()),
			verificationCompleter,
		)

		// send the presentation request to the responding user
		err = selfAccount.MessageSend(responderAddress, content)
		if err != nil {
			log.Fatal("failed to send credential request message", "error", err)
		}

		response := <-verificationCompleter

		log.Info("Response received with status", "status", response.Status())
		for _, c := range response.Credentials() {
			err = c.Validate()
			if err != nil {
				log.Warn("failed to validate credential", "error", err)
				continue
			}

			// check that the credential is not yet valid for use
			/*
				if c.ValidFrom().After(time.Now()) {
					log.Warn("credential is intended to be used in the future")
					continue
				}
			*/

			claims, err := c.CredentialSubjectClaims()
			if err != nil {
				log.Warn("failed to parse credential claims", "error", err)
				continue
			}

			parties, ok := claims["parties"].([]interface{})
			if !ok {
				log.Warn("parties claim is not an array")
				continue
			}

			var isIssued, isSigner bool

			for _, subject := range parties {
				subjectDetails, ok := subject.(map[string]interface{})
				if !ok {
					log.Warn("subject is not an object")
					continue
				}

				subjectType, ok := subjectDetails["type"].(string)
				if !ok || subjectType != "signatory" {
					continue
				}

				subjectID, ok := subjectDetails["id"].(string)
				if !ok {
					log.Warn("subject id is not a string")
					continue
				}

				// check if the agreement issuer (alice) is provided in the agreement
				if subjectID == inboxAddress.String() {
					isIssued = true
				}

				// check if the responder is included as a signer in the agreement
				if subjectID == responderAddress.String() {
					isSigner = true
				}
			}

			if isIssued && isSigner {
				log.Info("Agreement is valid and signed by both parties")
				selfAccount.CredentialStore(c)
			} else {
				log.Warn("Agreement is not valid or not signed by both parties")
			}
		}
		os.Exit(1)
	}
}
