package main

import (
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
)

func main() {
	cfg := &account.Config{
		StorageKey:  []byte("my secure random key"),
		StoragePath: "./storage",
		Environment: account.TargetSandbox,
		Callbacks: account.Callbacks{
			OnKeyPackage: account.DefaultKeyPackageAccept,
			OnWelcome:    account.DefaultWelcomeAccept,
			OnMessage: func(selfAccount *account.Account, msg *event.Message) {
				switch event.ContentTypeOf(msg) {
				case message.ContentTypeIntroduction:
					handleIntroduction(selfAccount, msg)
				case message.ContentTypeChat:
					handleChat(msg)
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

	// the account can now be connected with using the following inbox address.
	// any user can now negotiate an encrypted group using this address.
	log.Printf("waiting for connections. inboxAddress: %s", selfAccount.InboxDefault().String())

	runtime.Goexit()
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

	log.Printf("received introduction. documentAddress: %s", introduction.DocumentAddress())
}

func handleChat(msg *event.Message) {
	chat, err := message.DecodeChat(msg.Content())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("received message: %s", chat.Message())
}
