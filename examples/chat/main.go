package main

import (
	"encoding/hex"

	"github.com/charmbracelet/log"
	"github.com/joinself/self-go-sdk-next/account"
	"github.com/joinself/self-go-sdk-next/examples/common"
	"github.com/joinself/self-go-sdk-next/message"
)

func main() {
	discoveryManager := common.NewDiscoveryManager()
	// Load the default account config
	cfg := common.NewExamplesDefaultConfig()

	// Override the default OnMessagecallback to handle discovery responses and chat messages
	cfg.Callbacks.OnMessage = func(selfAccount *account.Account, msg *message.Message) {
		switch message.ContentType(msg) {
		case message.TypeDiscoveryResponse:
			discoveryManager.HandleDiscoveryResponse(selfAccount, msg)
		case message.TypeChat:
			chat, err := message.DecodeChat(msg)
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
	}

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
		response := <-completer
		log.Info(
			"received response to discovery request",
			"requestId", hex.EncodeToString(response.ID()),
		)

		// prepare a chat message to send to the user
		content, err := message.NewChat().
			Message("Hello!").
			Finish()
		if err != nil {
			log.Fatal("failed to encode chat message", "error", err)
		}

		// send the chat message to the user
		log.Info("sending message", "toAddress", response.FromAddress().String())
		err = selfAccount.MessageSend(
			response.FromAddress(),
			content,
		)
		if err != nil {
			log.Fatal("failed to send chat message", "error", err)
		}
		log.Info("sent message", "toAddress", response.FromAddress().String())
	}
}
