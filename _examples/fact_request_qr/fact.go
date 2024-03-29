// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/google/uuid"
	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/fact"
)

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
		log.Fatal(err)
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

	s := server{
		cid:  uuid.New().String(),
		fact: client.FactService(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/qr.png", s.qrcode)

	log.Println("starting server")

	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := http.Serve(l, mux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	openbrowser("http://localhost:9999/qr.png")

	log.Println("waiting for response")

	resp, err := s.fact.WaitForResponse(s.cid, time.Minute)
	if err != nil {
		log.Fatal("fact request returned with: ", err)
	}

	for _, f := range resp.Facts {
		log.Println(f.Fact, ":", f.AttestedValues())
	}
}

type server struct {
	cid  string
	fact *fact.Service
}

// serves the qr code image
func (s server) qrcode(w http.ResponseWriter, r *http.Request) {
	req := fact.QRFactRequest{
		ConversationID: s.cid,  // this is required. session ID maybe more sensible?
		Description:    "Info", // this is optional
		Facts: []fact.Fact{
			{
				Fact:    fact.FactEmailAddress,              // required
				Sources: []string{fact.SourceUserSpecified}, // required
			},
		},
		Expiry: time.Minute * 5, // this is required ?
		QRConfig: fact.QRConfig{
			Size:            400,       // this is optional/defaulted
			BackgroundColor: "#FFFFFF", // this is optional/defaulted
			ForegroundColor: "#000000", // this is optional/defaulted
		},
	}

	qrdata, _ := s.fact.GenerateQRCode(&req)

	w.Write(qrdata)
}

// ignore this stuff
func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Fatal("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
