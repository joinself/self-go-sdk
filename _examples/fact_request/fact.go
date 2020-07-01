package main

import (
	"log"
	"os"
	"time"

	selfsdk "github.com/selfid-net/self-go-sdk"
	"github.com/selfid-net/self-go-sdk/fact"
)

// expects 1 argument - the Self ID you want to authenticate
func main() {
	cfg := selfsdk.Config{
		SelfAppID:     os.Getenv("SELF_APP_ID"),
		SelfAppSecret: os.Getenv("SELF_APP_SECRET"),
		StorageKey:    "my-secret-crypto-storage-key",
	}

	if os.Getenv("SELF_ENV") != "" {
		cfg.Environment = os.Getenv("SELF_ENV")
	}

	client, err := selfsdk.New(cfg)
	if err != nil {
		panic(err)
	}

	err = client.MessagingService().PermitConnection("*")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		panic("you must specify a self id as an argument")
	}

	log.Println("requesting user information")

	req := fact.FactRequest{
		SelfID:      os.Args[1],
		Description: "info",
		Facts: []fact.Fact{
			{
				Fact:    fact.FactPhone,
				Sources: []string{fact.SourceUserSpecified},
			},
		},
		Expiry: time.Minute * 5,
	}

	factService := client.FactService()

	resp, err := factService.Request(&req)
	if err != nil {
		log.Fatal("fact request returned with: ", err)
	}

	for _, f := range resp.Facts {
		log.Println(f.Fact, ":", f.AttestedValues())
	}
}