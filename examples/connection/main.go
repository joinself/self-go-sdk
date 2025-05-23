package main

import (
	"encoding/hex"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
)

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
			// invoked when a user is attempting to establish an encrypted group with
			// this account.
			OnKeyPackage: func(selfAccount *account.Account, kpg *event.KeyPackage) {
				// establish a new encrypted group with the user using the key package
				// they have sent.
				groupAddress, err := selfAccount.ConnectionEstablish(
					kpg.ToAddress(),
					kpg.KeyPackage(),
				)

				if err != nil {
					log.Warn("failed to establish connection and encrypted group", "error", err.Error())
					return
				}

				log.Info(
					"established encrypted group",
					"with", kpg.FromAddress().String(),
					"group", groupAddress.String(),
				)

				content, err := message.NewChat().
					Message("Hello!").
					Finish()

				if err != nil {
					log.Fatal("failed to encode chat message", "error", err)
				}

				log.Info(
					"sending message",
					"toAddress", groupAddress.String(),
				)

				err = selfAccount.MessageSend(
					groupAddress,
					content,
				)

				if err != nil {
					log.Fatal("failed to send chat message", "error", err)
				}
			},
			// invoked when there is a message sent to an encrypted group we are subscribed to
			OnMessage: func(selfAccount *account.Account, msg *event.Message) {
				switch event.ContentTypeOf(msg) {
				case message.ContentTypeIntroduction:
					introduction, err := message.DecodeIntroduction(msg.Content())
					if err != nil {
						log.Warn("failed to decode introduction", "error", err)
						return
					}

					tokens, err := introduction.Tokens()
					if err != nil {
						log.Warn("failed to decode introduction tokens", "error", err)
						return
					}

					for _, token := range tokens {
						err = selfAccount.TokenStore(
							msg.FromAddress(),
							msg.ToAddress(),
							msg.ToAddress(),
							token,
						)

						if err != nil {
							log.Warn("failed to store introduction tokens", "error", err)
							return
						}
					}

					log.Info(
						"received introduction",
						"from", msg.FromAddress().String(),
						"messageId", hex.EncodeToString(msg.ID()),
						"tokens", len(tokens),
					)

				case message.ContentTypeChat:
					chat, err := message.DecodeChat(msg.Content())
					if err != nil {
						log.Warn("failed to decode chat", "error", err)
						return
					}

					log.Info(
						"received chat message",
						"from", msg.FromAddress().String(),
						"messageId", hex.EncodeToString(msg.ID()),
						"referencing", hex.EncodeToString(chat.Referencing()),
						"message", chat.Message(),
						"attachments", len(chat.Attachments()),
					)
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

	log.Info("initialized account success")

	inboxes, err := selfAccount.InboxList()
	if err != nil {
		log.Fatal("failed to list inboxes", "error", err)
	}

	// the account can now be connected with using the following inbox address.
	// any user can now negotiate an encrypted group using this address.
	log.Info("waiting for connections", "inbox_address", inboxes[0].String())

	runtime.Goexit()
}
