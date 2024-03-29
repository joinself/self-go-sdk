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
	"github.com/joinself/self-go-sdk/authentication"
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
		auth: client.AuthenticationService(),
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

	resp, err := s.auth.WaitForResponse(s.cid, time.Minute)
	if err != nil {
		log.Fatal("auth returned with: ", err)
	}

	if !resp.Accepted {
		log.Fatal("auth rejected:", resp.SelfID)
	}

	log.Println("authentication succeeded:", resp.SelfID)
}

type server struct {
	cid  string
	auth *authentication.Service
}

// serves the qr code image
func (s server) qrcode(w http.ResponseWriter, r *http.Request) {
	req := authentication.QRAuthenticationRequest{
		ConversationID: s.cid,
		Expiry:         time.Minute * 5,
		QRConfig: fact.QRConfig{
			Size:            400,
			BackgroundColor: "#FFFFFF",
			ForegroundColor: "#000000",
		},
	}

	qrdata, err := s.auth.GenerateQRCode(&req)
	if err != nil {
		log.Fatal(err)
	}

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
