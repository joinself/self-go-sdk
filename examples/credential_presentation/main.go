package main

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joinself/self-go-sdk-next/account"
	"github.com/joinself/self-go-sdk-next/credential"
	"github.com/joinself/self-go-sdk-next/examples/common"
	"github.com/joinself/self-go-sdk-next/message"
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
		case message.TypeCredentialPresentationResponse:
			log.Info(
				"received response to credential presentation request",
				"from", msg.FromAddress().String(),
				"requestId", hex.EncodeToString(msg.ID()),
			)

			credentialPresentationResponse, err := message.DecodeCredentialPresentationResponse(msg)
			if err != nil {
				log.Warn("failed to decode discovery response", "error", err)
				return
			}

			completer, ok := requests.LoadAndDelete(hex.EncodeToString(credentialPresentationResponse.ResponseTo()))
			if !ok {
				log.Warn(
					"received response to unknown request",
					"requestId", hex.EncodeToString(msg.ID()),
					"responseTo", hex.EncodeToString(credentialPresentationResponse.ResponseTo()),
				)
				return
			}

			completer.(chan *message.CredentialPresentationResponse) <- credentialPresentationResponse
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

		// create a new request and store a reference to it
		content, err := message.NewCredentialPresentationRequest().
			Type([]string{"VerifiablePresentation", "CustomPresentation"}).
			Details(credential.CredentialTypeLiveness, "livenessImageHash").
			Details(credential.CredentialTypeEmail, "emailAddress").
			Finish()

		if err != nil {
			log.Fatal("failed to encode credential request message", "error", err)
		}

		presentationCompleter := make(chan *message.CredentialPresentationResponse, 1)

		requests.Store(
			hex.EncodeToString(content.ID()),
			presentationCompleter,
		)

		// send the presentation request to the responding user
		err = selfAccount.MessageSend(responderAddress, content)
		if err != nil {
			log.Fatal("failed to send credential request message", "error", err)
		}

		response := <-presentationCompleter

		// validate the presentations and the
		for _, p := range response.Presentations() {
			err = p.Validate()
			if err != nil {
				log.Warn("failed to validate presentation", "error", err)
				continue
			}

			// check the presentation references the address we are communicating with
			if !p.Holder().Address().Matches(responderAddress) {
				log.Warn("recevied a presentation response for a different holder address")
				continue
			}

			for _, c := range p.Credentials() {
				err = c.Validate()
				if err != nil {
					log.Warn("failed to validate credential", "error", err)
					continue
				}

				// check that the credential is not yet valid for use
				if c.ValidFrom().After(time.Now()) {
					log.Warn("credential is intended to be used in the future")
					continue
				}

				// TODO check issuer identity, this is not working as keys are not setup
				// correctly on the verify service...
				/*
					if c.Issuer().Method() == credential.MethodAure {

						// check the issuer was valid at the time of issuance
						document, err := selfAccount.IdentityResolve(c.Issuer().Address())
						if err != nil {
							log.Warn("failed to resolve credential issuer", "error", err)
							continue
						}

						if !document.HasRolesAt(c.Issuer().SigningKey(), identity.RoleAssertion, c.Created()) {
							log.Warn("credential signing key was not valid at the time of issuance")
							continue
						}
					}
				*/

				claims, err := c.CredentialSubjectClaims()
				if err != nil {
					log.Warn("failed to parse credential claims", "error", err)
					continue
				}

				for k, v := range claims {
					log.Info(
						"credential value",
						"credentialType", c.CredentialType(),
						"field", k,
						"value", v,
					)
				}
			}
		}
	}
}
