// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"log"
	"os"

	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/authentication"
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

	log.Println("authenticating user")

	authService := client.AuthenticationService()

	err = authService.Request(authentication.AuthRequest{SelfID: os.Args[1]})
	if err != nil {
		log.Fatal("auth returned with: ", err)
	}

	log.Println("authentication succeeded")
}
