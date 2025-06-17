package main

import (
	"encoding/hex"
	"fmt"
	"log"
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
					handleDiscoveryResponse(msg)
				case message.ContentTypeIntroduction:
					handleIntroduction(selfAccount, msg)
				case message.ContentTypeChat:
					handleChat(msg)
				}
			},
		},
	}

	selfAccount, err := account.New(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		keyPackage, err := selfAccount.ConnectionNegotiateOutOfBand(selfAccount.InboxDefault(), time.Now().Add(time.Minute*5))
		if err != nil {
			log.Fatal(err.Error())
		}

		content, err := message.NewDiscoveryRequest().KeyPackage(keyPackage).Expires(time.Now().Add(time.Minute * 5)).Finish()
		if err != nil {
			log.Fatal(err.Error())
		}

		completer := make(chan *event.Message, 1)

		requests.Store(hex.EncodeToString(content.ID()), completer)

		qrCode, err := event.NewAnonymousMessage(content).SetFlags(event.MessageFlagTargetSandbox).EncodeToQR(event.QREncodingUnicode)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(string(qrCode))

		response := <-completer

		content, err = message.NewChat().Message("Hello!").Finish()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = selfAccount.MessageSend(response.FromAddress(), content)
		if err != nil {
			log.Fatal(err.Error())
		}

		summary, err := content.Summary()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = selfAccount.NotificationSend(response.FromAddress(), summary)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func handleDiscoveryResponse(msg *event.Message) {
	discoveryResponse, err := message.DecodeDiscoveryResponse(msg.Content())
	if err != nil {
		log.Println(err.Error())
		return
	}

	completer, ok := requests.LoadAndDelete(hex.EncodeToString(discoveryResponse.ResponseTo()))
	if !ok {
		log.Println("received response to an unknown discovery request")
		return
	}

	completer.(chan *event.Message) <- msg
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
