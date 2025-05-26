// Package main demonstrates simple credential exchange using the Self SDK.
//
// This is the BASIC level of credential exchange examples.
// Start here if you're new to credential exchange concepts.
//
// This example shows the basics of:
// - Setting up two clients (issuer and holder)
// - Creating a simple credential
// - Requesting and responding to credential exchanges
// - Understanding the exchange workflow
//
// 🎯 What you'll learn:
// • How credential exchange works between two parties
// • Basic request/response patterns
// • Simple credential creation and sharing
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

	// Step 2: Create a simple credential
	createSampleCredential(issuer, holder)

	// Step 3: Set up handlers for credential requests
	setupHandlers(issuer, holder)

	// Step 4: Demonstrate credential exchange
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
		log.Printf("Failed to send request: %v", err)
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
