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
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joinself/self-go-sdk/client"
)

const (
	// Configuration constants
	discoveryTimeout = 10 * time.Minute
	storageDir       = "./simple_chat_storage"

	// Chat commands
	helpCommand = "/help"
	quitCommand = "/quit"
)

func main() {
	fmt.Println("🚀 Self SDK Simple Chat Example")
	fmt.Println("===============================")

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
	fmt.Printf("📱 Your DID: %s\n", selfClient.DID())
	fmt.Println("   This is your unique identity on the Self network")

	// Set up chat functionality with enhanced features
	setupChatHandlers(selfClient)

	// Set up graceful shutdown handling
	ctx, cancel := setupGracefulShutdown()
	defer cancel()

	// Discover and connect to a peer using QR code
	peer, err := discoverPeer(selfClient, ctx)
	if err != nil {
		log.Fatal("❌ Failed to discover peer:", err)
	}

	fmt.Printf("✅ Connected to: %s\n", peer.DID())
	fmt.Println("🔐 All messages are end-to-end encrypted")

	// Send initial greeting with timestamp
	greeting := fmt.Sprintf("Hello! Connection established at %s", time.Now().Format("15:04:05"))
	err = selfClient.Chat().Send(peer.DID(), greeting)
	if err != nil {
		log.Fatal("❌ Failed to send greeting:", err)
	}

	fmt.Println("\n💬 Chat is now active!")
	fmt.Println("   • Messages will be echoed back")
	fmt.Println("   • Type '/help' for commands")
	fmt.Println("   • Press Ctrl+C to exit")
	fmt.Println()

	// Keep the program running to receive messages
	<-ctx.Done()
	fmt.Println("\n👋 Chat session ended. Goodbye!")
}

// setupChatHandlers configures enhanced chat message handling with timestamps and commands
func setupChatHandlers(selfClient *client.Client) {
	// Handle incoming chat messages with enhanced features
	selfClient.Chat().OnMessage(func(msg client.ChatMessage) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("\n[%s] %s: %s\n", timestamp, msg.From(), msg.Text())

		// Handle special commands
		if strings.HasPrefix(msg.Text(), helpCommand) {
			helpResponse := "Available commands:\n" +
				"  /help - Show this help message\n" +
				"  /quit - End the chat session\n" +
				"  Any other text will be echoed back"

			err := selfClient.Chat().Send(msg.From(), helpResponse)
			if err != nil {
				fmt.Printf("❌ Failed to send help response: %v\n", err)
			}
			return
		}

		if strings.HasPrefix(msg.Text(), quitCommand) {
			farewell := "Goodbye! Chat session ended by peer."
			err := selfClient.Chat().Send(msg.From(), farewell)
			if err != nil {
				fmt.Printf("❌ Failed to send farewell: %v\n", err)
			}
			return
		}

		// Echo regular messages back with timestamp
		echoMsg := fmt.Sprintf("Echo [%s]: %s", timestamp, msg.Text())
		err := selfClient.Chat().Send(msg.From(), echoMsg)
		if err != nil {
			fmt.Printf("❌ Failed to send echo: %v\n", err)
		}
	})

	// Handle discovery responses for subscription-based peer discovery
	// This allows multiple peers to connect by scanning the same QR code
	selfClient.Discovery().OnResponse(func(peer *client.Peer) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("\n[%s] 🔍 New peer discovered: %s\n", timestamp, peer.DID())

		// Send a welcome message to newly discovered peers
		welcome := fmt.Sprintf("Welcome! You connected at %s", timestamp)
		err := selfClient.Chat().Send(peer.DID(), welcome)
		if err != nil {
			fmt.Printf("❌ Failed to send welcome message: %v\n", err)
		}
	})
}

// discoverPeer handles the QR code-based peer discovery workflow
// This demonstrates how Self SDK enables secure peer-to-peer connections
func discoverPeer(selfClient *client.Client, ctx context.Context) (*client.Peer, error) {
	fmt.Println("\n🔍 Starting peer discovery...")
	fmt.Println("   Generating QR code for secure connection...")

	// Generate a QR code containing cryptographic material for secure connection
	// The QR code includes key exchange information for establishing E2E encryption
	qr, err := selfClient.Discovery().GenerateQR()
	if err != nil {
		return nil, fmt.Errorf("failed to generate discovery QR code: %w", err)
	}

	fmt.Println("\n📱 Scan this QR code with another Self client to connect:")
	fmt.Println("   The QR code contains cryptographic keys for secure connection")

	qrCode, err := qr.Unicode()
	if err != nil {
		return nil, fmt.Errorf("failed to render QR code: %w", err)
	}
	fmt.Println(qrCode)

	// Wait for someone to scan the QR code with timeout and cancellation support
	fmt.Printf("⏳ Waiting for connection (timeout: %v)...\n", discoveryTimeout)
	fmt.Println("   Press Ctrl+C to cancel")

	discoveryCtx, cancel := context.WithTimeout(ctx, discoveryTimeout)
	defer cancel()

	peer, err := qr.WaitForResponse(discoveryCtx)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("no peer connected within %v", discoveryTimeout)
		}
		if err == context.Canceled {
			return nil, fmt.Errorf("discovery cancelled by user")
		}
		return nil, fmt.Errorf("error during peer discovery: %w", err)
	}

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
