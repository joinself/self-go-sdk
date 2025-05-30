package main

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
)

var requests sync.Map

func main() {
	// initialize an account that will be used to interact with self and other entities on
	// the network. the account provides storage of all cryptographic key material, as well
	// as credentials and all state used for e2e encrypted messaging with other entitites
	cfg := &account.Config{
		// provide a secure storage key that will  be used to encrypt your local account
		// state. this should be replaced with a securely generated key!
		StorageKey: make([]byte, 32),
		// provide a storage path to the directory where your local account state will be
		// stored
		StoragePath: "./storage",
		// provide an environment to target [Develop, Sandbox]
		Environment: account.TargetSandbox,
		// provide the level of log granularity [Error, Warn, Info, Debug, Trace]
		LogLevel: account.LogWarn,
		// specify callbacks to handle events
		Callbacks: account.Callbacks{
			// invoked when the messaging socket connects
			OnConnect: func(selfAccount *account.Account) {
				log.Info("messaging socket connected")
			},
			// invoked when the messaging socket disconnects. if there is no error
			OnDisconnect: func(selfAccount *account.Account, err error) {
				if err != nil {
					log.Warn("messaging socket disconnected", "error", err)
				} else {
					log.Info("messaging socket disconnected")
				}
			},
			// invoked when there is a response to a discovery request from a new address.
			OnWelcome: func(selfAccount *account.Account, wlc *event.Welcome) {
				// we have received a response to our discovery request that is from a new
				// user/address that we do not have an  end to end encrypted session.
				// accept the invite to join the encrypted group created by the user.
				groupAddres, err := selfAccount.ConnectionAccept(
					wlc.ToAddress(),
					wlc.Welcome(),
				)

				if err != nil {
					log.Warn("failed to accept connection to encrypted group", "error", err.Error())
					return
				}

				log.Info(
					"accepted connection encrypted group",
					"from", wlc.FromAddress().String(),
					"group", groupAddres.String(),
				)
			},
			// invoked when there is a message sent to an encrypted group we are subscribed to
			OnMessage: func(selfAccount *account.Account, msg *event.Message) {
				switch event.ContentTypeOf(msg) {
				case message.ContentTypeDiscoveryResponse:
					log.Info(
						"received response to discovery request",
						"from", msg.FromAddress().String(),
						"requestId", hex.EncodeToString(msg.ID()),
					)

					discoveryResponse, err := message.DecodeDiscoveryResponse(msg.Content())
					if err != nil {
						log.Warn("failed to decode discovery response", "error", err)
						return
					}

					completer, ok := requests.LoadAndDelete(hex.EncodeToString(discoveryResponse.ResponseTo()))
					if !ok {
						log.Warn(
							"received response to unknown request",
							"requestId", hex.EncodeToString(msg.ID()),
						)
						return
					}

					completer.(chan *message.DiscoveryResponse) <- discoveryResponse
				}
			},
		},
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
		// to determine which user we are interacting with, we can generate a
		// discovery request and encode it to a qr code that your users can scan.
		// as we may not have interacted with this user before, we need to prepare
		// information they need to establish an encrypted group

		// get a key package that the responder can use to create an encryped group
		// for us, if there is not already an existing one.
		keyPackage, err := selfAccount.ConnectionNegotiateOutOfBand(
			inboxAddress,
			time.Now().Add(time.Minute*5),
		)

		if err != nil {
			log.Fatal("failed to generate key package", "error", err)
		}

		// build the key package into a discovery request
		content, err := message.NewDiscoveryRequest().
			KeyPackage(keyPackage).
			Expires(time.Now().Add(time.Minute * 5)).
			Finish()

		if err != nil {
			log.Fatal("failed to build discovery request", "error", err)
		}

		// create a channel to track the response from our qr code
		completer := make(chan *message.DiscoveryResponse, 1)

		requests.Store(
			hex.EncodeToString(content.ID()),
			completer,
		)

		// encode it as a QR code. This can be encoded as either an SVG
		// for use in rendering on a web page, or Unicode, for encoding
		// in text based environments like a terminal
		qrCode, err := event.NewAnonymousMessage(content).
			SetFlags(event.MessageFlagTargetSandbox).
			EncodeToQR(event.QREncodingUnicode)

		if err != nil {
			log.Fatal("failed to encode anonymous message", "error", err)
		}

		log.Info("scan the qr code to complete the discovery request")

		fmt.Println(string(qrCode))

		log.Info(
			"waiting for response to discovery request",
			"requestId", hex.EncodeToString(content.ID()),
		)

		<-completer

		log.Info(
			"received response to discovery request",
			"requestId", hex.EncodeToString(content.ID()),
		)
	}
}
