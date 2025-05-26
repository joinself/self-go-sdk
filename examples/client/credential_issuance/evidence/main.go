// Package main demonstrates credential issuance with evidence using the Self SDK.
//
// This is the EVIDENCE level of credential issuance examples.
// Prerequisites: Complete basic/main.go and multi_claim/main.go first.
//
// This example shows:
// - Creating custom credential types
// - Attaching file evidence to credentials
// - Asset management and upload functionality
// - Creating verifiable presentations
// - Linking evidence to credential claims
//
// 🎯 What you'll learn:
// • How to attach evidence files to credentials
// • Asset management and secure storage
// • Creating verifiable presentations
// • Linking evidence to claims with hashes
// • Custom credential types
//
// 📚 Next steps:
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
	fmt.Println("🎓 Evidence-Based Credential Issuance Demo")
	fmt.Println("===========================================")
	fmt.Println("This demo shows how to create credentials with file evidence.")
	fmt.Println("📚 This is the EVIDENCE level - adding proof to credentials.")
	fmt.Println()

	// Step 1: Create issuer and holder clients
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Create credentials with evidence
	createCertificationWithEvidence(issuer, holder)

	fmt.Println("✅ Evidence demo completed!")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run ../complex/main.go for complex nested data structures")
	fmt.Println("   • Run ../advanced/main.go for all features combined")
	fmt.Println()
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	// Create issuer client
	issuer, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("evidence_issuer"),
		StoragePath: "./evidence_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	// Create holder client
	holder, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("evidence_holder"),
		StoragePath: "./evidence_holder_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create holder:", err)
	}

	fmt.Println("✅ Clients created successfully")
	return issuer, holder
}

// createCertificationWithEvidence creates a certification credential with file evidence
func createCertificationWithEvidence(issuer, holder *client.Client) {
	fmt.Println("🎓 Creating certification credential with evidence...")
	fmt.Println("   This demonstrates advanced features: evidence, presentations")
	fmt.Println("   Evidence provides additional proof supporting credential claims")
	fmt.Println()

	// Step 1: Create evidence asset
	evidence := createEvidence(issuer)
	if evidence == nil {
		return
	}

	// Step 2: Create credential with evidence reference
	credential := createCredentialWithEvidence(issuer, holder, evidence)
	if credential == nil {
		return
	}

	// Step 3: Create verifiable presentation
	createPresentation(issuer, credential)
}

// createEvidence creates and uploads supporting documentation
func createEvidence(issuer *client.Client) *client.CredentialAsset {
	fmt.Println("📄 Creating evidence asset...")
	fmt.Println("   Evidence can be any file type: PDFs, images, documents, etc.")

	// Create mock certificate document
	certificateData := []byte(`Certificate of Completion
Advanced Go Programming Course

Student: John Doe
Course: Advanced Go Programming with Self SDK
Institution: Self SDK Academy
Grade: A+
Credits: 40 hours
Date: ` + time.Now().Format("2006-01-02") + `

This certificate verifies that the above-named student has
successfully completed the Advanced Go Programming course
with distinction.

Instructor: Jane Smith
Director: Dr. Alice Johnson`)

	evidence, err := issuer.Credentials().CreateAsset("certificate.pdf", "application/pdf", certificateData)
	if err != nil {
		log.Printf("Failed to create evidence: %v", err)
		return nil
	}

	fmt.Printf("   📄 Evidence created: %s\n", evidence.Name)
	fmt.Printf("   🔗 Asset ID: %x\n", evidence.ID())
	fmt.Printf("   🔐 Content Hash: %x\n", evidence.Hash())
	fmt.Printf("   📏 Size: %d bytes\n", len(certificateData))
	fmt.Println("   ✅ Evidence uploaded to secure storage")
	fmt.Println()

	return evidence
}

// createCredentialWithEvidence creates a custom credential with evidence reference
func createCredentialWithEvidence(issuer, holder *client.Client, evidence *client.CredentialAsset) *credential.VerifiableCredential {
	fmt.Println("🏗️ Building custom certification credential...")

	customCredential, err := issuer.Credentials().NewCredentialBuilder().
		Type([]string{"VerifiableCredential", "CertificationCredential"}). // Custom credential type
		Subject(holder.DID()).                                             // Credential subject
		Issuer(issuer.DID()).                                              // Credential issuer
		Claim("certificationType", "Professional Development").            // Type of certification
		Claim("courseName", "Advanced Go Programming").                    // Course name
		Claim("completionDate", time.Now().Format("2006-01-02")).          // Completion date
		Claim("certificateHash", fmt.Sprintf("%x", evidence.Hash())).      // Link to evidence
		Claim("grade", "A+").                                              // Achievement grade
		Claim("institution", "Self SDK Academy").                          // Issuing institution
		Claim("courseHours", 40).                                          // Course duration
		Claim("instructor", "Jane Smith").                                 // Instructor name
		Claim("evidenceId", fmt.Sprintf("%x", evidence.ID())).             // Evidence asset ID
		ValidFrom(time.Now()).                                             // Validity period
		SignWith(issuer.DID(), time.Now()).                                // Cryptographic signature
		Issue(issuer)                                                      // Issue credential

	if err != nil {
		log.Printf("Failed to create custom credential: %v", err)
		return nil
	}

	fmt.Printf("   ✅ Certification credential created successfully\n")
	fmt.Printf("   🎓 Course: Advanced Go Programming\n")
	fmt.Printf("   📅 Completed: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("   🏆 Grade: A+\n")
	fmt.Printf("   🏫 Institution: Self SDK Academy\n")
	fmt.Printf("   👨‍🏫 Instructor: Jane Smith\n")
	fmt.Printf("   ⏱️  Duration: 40 hours\n")
	fmt.Printf("   🔒 Type: %v\n", customCredential.CredentialType())
	fmt.Printf("   🔗 Evidence Hash: %x\n", evidence.Hash())
	fmt.Println()

	return customCredential
}

// createPresentation creates a verifiable presentation from the credential
func createPresentation(issuer *client.Client, cred *credential.VerifiableCredential) {
	fmt.Println("📋 Creating verifiable presentation...")
	fmt.Println("   Presentations package credentials for sharing with verifiers")

	presentation, err := issuer.Credentials().CreatePresentation(
		[]string{"VerifiablePresentation", "CertificationPresentation"}, // Presentation type
		[]*credential.VerifiableCredential{cred},                        // Credentials to include
	)
	if err != nil {
		log.Printf("Failed to create presentation: %v", err)
		return
	}

	fmt.Printf("   ✅ Presentation created successfully\n")
	fmt.Printf("   📋 Type: %v\n", presentation.PresentationType())
	fmt.Printf("   📄 Credentials included: %d\n", len(presentation.Credentials()))
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Created evidence asset (PDF certificate)")
	fmt.Println("   2. Uploaded evidence to secure storage")
	fmt.Println("   3. Created credential with evidence reference")
	fmt.Println("   4. Linked evidence using cryptographic hash")
	fmt.Println("   5. Created verifiable presentation for sharing")
	fmt.Println()
	fmt.Println("📚 Key Learning Points:")
	fmt.Println("   • Evidence provides additional verification material")
	fmt.Println("   • Asset management handles secure file storage")
	fmt.Println("   • Hash references link credentials to evidence")
	fmt.Println("   • Presentations package credentials for sharing")
	fmt.Println("   • Custom credential types support specific use cases")
	fmt.Println("   • Evidence integrity is cryptographically protected")
	fmt.Println()
}
