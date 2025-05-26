// Package main demonstrates multi-credential exchange using the Self SDK.
//
// This is the INTERMEDIATE level of credential exchange examples.
// Complete basic_exchange.go first before trying this example.
//
// This example shows:
// - Working with multiple credential types (email, profile, education)
// - Creating different types of credentials
// - Requesting multiple credentials in one exchange
// - Processing complex credential responses
//
// 🎯 What you'll learn:
// • How to handle multiple credential types
// • Creating credentials with different claim structures
// • Multi-credential request patterns
// • Processing complex responses
//
// 📚 Prerequisites: basic_exchange.go
// 📚 Next steps:
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
	fmt.Println("🔄 Multi-Credential Exchange Demo")
	fmt.Println("==================================")
	fmt.Println("This demo shows credential exchange with multiple credential types.")
	fmt.Println("📚 This is the INTERMEDIATE level - complete basic_exchange.go first.")
	fmt.Println()

	// Step 1: Create clients
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Create multiple types of credentials
	createMultipleCredentials(issuer, holder)

	// Step 3: Set up handlers for multi-credential requests
	setupMultiHandlers(issuer, holder)

	// Step 4: Demonstrate multi-credential exchange
	demonstrateMultiExchange(issuer, holder)

	fmt.Println("✅ Multi-credential demo completed!")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run advanced_exchange.go for complex parameters and verification")
	fmt.Println("   • Run discovery_exchange.go for QR code discovery integration")
	fmt.Println()
	fmt.Println("The clients will keep running. Press Ctrl+C to exit.")

	select {}
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	issuer, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("multi_issuer"),
		StoragePath: "./multi_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	holder, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("multi_holder"),
		StoragePath: "./multi_holder_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create holder:", err)
	}

	fmt.Println("✅ Clients created successfully")
	return issuer, holder
}

// createMultipleCredentials creates different types of credentials for demonstration
func createMultipleCredentials(issuer, holder *client.Client) {
	fmt.Println("📝 Creating multiple types of credentials...")

	// 1. Email credential
	fmt.Println("📧 Creating email credential...")
	_, err := issuer.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeEmail).
		Subject(holder.DID()).
		Issuer(issuer.DID()).
		Claim("emailAddress", "alice@example.com").
		Claim("verified", true).
		Claim("verificationDate", time.Now().Format("2006-01-02")).
		ValidFrom(time.Now()).
		SignWith(issuer.DID(), time.Now()).
		Issue(issuer)

	if err != nil {
		log.Printf("Failed to create email credential: %v", err)
	} else {
		fmt.Println("   ✅ Email credential created: alice@example.com")
	}

	// 2. Profile credential
	fmt.Println("👤 Creating profile credential...")
	_, err = issuer.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeProfileName).
		Subject(holder.DID()).
		Issuer(issuer.DID()).
		Claim("firstName", "Alice").
		Claim("lastName", "Johnson").
		Claim("displayName", "Alice Johnson").
		Claim("country", "Canada").
		ValidFrom(time.Now()).
		SignWith(issuer.DID(), time.Now()).
		Issue(issuer)

	if err != nil {
		log.Printf("Failed to create profile credential: %v", err)
	} else {
		fmt.Println("   ✅ Profile credential created: Alice Johnson")
	}

	// 3. Education credential (custom type)
	fmt.Println("🎓 Creating education credential...")
	_, err = issuer.Credentials().NewCredentialBuilder().
		Type([]string{"VerifiableCredential", "EducationCredential"}).
		Subject(holder.DID()).
		Issuer(issuer.DID()).
		Claim("degree", "Bachelor of Computer Science").
		Claim("institution", "Tech University").
		Claim("graduationYear", 2023).
		Claim("gpa", 3.8).
		ValidFrom(time.Now()).
		SignWith(issuer.DID(), time.Now()).
		Issue(issuer)

	if err != nil {
		log.Printf("Failed to create education credential: %v", err)
	} else {
		fmt.Println("   ✅ Education credential created: Bachelor of Computer Science")
	}

	fmt.Println("✅ All credentials created successfully")
	fmt.Println()
}

// setupMultiHandlers configures handlers for multi-credential requests
func setupMultiHandlers(issuer, holder *client.Client) {
	fmt.Println("🔧 Setting up multi-credential handlers...")

	// Holder responds to multi-credential requests
	holder.Credentials().OnPresentationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("📨 Holder received multi-credential request from: %s\n", req.From())
		fmt.Printf("   Request ID: %s\n", req.RequestID())
		fmt.Printf("   Requested types: %v\n", req.Type())

		// Show what's being requested in detail
		fmt.Println("   📋 Detailed request:")
		for i, detail := range req.Details() {
			fmt.Printf("     %d. Type: %v\n", i+1, detail.CredentialType)
			for j, param := range detail.Parameters {
				fmt.Printf("        Parameter %d: %s != empty\n", j+1, param.Field)
			}
		}

		// For demo, we'll reject but show we understand the multi-credential request
		fmt.Println("   ❌ Rejecting multi-credential request (demo)")
		fmt.Println("      In production: would check for each credential type and respond")
		req.Reject()
	})

	// Issuer processes multi-credential responses
	issuer.Credentials().OnPresentationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("📨 Issuer received multi-credential response from: %s\n", resp.From())
		fmt.Printf("   Status: %s\n", utils.ResponseStatusToString(resp.Status()))
		fmt.Printf("   Presentations received: %d\n", len(resp.Presentations()))

		// Process each presentation
		for i, presentation := range resp.Presentations() {
			fmt.Printf("   Presentation %d:\n", i+1)
			fmt.Printf("     Type: %v\n", presentation.PresentationType())
			fmt.Printf("     Credentials: %d\n", len(presentation.Credentials()))

			// Process each credential in the presentation
			for j, cred := range presentation.Credentials() {
				fmt.Printf("       Credential %d: %v\n", j+1, cred.CredentialType())
			}
		}
	})

	fmt.Println("✅ Multi-credential handlers configured")
	fmt.Println()
}

// demonstrateMultiExchange shows requesting multiple credential types
func demonstrateMultiExchange(issuer, holder *client.Client) {
	fmt.Println("🔄 Demonstrating multi-credential exchange...")

	// Create a request for multiple credential types
	details := []*client.CredentialDetail{
		{
			CredentialType: credential.CredentialTypeEmail,
			Parameters: []*client.CredentialParameter{
				{
					Operator: message.OperatorNotEquals,
					Field:    "emailAddress",
					Value:    "",
				},
			},
		},
		{
			CredentialType: credential.CredentialTypeProfileName,
			Parameters: []*client.CredentialParameter{
				{
					Operator: message.OperatorNotEquals,
					Field:    "firstName",
					Value:    "",
				},
			},
		},
		{
			CredentialType: []string{"VerifiableCredential", "EducationCredential"},
			Parameters: []*client.CredentialParameter{
				{
					Operator: message.OperatorNotEquals,
					Field:    "degree",
					Value:    "",
				},
			},
		},
	}

	fmt.Println("📤 Issuer requesting multiple credentials from holder:")
	fmt.Println("   📧 Email credential")
	fmt.Println("   👤 Profile credential")
	fmt.Println("   🎓 Education credential")

	// Send the multi-credential request
	req, err := issuer.Credentials().RequestPresentationWithTimeout(
		holder.DID(),
		details,
		15*time.Second,
	)
	if err != nil {
		log.Printf("Failed to send multi-credential request: %v", err)
		return
	}

	fmt.Printf("   Request sent with ID: %s\n", req.RequestID())

	// Wait for response
	fmt.Println("⏳ Waiting for multi-credential response...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
	fmt.Println("   1. Issuer requested 3 different credential types")
	fmt.Println("   2. Holder received the complex request and processed each type")
	fmt.Println("   3. Holder rejected the request (demo)")
	fmt.Println("   4. In real scenarios, holder would provide matching credentials")
	fmt.Println("   5. Multiple credentials can be bundled in a single response")
	fmt.Println()
}
