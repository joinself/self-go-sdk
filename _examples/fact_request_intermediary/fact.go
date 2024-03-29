// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"log"
	"os"
	"time"

	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/fact"
)

// expects 2 arguments - the Self ID you want to authenticate and the self ID of the intermediary to use
func main() {
	cfg := selfsdk.Config{
		SelfAppID:           os.Getenv("SELF_APP_ID"),
		SelfAppDeviceSecret: os.Getenv("SELF_APP_DEVICE_SECRET"),
		StorageKey:          "my-secret-crypto-storage-key",
		StorageDir:          "../.storage/",
	}

	if os.Getenv("SELF_ENV") != "" {
		cfg.Environment = os.Getenv("SELF_ENV")
	}

	client, err := selfsdk.New(cfg)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = client.Start()
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		panic("you must specify a self id as an argument")
	}

	var intermediary string

	if len(os.Args) > 3 {
		intermediary = os.Args[2]
	}

	log.Println("requesting user information by an intermediary")

	req := fact.IntermediaryFactRequest{
		SelfID:       os.Args[1],
		Intermediary: intermediary,
		Description:  "info",
		Facts: []fact.Fact{
			{
				Fact:          fact.FactEmailAddress,
				Sources:       []string{fact.SourceUserSpecified},
				Operator:      "==",
				ExpectedValue: "test@example.com",
			},
		},
		Expiry: time.Minute * 5,
	}

	factService := client.FactService()

	resp, err := factService.RequestViaIntermediary(&req)
	if err != nil {
		log.Fatal("fact request returned with: ", err)
	}

	for _, f := range resp.Facts {
		if f.Result() != true {
			log.Fatal("intermediary could not verify the required facts")
		}
		log.Printf("Your assertion that %s %s is %t\n", f.Fact, f.Operator, f.Result())
	}

	log.Println("intermediary verified the requried facts")
}
