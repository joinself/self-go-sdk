package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joinself/self-go-sdk/client"
)

func main() {
	// Check for required environment variables
	storageKey := os.Getenv("SELF_STORAGE_KEY")
	storagePath := os.Getenv("SELF_STORAGE_PATH")
	if storageKey == "" || storagePath == "" {
		log.Fatal("Please set SELF_STORAGE_KEY and SELF_STORAGE_PATH environment variables")
	}

	// Create client configuration
	config := client.Config{
		StorageKey:  []byte(storageKey),
		StoragePath: storagePath,
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	}

	// Create the client
	selfClient, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer selfClient.Close()

	fmt.Printf("🚀 Advanced Features Demo - Client DID: %s\n\n", selfClient.DID())

	// === STORAGE FEATURES ===
	fmt.Println("📦 Storage Features Demo")
	demonstrateStorage(selfClient)
	fmt.Println()

	// === NOTIFICATION FEATURES ===
	fmt.Println("🔔 Notification Features Demo")
	demonstrateNotifications(selfClient)
	fmt.Println()

	// === ACCOUNT PAIRING FEATURES ===
	fmt.Println("🔗 Account Pairing Features Demo")
	demonstratePairing(selfClient)
	fmt.Println()

	fmt.Println("✅ Advanced features demo completed!")
}

func demonstrateStorage(selfClient *client.Client) {
	storage := selfClient.Storage()

	// Basic storage operations
	fmt.Println("  • Basic storage operations")

	// Store simple values
	err := storage.StoreString("user:name", "Alice")
	if err != nil {
		log.Printf("Failed to store string: %v", err)
		return
	}

	// Store JSON data
	userData := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"email": "alice@example.com",
	}
	err = storage.StoreJSON("user:profile", userData)
	if err != nil {
		log.Printf("Failed to store JSON: %v", err)
		return
	}

	// Retrieve values
	name, err := storage.LookupString("user:name")
	if err != nil {
		log.Printf("Failed to lookup string: %v", err)
	} else {
		fmt.Printf("    Retrieved name: %s\n", name)
	}

	var profile map[string]interface{}
	err = storage.LookupJSON("user:profile", &profile)
	if err != nil {
		log.Printf("Failed to lookup JSON: %v", err)
	} else {
		fmt.Printf("    Retrieved profile: %+v\n", profile)
	}

	// Namespaced storage
	fmt.Println("  • Namespaced storage")
	userStorage := storage.Namespace("user")

	err = userStorage.StoreString("preferences", "dark_mode=true")
	if err != nil {
		log.Printf("Failed to store in namespace: %v", err)
	} else {
		fmt.Println("    Stored user preferences in namespace")
	}

	// Temporary storage with expiry
	fmt.Println("  • Temporary storage with expiry")
	err = storage.StoreTemporaryString("session:token", "abc123", 5*time.Second)
	if err != nil {
		log.Printf("Failed to store temporary value: %v", err)
	} else {
		fmt.Println("    Stored temporary session token (expires in 5 seconds)")
	}

	// Cache functionality
	fmt.Println("  • Cache functionality")
	cache := storage.Cache("api")

	err = cache.SetString("response:users", `[{"id":1,"name":"Alice"}]`)
	if err != nil {
		log.Printf("Failed to cache value: %v", err)
	} else {
		fmt.Println("    Cached API response")
	}

	if cache.Has("response:users") {
		cachedData, err := cache.GetString("response:users")
		if err != nil {
			log.Printf("Failed to get cached value: %v", err)
		} else {
			fmt.Printf("    Retrieved from cache: %s\n", cachedData)
		}
	}
}

func demonstrateNotifications(selfClient *client.Client) {
	notifications := selfClient.Notifications()

	// Register notification sent handler
	notifications.OnNotificationSent(func(peerDID string, summary *client.NotificationSummary) {
		fmt.Printf("    📤 Notification sent to %s: %s\n", peerDID, summary.Title)
	})

	// For demo purposes, we'll use the client's own DID as the target
	// In a real application, this would be another user's DID
	targetDID := selfClient.DID()

	// Send different types of notifications
	fmt.Println("  • Sending various notification types")

	// Chat notification
	err := notifications.SendChatNotification(targetDID, "Hello! This is a test message.")
	if err != nil {
		log.Printf("Failed to send chat notification: %v", err)
	}

	// Credential notification
	err = notifications.SendCredentialNotification(targetDID, "identity", "request")
	if err != nil {
		log.Printf("Failed to send credential notification: %v", err)
	}

	// Group invite notification
	err = notifications.SendGroupInviteNotification(targetDID, "Development Team", "Alice")
	if err != nil {
		log.Printf("Failed to send group invite notification: %v", err)
	}

	// Custom notification
	err = notifications.SendCustomNotification(
		targetDID,
		"System Alert",
		"Your account security settings have been updated",
		"security",
	)
	if err != nil {
		log.Printf("Failed to send custom notification: %v", err)
	}

	fmt.Println("    ✅ Notification examples completed")
}

func demonstratePairing(selfClient *client.Client) {
	pairing := selfClient.Pairing()

	// Get pairing code
	fmt.Println("  • Getting pairing code")
	pairingCode, err := pairing.GetPairingCode()
	if err != nil {
		log.Printf("Failed to get pairing code: %v", err)
		return
	}

	fmt.Printf("    Pairing Code: %s\n", pairingCode.Code)
	fmt.Printf("    Unpaired: %t\n", pairingCode.Unpaired)
	fmt.Printf("    Expires: %s\n", pairingCode.ExpiresAt.Format(time.RFC3339))

	// Generate QR code representation
	qrCode, err := pairing.GeneratePairingQR()
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
	} else {
		fmt.Printf("    QR Code: %s\n", qrCode)
	}

	// Check if paired
	isPaired, err := pairing.IsPaired()
	if err != nil {
		log.Printf("Failed to check pairing status: %v", err)
	} else {
		fmt.Printf("    Is Paired: %t\n", isPaired)
	}

	// Register pairing event handlers
	fmt.Println("  • Setting up pairing event handlers")

	pairing.OnPairingRequest(func(request *client.IncomingPairingRequest) {
		fmt.Printf("    📥 Pairing request from: %s\n", request.From())
		fmt.Printf("       Address: %s\n", request.Address().String())
		fmt.Printf("       Roles: %d\n", request.Roles())
		fmt.Printf("       Expires: %s\n", request.Expires().Format(time.RFC3339))

		// In a real application, you would prompt the user or check permissions
		// For demo purposes, we'll auto-reject to avoid actual pairing
		fmt.Println("       🚫 Auto-rejecting for demo purposes")
		err := request.Reject()
		if err != nil {
			log.Printf("Failed to reject pairing request: %v", err)
		}
	})

	pairing.OnPairingResponse(func(response *client.PairingResponse) {
		fmt.Printf("    📨 Pairing response from: %s\n", response.From())
		fmt.Printf("       Status: %d\n", response.Status())
		if response.Operation() != nil {
			fmt.Println("       ✅ Operation included")
		}
		if len(response.Assets()) > 0 {
			fmt.Printf("       📎 %d assets included\n", len(response.Assets()))
		}
	})

	// Example of sending a pairing request (commented out to avoid actual pairing)
	fmt.Println("  • Pairing request example (not executed)")
	fmt.Println("    // Create a signing key for the request")
	fmt.Println("    // signingKey, _ := signing.NewKey()")
	fmt.Println("    // request, err := pairing.RequestPairing(")
	fmt.Println("    //     \"target_did\",")
	fmt.Println("    //     signingKey.PublicKey(),")
	fmt.Println("    //     identity.RoleOwner,")
	fmt.Println("    // )")
	fmt.Println("    // if err == nil {")
	fmt.Println("    //     ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)")
	fmt.Println("    //     defer cancel()")
	fmt.Println("    //     response, err := request.WaitForResponse(ctx)")
	fmt.Println("    //     // Handle response...")
	fmt.Println("    // }")

	fmt.Println("    ✅ Pairing examples completed")
}

// Example of advanced storage patterns
func demonstrateAdvancedStoragePatterns(selfClient *client.Client) {
	storage := selfClient.Storage()

	// Session management
	sessionStorage := storage.Namespace("session")

	// Store session with expiry
	sessionData := map[string]interface{}{
		"user_id":     "12345",
		"created_at":  time.Now(),
		"permissions": []string{"read", "write"},
	}

	err := sessionStorage.StoreJSONWithExpiry("current", sessionData, time.Now().Add(24*time.Hour))
	if err != nil {
		log.Printf("Failed to store session: %v", err)
	}

	// User preferences with namespace
	userPrefs := storage.Namespace("user:12345")

	preferences := map[string]interface{}{
		"theme":         "dark",
		"notifications": true,
		"language":      "en",
	}

	err = userPrefs.StoreJSON("preferences", preferences)
	if err != nil {
		log.Printf("Failed to store preferences: %v", err)
	}

	// Cache with TTL for API responses
	apiCache := storage.Cache("api")

	// Cache user data for 1 hour
	userData := `{"id": 12345, "name": "Alice", "email": "alice@example.com"}`
	err = apiCache.SetWithTTL("user:12345", []byte(userData), time.Hour)
	if err != nil {
		log.Printf("Failed to cache user data: %v", err)
	}

	// Check cache before making API call
	if apiCache.Has("user:12345") {
		cachedUser, err := apiCache.GetString("user:12345")
		if err == nil {
			fmt.Printf("Using cached user data: %s\n", cachedUser)
		}
	}
}

// Example of notification integration with other components
func demonstrateNotificationIntegration(selfClient *client.Client) {
	notifications := selfClient.Notifications()
	chat := selfClient.Chat()

	// Send a chat message and notification
	targetDID := "example_peer_did"
	message := "Hello! How are you doing?"

	// Send the chat message
	err := chat.Send(targetDID, message)
	if err != nil {
		log.Printf("Failed to send chat message: %v", err)
		return
	}

	// Send a notification about the chat message
	err = notifications.SendChatNotification(targetDID, message)
	if err != nil {
		log.Printf("Failed to send chat notification: %v", err)
	}

	// For group chats, send group notifications
	groupChats := selfClient.GroupChats()
	groups := groupChats.ListGroups()

	for _, group := range groups {
		for _, member := range group.Members() {
			if member.DID != selfClient.DID() {
				err := notifications.SendGroupChatNotification(
					member.DID,
					group.Name(),
					"New message in group chat",
				)
				if err != nil {
					log.Printf("Failed to send group notification to %s: %v", member.DID, err)
				}
			}
		}
	}
}
