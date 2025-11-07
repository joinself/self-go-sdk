package main

import (
	"encoding/base64"
	"log"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/pairwise"
)

var (
	users    map[int]*user
	requests map[string]int
	mu       sync.Mutex
	encode   = base64.RawStdEncoding.EncodeToString
)

type user struct {
	ID        int
	Reference *pairwise.Identity
}

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
				case message.ContentTypeIntroduction:
					handleIntroduction(selfAccount, msg)
					requestUserAuthentication(0)
				case message.ContentTypeCredentialPresentationResponse:
					handleCredentialPresentationResponse(selfAccount, msg)
				default:
					log.Printf("received unhandled event")
				}
			},
		},
	}

	_, err := account.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("waiting for user registration")

	runtime.Goexit()
}

func handleIntroduction(selfAccount *account.Account, msg *event.Message) {
	introduction, err := message.DecodeIntroduction(msg.Content())
	if err != nil {
		log.Fatal("failed to decode introduction response", "error", err)
	}

	pairwiseIdentity, err := selfAccount.ConnectionPairwiseIntroduction(
		introduction.DocumentAddress(),
		introduction.Presentations(),
	)
	if err != nil {
		log.Fatal("failed to process introduction from new user")
	}

	mu.Lock()
	id := len(users)

	users[id] = &user{
		ID:        id,
		Reference: pairwiseIdentity,
	}

	mu.Unlock()

	log.Println("registered new user")
}

func handleCredentialPresentationResponse(selfAccount *account.Account, msg *event.Message) {
	credentialPresentationResponse, err := message.DecodeCredentialPresentationResponse(msg.Content())
	if err != nil {
		log.Fatal("failed to decode credential presentation response", "error", err)
	}

	mu.Lock()

	id, ok := requests[encode(msg.ID())]
	if !ok {
		mu.Unlock()
		return
	}

	user := users[id]

	mu.Unlock()

	userCredentials, err := selfAccount.CredentialGraphValidFor(
		user.Reference.DocumentAddress(),
		credential.SandboxTrustedIssuerRegistry(),
		credentialPresentationResponse.Presentations(),
	)

	if err != nil {
		log.Fatal("failed to validate credential presentation response credentials", "error", err)
	}

	var authenticated bool

	for _, c := range userCredentials {
		if slices.Contains(c.CredentialType(), credential.CredentialTypeLivenessAndFacialComparison) {
			continue
		}

		// TODO implement challenge and better APIS for checking

		authenticated = true
	}

	if !authenticated {
		log.Fatal("user not authenticated")
	}
}

func requestUserAuthentication(id int) {
	mu.Lock()
	user := users[id]
	mu.Unlock()

	content, err := message.NewCredentialPresentationRequest().
		PresentationType(credential.PresentationTypeLivenessAndFacialComparison).
		Holder(user.Reference.DocumentAddress()).
		BiometricAnchor(user.Reference.BiometricAnchor()).
		Expires(time.Now().Add(time.Minute * 5)).
		Finish()

	if err != nil {
		log.Fatal("failed to create credential presentation request", "error", err)
	}

	mu.Lock()
	requests[encode(content.ID())] = user.ID
	mu.Unlock()
}
