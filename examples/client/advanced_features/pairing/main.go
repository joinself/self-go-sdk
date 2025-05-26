// Package main demonstrates account pairing capabilities of the Self SDK.
//
// This is the PAIRING level of advanced features examples.
// Prerequisites: Complete ../storage/main.go and ../notifications/main.go first.
//
// This example shows:
// - Multi-device account synchronization
// - QR code-based device pairing
// - Secure cryptographic device verification
// - Cross-device state management
// - Pairing event handling and workflows
//
// 🎯 What you'll learn:
// • How to generate pairing codes and QR codes
// • Multi-device synchronization patterns
// • Secure device verification processes
// • Pairing event handling and management
// • Cross-device data synchronization
//
// 📚 Next steps:
// • ../production_patterns/main.go - Real-world storage patterns
// • ../integration/main.go - Component integration workflows
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("🔗 Account Pairing Demo")
	fmt.Println("=======================")
	fmt.Println("This demo showcases Self SDK account pairing capabilities.")
	fmt.Println("📚 This is the PAIRING level - multi-device synchronization.")
	fmt.Println()

	// Step 1: Create a Self client for pairing demonstrations
	pairingClient := createPairingClient()
	defer pairingClient.Close()

	fmt.Printf("🆔 Client DID: %s\n", pairingClient.DID())
	fmt.Println()

	// Step 2: Set up pairing event handlers
	setupPairingHandlers(pairingClient)

	// Step 3: Demonstrate pairing code generation
	demonstratePairingCodeGeneration(pairingClient)

	// Step 4: Show QR code generation for easy pairing
	demonstrateQRCodePairing(pairingClient)

	// Step 5: Explore pairing status and management
	demonstratePairingManagement(pairingClient)

	fmt.Println("✅ Account pairing demo completed!")
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Set up pairing event handlers for secure communication")
	fmt.Println("   2. Generated pairing codes for device authentication")
	fmt.Println("   3. Created QR codes for easy mobile device pairing")
	fmt.Println("   4. Demonstrated pairing status management and monitoring")
	fmt.Println()
	fmt.Println("🎯 Pairing benefits:")
	fmt.Println("   • Seamless multi-device experiences")
	fmt.Println("   • Secure cryptographic device verification")
	fmt.Println("   • Cross-device state synchronization")
	fmt.Println("   • Easy QR code-based device onboarding")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run ../production_patterns/main.go for real-world patterns")
	fmt.Println("   • Run ../integration/main.go for component integration")
}

// createPairingClient sets up a Self client for pairing demonstrations
func createPairingClient() *client.Client {
	fmt.Println("🔧 Setting up pairing client...")

	pairingClient, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("pairing_demo"),
		StoragePath: "./pairing_demo_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create pairing client:", err)
	}

	fmt.Println("✅ Pairing client created successfully")
	return pairingClient
}

// setupPairingHandlers configures pairing event handlers
func setupPairingHandlers(selfClient *client.Client) {
	fmt.Println("🔹 Setting Up Pairing Handlers")
	fmt.Println("==============================")
	fmt.Println("Configuring secure pairing event handling...")
	fmt.Println()

	pairing := selfClient.Pairing()

	// Handler for incoming pairing requests
	pairing.OnPairingRequest(func(request *client.IncomingPairingRequest) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("   📥 [%s] Pairing request received\n", timestamp)
		fmt.Printf("      👤 From: %s\n", request.From())
		fmt.Printf("      🆔 Request ID: %s\n", request.RequestID())
		fmt.Printf("      📍 Address: %s\n", request.Address().String())
		fmt.Printf("      👥 Roles: %d\n", request.Roles())
		fmt.Printf("      ⏰ Expires: %s\n", request.Expires().Format("15:04:05"))
		fmt.Println()

		// For demo safety, we'll auto-reject pairing requests
		// In a real application, you'd present this to the user for approval
		fmt.Println("      🚫 Auto-rejecting for demo safety")
		fmt.Println("      💡 In production: present to user for approval")
		err := request.Reject()
		if err != nil {
			log.Printf("Failed to reject pairing request: %v", err)
		} else {
			fmt.Println("      ✅ Pairing request rejected safely")
		}
		fmt.Println()
	})

	// Handler for pairing responses
	pairing.OnPairingResponse(func(response *client.PairingResponse) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("   📨 [%s] Pairing response received\n", timestamp)
		fmt.Printf("      👤 From: %s\n", response.From())
		fmt.Printf("      📊 Status: %d\n", response.Status())

		if response.Operation() != nil {
			fmt.Println("      ✅ Identity operation included")
		}

		if len(response.Assets()) > 0 {
			fmt.Printf("      📎 %d supporting assets included\n", len(response.Assets()))
		}

		fmt.Println("      🔐 Pairing response processed securely")
		fmt.Println()
	})

	fmt.Println("✅ Pairing event handlers configured")
	fmt.Println("   • Incoming pairing request handling")
	fmt.Println("   • Pairing response processing")
	fmt.Println("   • Secure cryptographic verification")
	fmt.Println("   • Automatic safety measures enabled")
	fmt.Println()
}

// demonstratePairingCodeGeneration shows pairing code creation
func demonstratePairingCodeGeneration(selfClient *client.Client) {
	fmt.Println("🔹 Pairing Code Generation")
	fmt.Println("==========================")
	fmt.Println("Creating secure pairing codes for device authentication...")
	fmt.Println()

	pairing := selfClient.Pairing()

	// Get current pairing code
	fmt.Println("🔑 Generating pairing code...")
	pairingCode, err := pairing.GetPairingCode()
	if err != nil {
		log.Printf("Failed to get pairing code: %v", err)
		return
	}

	fmt.Printf("   ✅ Pairing code generated successfully\n")
	fmt.Printf("   🔑 Code: %s\n", pairingCode.Code)
	fmt.Printf("   📱 Unpaired status: %t\n", pairingCode.Unpaired)
	fmt.Printf("   ⏰ Expires at: %s\n", pairingCode.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   🕐 Valid for: %v\n", time.Until(pairingCode.ExpiresAt).Round(time.Minute))
	fmt.Println()

	// Check current pairing status
	fmt.Println("🔍 Checking pairing status...")
	isPaired, err := pairing.IsPaired()
	if err != nil {
		log.Printf("Failed to check pairing status: %v", err)
	} else {
		fmt.Printf("   📊 Current pairing status: %t\n", isPaired)
		if isPaired {
			fmt.Println("   ✅ Device is already paired with other devices")
		} else {
			fmt.Println("   📱 Device is ready for initial pairing")
		}
	}

	fmt.Println("\n🔐 Security Features:")
	fmt.Println("   • Cryptographically secure pairing codes")
	fmt.Println("   • Time-limited validity for security")
	fmt.Println("   • Automatic expiration prevents misuse")
	fmt.Println("   • Unique codes for each pairing session")
	fmt.Println()
}

// demonstrateQRCodePairing shows QR code generation for easy pairing
func demonstrateQRCodePairing(selfClient *client.Client) {
	fmt.Println("🔹 QR Code Pairing")
	fmt.Println("==================")
	fmt.Println("Generating QR codes for easy mobile device pairing...")
	fmt.Println()

	pairing := selfClient.Pairing()

	// Generate QR code for pairing
	fmt.Println("📱 Generating QR code for mobile pairing...")
	qrCode, err := pairing.GeneratePairingQR()
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
		return
	}

	fmt.Printf("   ✅ QR code generated successfully\n")
	fmt.Printf("   📏 QR data length: %d characters\n", len(qrCode))
	fmt.Printf("   🔗 QR data preview: %s...\n", qrCode[:min(50, len(qrCode))])
	fmt.Println()

	// Simulate QR code usage scenarios
	fmt.Println("📱 QR Code Usage Scenarios:")
	fmt.Println()

	fmt.Println("   📲 Mobile App Pairing:")
	fmt.Println("      1. Open Self mobile app")
	fmt.Println("      2. Navigate to 'Add Device' or 'Pair Device'")
	fmt.Println("      3. Scan this QR code with your mobile camera")
	fmt.Println("      4. Confirm pairing on both devices")
	fmt.Println()

	fmt.Println("   💻 Desktop-to-Desktop Pairing:")
	fmt.Println("      1. Display QR code on first device")
	fmt.Println("      2. Use second device to scan QR code")
	fmt.Println("      3. Complete cryptographic handshake")
	fmt.Println("      4. Verify pairing success on both devices")
	fmt.Println()

	fmt.Println("   🔄 Cross-Platform Pairing:")
	fmt.Println("      1. Generate QR code on any Self SDK application")
	fmt.Println("      2. Scan with any other Self SDK application")
	fmt.Println("      3. Automatic protocol negotiation")
	fmt.Println("      4. Seamless cross-platform synchronization")
	fmt.Println()

	fmt.Println("🎯 QR Code Benefits:")
	fmt.Println("   • No manual code entry required")
	fmt.Println("   • Reduced user error in pairing process")
	fmt.Println("   • Fast and intuitive pairing experience")
	fmt.Println("   • Works across all device types and platforms")
	fmt.Println("   • Embedded security credentials for safe pairing")
	fmt.Println()
}

// demonstratePairingManagement shows pairing status and management
func demonstratePairingManagement(selfClient *client.Client) {
	fmt.Println("🔹 Pairing Management")
	fmt.Println("=====================")
	fmt.Println("Managing device pairing status and synchronization...")
	fmt.Println()

	pairing := selfClient.Pairing()

	// Demonstrate pairing status monitoring
	fmt.Println("📊 Pairing Status Monitoring:")

	// Check if device is paired
	isPaired, err := pairing.IsPaired()
	if err != nil {
		log.Printf("Failed to check pairing status: %v", err)
	} else {
		fmt.Printf("   🔄 Pairing Status: %t\n", isPaired)

		if isPaired {
			fmt.Println("   ✅ Device is part of a paired device network")
			fmt.Println("   🔄 Data synchronization is active")
			fmt.Println("   📱 Cross-device features are available")
		} else {
			fmt.Println("   📱 Device is standalone (not paired)")
			fmt.Println("   🔗 Ready to pair with other devices")
			fmt.Println("   💡 Use QR code or pairing code to connect")
		}
	}
	fmt.Println()

	// Demonstrate pairing code refresh
	fmt.Println("🔄 Pairing Code Management:")
	fmt.Println("   💡 Pairing codes automatically expire for security")
	fmt.Println("   🔄 New codes can be generated as needed")
	fmt.Println("   ⏰ Each code has a limited validity period")
	fmt.Println("   🔐 Expired codes cannot be used for pairing")
	fmt.Println()

	// Show pairing best practices
	fmt.Println("🎯 Pairing Best Practices:")
	fmt.Println()

	fmt.Println("   🔐 Security:")
	fmt.Println("      • Always verify device identity before pairing")
	fmt.Println("      • Use QR codes in secure environments")
	fmt.Println("      • Monitor pairing requests for unauthorized attempts")
	fmt.Println("      • Regularly review paired device list")
	fmt.Println()

	fmt.Println("   📱 User Experience:")
	fmt.Println("      • Provide clear pairing instructions")
	fmt.Println("      • Show pairing progress and status")
	fmt.Println("      • Offer multiple pairing methods (QR, code)")
	fmt.Println("      • Confirm successful pairing to users")
	fmt.Println()

	fmt.Println("   🔄 Synchronization:")
	fmt.Println("      • Ensure data consistency across devices")
	fmt.Println("      • Handle offline/online synchronization")
	fmt.Println("      • Provide conflict resolution mechanisms")
	fmt.Println("      • Monitor synchronization health")
	fmt.Println()

	fmt.Println("🔄 Integration Opportunities:")
	fmt.Println("   • Combine with storage for cross-device data sync")
	fmt.Println("   • Use with notifications for pairing alerts")
	fmt.Println("   • Integrate with chat for multi-device messaging")
	fmt.Println("   • Connect to credentials for identity synchronization")
	fmt.Println()
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
