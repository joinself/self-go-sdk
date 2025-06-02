// Package main demonstrates advanced storage capabilities of the Self SDK.
//
// This is the STORAGE level of advanced features examples.
// Start here to learn the foundation of advanced Self SDK storage patterns.
//
// This example shows:
// - Encrypted local storage with automatic security
// - Namespacing for organized data management
// - TTL (Time To Live) for automatic data expiry
// - Caching for performance optimization
// - Different data types and storage patterns
//
// 🎯 What you'll learn:
// • How to organize data with namespaces
// • Temporary storage with automatic expiry
// • Performance optimization with caching
// • Different storage patterns for various use cases
// • Encrypted storage security benefits
//
// 📚 Next steps:
// • ../notifications/main.go - Push notification system
// • ../pairing/main.go - Account pairing and multi-device sync
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
	fmt.Println("📦 Advanced Storage Demo")
	fmt.Println("========================")
	fmt.Println("This demo showcases advanced Self SDK storage capabilities.")
	fmt.Println("📚 This is the STORAGE level - foundation of advanced features.")
	fmt.Println()

	// Step 1: Create a Self client for storage demonstrations
	storageClient := createStorageClient()
	defer storageClient.Close()

	fmt.Printf("🆔 Client DID: %s\n", storageClient.DID())
	fmt.Println()

	// Step 2: Demonstrate basic storage operations
	demonstrateBasicStorage(storageClient)

	// Step 3: Show namespaced storage for organization
	demonstrateNamespacedStorage(storageClient)

	// Step 4: Explore temporary storage with TTL
	demonstrateTemporaryStorage(storageClient)

	// Step 5: Show cache management for performance
	demonstrateCacheManagement(storageClient)

	fmt.Println("✅ Storage capabilities demo completed!")
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Explored basic encrypted storage operations")
	fmt.Println("   2. Organized data using namespaces")
	fmt.Println("   3. Used TTL for automatic data expiry")
	fmt.Println("   4. Implemented caching for performance optimization")
	fmt.Println()
	fmt.Println("🎯 Storage benefits:")
	fmt.Println("   • Automatic encryption for all stored data")
	fmt.Println("   • Organized data management with namespaces")
	fmt.Println("   • Automatic cleanup with TTL expiry")
	fmt.Println("   • Performance optimization with intelligent caching")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run ../notifications/main.go to learn about push notifications")
	fmt.Println("   • Run ../pairing/main.go for account pairing capabilities")
	fmt.Println("   • Run ../production_patterns/main.go for real-world patterns")
	fmt.Println("   • Run ../integration/main.go for component integration")
}

// createStorageClient sets up a Self client for storage demonstrations
func createStorageClient() *client.Client {
	fmt.Println("🔧 Setting up storage client...")

	storageClient, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("storage_demo"),
		StoragePath: "./storage_demo_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create storage client:", err)
	}

	fmt.Println("✅ Storage client created successfully")
	return storageClient
}

// demonstrateBasicStorage shows fundamental storage operations
func demonstrateBasicStorage(selfClient *client.Client) {
	fmt.Println("🔹 Basic Storage Operations")
	fmt.Println("===========================")
	fmt.Println("Learning the fundamentals of Self SDK storage...")
	fmt.Println()

	storage := selfClient.Storage()

	// Store different data types
	fmt.Println("📝 Storing different data types...")

	// String storage
	err := storage.StoreString("user:name", "Alice Johnson")
	if err != nil {
		log.Printf("Failed to store string: %v", err)
		return
	}
	fmt.Println("   ✅ Stored string: user name")

	// JSON storage
	userData := map[string]interface{}{
		"name":  "Alice Johnson",
		"age":   30,
		"email": "alice@example.com",
		"role":  "developer",
		"settings": map[string]interface{}{
			"theme":         "dark",
			"notifications": true,
			"language":      "en",
		},
	}
	err = storage.StoreJSON("user:profile", userData)
	if err != nil {
		log.Printf("Failed to store JSON: %v", err)
		return
	}
	fmt.Println("   ✅ Stored JSON: user profile with nested data")

	// Retrieve and verify data
	fmt.Println("\n📖 Retrieving stored data...")

	name, err := storage.LookupString("user:name")
	if err != nil {
		log.Printf("Failed to lookup string: %v", err)
	} else {
		fmt.Printf("   ✅ Retrieved name: %s\n", name)
	}

	var profile map[string]interface{}
	err = storage.LookupJSON("user:profile", &profile)
	if err != nil {
		log.Printf("Failed to lookup JSON: %v", err)
	} else {
		fmt.Printf("   ✅ Retrieved profile: %s (%s)\n", profile["name"], profile["role"])
		if settings, ok := profile["settings"].(map[string]interface{}); ok {
			fmt.Printf("   🎨 Theme: %s, Language: %s\n", settings["theme"], settings["language"])
		}
	}

	fmt.Println("\n🔐 Security Note:")
	fmt.Println("   • All data is automatically encrypted at rest")
	fmt.Println("   • No additional encryption setup required")
	fmt.Println("   • Data is tied to your specific client identity")
	fmt.Println()
}

// demonstrateNamespacedStorage shows organized storage patterns
func demonstrateNamespacedStorage(selfClient *client.Client) {
	fmt.Println("🔹 Namespaced Storage")
	fmt.Println("=====================")
	fmt.Println("Organizing data with logical namespaces...")
	fmt.Println()

	storage := selfClient.Storage()

	// User-specific namespace
	fmt.Println("👤 User namespace - personal data:")
	userStorage := storage.Namespace("user")
	err := userStorage.StoreString("preferences", "dark_mode=true,lang=en,notifications=true")
	if err != nil {
		log.Printf("Failed to store in user namespace: %v", err)
	} else {
		fmt.Println("   ✅ Stored user preferences")
	}

	userSettings := map[string]interface{}{
		"dashboard_layout": "grid",
		"auto_save":        true,
		"privacy_level":    "standard",
		"last_login":       time.Now().Format("2006-01-02 15:04:05"),
	}
	err = userStorage.StoreJSON("settings", userSettings)
	if err != nil {
		log.Printf("Failed to store user settings: %v", err)
	} else {
		fmt.Println("   ✅ Stored user settings")
	}

	// Application settings namespace
	fmt.Println("\n🔧 Application namespace - app configuration:")
	appStorage := storage.Namespace("app")
	appSettings := map[string]interface{}{
		"version":     "1.2.3",
		"environment": "development",
		"debug_mode":  true,
		"features": map[string]bool{
			"chat":          true,
			"groups":        true,
			"credentials":   true,
			"notifications": true,
			"advanced_ui":   false,
		},
		"api_endpoints": map[string]string{
			"auth":    "https://auth-dev.example.com",
			"api":     "https://api-dev.example.com",
			"storage": "https://storage-dev.example.com",
		},
	}
	err = appStorage.StoreJSON("config", appSettings)
	if err != nil {
		log.Printf("Failed to store app settings: %v", err)
	} else {
		fmt.Println("   ✅ Stored application configuration")
	}

	// Session namespace for temporary data
	fmt.Println("\n🔑 Session namespace - temporary session data:")
	sessionStorage := storage.Namespace("session")
	sessionData := map[string]interface{}{
		"session_id":    "sess_abc123def456",
		"user_id":       "user_12345",
		"created_at":    time.Now(),
		"last_activity": time.Now(),
		"permissions":   []string{"read", "write", "admin"},
		"device_info": map[string]string{
			"type":      "desktop",
			"browser":   "chrome",
			"os":        "macos",
			"ip_region": "US-West",
		},
	}
	err = sessionStorage.StoreJSON("current", sessionData)
	if err != nil {
		log.Printf("Failed to store session: %v", err)
	} else {
		fmt.Println("   ✅ Stored session data")
	}

	// Demonstrate namespace isolation
	fmt.Println("\n🔍 Verifying namespace isolation:")
	userPrefs, err := userStorage.LookupString("preferences")
	if err == nil {
		fmt.Printf("   ✅ User namespace accessible: %s\n", userPrefs[:30]+"...")
	}

	var appConfig map[string]interface{}
	err = appStorage.LookupJSON("config", &appConfig)
	if err == nil {
		fmt.Printf("   ✅ App namespace accessible: version %s\n", appConfig["version"])
	}

	fmt.Println("\n💡 Namespace Benefits:")
	fmt.Println("   • Logical separation of different data types")
	fmt.Println("   • Prevents naming conflicts between components")
	fmt.Println("   • Easier data management and cleanup")
	fmt.Println("   • Clear data ownership and access patterns")
	fmt.Println()
}

// demonstrateTemporaryStorage shows TTL-based storage
func demonstrateTemporaryStorage(selfClient *client.Client) {
	fmt.Println("🔹 Temporary Storage with TTL")
	fmt.Println("=============================")
	fmt.Println("Automatic data expiry for temporary information...")
	fmt.Println()

	storage := selfClient.Storage()

	// Short-lived session token
	fmt.Println("⏰ Creating short-lived data:")
	err := storage.StoreTemporaryString("session:token", "temp_token_abc123", 10*time.Second)
	if err != nil {
		log.Printf("Failed to store temporary token: %v", err)
	} else {
		fmt.Println("   ✅ Stored temporary session token (expires in 10 seconds)")
	}

	// Temporary verification code
	err = storage.StoreTemporaryString("verification:code", "987654", 5*time.Minute)
	if err != nil {
		log.Printf("Failed to store verification code: %v", err)
	} else {
		fmt.Println("   ✅ Stored verification code (expires in 5 minutes)")
	}

	// Temporary user state
	tempState := map[string]interface{}{
		"onboarding_step":   3,
		"tutorial_shown":    true,
		"temp_preferences":  map[string]string{"theme": "auto", "lang": "en"},
		"wizard_progress":   75,
		"unsaved_changes":   true,
		"last_auto_save":    time.Now(),
		"temp_data_expires": time.Now().Add(1 * time.Hour),
	}
	err = storage.StoreTemporaryJSON("user:temp_state", tempState, 1*time.Hour)
	if err != nil {
		log.Printf("Failed to store temporary state: %v", err)
	} else {
		fmt.Println("   ✅ Stored temporary user state (expires in 1 hour)")
	}

	// Demonstrate immediate retrieval
	fmt.Println("\n📖 Verifying temporary data exists:")
	token, err := storage.LookupString("session:token")
	if err == nil {
		fmt.Printf("   ✅ Token retrieved: %s...\n", token[:10])
	} else {
		fmt.Printf("   ⏰ Token may have expired: %v\n", err)
	}

	code, err := storage.LookupString("verification:code")
	if err == nil {
		fmt.Printf("   ✅ Verification code: %s\n", code)
	}

	var tempUserState map[string]interface{}
	err = storage.LookupJSON("user:temp_state", &tempUserState)
	if err == nil {
		fmt.Printf("   ✅ Temp state: step %v, progress %v%%\n",
			tempUserState["onboarding_step"], tempUserState["wizard_progress"])
	}

	fmt.Println("\n⚡ TTL Benefits:")
	fmt.Println("   • Automatic cleanup prevents storage bloat")
	fmt.Println("   • Perfect for session tokens and temporary data")
	fmt.Println("   • No manual cleanup required")
	fmt.Println("   • Configurable expiry times for different use cases")
	fmt.Println()
}

// demonstrateCacheManagement shows caching patterns
func demonstrateCacheManagement(selfClient *client.Client) {
	fmt.Println("🔹 Cache Management")
	fmt.Println("===================")
	fmt.Println("Performance optimization with intelligent caching...")
	fmt.Println()

	storage := selfClient.Storage()

	// API response cache
	fmt.Println("🗄️ Setting up API response cache:")
	apiCache := storage.Cache("api")

	// Cache user list
	userListJSON := `[
		{"id": 1, "name": "Alice", "role": "admin", "last_seen": "2024-01-15T10:30:00Z"},
		{"id": 2, "name": "Bob", "role": "user", "last_seen": "2024-01-15T09:15:00Z"},
		{"id": 3, "name": "Charlie", "role": "moderator", "last_seen": "2024-01-15T11:45:00Z"},
		{"id": 4, "name": "Diana", "role": "user", "last_seen": "2024-01-15T08:20:00Z"}
	]`
	err := apiCache.SetString("users:list", userListJSON)
	if err != nil {
		log.Printf("Failed to cache user list: %v", err)
	} else {
		fmt.Println("   ✅ Cached user list (no expiry)")
	}

	// Cache with TTL
	profileData := `{
		"id": 123,
		"name": "Alice Johnson",
		"email": "alice@example.com",
		"avatar": "https://example.com/avatar.jpg",
		"last_login": "2024-01-15T10:30:00Z",
		"preferences": {"theme": "dark", "language": "en"},
		"stats": {"messages_sent": 1247, "groups_joined": 8}
	}`
	err = apiCache.SetWithTTL("profile:123", []byte(profileData), 30*time.Minute)
	if err != nil {
		log.Printf("Failed to cache profile: %v", err)
	} else {
		fmt.Println("   ✅ Cached user profile (expires in 30 minutes)")
	}

	// Search results cache
	searchResults := `{
		"query": "development team",
		"results": [
			{"id": 1, "name": "Dev Team Alpha", "members": 12, "active": true},
			{"id": 2, "name": "Dev Team Beta", "members": 8, "active": true},
			{"id": 3, "name": "Dev Team Gamma", "members": 15, "active": false}
		],
		"total": 3,
		"cached_at": "2024-01-15T10:30:00Z",
		"search_time_ms": 45
	}`
	err = apiCache.SetWithTTL("search:development_team", []byte(searchResults), 15*time.Minute)
	if err != nil {
		log.Printf("Failed to cache search results: %v", err)
	} else {
		fmt.Println("   ✅ Cached search results (expires in 15 minutes)")
	}

	// Settings cache
	fmt.Println("\n⚙️ Setting up application cache:")
	settingsCache := storage.Cache("settings")
	err = settingsCache.SetString("ui:theme", "dark")
	if err != nil {
		log.Printf("Failed to cache theme: %v", err)
	} else {
		fmt.Println("   ✅ Cached UI theme setting")
	}

	err = settingsCache.SetString("ui:language", "en")
	if err != nil {
		log.Printf("Failed to cache language: %v", err)
	} else {
		fmt.Println("   ✅ Cached language setting")
	}

	// Demonstrate cache retrieval
	fmt.Println("\n📖 Demonstrating cache retrieval:")
	if apiCache.Has("users:list") {
		cachedUsers, err := apiCache.GetString("users:list")
		if err == nil {
			fmt.Printf("   ⚡ Cache hit: Retrieved user list (%d bytes)\n", len(cachedUsers))
		}
	}

	if apiCache.Has("profile:123") {
		cachedProfile, err := apiCache.Get("profile:123")
		if err == nil {
			fmt.Printf("   ⚡ Cache hit: Retrieved user profile (%d bytes)\n", len(cachedProfile))
		}
	}

	if settingsCache.Has("ui:theme") {
		theme, err := settingsCache.GetString("ui:theme")
		if err == nil {
			fmt.Printf("   ⚡ Cache hit: UI theme is '%s'\n", theme)
		}
	}

	fmt.Println("\n🚀 Cache Benefits:")
	fmt.Println("   • Dramatically faster data access")
	fmt.Println("   • Reduced API calls and network traffic")
	fmt.Println("   • Configurable TTL for different data types")
	fmt.Println("   • Automatic memory management")
	fmt.Println("   • Perfect for frequently accessed data")
	fmt.Println()
}
