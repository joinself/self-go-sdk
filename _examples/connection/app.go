// Copyright 2020 Self Group Ltd. All Rights Reserved.

package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	selfsdk "github.com/joinself/self-go-sdk"
	"github.com/joinself/self-go-sdk/chat"
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

	err = client.Start()
	if err != nil {
		panic(err)
	}

	s := server{
		cid:  uuid.New().String(),
		chat: client.ChatService(),
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
	link, err := s.chat.GenerateConnectionDeepLink(chat.ConnectionConfig{
		Expiry: time.Minute * 5, // this is required ?
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("or click on " + link)

	log.Println("waiting for response")

	s.chat.OnConnection(func(iss, status string) {
		log.Println("Response received from " + iss + " with status " + status)
		parts := strings.Split(iss, ":")
		s.chat.Message([]string{parts[0]}, "Hi there!")
	})

	defer func() {
		err = client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(10 * time.Minute)
}

type server struct {
	cid  string
	chat *chat.Service
}

// serves the qr code image
func (s server) qrcode(w http.ResponseWriter, r *http.Request) {
	qrdata, _ := s.chat.GenerateConnectionQR(chat.ConnectionConfig{
		Expiry: time.Minute * 5, // this is required ?
	})

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
