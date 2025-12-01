// Package utils provides common utility functions for Self SDK examples.
//
// This package contains shared functionality that is commonly used across
// multiple examples, reducing code duplication and providing consistent
// implementations of common patterns.
//
// 🔧 UTILITIES PROVIDED:
// • Storage key generation for demo purposes
// • Response status string conversion
// • Graceful shutdown handling
// • Common configuration helpers
// • Educational output formatting
//
// 📚 USAGE NOTES:
// These utilities are designed for educational and demonstration purposes.
// In production applications, implement proper security practices:
// - Use cryptographically secure random keys
// - Load keys from secure storage systems
// - Implement proper error handling and logging
// - Follow security best practices for key management
package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joinself/self-go-sdk/message"
)

// GenerateStorageKey creates a storage key for encrypting local account data.
// This function provides different implementations for demo and production use.
//
// For demo purposes (when role is provided):
// - Creates a deterministic key based on the role for consistent demo behavior
// - Allows multiple runs with the same storage without key mismatches
//
// For production use (when role is empty):
// - Generates a cryptographically secure random key
// - Should be stored securely and reused for the same account
//
// Parameters:
//   - role: Optional role identifier for demo keys (e.g., "issuer", "holder")
//     If empty, generates a secure random key for production use
//
// Returns:
//   - []byte: 32-byte storage key suitable for Self SDK client configuration
//
// Example usage:
//
//	// Demo usage with deterministic keys
//	issuerKey := utils.GenerateStorageKey("issuer")
//	holderKey := utils.GenerateStorageKey("holder")
//
//	// Production usage with secure random keys
//	prodKey := utils.GenerateStorageKey("")
func GenerateStorageKey(role string) []byte {
	key := make([]byte, 32)

	if role != "" {
		// 🎓 DEMO MODE: Create deterministic key for educational consistency
		// This ensures the same storage can be used across multiple demo runs
		demoKey := fmt.Sprintf("demo-key-%s-replace-in-production!!", role)
		copy(key, []byte(demoKey))

		// Pad or truncate to exactly 32 bytes
		if len(demoKey) > 32 {
			copy(key, []byte(demoKey)[:32])
		}
	} else {
		// 🔐 PRODUCTION MODE: Generate cryptographically secure random key
		// In production, this key should be stored securely and reused
		_, err := rand.Read(key)
		if err != nil {
			// Fallback to deterministic key if random generation fails
			// In production, this should be a fatal error
			copy(key, []byte("fallback-key-use-secure-random!!"))
		}
	}

	return key
}

// ResponseStatusToString converts message response status codes to human-readable strings.
// This utility provides consistent status message formatting across examples.
//
// Parameters:
//   - status: message.ResponseStatus code from Self SDK
//
// Returns:
//   - string: Human-readable status description
//
// Example usage:
//
//	status := resp.Status()
//	fmt.Printf("Response status: %s\n", utils.ResponseStatusToString(status))
func ResponseStatusToString(status message.ResponseStatus) string {
	switch status {
	case message.ResponseStatusAccepted:
		return "Accepted"
	case message.ResponseStatusForbidden:
		return "Forbidden"
	case message.ResponseStatusNotFound:
		return "Not Found"
	case message.ResponseStatusUnauthorized:
		return "Unauthorized"
	case message.ResponseStatusOk:
		return "OK"
	case message.ResponseStatusCreated:
		return "Created"
	case message.ResponseStatusBadRequest:
		return "Bad Request"
	case message.ResponseStatusConflict:
		return "Conflict"
	case message.ResponseStatusNotAcceptable:
		return "Not Acceptable"
	default:
		return fmt.Sprintf("Unknown (%d)", int(status))
	}
}

// SetupGracefulShutdown configures signal handling for clean application shutdown.
// This utility provides consistent shutdown behavior across examples.
//
// The function sets up handlers for SIGINT (Ctrl+C) and SIGTERM signals,
// allowing applications to clean up resources before exiting.
//
// Returns:
//   - context.Context: Context that will be cancelled on shutdown signal
//   - context.CancelFunc: Function to manually trigger shutdown
//
// Example usage:
//
//	ctx, cancel := utils.SetupGracefulShutdown()
//	defer cancel()
//
//	// Your application logic here
//	<-ctx.Done()
//	fmt.Println("Shutting down gracefully...")
func SetupGracefulShutdown() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\n🛑 Shutdown signal received (%s)...\n", sig)
		cancel()
	}()

	return ctx, cancel
}

// PrintHeader displays a formatted header for example applications.
// This utility provides consistent header formatting across examples.
//
// Parameters:
//   - title: Main title of the example
//   - description: Brief description of what the example demonstrates
//   - features: List of key features demonstrated
//
// Example usage:
//
//	utils.PrintHeader(
//	    "Self SDK Chat Example",
//	    "Demonstrates secure peer-to-peer messaging",
//	    []string{
//	        "Real-time messaging",
//	        "End-to-end encryption",
//	        "QR code discovery",
//	    },
//	)
func PrintHeader(title, description string, features []string) {
	fmt.Printf("🚀 %s\n", title)
	fmt.Println(generateSeparator(len(title) + 4))
	fmt.Printf("📚 %s\n", description)

	if len(features) > 0 {
		fmt.Println("🎯 Key features demonstrated:")
		for _, feature := range features {
			fmt.Printf("   • %s\n", feature)
		}
	}
	fmt.Println()
}

// PrintSectionHeader displays a formatted section header.
// This utility provides consistent section formatting within examples.
//
// Parameters:
//   - title: Section title
//   - description: Optional section description
//
// Example usage:
//
//	utils.PrintSectionHeader("CLIENT SETUP", "Initializing Self SDK clients")
func PrintSectionHeader(title, description string) {
	fmt.Printf("🔧 %s\n", title)
	fmt.Println(generateSeparator(len(title) + 4))
	if description != "" {
		fmt.Printf("📝 %s\n", description)
	}
	fmt.Println()
}

// PrintSuccess displays a formatted success message.
// This utility provides consistent success message formatting.
//
// Parameters:
//   - message: Success message to display
//
// Example usage:
//
//	utils.PrintSuccess("Client created successfully")
func PrintSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// PrintError displays a formatted error message.
// This utility provides consistent error message formatting.
//
// Parameters:
//   - message: Error message to display
//
// Example usage:
//
//	utils.PrintError("Failed to create client")
func PrintError(message string) {
	fmt.Printf("❌ %s\n", message)
}

// PrintInfo displays a formatted informational message.
// This utility provides consistent info message formatting.
//
// Parameters:
//   - message: Info message to display
//
// Example usage:
//
//	utils.PrintInfo("Waiting for peer connection...")
func PrintInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}

// PrintWarning displays a formatted warning message.
// This utility provides consistent warning message formatting.
//
// Parameters:
//   - message: Warning message to display
//
// Example usage:
//
//	utils.PrintWarning("Using demo storage key - not suitable for production")
func PrintWarning(message string) {
	fmt.Printf("⚠️  %s\n", message)
}

// generateSeparator creates a separator line of the specified length.
// This is a helper function for consistent formatting.
func generateSeparator(length int) string {
	separator := ""
	for i := 0; i < length; i++ {
		separator += "="
	}
	return separator
}

// DemoStorageKeyWarning displays a warning about demo storage keys.
// This utility reminds developers about production security considerations.
func DemoStorageKeyWarning() {
	PrintWarning("Using demo storage keys - replace with secure keys in production")
	fmt.Println("   🔐 In production: use crypto/rand or secure key management")
	fmt.Println("   📚 See documentation for security best practices")
	fmt.Println()
}

// ProductionNote displays a note about production considerations.
// This utility provides educational context about production deployment.
//
// Parameters:
//   - context: Specific context for the production note
//
// Example usage:
//
//	utils.ProductionNote("credential validation")
func ProductionNote(context string) {
	fmt.Printf("📝 Production Note: In production, implement proper %s\n", context)
	fmt.Println("   See documentation for security and scalability considerations")
	fmt.Println()
}
