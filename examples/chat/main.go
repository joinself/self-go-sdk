package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

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
			OnWelcome: account.DefaultWelcomeAccept,
			OnMessage: func(selfAccount *account.Account, msg *event.Message) {
				switch event.ContentTypeOf(msg) {
				case message.ContentTypeDiscoveryResponse:
					handleDiscoveryResponse(selfAccount, msg)
				case message.ContentTypeChat:
					handleChat(msg)
				}
			},
		},
	}

	selfAccount, err := account.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	expires := time.Now().Add(time.Minute * 5)

	// create a new discovery request that we can encode as a QR code
	content, err := message.NewDiscoveryRequest().InboxAddress(selfAccount.InboxDefault()).Expires(expires).Finish()
	if err != nil {
		log.Fatal(err)
	}

	// format the discovery request as a QR code
	qrCode, err := event.NewAnonymousMessage(content).SetFlags(event.MessageFlagTargetSandbox).EncodeToQR(event.QREncodingUnicode)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(qrCode))

	runtime.Goexit()
}

func handleDiscoveryResponse(selfAccount *account.Account, msg *event.Message) {
	content, err := message.NewChat().Message("Hello!").Finish()
	if err != nil {
		log.Fatal(err)
	}

	err = selfAccount.MessageSend(msg.FromAddress(), content)
	if err != nil {
		log.Fatal(err)
	}
}

func handleChat(msg *event.Message) {
	chat, err := message.DecodeChat(msg.Content())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("received message:", chat.Message())
}
