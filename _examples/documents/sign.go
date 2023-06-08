// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/documents"
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

	defer client.Close()

	err = client.Start()
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		panic("you must specify a self id as an argument")
	}

	ds := client.DocsService()

	content, err := ioutil.ReadFile("./sample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	objects := make([]documents.InputObject, 0)
	objects = append(objects, documents.InputObject{
		Name: "Terms and conditions",
		Data: content,
		Mime: "application/pdf",
	})

	log.Println("sending document sign request")
	resp, err := ds.RequestSignature(os.Args[1], "Read and sign this documents", objects)
	if err != nil {
		log.Println(err.Error())
	}

	if resp.Status == "accepted" {
		fmt.Println("Document has been signed")
		fmt.Println("")
		fmt.Println("signed documents:")
		spew.Dump(resp.SignedObjects)
		for _, o := range resp.SignedObjects {
			fmt.Println("- Name: " + o.Name)
			fmt.Println("  Link: " + o.Link)
			fmt.Println("  Hash: " + o.Hash)
		}
		fmt.Println("")
		fmt.Println("full signature:")
		fmt.Println(resp.Signature)
	} else {
		fmt.Println("Document signature has been rejected")
	}

}
