// Package main demonstrates simple credential exchange using the Self SDK.
//
// This is the BASIC level of credential exchange examples.
// Start here if you're new to credential exchange concepts.
//
// This example shows the basics of:
// - Setting up two clients (issuer and holder)
// - Connecting clients through discovery
// - Creating a simple credential
// - Requesting and responding to credential exchanges
// - Understanding the exchange workflow
//
// 🎯 What you'll learn:
// • How credential exchange works between two parties
// • Basic request/response patterns
// • Simple credential creation and sharing
// • Client connection establishment
//
// 📚 Next steps:
// • multi_credential_exchange.go - Multiple credential types
// • advanced_exchange.go - Complex parameters and verification
// • discovery_exchange.go - QR code discovery integration
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/examples/utils"
	"github.com/joinself/self-go-sdk/message"
)

func main() {
	fmt.Println("🔄 Basic Credential Exchange Demo")
	fmt.Println("==================================")
	fmt.Println("This demo shows basic credential exchange between two parties.")
	fmt.Println("📚 This is the BASIC level - start here if you're new to credential exchange.")
	fmt.Println()

	// Step 1: Create two clients - one issuer, one holder
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Connect the clients through discovery
	connectClients(issuer, holder)

	// Step 3: Create a simple credential
	createSampleCredential(issuer, holder)

	// Step 4: Set up handlers for credential requests
	setupHandlers(issuer, holder)

	// Step 5: Demonstrate credential exchange
	demonstrateExchange(issuer, holder)

	fmt.Println("✅ Basic demo completed!")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run multi_credential_exchange.go to learn about multiple credential types")
	fmt.Println("   • Run advanced_exchange.go for complex parameters and verification")
	fmt.Println("   • Run discovery_exchange.go for QR code discovery integration")
	fmt.Println()
	fmt.Println("The clients will keep running to show ongoing exchange capabilities.")
	fmt.Println("Press Ctrl+C to exit.")

	// Keep running to demonstrate exchange capabilities
	select {}
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	// Create issuer client
	issuer, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("basic_issuer"),
		StoragePath: "./basic_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	// Create holder client
	holder, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("basic_holder"),
		StoragePath: "./basic_holder_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create holder:", err)
	}

	fmt.Println("✅ Clients created successfully")
	return issuer, holder
}

// connectClients establishes a connection between the two clients programmatically
func connectClients(issuer, holder *client.Client) {
	fmt.Println("🔗 Connecting clients programmatically...")
	fmt.Println("   This establishes a secure connection without QR codes")
	fmt.Println("   Using the new Connection component for easy peer-to-peer connectivity")

	fmt.Println("   📡 Initiating connection negotiation...")

	// Use the new Connection component to establish connection
	err := client.ConnectTwoClientsWithTimeout(issuer, holder, 10*time.Second)
	if err != nil {
		fmt.Printf("   ❌ Connection failed: %v\n", err)
		fmt.Println("   💡 This may happen in demo environments")
		fmt.Println("   🔗 In production, ensure both clients are connected to the messaging service")
		return
	}

	fmt.Println("   ✅ Connection established successfully!")
	fmt.Println("   🔐 Clients can now exchange messages securely")
	fmt.Println("   🎉 Ready for credential exchange!")

	// Give a moment for the connection to fully establish
	time.Sleep(1 * time.Second)
	fmt.Println()
}

// createSampleCredential creates a simple email credential for demonstration
func createSampleCredential(issuer, holder *client.Client) {
	fmt.Println("📧 Creating sample email credential...")

	// Create a simple email credential
	_, err := issuer.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeEmail).
		Subject(holder.DID()).
		Issuer(issuer.DID()).
		Claim("emailAddress", "demo@example.com").
		Claim("verified", true).
		ValidFrom(time.Now()).
		SignWith(issuer.DID(), time.Now()).
		Issue(issuer)

	if err != nil {
		log.Printf("Failed to create credential: %v", err)
		return
	}

	fmt.Println("✅ Sample credential created: demo@example.com")
	fmt.Println()
}

// setupHandlers configures how clients respond to credential requests
func setupHandlers(issuer, holder *client.Client) {
	fmt.Println("🔧 Setting up exchange handlers...")

	// When someone asks the holder for credentials
	holder.Credentials().OnPresentationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("📨 Holder received request from: %s\n", req.From())
		fmt.Printf("   Requested: %v\n", req.Type())

		// For this demo, we'll reject the request
		// In a real app, you'd check if you have the credential and respond accordingly
		fmt.Println("   ❌ Rejecting request (demo)")
		req.Reject()
	})

	// When the issuer gets a response to their request
	issuer.Credentials().OnPresentationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("📨 Issuer received response from: %s\n", resp.From())
		fmt.Printf("   Status: %s\n", utils.ResponseStatusToString(resp.Status()))
		fmt.Printf("   Presentations: %d\n", len(resp.Presentations()))
	})

	fmt.Println("✅ Handlers configured")
	fmt.Println()
}

// demonstrateExchange shows a simple credential exchange request
func demonstrateExchange(issuer, holder *client.Client) {
	fmt.Println("🔄 Demonstrating credential exchange...")
	fmt.Println("   ⚠️  Note: This demo shows the request/response pattern")
	fmt.Println("   💡 In a real scenario, clients would be connected via QR code discovery")

	// Create a simple request for email credentials
	details := []*client.CredentialDetail{
		{
			CredentialType: credential.CredentialTypeEmail,
			Parameters: []*client.CredentialParameter{
				{
					Operator: message.OperatorNotEquals,
					Field:    "emailAddress",
					Value:    "", // Looking for any non-empty email
				},
			},
		},
	}

	fmt.Println("📤 Issuer requesting email credential from holder...")

	// Send the request
	req, err := issuer.Credentials().RequestPresentationWithTimeout(
		holder.DID(),
		details,
		10*time.Second,
	)
	if err != nil {
		fmt.Printf("   ❌ Request failed: %v\n", err)
		fmt.Println("   💡 This is expected in the demo since clients aren't actually connected")
		fmt.Println("   🔗 In a real app, ensure clients are connected via discovery first")
		fmt.Println()

		// Show what would happen in a real scenario
		demonstrateHypotheticalExchange()
		return
	}

	fmt.Printf("   Request sent with ID: %s\n", req.RequestID())

	// Wait for response
	fmt.Println("⏳ Waiting for response...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := req.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("   ⏰ Request timed out (expected in demo)")
		} else {
			fmt.Printf("   ❌ Request failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Response received: %s\n", utils.ResponseStatusToString(resp.Status()))
	}

	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Issuer requested email credential from holder")
	fmt.Println("   2. Holder received the request and rejected it (demo)")
	fmt.Println("   3. Issuer received the rejection response")
	fmt.Println("   4. In real scenarios, holder would share actual credentials")
	fmt.Println()
}

// demonstrateHypotheticalExchange shows what would happen with connected clients
func demonstrateHypotheticalExchange() {
	fmt.Println("📚 WHAT WOULD HAPPEN WITH CONNECTED CLIENTS:")
	fmt.Println("============================================")
	fmt.Println("🔗 If clients were properly connected via discovery:")
	fmt.Println()
	fmt.Println("   1. 📱 Holder scans issuer's QR code")
	fmt.Println("   2. 🔐 Secure encrypted connection established")
	fmt.Println("   3. 📤 Issuer sends credential request")
	fmt.Println("   4. 📨 Holder receives request instantly")
	fmt.Println("   5. 📋 Holder checks available credentials")
	fmt.Println("   6. ✅ Holder responds with matching credentials")
	fmt.Println("   7. 🎉 Issuer receives and validates credentials")
	fmt.Println()
	fmt.Println("🔧 To see this in action:")
	fmt.Println("   • Run the discovery_exchange.go example")
	fmt.Println("   • Use two separate devices/terminals")
	fmt.Println("   • Scan QR codes to establish real connections")
	fmt.Println()
}
