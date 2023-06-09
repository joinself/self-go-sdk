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

	err = client.MessagingService().PermitConnection("*")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		panic("you must specify a self id as an argument")
	}

	selfID := os.Args[1]
	cert := fact.Delegation{
		Subjects:  []string{selfID},
		Actions:   []string{"authenticate"},
		Effect:    "allow",
		Resources: []string{"resources:appID2:authentication"},
	}

	encodedCert, err := cert.Encode()
	if err != nil {
		panic(err)
	}

	source := "supu"

	g := fact.FactGroup{
		Name: "group name",
		Icon: "wifi_find",
	}
	f := []fact.FactToIssue{
		fact.FactToIssue{
			Key:    "cert",
			Value:  encodedCert,
			Source: source,
			Group:  &g,
			Type:   "delegation_certificate",
		},
	}

	client.FactService().Issue(selfID, f, []string{})
	time.Sleep(5 * time.Second)

	log.Println("requesting certificate")
	req := fact.FactRequest{
		SelfID:      selfID,
		Description: "We need access to your flight confirmation codes to reschedule your flights",
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

	for _, f := range resp.Facts {
		returnedCert, err := fact.ParseDelegationCertificate(f.AttestedValues()[0])
		if err != nil {
			panic(err)
		}

		log.Println(f.Fact, ":", returnedCert)
	}
}
