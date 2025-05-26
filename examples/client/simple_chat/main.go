// Package main demonstrates simple chat messaging using the Self SDK.
//
// This example shows the basics of:
// - Setting up a Self client for messaging
// - Establishing peer connections via QR codes
// - Sending and receiving real-time messages
// - Understanding the chat workflow
//
// 🎯 What you'll learn:
// • How peer-to-peer chat works with Self SDK
// • Basic message sending and receiving patterns
// • QR code-based peer discovery
// • Real-time encrypted messaging
//
// 💬 CHAT CAPABILITIES DEMONSTRATED:
// • Secure peer-to-peer messaging
// • End-to-end encryption (automatic)
// • Real-time bidirectional communication
// • Simple message echo functionality
// • Multi-peer support
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("💬 Simple Chat Demo")
	fmt.Println("===================")
	fmt.Println("This demo shows basic chat messaging between peers.")
	fmt.Println()

	// Step 1: Create a Self client
	chatClient := createClient()
	defer chatClient.Close()

	fmt.Printf("🆔 Your DID: %s\n", chatClient.DID())
	fmt.Println()

	// Step 2: Set up message handlers
	setupChatHandlers(chatClient)

	// Step 3: Discover and connect to a peer
	peer := discoverPeer(chatClient)

	// Step 4: Demonstrate chat messaging
	demonstrateChat(chatClient, peer)

	fmt.Println("✅ Basic chat demo completed!")
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Created a Self client for messaging")
	fmt.Println("   2. Set up handlers to receive and process messages")
	fmt.Println("   3. Used QR code to discover and connect to a peer")
	fmt.Println("   4. Exchanged encrypted messages in real-time")
	fmt.Println("   5. Demonstrated echo functionality")
	fmt.Println()
	fmt.Println("The client will keep running to show ongoing chat capabilities.")
	fmt.Println("Send messages from another instance to see real-time messaging!")
	fmt.Println("Press Ctrl+C to exit.")

	// Keep running to demonstrate ongoing chat capabilities
	select {}
}

// createClient sets up a Self client for chat messaging
func createClient() *client.Client {
	fmt.Println("🔧 Setting up chat client...")

	chatClient, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("simple_chat"),
		StoragePath: "./simple_chat_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create chat client:", err)
	}

	fmt.Println("✅ Chat client created successfully")
	return chatClient
}

// setupChatHandlers configures how the client handles incoming messages
func setupChatHandlers(chatClient *client.Client) {
	fmt.Println("📨 Setting up message handlers...")

	// Handle incoming chat messages
	chatClient.Chat().OnMessage(func(msg client.ChatMessage) {
		timestamp := time.Now().Format("15:04:05")

		fmt.Printf("\n📨 [%s] Message received from %s:\n", timestamp, msg.From())
		fmt.Printf("   💬 \"%s\"\n", msg.Text())

		// Demonstrate different types of responses based on message content
		response := generateResponse(msg.Text(), timestamp)

		fmt.Printf("📤 [%s] Sending response...\n", timestamp)
		err := chatClient.Chat().Send(msg.From(), response)
		if err != nil {
			fmt.Printf("❌ Failed to send response: %v\n", err)
		} else {
			fmt.Printf("✅ Response sent: \"%s\"\n", response)
		}
		fmt.Println()
	})

	// Handle new peer connections
	chatClient.Discovery().OnResponse(func(peer *client.Peer) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("\n🔍 [%s] New peer connected: %s\n", timestamp, peer.DID())

		// Send welcome message to new peers
		welcome := fmt.Sprintf("🎉 Welcome to Self SDK Chat! Connected at %s. Try sending me a message!", timestamp)
		fmt.Printf("📤 [%s] Sending welcome message...\n", timestamp)
		err := chatClient.Chat().Send(peer.DID(), welcome)
		if err != nil {
			fmt.Printf("❌ Failed to send welcome: %v\n", err)
		} else {
			fmt.Printf("✅ Welcome message sent\n")
		}
		fmt.Println()
	})

	fmt.Println("✅ Message handlers configured")
	fmt.Println()
}

// generateResponse creates different responses based on the incoming message
func generateResponse(message, timestamp string) string {
	message = strings.ToLower(strings.TrimSpace(message))

	switch {
	case strings.Contains(message, "hello") || strings.Contains(message, "hi"):
		return fmt.Sprintf("👋 Hello there! Message received at %s", timestamp)
	case strings.Contains(message, "how are you"):
		return "🤖 I'm doing great! Thanks for asking. I'm a Self SDK chat demo."
	case strings.Contains(message, "help"):
		return "💡 This is a chat demo. Try saying 'hello', 'how are you', or just send any message!"
	case strings.Contains(message, "time"):
		return fmt.Sprintf("🕐 Current time is %s", timestamp)
	default:
		return fmt.Sprintf("🔄 Echo [%s]: %s", timestamp, message)
	}
}

// discoverPeer establishes a connection with another peer via QR code
func discoverPeer(chatClient *client.Client) *client.Peer {
	fmt.Println("🔍 Discovering peer for chat...")
	fmt.Println("🔑 Generating QR code for secure connection...")

	// Generate QR code for peer discovery
	qr, err := chatClient.Discovery().GenerateQR()
	if err != nil {
		log.Fatal("Failed to generate QR code:", err)
	}

	fmt.Println("\n📱 SCAN THIS QR CODE with another Self client:")
	fmt.Println("   • Run another instance of this program")
	fmt.Println("   • Use the Self mobile app")
	fmt.Println("   • Any Self SDK application")

	qrCode, err := qr.Unicode()
	if err != nil {
		log.Fatal("Failed to render QR code:", err)
	}
	fmt.Println(qrCode)

	fmt.Println("⏳ Waiting for peer to scan QR code...")
	fmt.Println("   🔐 Secure encrypted connection will be established")
	fmt.Println("   🛑 Press Ctrl+C to cancel")
	fmt.Println()

	// Wait for peer connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	peer, err := qr.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Fatal("❌ No peer connected within timeout. Try running another instance of this program.")
		}
		log.Fatal("❌ Failed to connect to peer:", err)
	}

	fmt.Printf("✅ Peer connected: %s\n", peer.DID())
	fmt.Println("🔐 Secure encrypted channel established")
	fmt.Println()

	return peer
}

// demonstrateChat shows basic chat functionality with the connected peer
func demonstrateChat(chatClient *client.Client, peer *client.Peer) {
	fmt.Println("💬 Demonstrating chat messaging...")

	// Send initial greeting
	greeting := fmt.Sprintf("🎉 Hello! Chat demo started at %s. This message is end-to-end encrypted!",
		time.Now().Format("15:04:05"))

	fmt.Println("📤 Sending initial greeting...")
	err := chatClient.Chat().Send(peer.DID(), greeting)
	if err != nil {
		log.Printf("Failed to send greeting: %v", err)
		return
	}
	fmt.Printf("✅ Greeting sent: \"%s\"\n", greeting)

	// Send a few demo messages to showcase different responses
	demoMessages := []string{
		"Hello there!",
		"How are you?",
		"What time is it?",
		"This is a test message",
	}

	fmt.Println("\n📤 Sending demo messages...")
	for i, msg := range demoMessages {
		time.Sleep(2 * time.Second) // Small delay between messages

		fmt.Printf("📤 [%d/%d] Sending: \"%s\"\n", i+1, len(demoMessages), msg)
		err := chatClient.Chat().Send(peer.DID(), msg)
		if err != nil {
			fmt.Printf("❌ Failed to send message: %v\n", err)
		} else {
			fmt.Printf("✅ Message sent successfully\n")
		}
	}

	fmt.Println("\n🎯 Demo messages sent!")
	fmt.Println("   • Each message is automatically encrypted")
	fmt.Println("   • Responses are generated based on message content")
	fmt.Println("   • Try sending messages from the other client to see real-time chat")
	fmt.Println()
}
