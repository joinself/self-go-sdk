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
		LogLevel:    account.LogWarn,
		Callbacks: account.Callbacks{
			OnWelcome: account.DefaultWelcomeAccept,
			OnMessage: func(selfAccount *account.Account, msg *event.Message) {
				switch event.ContentTypeOf(msg) {
				case message.ContentTypeDiscoveryResponse:
					handleDiscoveryResponse(msg)
				default:
					log.Printf("received unhandled event")
				}
			},
		},
	}

	selfAccount, err := account.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	expires := time.Now().Add(time.Minute * 5)

	keyPackage, err := selfAccount.ConnectionNegotiateOutOfBand(selfAccount.InboxDefault(), expires)
	if err != nil {
		log.Fatal(err)
	}

	content, err := message.NewDiscoveryRequest().KeyPackage(keyPackage).Expires(expires).Finish()
	if err != nil {
		log.Fatal(err)
	}

	discoveryCompleter := make(chan *event.Message, 1)
	requests.Store(hex.EncodeToString(content.ID()), discoveryCompleter)

	qrCode, err := event.NewAnonymousMessage(content).SetFlags(event.MessageFlagTargetSandbox).EncodeToQR(event.QREncodingUnicode)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(qrCode))

	discoveryResponse := <-discoveryCompleter

	log.Printf("received discovery response from: %s", discoveryResponse.FromAddress().String())
}

func handleDiscoveryResponse(msg *event.Message) {
	discoveryResponse, err := message.DecodeDiscoveryResponse(msg.Content())
	if err != nil {
		log.Fatal(err)
	}

	completer, ok := requests.LoadAndDelete(hex.EncodeToString(discoveryResponse.ResponseTo()))
	if !ok {
		log.Fatal("received response to an unknown discovery request")
	}

	completer.(chan *event.Message) <- msg
}
