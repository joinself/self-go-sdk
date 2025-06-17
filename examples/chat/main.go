package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
)

var requests sync.Map

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

	keyPackage, err := selfAccount.ConnectionNegotiateOutOfBand(selfAccount.InboxDefault(), time.Now().Add(time.Minute*5))
	if err != nil {
		log.Fatal(err.Error())
	}

	content, err := message.NewDiscoveryRequest().KeyPackage(keyPackage).Expires(time.Now().Add(time.Minute * 5)).Finish()
	if err != nil {
		log.Fatal(err.Error())
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
		log.Println(err.Error())
		return
	}

	_, ok := requests.LoadAndDelete(hex.EncodeToString(discoveryResponse.ResponseTo()))
	if !ok {
		log.Println("received response to an unknown discovery request")
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
