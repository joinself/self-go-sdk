// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/fact"
)

// expects 1 argument - the Self ID you want to authenticate
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

	log.Println("requesting user information")

	req := fact.FactRequest{
		SelfID:      os.Args[1],
		Description: "info",
		Facts: []fact.Fact{
			{
				Fact:    fact.FactPhoto,
				Sources: []string{fact.SourcePassport},
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
		v := f.AttestedValues()[0]
		o := resp.Objects[v]

		ct, err := o.GetContent()
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("/tmp/output.jpg", ct, 0644)
		if err != nil {
			log.Fatal(err)
		}
		println("file stored on /tmp/output.jpg")
	}
}
