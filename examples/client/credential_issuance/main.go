// Package main demonstrates simple credential issuance using the Self SDK.
//
// This is a simplified version of the credential issuance tutorial.
// For the complete educational progression, see the individual tutorial files:
//
// 📚 Educational Progression:
// 1. basic/main.go - Foundation concepts (start here)
// 2. multi_claim/main.go - Multiple claims in credentials
// 3. evidence/main.go - Evidence and asset management
// 4. complex/main.go - Complex nested data structures
// 5. advanced/main.go - All features combined
//
// This example shows the basics of:
// - Setting up issuer and holder clients
// - Creating a simple credential
// - Understanding the issuance workflow
// - Basic credential builder usage
//
// 🎯 What you'll learn:
// • How credential issuance works
// • Basic credential creation patterns
// • Simple claim addition
// • Client setup and configuration
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
	fmt.Println("🎓 Simple Credential Issuance Demo")
	fmt.Println("===================================")
	fmt.Println("This demo shows basic credential issuance between issuer and holder.")
	fmt.Println()

	// Step 1: Create issuer and holder clients
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Create a simple credential
	createSimpleCredential(issuer, holder)

	fmt.Println("✅ Basic demo completed!")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • cd basic && go run main.go - Foundation concepts (start here)")
	fmt.Println("   • cd multi_claim && go run main.go - Multiple claims in credentials")
	fmt.Println("   • cd evidence && go run main.go - Evidence and asset management")
	fmt.Println("   • cd complex && go run main.go - Complex nested data structures")
	fmt.Println("   • cd advanced && go run main.go - All features combined")
	fmt.Println()
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	// Create issuer client
	issuer, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("simple_issuer"),
		StoragePath: "./simple_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	// Create holder client
	holder, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("simple_holder"),
		StoragePath: "./simple_holder_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create holder:", err)
	}

	fmt.Println("✅ Clients created successfully")
	return issuer, holder
}

// createSimpleCredential creates a basic email credential for demonstration
func createSimpleCredential(issuer, holder *client.Client) {
	fmt.Println("📧 Creating simple email credential...")

	// Create a basic email credential using the builder pattern
	emailCredential, err := issuer.Credentials().NewCredentialBuilder().
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

	fmt.Printf("✅ Email credential created successfully\n")
	fmt.Printf("   📧 Email: demo@example.com\n")
	fmt.Printf("   ✔️  Verified: true\n")
	fmt.Printf("   🔒 Type: %v\n", emailCredential.CredentialType())
	fmt.Printf("   🆔 Subject: %s\n", emailCredential.CredentialSubject().String())
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Issuer created a verifiable credential")
	fmt.Println("   2. Added claims (email address, verification status)")
	fmt.Println("   3. Signed with cryptographic key")
	fmt.Println("   4. Credential is now ready for sharing")
	fmt.Println()
}
