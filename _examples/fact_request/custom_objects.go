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

	log.Println("issuing custom facts with objects attached")
	selfid := os.Args[1]
	source := "Custom Objects"
	content, err := ioutil.ReadFile("./my_image.png")
	if err != nil {
		log.Fatal(err)
	}
	obj, err := client.NewObject(content, "test", "image/png")
	if err != nil {
		log.Fatal(err)
	}

	f := []fact.FactToIssue{
		fact.FactToIssue{
			Key:    "my-image",
			Value:  obj.Digest,
			Source: source,
			Object: obj,
		},
	}

	client.FactService().Issue(selfid, f, []string{})
	time.Sleep(5 * time.Second)

	log.Println("requesting custom facts with objects")
	req := fact.FactRequest{
		SelfID:      selfid,
		Description: "This is a sample of issuing custom objects",
		Facts: []fact.Fact{
			{
				Fact:    f[0].Key,
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

	digests := []string{}
	for _, f := range resp.Facts {
		log.Println(f.Fact, ":", f.AttestedValues())
		digests = append(digests, f.AttestedValues()...)
	}

	for _, d := range digests {
		if _, ok := resp.Objects[d]; ok {
			ct, err := resp.Objects[d].GetContent()
			if err != nil {
				log.Fatal(err)
			}
			err = ioutil.WriteFile("/tmp/output-go.jpg", ct, 0644)
			if err != nil {
				log.Fatal(err)
			}
			println("file stored on /tmp/output-go.jpg")
		}
	}

}
