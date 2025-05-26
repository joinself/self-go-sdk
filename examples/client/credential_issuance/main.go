package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/message"
)

func main() {
	// Create a new Self client (issuer)
	issuerClient, err := client.NewClient(client.Config{
		StorageKey:  make([]byte, 32), // In production, use a secure key
		StoragePath: "./issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer client:", err)
	}
	defer issuerClient.Close()

	// Create another client (holder)
	holderClient, err := client.NewClient(client.Config{
		StorageKey:  make([]byte, 32), // In production, use a secure key
		StoragePath: "./holder_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create holder client:", err)
	}
	defer holderClient.Close()

	fmt.Printf("Issuer DID: %s\n", issuerClient.DID())
	fmt.Printf("Holder DID: %s\n", holderClient.DID())

	// Set up credential request handlers for the holder
	holderClient.Credentials().OnVerificationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("\n🔍 Holder received verification request from: %s\n", req.From())
		fmt.Printf("   Type: %v\n", req.Type())
		fmt.Printf("   Evidence: %d items\n", len(req.Evidence()))
		fmt.Printf("   Proof: %d presentations\n", len(req.Proof()))

		// For demo purposes, we'll reject the request
		fmt.Println("   ❌ Rejecting request (demo)")
		err := req.Reject()
		if err != nil {
			fmt.Printf("   Failed to reject request: %v\n", err)
		}
	})

	// Set up response handlers for the issuer
	issuerClient.Credentials().OnVerificationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("\n📨 Issuer received verification response from: %s\n", resp.From())
		fmt.Printf("   Status: %s\n", responseStatusToString(resp.Status()))
		fmt.Printf("   Credentials: %d\n", len(resp.Credentials()))
	})

	// Example 1: Create a custom credential with evidence
	fmt.Println("\n📋 Creating custom credential with evidence...")

	// Create some evidence (e.g., a PDF document)
	pdfData := []byte("This is a mock PDF document for demonstration purposes")
	evidence, err := issuerClient.Credentials().CreateAsset("agreement.pdf", "application/pdf", pdfData)
	if err != nil {
		log.Printf("Failed to create evidence: %v", err)
		return
	}

	fmt.Printf("   Created evidence asset: %s (ID: %x)\n", evidence.Name, evidence.ID())

	// Create a custom credential using the builder
	customCredential, err := issuerClient.Credentials().NewCredentialBuilder().
		Type([]string{"VerifiableCredential", "AgreementCredential"}).
		Subject(holderClient.DID()).
		Issuer(issuerClient.DID()).
		Claim("agreementType", "Service Agreement").
		Claim("agreementVersion", "1.0").
		Claim("termsHash", fmt.Sprintf("%x", evidence.Hash())).
		Claims(map[string]interface{}{
			"parties": []map[string]interface{}{
				{
					"type": "issuer",
					"id":   issuerClient.DID(),
				},
				{
					"type": "holder",
					"id":   holderClient.DID(),
				},
			},
			"effectiveDate": time.Now().Format("2006-01-02"),
		}).
		ValidFrom(time.Now()).
		SignWith(issuerClient.DID(), time.Now()).
		Issue(issuerClient)

	if err != nil {
		log.Printf("Failed to create custom credential: %v", err)
		return
	}

	fmt.Printf("   ✅ Created custom credential: %v\n", customCredential.CredentialType())

	// Create a presentation with the credential
	presentation, err := createPresentation(issuerClient, customCredential)
	if err != nil {
		log.Printf("Failed to create presentation: %v", err)
		return
	}

	fmt.Printf("   ✅ Created presentation: %v\n", presentation.PresentationType())

	// Example 2: Discovery and connection
	fmt.Println("\n📱 Generating QR code for discovery...")
	qr, err := issuerClient.Discovery().GenerateQR()
	if err != nil {
		log.Fatal("Failed to generate QR code:", err)
	}

	qrCode, err := qr.Unicode()
	if err != nil {
		log.Fatal("Failed to get QR code:", err)
	}
	fmt.Println("Scan this QR code with the holder client to connect:")
	fmt.Println(qrCode)

	// Wait for someone to scan the QR code
	fmt.Println("\n⏳ Waiting for connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	peer, err := qr.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("No one connected within 2 minutes. Continuing with demo...")
			// For demo purposes, we'll simulate a connection
			peer = &client.Peer{}
		} else {
			log.Fatal("Error waiting for connection:", err)
		}
	} else {
		fmt.Printf("✅ Connected to: %s\n", peer.DID())
		// Wait a moment for the connection to be fully established
		time.Sleep(2 * time.Second)
	}

	// Example 3: Request verification with evidence and proof
	fmt.Println("\n🔍 Requesting credential verification with evidence...")

	// Create evidence for the verification request
	credentialEvidence := []*client.CredentialEvidence{
		{
			Type:   "terms",
			Object: evidence.Object(),
		},
	}

	// Create proof presentations
	proofPresentations := []*credential.VerifiablePresentation{presentation}

	// For demo purposes, we'll use the holder's DID as the target
	// In a real scenario, this would be the connected peer's DID
	targetDID := holderClient.DID()

	verificationReq, err := issuerClient.Credentials().RequestVerificationWithEvidence(
		targetDID,
		[]string{"VerifiableCredential", "AgreementCredential"},
		credentialEvidence,
		proofPresentations,
	)
	if err != nil {
		log.Printf("Failed to request verification: %v", err)
	} else {
		fmt.Printf("   Request ID: %s\n", verificationReq.RequestID())
		fmt.Printf("   Evidence: %d items\n", len(credentialEvidence))
		fmt.Printf("   Proof: %d presentations\n", len(proofPresentations))

		// Wait for response
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		resp, err := verificationReq.WaitForResponse(ctx)
		if err != nil {
			fmt.Printf("   ❌ Request failed or timed out: %v\n", err)
		} else {
			fmt.Printf("   ✅ Received response with status: %s\n", responseStatusToString(resp.Status()))
		}
	}

	// Example 4: Create and issue multiple credential types
	fmt.Println("\n📋 Creating multiple credential types...")

	// Email credential
	emailCredential, err := issuerClient.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeEmail).
		Subject(holderClient.DID()).
		Issuer(issuerClient.DID()).
		Claim("emailAddress", "user@example.com").
		Claim("verified", true).
		Claim("verificationDate", time.Now().Format("2006-01-02")).
		ValidFrom(time.Now()).
		SignWith(issuerClient.DID(), time.Now()).
		Issue(issuerClient)

	if err != nil {
		log.Printf("Failed to create email credential: %v", err)
	} else {
		fmt.Printf("   ✅ Created email credential: %v\n", emailCredential.CredentialType())
	}

	// Profile credential with custom claims
	profileCredential, err := issuerClient.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeProfileName).
		Subject(holderClient.DID()).
		Issuer(issuerClient.DID()).
		Claim("firstName", "John").
		Claim("lastName", "Doe").
		Claim("displayName", "John Doe").
		Claim("profileLevel", "verified").
		ValidFrom(time.Now()).
		SignWith(issuerClient.DID(), time.Now()).
		Issue(issuerClient)

	if err != nil {
		log.Printf("Failed to create profile credential: %v", err)
	} else {
		fmt.Printf("   ✅ Created profile credential: %v\n", profileCredential.CredentialType())
	}

	// Organization credential with complex claims
	orgCredential, err := issuerClient.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeOrganisation).
		Subject(holderClient.DID()).
		Issuer(issuerClient.DID()).
		Claims(map[string]interface{}{
			"organizationName": "Example Corp",
			"role":             "Software Engineer",
			"department":       "Engineering",
			"employeeId":       "EMP001",
			"startDate":        "2023-01-01",
			"permissions": []string{
				"read:documents",
				"write:code",
				"deploy:staging",
			},
			"manager": map[string]interface{}{
				"name":  "Jane Smith",
				"email": "jane.smith@example.com",
			},
		}).
		ValidFrom(time.Now()).
		SignWith(issuerClient.DID(), time.Now()).
		Issue(issuerClient)

	if err != nil {
		log.Printf("Failed to create organization credential: %v", err)
	} else {
		fmt.Printf("   ✅ Created organization credential: %v\n", orgCredential.CredentialType())
	}

	fmt.Println("\n🎉 Enhanced credential features demo completed!")
	fmt.Println("Features demonstrated:")
	fmt.Println("  ✅ Custom credential creation with builder pattern")
	fmt.Println("  ✅ Evidence/proof attachments (file uploads)")
	fmt.Println("  ✅ Custom credential schemas with flexible claims")
	fmt.Println("  ✅ Multiple credential types (Email, Profile, Organization)")
	fmt.Println("  ✅ Complex nested claims and arrays")
	fmt.Println("  ✅ Asset management (upload/download)")
}

// Helper function to create a presentation
func createPresentation(client *client.Client, cred *credential.VerifiableCredential) (*credential.VerifiablePresentation, error) {
	return client.Credentials().CreatePresentation(
		[]string{"VerifiablePresentation", "AgreementPresentation"},
		[]*credential.VerifiableCredential{cred},
	)
}

// Helper function to convert response status to string
func responseStatusToString(status message.ResponseStatus) string {
	switch status {
	case message.ResponseStatusAccepted:
		return "Accepted"
	case message.ResponseStatusForbidden:
		return "Forbidden"
	case message.ResponseStatusNotFound:
		return "Not Found"
	case message.ResponseStatusUnauthorized:
		return "Unauthorized"
	default:
		return "Unknown"
	}
}
