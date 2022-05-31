// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"log"
	"os"
	"time"

	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/fact"
)

// expects 1 argument - the Self ID you want to authenticate
func main() {
	appID := os.Getenv("SELF_APP_ID")
	cfg := selfsdk.Config{
		SelfAppID:           appID,
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

	err = client.MessagingService().PermitConnection("*")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		panic("you must specify a self id as an argument")
	}

	log.Println("issuing custom facts")
	selfid := os.Args[1]
	source := "supu"
	factName := "fooGo"

	f := fact.FactToIssue{
		Key:   factName,
		Value: "bar",
	}
	client.FactService().Issue(selfid, source, []fact.FactToIssue{f}, []string{})
	time.Sleep(5 * time.Second)

	log.Println("requesting custom facts")
	req := fact.FactRequest{
		SelfID:      selfid,
		Description: "info",
		Facts: []fact.Fact{
			{
				Fact:    factName,
				Sources: []string{source},
				Issuers: []string{appID},
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

	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
