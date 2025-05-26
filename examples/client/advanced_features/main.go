// Package main provides an overview of the Self SDK's advanced features.
//
// This is the entry point for learning advanced Self SDK capabilities.
// For focused learning, explore the individual tutorial directories:
//
// 📚 Educational Progression:
// 1. storage/ - Advanced storage capabilities (start here)
// 2. notifications/ - Push notification system
// 3. pairing/ - Account pairing and multi-device sync
// 4. production_patterns/ - Real-world storage and session management
// 5. integration/ - Component integration and workflows
//
// This overview shows a quick demonstration of each capability.
// For deep learning, visit each subdirectory for focused examples.
//
// 🎯 What you'll learn across all examples:
// • Advanced storage with namespacing, TTL, and caching
// • Push notification system for user engagement
// • Account pairing for multi-device experiences
// • Production-ready storage and session patterns
// • Integration between different SDK components
//
// 🚀 ADVANCED CAPABILITIES OVERVIEW:
// • Encrypted local storage with namespacing
// • Cache management with TTL support
// • Push notification delivery system
// • Account pairing and synchronization
// • Production storage patterns
// • Multi-component integration workflows
package main

import (
	"fmt"
	"log"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("🚀 Advanced Features Overview")
	fmt.Println("=============================")
	fmt.Println("Welcome to the Self SDK advanced features learning path!")
	fmt.Println()
	fmt.Println("This overview demonstrates the breadth of advanced capabilities.")
	fmt.Println("For focused learning, explore each subdirectory:")
	fmt.Println()

	// Create a client for quick demonstrations
	advancedClient := createClient()
	defer advancedClient.Close()

	fmt.Printf("🆔 Client DID: %s\n", advancedClient.DID())
	fmt.Println()

	// Quick overview of each capability
	demonstrateCapabilityOverview(advancedClient)

	fmt.Println("📚 Educational Learning Path")
	fmt.Println("============================")
	fmt.Println()
	fmt.Println("Follow this progression for optimal learning:")
	fmt.Println()
	fmt.Println("1️⃣  STORAGE (Complexity: 5/10)")
	fmt.Println("   📁 cd storage && go run main.go")
	fmt.Println("   🎯 Learn: Namespacing, TTL, caching, encrypted storage")
	fmt.Println("   ⏱️  Time: 15-20 minutes")
	fmt.Println()
	fmt.Println("2️⃣  NOTIFICATIONS (Complexity: 4/10)")
	fmt.Println("   🔔 cd notifications && go run main.go")
	fmt.Println("   🎯 Learn: Push notifications, event handling, user engagement")
	fmt.Println("   ⏱️  Time: 10-15 minutes")
	fmt.Println()
	fmt.Println("3️⃣  PAIRING (Complexity: 5/10)")
	fmt.Println("   🔗 cd pairing && go run main.go")
	fmt.Println("   🎯 Learn: Multi-device sync, QR pairing, device management")
	fmt.Println("   ⏱️  Time: 15-20 minutes")
	fmt.Println()
	fmt.Println("4️⃣  PRODUCTION PATTERNS (Complexity: 6/10)")
	fmt.Println("   🏭 cd production_patterns && go run main.go")
	fmt.Println("   🎯 Learn: Session management, state persistence, optimization")
	fmt.Println("   ⏱️  Time: 20-25 minutes")
	fmt.Println()
	fmt.Println("5️⃣  INTEGRATION (Complexity: 7/10)")
	fmt.Println("   🔄 cd integration && go run main.go")
	fmt.Println("   🎯 Learn: Multi-component workflows, coordinated features")
	fmt.Println("   ⏱️  Time: 20-30 minutes")
	fmt.Println()
	fmt.Println("📊 Total Learning Time: ~80-110 minutes")
	fmt.Println()
	fmt.Println("🎓 Prerequisites:")
	fmt.Println("   • Complete simple_chat, group_chat examples first")
	fmt.Println("   • Basic understanding of Go and Self SDK concepts")
	fmt.Println("   • Familiarity with storage and caching concepts")
	fmt.Println()
	fmt.Println("💡 Pro Tips:")
	fmt.Println("   • Follow the numbered progression for best learning")
	fmt.Println("   • Each example builds on previous concepts")
	fmt.Println("   • Take time to understand each pattern before moving on")
	fmt.Println("   • Experiment with the code to deepen understanding")
	fmt.Println()
	fmt.Println("🚀 Ready to start? Begin with: cd storage && go run main.go")
}

// createClient sets up a Self client for demonstrations
func createClient() *client.Client {
	fmt.Println("🔧 Setting up overview client...")

	client, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("advanced_overview"),
		StoragePath: "./advanced_overview_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fmt.Println("✅ Overview client created successfully")
	return client
}

// demonstrateCapabilityOverview shows a quick preview of each advanced capability
func demonstrateCapabilityOverview(selfClient *client.Client) {
	fmt.Println("🔍 Quick Capability Overview")
	fmt.Println("============================")
	fmt.Println()

	// Storage overview
	fmt.Println("📦 STORAGE CAPABILITIES")
	fmt.Println("   • Encrypted local storage with automatic security")
	fmt.Println("   • Namespacing for organized data management")
	fmt.Println("   • TTL (Time To Live) for automatic data expiry")
	fmt.Println("   • Caching for performance optimization")
	storage := selfClient.Storage()
	err := storage.StoreString("overview:demo", "Advanced storage working!")
	if err == nil {
		fmt.Println("   ✅ Storage system operational")
	}
	fmt.Println()

	// Notifications overview
	fmt.Println("🔔 NOTIFICATION SYSTEM")
	fmt.Println("   • Push notifications for real-time user engagement")
	fmt.Println("   • Multiple notification types (chat, credential, custom)")
	fmt.Println("   • Event-driven notification handling")
	fmt.Println("   • Delivery tracking and status management")
	notifications := selfClient.Notifications()
	fmt.Println("   ✅ Notification system available")
	_ = notifications // Demonstrate availability
	fmt.Println()

	// Pairing overview
	fmt.Println("🔗 ACCOUNT PAIRING")
	fmt.Println("   • Multi-device account synchronization")
	fmt.Println("   • QR code-based device pairing")
	fmt.Println("   • Secure cryptographic device verification")
	fmt.Println("   • Cross-device state management")
	pairing := selfClient.Pairing()
	fmt.Println("   ✅ Pairing system ready")
	_ = pairing // Demonstrate availability
	fmt.Println()

	// Production patterns overview
	fmt.Println("🏭 PRODUCTION PATTERNS")
	fmt.Println("   • Session management with automatic expiry")
	fmt.Println("   • Application state persistence")
	fmt.Println("   • Performance optimization strategies")
	fmt.Println("   • Scalable data access patterns")
	fmt.Println("   ✅ Production patterns demonstrated in examples")
	fmt.Println()

	// Integration overview
	fmt.Println("🔄 COMPONENT INTEGRATION")
	fmt.Println("   • Coordinated workflows between SDK components")
	fmt.Println("   • Storage + Chat + Notifications integration")
	fmt.Println("   • Complex multi-feature applications")
	fmt.Println("   • Real-world application architecture patterns")
	fmt.Println("   ✅ Integration patterns ready for exploration")
	fmt.Println()
}
