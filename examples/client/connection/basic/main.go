// Package main demonstrates connection establishment between Self SDK clients.
//
// This example focuses specifically on the different ways to establish
// secure connections between two Self SDK clients:
// - Programmatic connections (for demos and testing)
// - QR code discovery (for real-world scenarios)
// - Connection status checking and management
//
// 🎯 What you'll learn:
// • How to establish programmatic connections without QR codes
// • How to use QR code discovery for real connections
// • How to check connection status and manage peers
// • Understanding the connection lifecycle
//
// 📚 Use cases:
// • Demo applications and testing
// • Automated connection establishment
// • Real-world peer discovery
// • Connection troubleshooting
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("🔗 Self SDK Connection Examples")
	fmt.Println("===============================")
	fmt.Println("This demo shows different ways to establish connections between Self clients.")
	fmt.Println()

	// Create two clients for demonstration
	client1, client2 := createClients()
	defer client1.Close()
	defer client2.Close()

	fmt.Printf("🏢 Client 1: %s\n", client1.DID())
	fmt.Printf("👤 Client 2: %s\n", client2.DID())
	fmt.Println()

	// Show the menu and let user choose
	showMenu()

	for {
		fmt.Print("Choose an option (1-5): ")
		var choice int
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			demonstrateProgrammaticConnection(client1, client2)
		case 2:
			demonstrateQRCodeDiscovery(client1, client2)
		case 3:
			demonstrateConnectionStatus(client1, client2)
		case 4:
			demonstrateConnectionTroubleshooting(client1, client2)
		case 5:
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice. Please try again.")
		}

		fmt.Println()
		showMenu()
	}
}

func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up two clients...")

	// Create first client
	client1, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("connection_client1"),
		StoragePath: "./connection_client1_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create client1:", err)
	}

	// Create second client
	client2, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("connection_client2"),
		StoragePath: "./connection_client2_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create client2:", err)
	}

	fmt.Println("✅ Clients created successfully")
	return client1, client2
}

func showMenu() {
	fmt.Println("📋 Connection Examples Menu:")
	fmt.Println("1. 🤖 Programmatic Connection (for demos/testing)")
	fmt.Println("2. 📱 QR Code Discovery (real-world scenario)")
	fmt.Println("3. 📊 Connection Status & Management")
	fmt.Println("4. 🔧 Connection Troubleshooting")
	fmt.Println("5. 🚪 Exit")
	fmt.Println()
}

func demonstrateProgrammaticConnection(client1, client2 *client.Client) {
	fmt.Println("🤖 PROGRAMMATIC CONNECTION DEMO")
	fmt.Println("================================")
	fmt.Println("This method establishes connections without QR codes.")
	fmt.Println("Perfect for demos, testing, and same-process scenarios.")
	fmt.Println()

	// Method 1: Using ConnectTwoClients utility
	fmt.Println("📡 Method 1: ConnectTwoClients utility")
	fmt.Println("   This is the simplest way to connect two clients in the same process")

	err := client.ConnectTwoClientsWithTimeout(client1, client2, 10*time.Second)
	if err != nil {
		fmt.Printf("   ❌ Connection failed: %v\n", err)
		fmt.Println("   💡 This may happen in demo environments")
	} else {
		fmt.Println("   ✅ Connection established successfully!")
	}
	fmt.Println()

	// Method 2: Using individual ConnectToPeer
	fmt.Println("📡 Method 2: Individual ConnectToPeer")
	fmt.Println("   This shows how to connect to a specific peer by DID")

	result, err := client1.Connection().ConnectToPeerWithTimeout(client2.DID(), 10*time.Second)
	if err != nil {
		fmt.Printf("   ❌ Connection attempt failed: %v\n", err)
	} else {
		fmt.Printf("   📊 Connection result:\n")
		fmt.Printf("      Peer DID: %s\n", result.PeerDID)
		fmt.Printf("      Connected: %v\n", result.Connected)
		if result.Error != nil {
			fmt.Printf("      Error: %v\n", result.Error)
		}
	}
	fmt.Println()

	fmt.Println("🎓 Key Benefits of Programmatic Connections:")
	fmt.Println("   • No QR code scanning required")
	fmt.Println("   • Perfect for automated testing")
	fmt.Println("   • Great for demo applications")
	fmt.Println("   • Enables same-process client connections")
	fmt.Println()

	fmt.Println("⚠️  Important Notes:")
	fmt.Println("   • Both clients must be connected to the messaging service")
	fmt.Println("   • In production, prefer QR code discovery for security")
	fmt.Println("   • Connection may timeout in demo environments")
	fmt.Println()
}

func demonstrateQRCodeDiscovery(client1, client2 *client.Client) {
	fmt.Println("📱 QR CODE DISCOVERY DEMO")
	fmt.Println("=========================")
	fmt.Println("This is the standard way to establish connections in real applications.")
	fmt.Println("One client generates a QR code, the other scans it.")
	fmt.Println()

	// Set up discovery response handler for client1
	client1.Discovery().OnResponse(func(peer *client.Peer) {
		fmt.Printf("🎉 Client1 discovered peer: %s\n", peer.DID())
	})

	// Generate QR code from client1
	fmt.Println("📱 Step 1: Client1 generates QR code")
	qr, err := client1.Discovery().GenerateQRWithTimeout(30 * time.Second)
	if err != nil {
		fmt.Printf("❌ Failed to generate QR code: %v\n", err)
		return
	}

	// Display the QR code
	qrCode, err := qr.Unicode()
	if err != nil {
		fmt.Printf("❌ Failed to get QR code: %v\n", err)
		return
	}

	fmt.Println("📱 QR Code for Client1:")
	fmt.Println(qrCode)
	fmt.Printf("🔗 Request ID: %s\n", qr.RequestID())
	fmt.Println()

	fmt.Println("📱 Step 2: In a real scenario, Client2 would scan this QR code")
	fmt.Println("   • Use a QR code scanner app or camera")
	fmt.Println("   • The QR contains Client1's discovery information")
	fmt.Println("   • Scanning establishes a secure encrypted connection")
	fmt.Println()

	fmt.Println("⏳ Waiting for QR code to be scanned (30 seconds)...")
	fmt.Println("   💡 In this demo, no one will scan it, so it will timeout")

	// Wait for response (will timeout in demo)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	peer, err := qr.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("   ⏰ QR code timed out (expected in demo)")
		} else {
			fmt.Printf("   ❌ Discovery failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Connected to peer: %s\n", peer.DID())
	}
	fmt.Println()

	fmt.Println("🎓 QR Code Discovery Benefits:")
	fmt.Println("   • Secure peer-to-peer discovery")
	fmt.Println("   • Works across different devices/networks")
	fmt.Println("   • User-friendly for mobile apps")
	fmt.Println("   • Industry standard for Self SDK")
	fmt.Println()

	fmt.Println("📱 Real-World Usage:")
	fmt.Println("   1. App A generates QR code")
	fmt.Println("   2. User scans QR with App B")
	fmt.Println("   3. Secure connection established")
	fmt.Println("   4. Apps can now exchange messages/credentials")
	fmt.Println()
}

func demonstrateConnectionStatus(client1, client2 *client.Client) {
	fmt.Println("📊 CONNECTION STATUS & MANAGEMENT")
	fmt.Println("=================================")
	fmt.Println("Learn how to check and manage connection status.")
	fmt.Println()

	// Check if clients are connected to each other
	fmt.Println("🔍 Checking connection status...")

	isConnected1to2 := client1.Connection().IsConnectedTo(client2.DID())
	isConnected2to1 := client2.Connection().IsConnectedTo(client1.DID())

	fmt.Printf("   Client1 → Client2: %v\n", isConnected1to2)
	fmt.Printf("   Client2 → Client1: %v\n", isConnected2to1)
	fmt.Println()

	// List connected peers
	fmt.Println("👥 Connected peers:")

	peers1 := client1.Connection().ListConnectedPeers()
	peers2 := client2.Connection().ListConnectedPeers()

	fmt.Printf("   Client1 connected to %d peers: %v\n", len(peers1), peers1)
	fmt.Printf("   Client2 connected to %d peers: %v\n", len(peers2), peers2)
	fmt.Println()

	// Attempt to establish connection and check status
	fmt.Println("🔗 Attempting connection and monitoring status...")

	result, err := client1.Connection().ConnectToPeerWithTimeout(client2.DID(), 5*time.Second)
	if err != nil {
		fmt.Printf("   ❌ Connection attempt failed: %v\n", err)
	} else {
		fmt.Printf("   📊 Connection attempt result:\n")
		fmt.Printf("      Target: %s\n", result.PeerDID)
		fmt.Printf("      Success: %v\n", result.Connected)
		if result.Error != nil {
			fmt.Printf("      Error: %v\n", result.Error)
		}
	}
	fmt.Println()

	fmt.Println("💡 Connection Status Tips:")
	fmt.Println("   • Connection status may not update immediately")
	fmt.Println("   • In demo environments, connections often timeout")
	fmt.Println("   • Real connections persist across app restarts")
	fmt.Println("   • Use IsConnectedTo() before sending messages")
	fmt.Println()

	fmt.Println("🔧 Connection Management Best Practices:")
	fmt.Println("   • Check connection status before operations")
	fmt.Println("   • Handle connection failures gracefully")
	fmt.Println("   • Implement reconnection logic for critical apps")
	fmt.Println("   • Monitor peer lists for active connections")
	fmt.Println()
}

func demonstrateConnectionTroubleshooting(client1, client2 *client.Client) {
	fmt.Println("🔧 CONNECTION TROUBLESHOOTING")
	fmt.Println("=============================")
	fmt.Println("Common connection issues and how to resolve them.")
	fmt.Println()

	fmt.Println("❓ Common Connection Problems:")
	fmt.Println()

	fmt.Println("1. 🚫 'Connection timeout' errors")
	fmt.Println("   Causes:")
	fmt.Println("   • Clients not connected to messaging service")
	fmt.Println("   • Network connectivity issues")
	fmt.Println("   • Firewall blocking connections")
	fmt.Println("   Solutions:")
	fmt.Println("   • Check internet connectivity")
	fmt.Println("   • Verify messaging service is reachable")
	fmt.Println("   • Try increasing timeout duration")
	fmt.Println()

	fmt.Println("2. 🔑 'Keypair not found' errors")
	fmt.Println("   Causes:")
	fmt.Println("   • Invalid DID format")
	fmt.Println("   • Client not properly initialized")
	fmt.Println("   • Storage corruption")
	fmt.Println("   Solutions:")
	fmt.Println("   • Verify DID format is correct")
	fmt.Println("   • Recreate client if needed")
	fmt.Println("   • Check storage permissions")
	fmt.Println()

	fmt.Println("3. 📡 'Failed to find sender address' warnings")
	fmt.Println("   Causes:")
	fmt.Println("   • Clients not connected to each other")
	fmt.Println("   • Message sent before connection established")
	fmt.Println("   Solutions:")
	fmt.Println("   • Establish connection first")
	fmt.Println("   • Wait for connection confirmation")
	fmt.Println("   • Use connection status checks")
	fmt.Println()

	// Demonstrate diagnostic checks
	fmt.Println("🔍 Running diagnostic checks...")
	fmt.Println()

	// Check client DIDs
	fmt.Printf("✅ Client1 DID: %s\n", client1.DID())
	fmt.Printf("✅ Client2 DID: %s\n", client2.DID())

	if client1.DID() == "" || client2.DID() == "" {
		fmt.Println("❌ Invalid DID detected!")
	}
	fmt.Println()

	// Test connection attempt with detailed error handling
	fmt.Println("🧪 Testing connection with detailed error handling...")

	result, err := client1.Connection().ConnectToPeerWithTimeout(client2.DID(), 3*time.Second)
	if err != nil {
		fmt.Printf("❌ Connection error: %v\n", err)
		fmt.Println("   💡 This is expected in demo environments")
	} else if result.Error != nil {
		fmt.Printf("❌ Connection failed: %v\n", result.Error)

		// Provide specific troubleshooting based on error
		errorStr := result.Error.Error()
		if contains(errorStr, "timeout") {
			fmt.Println("   🔧 Troubleshooting: Try increasing timeout or check network")
		} else if contains(errorStr, "invalid") {
			fmt.Println("   🔧 Troubleshooting: Check DID format and client initialization")
		} else {
			fmt.Println("   🔧 Troubleshooting: Check logs for more details")
		}
	} else if result.Connected {
		fmt.Println("✅ Connection successful!")
	} else {
		fmt.Println("⚠️  Connection attempt completed but not connected")
	}
	fmt.Println()

	fmt.Println("🛠️  Debugging Tips:")
	fmt.Println("   • Enable debug logging: LogLevel: client.LogDebug")
	fmt.Println("   • Check network connectivity")
	fmt.Println("   • Verify both clients are running")
	fmt.Println("   • Test with longer timeouts")
	fmt.Println("   • Try QR code discovery instead")
	fmt.Println()

	fmt.Println("📞 Getting Help:")
	fmt.Println("   • Check Self SDK documentation")
	fmt.Println("   • Review example code")
	fmt.Println("   • Enable verbose logging")
	fmt.Println("   • Test in different environments")
	fmt.Println()
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && contains(s[1:], substr))
}
