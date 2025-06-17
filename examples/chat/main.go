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

func handleDiscoveryResponse(selfAccount *account.Account, msg *event.Message) {
	_, err := message.DecodeDiscoveryResponse(msg.Content())
	if err != nil {
		log.Println(err.Error())
		return
	}

	sendMessage(selfAccount, msg)
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

func handleChat(msg *event.Message) {
	chat, err := message.DecodeChat(msg.Content())
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("received message:", chat.Message())
}

func sendMessage(selfAccount *account.Account, msg *event.Message) {
	content, err := message.NewChat().Message("Hello!").Finish()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = selfAccount.MessageSend(msg.FromAddress(), content)
	if err != nil {
		log.Fatal(err.Error())
	}

	summary, err := content.Summary()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = selfAccount.NotificationSend(msg.FromAddress(), summary)
	if err != nil {
		log.Fatal(err.Error())
	}
}
