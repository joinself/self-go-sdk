// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/documents"
	"github.com/joinself/self-go-sdk/messaging"
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

	client.MessagingService().Subscribe("document.sign.resp", func(m *messaging.Message) {
		var payload map[string]interface{}
		err := json.Unmarshal(m.Payload, &payload)
		if err != nil {
			log.Printf("failed to decode message payload: %s", err.Error())
			return
		}

		var resp documents.Response
		err = json.Unmarshal(m.Payload, &resp)
		if err != nil {
			log.Printf("failed unmarshalling response: %s", err.Error())
			return
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
			fmt.Println(m.Signature)
		} else {
			fmt.Println("Document signature has been rejected")
		}
	})

	log.Println("sending document sign request")
	cid := uuid.New().String()
	err = ds.RequestSignatureAsync(cid, os.Args[1], "Read and sign this documents", objects)
	if err != nil {
		log.Println(err.Error())
	}

	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Notify the channel on SIGINT (Ctrl+C) and SIGTERM
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to block the main goroutine until a signal is received
	done := make(chan bool, 1)

	// Goroutine to handle received signals
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("Received signal:", sig)
		done <- true
	}()

	fmt.Println("Press Ctrl+C to exit")

	// Wait for a signal to be handled
	<-done

	fmt.Println("Exiting program")

}
