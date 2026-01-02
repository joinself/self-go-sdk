package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"runtime"
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
	requests map[string]*request
	mu       sync.Mutex
	encode   = base64.RawStdEncoding.EncodeToString
)

type user struct {
	ID        int
	Reference *pairwise.Identity
}

type request struct {
	UserID    int
	Challenge []byte
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
					userID := handleIntroduction(selfAccount, msg)
					requestUserAuthentication(userID)
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

func handleIntroduction(selfAccount *account.Account, msg *event.Message) int {
	introduction, err := message.DecodeIntroduction(msg.Content())
	if err != nil {
		log.Fatal("failed to decode introduction response", "error", err)
	}

	pairwiseIntroduction, err := introduction.Introduction()
	if err != nil {
		log.Fatal("invalid pairwise introduction", "error", err)
	}

	pairwiseIdentity, err := selfAccount.ConnectionPairwiseIntroductionValidate(
		msg.FromAddress(),
		pairwiseIntroduction,
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

	return id
}

func handleCredentialPresentationResponse(selfAccount *account.Account, msg *event.Message) {
	credentialPresentationResponse, err := message.DecodeCredentialPresentationResponse(msg.Content())
	if err != nil {
		log.Fatal("failed to decode credential presentation response", "error", err)
	}

	mu.Lock()

	request, ok := requests[encode(msg.ID())]
	if !ok {
		mu.Unlock()
		return
	}

	user := users[request.UserID]

	mu.Unlock()

	graph, err := selfAccount.CredentialGraphCreate(
		credential.SandboxTrustedIssuerRegistry(),
		credentialPresentationResponse.Presentations(),
	)

	if err != nil {
		log.Fatal("failed to validate credential presentation response credentials", "error", err)
	}

	if !graph.ValidAuthenticationFor(user.Reference, request.Challenge) {
		log.Fatal("user not authenticated")
	}

	log.Println("user authenticated successfully")
}

func requestUserAuthentication(id int) {
	mu.Lock()
	user := users[id]
	mu.Unlock()

	challenge := make([]byte, 32)
	rand.Read(challenge)

	content, err := message.NewCredentialPresentationRequest().
		PresentationType(credential.PresentationTypeLivenessAndFacialComparison).
		Authenticate(user.Reference, challenge).
		Expires(time.Now().Add(time.Minute * 5)).
		Finish()

	if err != nil {
		log.Fatal("failed to create credential presentation request", "error", err)
	}

	mu.Lock()
	requests[encode(content.ID())] = &request{
		UserID:    user.ID,
		Challenge: challenge,
	}
	mu.Unlock()
}
