// Package main demonstrates discovery subscription capabilities using the Self SDK.
//
// This example showcases the power of subscription-based peer discovery,
// where a single client can generate multiple QR codes and receive real-time
// notifications as different peers discover and connect to it.
//
// 🎯 What you'll learn:
// • How to set up discovery subscription handlers
// • Generating multiple QR codes for discovery
// • Real-time peer discovery notifications
// • Understanding the discovery workflow
//
// 🔄 Discovery Subscription Flow:
// 1. Client generates multiple QR codes
// 2. Each QR code can be scanned by different peers
// 3. Client receives real-time notifications for each discovery
// 4. Multiple peers can discover simultaneously
//
// 📱 Try this demo:
// • Run this program to generate QR codes
// • Use multiple Self clients to scan different QR codes
// • Watch real-time discovery notifications
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("🔍 Discovery Subscription Demo")
	fmt.Println("==============================")
	fmt.Println("This demo shows how to use subscription-based peer discovery.")
	fmt.Println("Generate multiple QR codes and get notified when peers discover you!")
	fmt.Println()

	// Step 1: Create and configure the discovery client
	discoveryClient := createDiscoveryClient()
	defer discoveryClient.Close()

	fmt.Printf("🆔 My DID: %s\n", discoveryClient.DID())
	fmt.Println()

	// Step 2: Set up discovery subscription handler
	setupDiscoveryHandler(discoveryClient)

	// Step 3: Generate multiple QR codes for discovery
	generateDiscoveryQRCodes(discoveryClient)

	// Step 4: Keep listening for discovery responses
	keepListening()
}

// createDiscoveryClient sets up a Self client configured for discovery
func createDiscoveryClient() *client.Client {
	fmt.Println("🔧 Setting up discovery client...")

	discoveryClient, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("discovery_demo"),
		StoragePath: "./discovery_demo_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create discovery client:", err)
	}

	fmt.Println("✅ Discovery client created successfully")
	return discoveryClient
}

// setupDiscoveryHandler configures the subscription-based discovery response handler
func setupDiscoveryHandler(discoveryClient *client.Client) {
	fmt.Println("🔧 Setting up discovery subscription handler...")

	// This is the key to subscription-based discovery!
	// The handler will be called for EVERY peer that discovers us
	discoveryClient.Discovery().OnResponse(func(peer *client.Peer) {
		fmt.Printf("\n🎉 NEW PEER DISCOVERED!\n")
		fmt.Printf("   DID: %s\n", peer.DID())
		fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05"))
		fmt.Printf("   Status: Ready for communication\n")

		// In a real application, you could:
		// - Store the peer information
		// - Initiate a chat session
		// - Request credentials
		// - Send a welcome message
		fmt.Printf("   💡 You can now communicate with this peer!\n")
		fmt.Println("   ────────────────────────────────────────")
	})

	fmt.Println("✅ Discovery handler configured")
	fmt.Println("   📡 Now listening for peer discoveries...")
	fmt.Println()
}

// generateDiscoveryQRCodes creates multiple QR codes to demonstrate subscription
func generateDiscoveryQRCodes(discoveryClient *client.Client) {
	fmt.Println("📱 Generating QR codes for discovery...")
	fmt.Println("Each QR code can be scanned by different peers.")
	fmt.Println("You'll receive notifications for each discovery!")
	fmt.Println()

	// Generate multiple QR codes to show subscription capabilities
	qrCodes := []struct {
		name    string
		timeout time.Duration
	}{
		{"Quick Discovery", 15 * time.Minute},
		{"Standard Discovery", 30 * time.Minute},
		{"Extended Discovery", 60 * time.Minute},
	}

	for i, qrConfig := range qrCodes {
		fmt.Printf("🔄 Generating %s QR code...\n", qrConfig.name)

		qr, err := discoveryClient.Discovery().GenerateQRWithTimeout(qrConfig.timeout)
		if err != nil {
			log.Printf("❌ Failed to generate QR code %d: %v", i+1, err)
			continue
		}

		fmt.Printf("\n--- %s (QR #%d) ---\n", qrConfig.name, i+1)
		fmt.Printf("Request ID: %s\n", qr.RequestID())
		fmt.Printf("Valid for: %v\n", qrConfig.timeout)

		qrCode, err := qr.Unicode()
		if err != nil {
			log.Printf("❌ Failed to get QR code %d: %v", i+1, err)
			continue
		}
		fmt.Println(qrCode)
		fmt.Println()
	}

	fmt.Println("✅ All QR codes generated successfully!")
	fmt.Println()
	fmt.Println("🎓 What's happening:")
	fmt.Println("   1. Three QR codes with different timeouts are active")
	fmt.Println("   2. Each can be scanned by different Self clients")
	fmt.Println("   3. You'll get real-time notifications for each discovery")
	fmt.Println("   4. Multiple peers can discover you simultaneously")
	fmt.Println()
}

// keepListening maintains the program to receive discovery responses
func keepListening() {
	fmt.Println("🔍 Listening for discovery responses...")
	fmt.Println("📱 Scan any QR code above with Self clients to see subscription in action!")
	fmt.Println()
	fmt.Println("💡 Try this:")
	fmt.Println("   • Use multiple devices to scan different QR codes")
	fmt.Println("   • Notice how each discovery triggers a separate notification")
	fmt.Println("   • See the real-time nature of the subscription system")
	fmt.Println()
	fmt.Println("Press Ctrl+C to exit.")

	// Keep the program running to receive discovery responses
	// This demonstrates the subscription nature - the client stays active
	// and receives notifications as they happen
	select {}
}
