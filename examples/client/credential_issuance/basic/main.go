// Package main demonstrates basic credential issuance using the Self SDK.
//
// This is the BASIC level of credential issuance examples.
// Start here if you're new to credential issuance concepts.
//
// This example shows the basics of:
// - Setting up issuer and holder clients
// - Creating a simple email credential
// - Understanding the credential builder pattern
// - Basic claim addition and signing
//
// 🎯 What you'll learn:
// • How credential issuance works
// • Basic credential creation patterns
// • Simple claim addition
// • Client setup and configuration
// • Cryptographic signing basics
//
// 📚 Next steps:
// • multi_claim/main.go - Multiple claims in credentials
// • evidence/main.go - Evidence and asset management
// • complex/main.go - Complex nested data structures
// • advanced/main.go - All features combined
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("🎓 Basic Credential Issuance Demo")
	fmt.Println("==================================")
	fmt.Println("This demo shows basic credential issuance between issuer and holder.")
	fmt.Println("📚 This is the BASIC level - start here if you're new to credential issuance.")
	fmt.Println()

	// Step 1: Create issuer and holder clients
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Create a simple credential
	createEmailCredential(issuer, holder)

	fmt.Println("✅ Basic demo completed!")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run ../multi_claim/main.go to learn about multiple claims")
	fmt.Println("   • Run ../evidence/main.go for evidence and asset management")
	fmt.Println("   • Run ../complex/main.go for complex nested data structures")
	fmt.Println("   • Run ../advanced/main.go for all features combined")
	fmt.Println()
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	// Create issuer client
	issuer, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("basic_issuer"),
		StoragePath: "./basic_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	// Create holder client
	holder, err := client.NewClient(client.Config{
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

// createEmailCredential creates a basic email credential for demonstration
func createEmailCredential(issuer, holder *client.Client) {
	fmt.Println("📧 Creating basic email credential...")
	fmt.Println("   This demonstrates the foundation of credential issuance")
	fmt.Println("   Key concepts: builder pattern, claims, signing")
	fmt.Println()

	// Create a basic email credential using the builder pattern
	emailCredential, err := issuer.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeEmail).                       // Set credential type
		Subject(holder.DID()).                                      // Who the credential is about
		Issuer(issuer.DID()).                                       // Who is issuing the credential
		Claim("emailAddress", "john.doe@example.com").              // Add email address claim
		Claim("verified", true).                                    // Add verification status
		Claim("verificationDate", time.Now().Format("2006-01-02")). // Add verification date
		ValidFrom(time.Now()).                                      // Set validity start time
		SignWith(issuer.DID(), time.Now()).                         // Sign with issuer's key
		Issue(issuer)                                               // Issue the credential

	if err != nil {
		log.Printf("Failed to create credential: %v", err)
		return
	}

	fmt.Printf("✅ Email credential created successfully\n")
	fmt.Printf("   📧 Email: john.doe@example.com\n")
	fmt.Printf("   ✔️  Verified: true\n")
	fmt.Printf("   📅 Date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("   🔒 Type: %v\n", emailCredential.CredentialType())
	fmt.Printf("   🆔 Subject: %s\n", emailCredential.CredentialSubject().String())
	fmt.Printf("   🏢 Issuer: %s\n", emailCredential.Issuer().String())
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Issuer created a verifiable credential")
	fmt.Println("   2. Added claims (email, verification status, date)")
	fmt.Println("   3. Signed with cryptographic key for integrity")
	fmt.Println("   4. Credential is now ready for sharing or verification")
	fmt.Println()
	fmt.Println("📚 Key Learning Points:")
	fmt.Println("   • Credentials contain claims about a subject")
	fmt.Println("   • Builder pattern provides clean, readable construction")
	fmt.Println("   • Cryptographic signatures ensure integrity")
	fmt.Println("   • Timestamps establish validity periods")
	fmt.Println()
}
