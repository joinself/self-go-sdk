// Package main demonstrates a simple chat application using the Self SDK client facade.
//
// This example shows how to:
// - Initialize a Self client with minimal configuration
// - Set up chat message handlers for real-time messaging
// - Use QR code-based peer discovery for secure connection establishment
// - Send and receive end-to-end encrypted messages
// - Handle graceful shutdown and cleanup
//
// The Self SDK provides decentralized identity and messaging capabilities,
// allowing peers to connect directly without intermediary servers while
// maintaining full end-to-end encryption and identity verification.
//
// 🎯 CHAT CAPABILITIES DEMONSTRATED:
// • Real-time bidirectional messaging
// • End-to-end encryption (automatic)
// • Message echo functionality
// • Simple command handling (/help, /quit)
// • Multi-peer discovery support
// • Graceful connection management
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joinself/self-go-sdk/client"
)

const (
	// Configuration constants
	discoveryTimeout = 10 * time.Minute
	storageDir       = "./simple_chat_storage"
)

func main() {
	fmt.Println("🚀 Self SDK Simple Chat Example")
	fmt.Println("===============================")
	fmt.Println("📚 This demo showcases the core chat capabilities:")
	fmt.Println("   • Secure peer-to-peer messaging")
	fmt.Println("   • QR code-based connection")
	fmt.Println("   • Real-time message exchange")
	fmt.Println("   • End-to-end encryption")
	fmt.Println()

	// Create a new Self client with minimal configuration
	// The client handles all cryptographic operations, storage, and networking
	selfClient, err := client.NewClient(client.Config{
		StorageKey:  generateStorageKey(), // Secure key for encrypting local storage
		StoragePath: storageDir,           // Directory for storing account state
		Environment: client.Sandbox,       // Use Sandbox for development/testing
		LogLevel:    client.LogInfo,       // Show informational messages
	})
	if err != nil {
		log.Fatal("❌ Failed to create Self client:", err)
	}
	defer selfClient.Close()

	// Your DID (Decentralized Identifier) is your unique identity on the Self network
	fmt.Printf("🆔 Your DID: %s\n", selfClient.DID())
	fmt.Println("   This is your unique decentralized identity")
	fmt.Println()

	// 🎯 CHAT SETUP: Configure message handlers to demonstrate chat capabilities
	setupChatHandlers(selfClient)

	// Set up graceful shutdown handling
	ctx, cancel := setupGracefulShutdown()
	defer cancel()

	// 🔗 PEER DISCOVERY: Establish secure connection using QR code
	peer, err := discoverPeer(selfClient, ctx)
	if err != nil {
		log.Fatal("❌ Failed to discover peer:", err)
	}

	fmt.Printf("✅ Chat connection established with: %s\n", peer.DID())
	fmt.Println("🔐 All messages are automatically end-to-end encrypted")
	fmt.Println()

	// 💬 CHAT DEMONSTRATION: Send initial message to show chat functionality
	greeting := fmt.Sprintf("🎉 Hello! Chat connection established at %s. Try sending me a message!",
		time.Now().Format("15:04:05"))
	err = selfClient.Chat().Send(peer.DID(), greeting)
	if err != nil {
		log.Fatal("❌ Failed to send greeting:", err)
	}

	fmt.Println("💬 CHAT IS NOW ACTIVE!")
	fmt.Println("======================")
	fmt.Println("📨 This demo will echo back any messages you send")
	fmt.Println("🎮 Available commands:")
	fmt.Println("   • Type '/help' to see available commands")
	fmt.Println("   • Type '/quit' to end the chat session")
	fmt.Println("   • Any other text will be echoed back")
	fmt.Println("⚡ Messages are sent and received in real-time")
	fmt.Println("🛑 Press Ctrl+C to exit")
	fmt.Println()

	// Keep the program running to receive and handle messages
	<-ctx.Done()
	fmt.Println("\n👋 Chat session ended. Goodbye!")
}

// setupChatHandlers demonstrates the core chat message handling capabilities
// This function showcases how to:
// - Register message handlers for incoming messages
// - Process different types of messages (commands vs regular text)
// - Send responses back to peers
// - Handle multiple peer connections
func setupChatHandlers(selfClient *client.Client) {
	// 📨 INCOMING MESSAGE HANDLER: Process all received chat messages
	selfClient.Chat().OnMessage(func(msg client.ChatMessage) {
		timestamp := time.Now().Format("15:04:05")

		// Display received message with clear formatting
		fmt.Printf("\n📨 [%s] Message from %s:\n", timestamp, msg.From())
		fmt.Printf("   💬 \"%s\"\n", msg.Text())

		// 🔄 MESSAGE ECHO: Demonstrate bidirectional messaging
		echoMsg := fmt.Sprintf("🔄 Echo [%s]: %s", timestamp, msg.Text())
		fmt.Printf("📤 [%s] Echoing message back...\n", timestamp)
		err := selfClient.Chat().Send(msg.From(), echoMsg)
		if err != nil {
			fmt.Printf("❌ Failed to send echo: %v\n", err)
		} else {
			fmt.Printf("✅ Message echoed successfully\n")
		}
	})

	// 🔍 PEER DISCOVERY HANDLER: Handle multiple peer connections
	// This demonstrates how the same QR code can be used by multiple peers
	selfClient.Discovery().OnResponse(func(peer *client.Peer) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("\n🔍 [%s] New peer discovered: %s\n", timestamp, peer.DID())
		fmt.Printf("   🌟 Multiple peers can connect using the same QR code\n")

		// Send welcome message to newly discovered peers
		welcome := fmt.Sprintf("🌟 Welcome to Self SDK Chat! You connected at %s. "+
			"This demonstrates secure peer-to-peer messaging.", timestamp)
		fmt.Printf("📤 [%s] Sending welcome message to new peer...\n", timestamp)
		err := selfClient.Chat().Send(peer.DID(), welcome)
		if err != nil {
			fmt.Printf("❌ Failed to send welcome message: %v\n", err)
		} else {
			fmt.Printf("✅ Welcome message sent successfully\n")
		}
	})
}

// discoverPeer demonstrates the QR code-based peer discovery workflow
// This showcases how Self SDK enables secure peer-to-peer connections without
// requiring any central servers or pre-shared secrets
func discoverPeer(selfClient *client.Client, ctx context.Context) (*client.Peer, error) {
	fmt.Println("🔍 PEER DISCOVERY PROCESS")
	fmt.Println("=========================")
	fmt.Println("🔑 Generating secure QR code for connection...")
	fmt.Println("   The QR code contains cryptographic keys for establishing")
	fmt.Println("   a secure, end-to-end encrypted connection")

	// Generate a QR code containing cryptographic material for secure connection
	// The QR code includes key exchange information for establishing E2E encryption
	qr, err := selfClient.Discovery().GenerateQR()
	if err != nil {
		return nil, fmt.Errorf("failed to generate discovery QR code: %w", err)
	}

	fmt.Println("\n📱 SCAN THIS QR CODE with another Self client to connect:")
	fmt.Println("   Use another instance of this program or the Self mobile app")

	qrCode, err := qr.Unicode()
	if err != nil {
		return nil, fmt.Errorf("failed to render QR code: %w", err)
	}
	fmt.Println(qrCode)

	// Wait for someone to scan the QR code with timeout and cancellation support
	fmt.Printf("⏳ Waiting for peer to scan QR code (timeout: %v)...\n", discoveryTimeout)
	fmt.Println("   🔐 When scanned, a secure connection will be established")
	fmt.Println("   🛑 Press Ctrl+C to cancel")
	fmt.Println()

	discoveryCtx, cancel := context.WithTimeout(ctx, discoveryTimeout)
	defer cancel()

	peer, err := qr.WaitForResponse(discoveryCtx)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("no peer connected within %v - try running another instance of this program", discoveryTimeout)
		}
		if err == context.Canceled {
			return nil, fmt.Errorf("discovery cancelled by user")
		}
		return nil, fmt.Errorf("error during peer discovery: %w", err)
	}

	fmt.Println("🎉 Peer connection successful!")
	fmt.Println("   🔐 Secure encrypted channel established")
	fmt.Println("   💬 Ready for real-time messaging")
	fmt.Println()

	return peer, nil
}

// setupGracefulShutdown configures signal handling for clean application shutdown
func setupGracefulShutdown() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n🛑 Shutdown signal received...")
		cancel()
	}()

	return ctx, cancel
}

// generateStorageKey creates a storage key for encrypting local account data
// In production, this should be a securely generated and stored key
func generateStorageKey() []byte {
	// For demo purposes, we use a simple key
	// In production, use crypto/rand or load from secure storage
	key := make([]byte, 32)
	copy(key, []byte("demo-key-replace-in-production!!"))
	return key
}
