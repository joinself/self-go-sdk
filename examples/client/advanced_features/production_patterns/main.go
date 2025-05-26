// Package main demonstrates production-ready patterns using the Self SDK.
//
// This is the PRODUCTION PATTERNS level of advanced features examples.
// Prerequisites: Complete ../storage/main.go, ../notifications/main.go, and ../pairing/main.go first.
//
// This example shows:
// - Real-world session management with automatic expiry
// - Application state persistence patterns
// - Performance optimization strategies
// - Scalable data access patterns
// - Production-ready error handling and recovery
//
// 🎯 What you'll learn:
// • Production session management patterns
// • Application state persistence strategies
// • Performance optimization with caching
// • Error handling and recovery mechanisms
// • Scalable data architecture patterns
//
// 📚 Next steps:
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
	fmt.Println("🏭 Production Patterns Demo")
	fmt.Println("===========================")
	fmt.Println("This demo showcases production-ready Self SDK patterns.")
	fmt.Println("📚 This is the PRODUCTION PATTERNS level - real-world applications.")
	fmt.Println()

	// Step 1: Create a Self client for production pattern demonstrations
	productionClient := createProductionClient()
	defer productionClient.Close()

	fmt.Printf("🆔 Client DID: %s\n", productionClient.DID())
	fmt.Println()

	// Step 2: Demonstrate session management patterns
	demonstrateSessionManagement(productionClient)

	// Step 3: Show application state persistence
	demonstrateStatePersistence(productionClient)

	// Step 4: Explore performance optimization strategies
	demonstratePerformanceOptimization(productionClient)

	// Step 5: Show error handling and recovery patterns
	demonstrateErrorHandlingPatterns(productionClient)

	fmt.Println("✅ Production patterns demo completed!")
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Implemented robust session management with automatic expiry")
	fmt.Println("   2. Demonstrated application state persistence strategies")
	fmt.Println("   3. Applied performance optimization with intelligent caching")
	fmt.Println("   4. Showed error handling and recovery mechanisms")
	fmt.Println()
	fmt.Println("🎯 Production benefits:")
	fmt.Println("   • Scalable and maintainable application architecture")
	fmt.Println("   • Robust error handling and recovery mechanisms")
	fmt.Println("   • Optimized performance for real-world usage")
	fmt.Println("   • Production-ready session and state management")
	fmt.Println()
	fmt.Println("📚 Ready for the final level?")
	fmt.Println("   • Run ../integration/main.go for complete component integration")
}

// createProductionClient sets up a Self client for production pattern demonstrations
func createProductionClient() *client.Client {
	fmt.Println("🔧 Setting up production client...")

	productionClient, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("production_demo"),
		StoragePath: "./production_demo_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create production client:", err)
	}

	fmt.Println("✅ Production client created successfully")
	return productionClient
}

// demonstrateSessionManagement shows production session handling patterns
func demonstrateSessionManagement(selfClient *client.Client) {
	fmt.Println("🔹 Session Management Patterns")
	fmt.Println("==============================")
	fmt.Println("Implementing robust session management for production applications...")
	fmt.Println()

	storage := selfClient.Storage()
	sessionStorage := storage.Namespace("session")

	// Create comprehensive session data
	fmt.Println("🔑 Creating production session...")
	sessionData := map[string]interface{}{
		"session_id":     generateSessionID(),
		"user_id":        "user_12345",
		"created_at":     time.Now(),
		"last_activity":  time.Now(),
		"expires_at":     time.Now().Add(24 * time.Hour),
		"permissions":    []string{"read", "write", "admin", "manage_users"},
		"security_level": "high",
		"device_info": map[string]string{
			"type":       "desktop",
			"browser":    "chrome",
			"os":         "macos",
			"ip_region":  "US-West",
			"user_agent": "Self-SDK-Client/1.0",
			"device_id":  "device_abc123",
		},
		"feature_flags": map[string]bool{
			"advanced_ui":   true,
			"beta_features": false,
			"debug_mode":    false,
			"analytics":     true,
		},
		"preferences": map[string]interface{}{
			"theme":         "dark",
			"language":      "en",
			"timezone":      "America/Los_Angeles",
			"notifications": true,
			"auto_save":     true,
		},
	}

	// Store session with automatic expiry
	err := sessionStorage.StoreJSONWithExpiry("current", sessionData, time.Now().Add(24*time.Hour))
	if err != nil {
		log.Printf("Failed to store session: %v", err)
	} else {
		fmt.Println("   ✅ Production session stored with 24-hour expiry")
	}

	// Store refresh token separately with longer expiry
	refreshToken := generateRefreshToken()
	err = storage.StoreTemporaryString("session:refresh_token", refreshToken, 7*24*time.Hour)
	if err != nil {
		log.Printf("Failed to store refresh token: %v", err)
	} else {
		fmt.Println("   ✅ Refresh token stored with 7-day expiry")
	}

	// Demonstrate session validation
	fmt.Println("\n🔍 Session validation and management:")
	isValid := validateSession(sessionStorage)
	if isValid {
		fmt.Println("   ✅ Session is valid and active")

		// Update last activity
		updateSessionActivity(sessionStorage)
		fmt.Println("   🔄 Session activity updated")
	} else {
		fmt.Println("   ❌ Session is invalid or expired")
		fmt.Println("   🔄 Session refresh or re-authentication required")
	}

	// Demonstrate session cleanup
	fmt.Println("\n🧹 Session cleanup patterns:")
	fmt.Println("   • Automatic expiry prevents stale sessions")
	fmt.Println("   • Refresh tokens enable seamless renewal")
	fmt.Println("   • Activity tracking for security monitoring")
	fmt.Println("   • Graceful session termination on logout")
	fmt.Println()
}

// demonstrateStatePersistence shows application state management
func demonstrateStatePersistence(selfClient *client.Client) {
	fmt.Println("🔹 Application State Persistence")
	fmt.Println("================================")
	fmt.Println("Managing application state for production reliability...")
	fmt.Println()

	storage := selfClient.Storage()

	// Application configuration management
	fmt.Println("⚙️ Application configuration management:")
	appStorage := storage.Namespace("app")
	appConfig := map[string]interface{}{
		"version":          "2.1.0",
		"environment":      "production",
		"debug_enabled":    false,
		"maintenance_mode": false,
		"feature_flags": map[string]bool{
			"new_ui":          true,
			"beta_features":   false,
			"advanced_search": true,
			"real_time_sync":  true,
			"analytics":       true,
		},
		"api_endpoints": map[string]string{
			"auth":          "https://auth.example.com",
			"api":           "https://api.example.com",
			"storage":       "https://storage.example.com",
			"notifications": "https://notifications.example.com",
			"websocket":     "wss://ws.example.com",
		},
		"limits": map[string]int{
			"max_file_size":   10485760, // 10MB
			"max_connections": 1000,
			"rate_limit_rpm":  60,
			"session_timeout": 3600,
		},
		"security": map[string]interface{}{
			"encryption_enabled":  true,
			"two_factor_required": true,
			"password_min_length": 12,
			"session_rotation":    true,
		},
	}

	err := appStorage.StoreJSON("config", appConfig)
	if err != nil {
		log.Printf("Failed to store app config: %v", err)
	} else {
		fmt.Println("   ✅ Application configuration persisted")
	}

	// User-specific state management
	fmt.Println("\n👤 User state management:")
	userStorage := storage.Namespace("user:12345")
	userState := map[string]interface{}{
		"profile": map[string]interface{}{
			"name":        "Alice Johnson",
			"email":       "alice@example.com",
			"role":        "admin",
			"department":  "Engineering",
			"last_login":  time.Now().Format("2006-01-02 15:04:05"),
			"login_count": 247,
		},
		"preferences": map[string]interface{}{
			"theme":            "dark",
			"language":         "en",
			"timezone":         "America/Los_Angeles",
			"notifications":    true,
			"auto_save":        true,
			"privacy_level":    "standard",
			"dashboard_layout": "grid",
		},
		"activity": map[string]interface{}{
			"messages_sent":      1247,
			"credentials_issued": 23,
			"groups_joined":      8,
			"files_uploaded":     156,
			"last_active":        time.Now(),
		},
		"settings": map[string]interface{}{
			"two_factor_enabled": true,
			"backup_enabled":     true,
			"sync_enabled":       true,
			"analytics_opt_in":   true,
		},
	}

	err = userStorage.StoreJSON("state", userState)
	if err != nil {
		log.Printf("Failed to store user state: %v", err)
	} else {
		fmt.Println("   ✅ User state persisted")
	}

	// Application metrics and analytics
	fmt.Println("\n📊 Metrics and analytics storage:")
	metricsStorage := storage.Namespace("metrics")
	metrics := map[string]interface{}{
		"daily_active_users":  1250,
		"messages_sent_today": 5670,
		"credentials_issued":  89,
		"storage_usage_mb":    2048,
		"api_calls_today":     12450,
		"error_rate_percent":  0.02,
		"average_response_ms": 145,
		"uptime_percent":      99.98,
		"last_updated":        time.Now(),
	}

	err = metricsStorage.StoreJSON("daily", metrics)
	if err != nil {
		log.Printf("Failed to store metrics: %v", err)
	} else {
		fmt.Println("   ✅ Application metrics stored")
	}

	fmt.Println("\n🎯 State persistence benefits:")
	fmt.Println("   • Reliable application configuration management")
	fmt.Println("   • User preferences survive application restarts")
	fmt.Println("   • Metrics tracking for performance monitoring")
	fmt.Println("   • Disaster recovery and backup capabilities")
	fmt.Println()
}

// demonstratePerformanceOptimization shows caching and optimization strategies
func demonstratePerformanceOptimization(selfClient *client.Client) {
	fmt.Println("🔹 Performance Optimization")
	fmt.Println("===========================")
	fmt.Println("Implementing production-grade performance optimizations...")
	fmt.Println()

	storage := selfClient.Storage()

	// Multi-tier caching strategy
	fmt.Println("🚀 Multi-tier caching implementation:")

	// L1 Cache: Frequently accessed data (short TTL)
	l1Cache := storage.Cache("l1")

	// Cache user sessions for quick access
	userSession := `{
		"user_id": "12345",
		"name": "Alice Johnson",
		"role": "admin",
		"permissions": ["read", "write", "admin"],
		"last_activity": "2024-01-15T10:30:00Z"
	}`
	err := l1Cache.SetWithTTL("session:12345", []byte(userSession), 5*time.Minute)
	if err == nil {
		fmt.Println("   ✅ L1 Cache: User session (5min TTL)")
	}

	// L2 Cache: API responses (medium TTL)
	l2Cache := storage.Cache("l2")

	// Cache API responses
	apiResponse := `{
		"users": [
			{"id": 1, "name": "Alice", "role": "admin", "active": true},
			{"id": 2, "name": "Bob", "role": "user", "active": true},
			{"id": 3, "name": "Charlie", "role": "moderator", "active": false}
		],
		"total": 3,
		"cached_at": "2024-01-15T10:30:00Z"
	}`
	err = l2Cache.SetWithTTL("api:users:list", []byte(apiResponse), 30*time.Minute)
	if err == nil {
		fmt.Println("   ✅ L2 Cache: API response (30min TTL)")
	}

	// L3 Cache: Static data (long TTL)
	l3Cache := storage.Cache("l3")

	// Cache configuration data
	configData := `{
		"app_name": "Self SDK Demo",
		"version": "2.1.0",
		"features": ["chat", "credentials", "groups"],
		"supported_languages": ["en", "es", "fr", "de"],
		"max_file_size": 10485760
	}`
	err = l3Cache.SetWithTTL("config:app", []byte(configData), 24*time.Hour)
	if err == nil {
		fmt.Println("   ✅ L3 Cache: Configuration (24hr TTL)")
	}

	// Demonstrate cache hit performance
	fmt.Println("\n⚡ Cache performance demonstration:")

	// Simulate cache hits
	start := time.Now()
	if l1Cache.Has("session:12345") {
		data, _ := l1Cache.Get("session:12345")
		elapsed := time.Since(start)
		fmt.Printf("   🎯 L1 Cache hit: %d bytes in %v\n", len(data), elapsed)
	}

	start = time.Now()
	if l2Cache.Has("api:users:list") {
		data, _ := l2Cache.Get("api:users:list")
		elapsed := time.Since(start)
		fmt.Printf("   🎯 L2 Cache hit: %d bytes in %v\n", len(data), elapsed)
	}

	start = time.Now()
	if l3Cache.Has("config:app") {
		data, _ := l3Cache.Get("config:app")
		elapsed := time.Since(start)
		fmt.Printf("   🎯 L3 Cache hit: %d bytes in %v\n", len(data), elapsed)
	}

	// Database query optimization simulation
	fmt.Println("\n🗄️ Database optimization patterns:")
	dbCache := storage.Cache("database")

	// Cache expensive query results
	queryResult := `{
		"query": "SELECT * FROM users WHERE active = true ORDER BY last_login DESC",
		"results": [
			{"id": 1, "name": "Alice", "last_login": "2024-01-15T10:30:00Z"},
			{"id": 2, "name": "Bob", "last_login": "2024-01-15T09:15:00Z"}
		],
		"execution_time_ms": 245,
		"cached_at": "2024-01-15T10:30:00Z"
	}`
	err = dbCache.SetWithTTL("query:active_users", []byte(queryResult), 15*time.Minute)
	if err == nil {
		fmt.Println("   ✅ Database query cached (15min TTL)")
	}

	fmt.Println("\n🎯 Performance optimization benefits:")
	fmt.Println("   • Dramatically reduced response times")
	fmt.Println("   • Lower database and API load")
	fmt.Println("   • Improved user experience")
	fmt.Println("   • Scalable architecture for high traffic")
	fmt.Println("   • Intelligent cache invalidation strategies")
	fmt.Println()
}

// demonstrateErrorHandlingPatterns shows robust error handling
func demonstrateErrorHandlingPatterns(selfClient *client.Client) {
	fmt.Println("🔹 Error Handling & Recovery")
	fmt.Println("============================")
	fmt.Println("Implementing production-grade error handling and recovery...")
	fmt.Println()

	storage := selfClient.Storage()

	// Error logging and tracking
	fmt.Println("📝 Error logging and tracking:")
	errorStorage := storage.Namespace("errors")

	// Simulate error tracking
	errorLog := map[string]interface{}{
		"error_id":    generateErrorID(),
		"timestamp":   time.Now(),
		"level":       "error",
		"component":   "storage",
		"message":     "Failed to connect to storage backend",
		"stack_trace": "storage.go:123 -> client.go:456 -> main.go:789",
		"user_id":     "12345",
		"session_id":  "sess_abc123",
		"request_id":  "req_def456",
		"metadata": map[string]interface{}{
			"retry_count":    3,
			"last_retry":     time.Now().Add(-5 * time.Minute),
			"error_category": "network",
			"severity":       "high",
		},
	}

	err := errorStorage.StoreJSON(fmt.Sprintf("log_%d", time.Now().Unix()), errorLog)
	if err == nil {
		fmt.Println("   ✅ Error logged for analysis and monitoring")
	}

	// Circuit breaker pattern simulation
	fmt.Println("\n🔄 Circuit breaker pattern:")
	circuitStorage := storage.Namespace("circuit_breaker")

	circuitState := map[string]interface{}{
		"service_name":      "external_api",
		"state":             "closed", // closed, open, half_open
		"failure_count":     0,
		"success_count":     150,
		"last_failure":      nil,
		"next_retry":        nil,
		"failure_threshold": 5,
		"timeout_duration":  "30s",
		"reset_timeout":     "60s",
	}

	err = circuitStorage.StoreJSON("external_api", circuitState)
	if err == nil {
		fmt.Println("   ✅ Circuit breaker state tracked")
	}

	// Retry mechanism with exponential backoff
	fmt.Println("\n🔁 Retry mechanism with exponential backoff:")
	retryStorage := storage.Namespace("retry")

	retryConfig := map[string]interface{}{
		"operation":      "api_call",
		"max_retries":    5,
		"base_delay_ms":  100,
		"max_delay_ms":   30000,
		"backoff_factor": 2.0,
		"jitter_enabled": true,
		"retry_count":    0,
		"last_attempt":   time.Now(),
		"next_attempt":   time.Now().Add(100 * time.Millisecond),
	}

	err = retryStorage.StoreJSON("api_call_config", retryConfig)
	if err == nil {
		fmt.Println("   ✅ Retry configuration stored")
	}

	// Health check and monitoring
	fmt.Println("\n🏥 Health check and monitoring:")
	healthStorage := storage.Namespace("health")

	healthStatus := map[string]interface{}{
		"overall_status": "healthy",
		"last_check":     time.Now(),
		"components": map[string]interface{}{
			"storage": map[string]interface{}{
				"status":        "healthy",
				"response_time": "5ms",
				"last_error":    nil,
			},
			"notifications": map[string]interface{}{
				"status":        "healthy",
				"response_time": "12ms",
				"last_error":    nil,
			},
			"pairing": map[string]interface{}{
				"status":        "healthy",
				"response_time": "8ms",
				"last_error":    nil,
			},
		},
		"metrics": map[string]interface{}{
			"uptime_seconds":     86400,
			"memory_usage_mb":    256,
			"cpu_usage_percent":  15.5,
			"disk_usage_percent": 45.2,
		},
	}

	err = healthStorage.StoreJSON("status", healthStatus)
	if err == nil {
		fmt.Println("   ✅ Health status monitored and stored")
	}

	fmt.Println("\n🛡️ Error handling benefits:")
	fmt.Println("   • Comprehensive error logging and analysis")
	fmt.Println("   • Automatic recovery with circuit breakers")
	fmt.Println("   • Intelligent retry mechanisms")
	fmt.Println("   • Real-time health monitoring")
	fmt.Println("   • Proactive issue detection and resolution")
	fmt.Println()
}

// Helper functions for production patterns

func generateSessionID() string {
	return fmt.Sprintf("sess_%d_%s", time.Now().Unix(), utils.GenerateStorageKey("session")[:8])
}

func generateRefreshToken() string {
	return fmt.Sprintf("refresh_%d_%s", time.Now().Unix(), utils.GenerateStorageKey("refresh")[:16])
}

func generateErrorID() string {
	return fmt.Sprintf("err_%d_%s", time.Now().Unix(), utils.GenerateStorageKey("error")[:8])
}

func validateSession(sessionStorage *client.StorageNamespace) bool {
	// Simulate session validation
	var session map[string]interface{}
	err := sessionStorage.LookupJSON("current", &session)
	if err != nil {
		return false
	}

	// Check if session has expired
	if expiresAt, ok := session["expires_at"].(string); ok {
		if expiry, err := time.Parse(time.RFC3339, expiresAt); err == nil {
			return time.Now().Before(expiry)
		}
	}

	return true
}

func updateSessionActivity(sessionStorage *client.StorageNamespace) {
	// Update last activity timestamp
	var session map[string]interface{}
	err := sessionStorage.LookupJSON("current", &session)
	if err == nil {
		session["last_activity"] = time.Now()
		sessionStorage.StoreJSON("current", session)
	}
}
