// Package main demonstrates multi-claim credential issuance using the Self SDK.
//
// This is the MULTI-CLAIM level of credential issuance examples.
// Prerequisites: Complete basic/main.go first.
//
// This example shows:
// - Creating credentials with multiple claims
// - Different data types in claims
// - Organizing related information in one credential
// - Building upon basic credential concepts
//
// 🎯 What you'll learn:
// • How to add multiple claims to a single credential
// • Different data types in claims (strings, booleans, numbers)
// • Organizing related identity information
// • Efficient credential structuring
//
// 📚 Next steps:
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
	fmt.Println("🎓 Multi-Claim Credential Issuance Demo")
	fmt.Println("========================================")
	fmt.Println("This demo shows how to create credentials with multiple claims.")
	fmt.Println("📚 This is the MULTI-CLAIM level - building on basic concepts.")
	fmt.Println()

	// Step 1: Create issuer and holder clients
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Create credentials with multiple claims
	createProfileCredential(issuer, holder)
	createEducationCredential(issuer, holder)

	fmt.Println("✅ Multi-claim demo completed!")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run ../evidence/main.go for evidence and asset management")
	fmt.Println("   • Run ../complex/main.go for complex nested data structures")
	fmt.Println("   • Run ../advanced/main.go for all features combined")
	fmt.Println()
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	// Create issuer client
	issuer, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("multi_issuer"),
		StoragePath: "./multi_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	// Create holder client
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

// createProfileCredential creates a profile credential with multiple claims
func createProfileCredential(issuer, holder *client.Client) {
	fmt.Println("👤 Creating profile credential with multiple claims...")
	fmt.Println("   This demonstrates grouping related information in one credential")
	fmt.Println("   Multiple claims can contain different data types")
	fmt.Println()

	// Create a profile credential with multiple related claims
	profileCredential, err := issuer.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeProfileName).                 // Profile credential type
		Subject(holder.DID()).                                      // Subject of the credential
		Issuer(issuer.DID()).                                       // Credential issuer
		Claim("firstName", "John").                                 // First name (string)
		Claim("lastName", "Doe").                                   // Last name (string)
		Claim("displayName", "John Doe").                           // Display name (string)
		Claim("profileLevel", "verified").                          // Verification level (string)
		Claim("country", "United States").                          // Country (string)
		Claim("age", 30).                                           // Age (number)
		Claim("isActive", true).                                    // Active status (boolean)
		Claim("registrationDate", time.Now().Format("2006-01-02")). // Registration date (string)
		ValidFrom(time.Now()).                                      // Validity start time
		SignWith(issuer.DID(), time.Now()).                         // Cryptographic signature
		Issue(issuer)                                               // Issue the credential

	if err != nil {
		log.Printf("Failed to create profile credential: %v", err)
		return
	}

	fmt.Printf("✅ Profile credential created successfully\n")
	fmt.Printf("   👤 Name: John Doe\n")
	fmt.Printf("   🌍 Country: United States\n")
	fmt.Printf("   🎂 Age: 30\n")
	fmt.Printf("   ⭐ Level: verified\n")
	fmt.Printf("   ✅ Active: true\n")
	fmt.Printf("   📅 Registration: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("   🔒 Type: %v\n", profileCredential.CredentialType())
	fmt.Println()
}

// createEducationCredential creates an education credential with academic claims
func createEducationCredential(issuer, holder *client.Client) {
	fmt.Println("🎓 Creating education credential with academic claims...")
	fmt.Println("   This shows how to structure educational achievements")
	fmt.Println("   Different claim types for academic information")
	fmt.Println()

	// Create an education credential with academic information
	educationCredential, err := issuer.Credentials().NewCredentialBuilder().
		Type([]string{"VerifiableCredential", "EducationCredential"}). // Education credential type
		Subject(holder.DID()).                                         // Subject of the credential
		Issuer(issuer.DID()).                                          // Credential issuer
		Claim("institution", "University of Technology").              // Institution name (string)
		Claim("degree", "Bachelor of Science").                        // Degree type (string)
		Claim("major", "Computer Science").                            // Major field (string)
		Claim("graduationYear", 2020).                                 // Graduation year (number)
		Claim("gpa", 3.8).                                             // GPA (float as number)
		Claim("honors", true).                                         // Honors status (boolean)
		Claim("creditsCompleted", 120).                                // Credits (number)
		Claim("thesis", "Machine Learning Applications").              // Thesis title (string)
		Claim("graduationDate", "2020-05-15").                         // Graduation date (string)
		ValidFrom(time.Now()).                                         // Validity start time
		SignWith(issuer.DID(), time.Now()).                            // Cryptographic signature
		Issue(issuer)                                                  // Issue the credential

	if err != nil {
		log.Printf("Failed to create education credential: %v", err)
		return
	}

	fmt.Printf("✅ Education credential created successfully\n")
	fmt.Printf("   🏫 Institution: University of Technology\n")
	fmt.Printf("   🎓 Degree: Bachelor of Science in Computer Science\n")
	fmt.Printf("   📅 Graduated: 2020-05-15\n")
	fmt.Printf("   📊 GPA: 3.8\n")
	fmt.Printf("   🏆 Honors: true\n")
	fmt.Printf("   📚 Credits: 120\n")
	fmt.Printf("   📝 Thesis: Machine Learning Applications\n")
	fmt.Printf("   🔒 Type: %v\n", educationCredential.CredentialType())
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Created two credentials with multiple claims each")
	fmt.Println("   2. Used different data types: strings, numbers, booleans")
	fmt.Println("   3. Grouped related information in single credentials")
	fmt.Println("   4. Each credential maintains cryptographic integrity")
	fmt.Println()
	fmt.Println("📚 Key Learning Points:")
	fmt.Println("   • Multiple related claims can be grouped in one credential")
	fmt.Println("   • Claims support different data types (strings, numbers, booleans)")
	fmt.Println("   • Grouping related information improves efficiency")
	fmt.Println("   • Each claim is individually verifiable")
	fmt.Println("   • Credential types help organize different kinds of information")
	fmt.Println()
}
