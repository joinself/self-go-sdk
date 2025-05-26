package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
)

func main() {
	// Create a new Self client
	selfClient, err := client.NewClient(client.Config{
		StorageKey:  make([]byte, 32), // In production, use a secure key
		StoragePath: "./discovery_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer selfClient.Close()

	fmt.Printf("My DID: %s\n", selfClient.DID())

	// Set up discovery response handler for subscription-based discovery
	selfClient.Discovery().OnResponse(func(peer *client.Peer) {
		fmt.Printf("🎉 New peer discovered: %s\n", peer.DID())
		fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05"))

		// You could initiate chat or other interactions here
		fmt.Println("   Ready to start communication with this peer!")
	})

	// Generate multiple QR codes to demonstrate subscription
	fmt.Println("\n📱 Generating QR codes for discovery...")
	fmt.Println("Multiple people can scan these codes and you'll be notified of each connection.")

	for i := 1; i <= 3; i++ {
		qr, err := selfClient.Discovery().GenerateQRWithTimeout(30 * time.Minute)
		if err != nil {
			log.Printf("Failed to generate QR code %d: %v", i, err)
			continue
		}

		fmt.Printf("\n--- QR Code #%d (Request ID: %s) ---\n", i, qr.RequestID())
		qrCode, err := qr.Unicode()
		if err != nil {
			log.Printf("Failed to get QR code %d: %v", i, err)
			continue
		}
		fmt.Println(qrCode)
	}

	fmt.Println("\n🔍 Listening for discovery responses...")
	fmt.Println("Scan any of the QR codes above with Self clients to see subscription in action.")
	fmt.Println("Press Ctrl+C to exit.")

	// Keep the program running to receive discovery responses
	select {}
}
