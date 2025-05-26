// Package main demonstrates credential exchange with QR code discovery using the Self SDK.
//
// This is the EXPERT level of credential exchange examples.
// Complete all previous examples before trying this one.
//
// This example shows:
// - QR code generation for peer discovery
// - Real-time peer-to-peer credential exchange
// - Live credential sharing workflows
// - Production-ready discovery patterns
//
// 🎯 What you'll learn:
// • QR code-based peer discovery
// • Real-time credential exchange with live peers
// • Production discovery workflows
// • Integration of discovery with credential exchange
//
// 📚 Prerequisites: basic_exchange.go, multi_credential_exchange.go, advanced_exchange.go
// 📚 This is the final level of the credential exchange tutorial series
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/examples/utils"
	"github.com/joinself/self-go-sdk/message"
)

func main() {
	fmt.Println("🔄 Discovery-Based Credential Exchange Demo")
	fmt.Println("===========================================")
	fmt.Println("This demo shows credential exchange with QR code discovery.")
	fmt.Println("📚 This is the EXPERT level - complete all previous examples first.")
	fmt.Println()

	// Step 1: Create clients
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Create credentials for exchange
	createDiscoveryCredentials(issuer, holder)

	// Step 3: Set up discovery handlers
	setupDiscoveryHandlers(issuer, holder)

	// Step 4: Demonstrate discovery-based exchange
	demonstrateDiscoveryExchange(issuer, holder)

	fmt.Println("✅ Discovery exchange demo completed!")
	fmt.Println()
	fmt.Println("🎓 Congratulations! You've completed the credential exchange tutorial series:")
	fmt.Println("   ✅ basic_exchange.go - Foundation concepts")
	fmt.Println("   ✅ multi_credential_exchange.go - Multiple credential types")
	fmt.Println("   ✅ advanced_exchange.go - Complex parameters and verification")
	fmt.Println("   ✅ discovery_exchange.go - QR code discovery integration")
	fmt.Println()
	fmt.Println("🚀 You're now ready to build production credential exchange applications!")
	fmt.Println()
	fmt.Println("The clients will keep running. Press Ctrl+C to exit.")

	select {}
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	issuer, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("discovery_issuer"),
		StoragePath: "./discovery_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	holder, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("discovery_holder"),
		StoragePath: "./discovery_holder_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create holder:", err)
	}

	fmt.Println("✅ Clients created successfully")
	return issuer, holder
}

// createDiscoveryCredentials creates credentials for discovery-based exchange
func createDiscoveryCredentials(issuer, holder *client.Client) {
	fmt.Println("📝 Creating credentials for discovery exchange...")

	// Create a professional credential for discovery demo
	fmt.Println("💼 Creating professional credential...")
	_, err := issuer.Credentials().NewCredentialBuilder().
		Type([]string{"VerifiableCredential", "ProfessionalCredential"}).
		Subject(holder.DID()).
		Issuer(issuer.DID()).
		Claims(map[string]interface{}{
			"professionalId":    "PROF-2024-001",
			"certificationName": "Self SDK Expert",
			"level":             "Advanced",
			"issueDate":         time.Now().Format("2006-01-02"),
			"skills": []string{
				"Credential Exchange",
				"QR Code Discovery",
				"Decentralized Identity",
				"Self SDK",
			},
			"verified": true,
		}).
		ValidFrom(time.Now()).
		SignWith(issuer.DID(), time.Now()).
		Issue(issuer)

	if err != nil {
		log.Printf("Failed to create professional credential: %v", err)
	} else {
		fmt.Println("   ✅ Professional credential created: Self SDK Expert")
	}

	fmt.Println("✅ Discovery credentials created successfully")
	fmt.Println()
}

// setupDiscoveryHandlers configures handlers for discovery-based exchange
func setupDiscoveryHandlers(issuer, holder *client.Client) {
	fmt.Println("🔧 Setting up discovery exchange handlers...")

	// Discovery-aware presentation handler
	holder.Credentials().OnPresentationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("📨 Holder received discovery-based request from: %s\n", req.From())
		fmt.Printf("   Request ID: %s\n", req.RequestID())
		fmt.Printf("   Requested types: %v\n", req.Type())
		fmt.Printf("   🔗 This request came through discovery connection\n")

		// Show discovery context
		fmt.Println("   🌐 Discovery context:")
		fmt.Println("      • Peer connected via QR code scan")
		fmt.Println("      • Secure encrypted channel established")
		fmt.Println("      • Real-time credential exchange enabled")

		// Process the request
		fmt.Println("   📋 Processing discovery-based credential request...")
		for i, detail := range req.Details() {
			fmt.Printf("     Detail %d - Type: %v\n", i+1, detail.CredentialType)
		}

		// This handler shows what would happen if this client received the request
		// In your case, the mobile app is handling the actual request
		fmt.Println("   📱 Note: In this demo, the mobile app is handling the actual request")
		fmt.Println("      This handler shows what would happen if this client received it")
		fmt.Println()
		fmt.Println("   🔍 What the mobile app should do:")
		fmt.Println("      1. Check if it has email credentials")
		fmt.Println("      2. Verify the requester's identity/permissions")
		fmt.Println("      3. Respond with email credentials (if available)")
		fmt.Println()
		fmt.Println("   ❌ Rejecting request (this client is just for demo)")
		fmt.Println("      The mobile app will handle the real request")
		req.Reject()
	})

	// Discovery-aware response handler
	issuer.Credentials().OnPresentationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("📨 Issuer received discovery-based response from: %s\n", resp.From())
		fmt.Printf("   Status: %s\n", utils.ResponseStatusToString(resp.Status()))
		fmt.Printf("   🔗 Response from discovered peer\n")

		// Process discovery response
		if len(resp.Presentations()) > 0 {
			fmt.Println("   🎉 Successfully received credentials from discovered peer!")
			for i, presentation := range resp.Presentations() {
				fmt.Printf("     Presentation %d: %v\n", i+1, presentation.PresentationType())

				// Display the actual credential data
				for j, credential := range presentation.Credentials() {
					fmt.Printf("       Credential %d:\n", j+1)
					fmt.Printf("         Type: %v\n", credential.CredentialType())
					fmt.Printf("         Subject: %s\n", credential.CredentialSubject().String())
					fmt.Printf("         Issuer: %s\n", credential.Issuer().String())

					// Display claims/data
					fmt.Println("         Claims:")
					claims, err := credential.CredentialSubjectClaims()
					if err != nil {
						fmt.Printf("           Error reading claims: %v\n", err)
					} else {
						for key, value := range claims {
							fmt.Printf("           %s: %v\n", key, value)
						}
					}
				}
			}
		}
	})

	fmt.Println("✅ Discovery handlers configured")
	fmt.Println()
}

// demonstrateDiscoveryExchange shows QR code discovery and credential exchange
func demonstrateDiscoveryExchange(issuer, holder *client.Client) {
	fmt.Println("🔗 DISCOVERY-BASED CREDENTIAL EXCHANGE")
	fmt.Println("======================================")
	fmt.Println("📱 Demonstrating QR code discovery and live credential exchange...")

	// Generate QR code for discovery
	fmt.Println("📱 Generating QR code for peer discovery...")
	qr, err := issuer.Discovery().GenerateQR()
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
		return
	}

	qrCode, err := qr.Unicode()
	if err != nil {
		log.Printf("Failed to get QR code: %v", err)
		return
	}

	fmt.Println("📱 QR CODE FOR CREDENTIAL EXCHANGE:")
	fmt.Println("   Scan this with another Self client to initiate credential exchange")
	fmt.Println(qrCode)
	fmt.Println()
	fmt.Println("🔐 QR Code Features:")
	fmt.Println("   • Contains cryptographic keys for secure connection")
	fmt.Println("   • Enables peer-to-peer credential exchange")
	fmt.Println("   • Compatible with Self mobile apps and SDK clients")
	fmt.Println("   • Establishes encrypted communication channel")
	fmt.Println()

	// Wait for peer connection
	fmt.Println("⏳ Waiting for peer to scan QR code and connect (30 seconds)...")
	fmt.Println("   💡 In a real scenario:")
	fmt.Println("      1. Another user scans this QR code with their Self app")
	fmt.Println("      2. Secure connection is established automatically")
	fmt.Println("      3. Credential exchange can begin immediately")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	peer, err := qr.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("⏰ No peer connected (normal for demo)")
			fmt.Println("   📚 What would happen with a real connection:")
			demonstrateHypotheticalExchange(issuer, holder)
		} else {
			log.Printf("❌ Discovery error: %v", err)
		}
		return
	}

	// Handle successful connection
	fmt.Printf("✅ Peer connected: %s\n", peer.DID())
	fmt.Println("   🔐 Secure encrypted channel established")
	fmt.Println("   🔄 Ready for real-time credential exchange")
	fmt.Println()

	// Demonstrate live credential exchange
	demonstrateLiveExchange(issuer, peer)
}

// demonstrateHypotheticalExchange shows what would happen with a real connection
func demonstrateHypotheticalExchange(issuer, holder *client.Client) {
	fmt.Println("📚 HYPOTHETICAL LIVE EXCHANGE SCENARIO")
	fmt.Println("======================================")
	fmt.Println("🎭 Simulating what would happen with a real peer connection...")

	fmt.Println("📤 Would request email credentials from connected peer:")
	fmt.Println("   📧 Email verification credentials")
	fmt.Println("   🎯 With non-empty email address")
	fmt.Println()

	fmt.Println("🔄 Live exchange workflow would be:")
	fmt.Println("   1. 📱 Peer scans QR code → secure connection established")
	fmt.Println("   2. 📤 Issuer sends credential request to peer")
	fmt.Println("   3. 📨 Peer receives request and processes it")
	fmt.Println("   4. 📋 Peer responds with matching credentials")
	fmt.Println("   5. ✅ Issuer receives and validates credentials")
	fmt.Println("   6. 🎉 Successful real-time credential exchange!")
	fmt.Println()

	fmt.Println("💡 Benefits of discovery-based exchange:")
	fmt.Println("   • 🚀 Instant peer-to-peer connections")
	fmt.Println("   • 🔐 End-to-end encryption")
	fmt.Println("   • 📱 Mobile-friendly QR code interface")
	fmt.Println("   • 🌐 No central authority required")
	fmt.Println("   • ⚡ Real-time credential sharing")
	fmt.Println()
}

// demonstrateLiveExchange shows live credential exchange with a connected peer
func demonstrateLiveExchange(issuer *client.Client, peer *client.Peer) {
	fmt.Println("🔄 LIVE CREDENTIAL EXCHANGE")
	fmt.Println("===========================")
	fmt.Println("🎉 Demonstrating live exchange with connected peer...")

	// Create request for the connected peer
	details := []*client.CredentialDetail{
		{
			CredentialType: []string{"VerifiableCredential", "EmailCredential"},
			Parameters: []*client.CredentialParameter{
				{
					Operator: message.OperatorNotEquals,
					Field:    "emailAddress",
					Value:    "",
				},
			},
		},
	}

	fmt.Printf("📤 Sending live credential request to peer: %s\n", peer.DID())
	fmt.Println("   🔍 Requesting: Email credentials with non-empty email address")

	// Send request to the live peer
	req, err := issuer.Credentials().RequestPresentationWithTimeout(
		peer.DID(),
		details,
		30*time.Second,
	)
	if err != nil {
		log.Printf("Failed to send live request: %v", err)
		return
	}

	fmt.Printf("   Request sent with ID: %s\n", req.RequestID())
	fmt.Println("   ⏳ Waiting for live peer response...")

	// Wait for live response
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	resp, err := req.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("   ⏰ Live request timed out")
		} else {
			fmt.Printf("   ❌ Live request failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Live response received: %s\n", utils.ResponseStatusToString(resp.Status()))

		// Display received credential data
		if len(resp.Presentations()) > 0 {
			fmt.Println("   📋 Received credential presentations:")
			for i, presentation := range resp.Presentations() {
				fmt.Printf("     Presentation %d: %v\n", i+1, presentation.PresentationType())

				// Display the actual credential data
				for j, credential := range presentation.Credentials() {
					fmt.Printf("       Credential %d:\n", j+1)
					fmt.Printf("         Type: %v\n", credential.CredentialType())
					fmt.Printf("         Subject: %s\n", credential.CredentialSubject().String())
					fmt.Printf("         Issuer: %s\n", credential.Issuer().String())

					// Display claims/data
					fmt.Println("         Claims:")
					claims, err := credential.CredentialSubjectClaims()
					if err != nil {
						fmt.Printf("           Error reading claims: %v\n", err)
					} else {
						for key, value := range claims {
							fmt.Printf("           %s: %v\n", key, value)
						}
					}
				}
			}
		}

		fmt.Println("   🎉 Successful live credential exchange!")
	}

	fmt.Println()
	fmt.Println("🎓 Live exchange completed!")
	fmt.Println("   • Real peer-to-peer connection established")
	fmt.Println("   • Credentials exchanged in real-time")
	fmt.Println("   • Secure encrypted communication")
	fmt.Println("   • Production-ready discovery workflow")
	fmt.Println()
}
