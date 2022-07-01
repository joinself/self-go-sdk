// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"fmt"
	"log"
	"os"

	selfsdk "github.com/joinself/self-go-sdk"
)

// expects 1 argument - the Self ID you want to lookup
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

	if len(os.Args) < 2 {
		panic("you must specify a self id as an argument")
	}

	log.Println("looking up identity")

	identityService := client.IdentityService()

	identity, err := identityService.GetIdentity(os.Args[1])
	if err != nil {
		log.Fatal("identity lookup returned with: ", err)
	}

	log.Println("identity lookup succeeded")

	for _, o := range identity.History {
		fmt.Println(string(o))
	}

	fmt.Println(identity.Proofs)

	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
