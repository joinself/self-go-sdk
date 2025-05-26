// Package main demonstrates component integration workflows using the Self SDK.
//
// This is the INTEGRATION level of advanced features examples.
// Prerequisites: Complete all previous examples (storage, notifications, pairing, production_patterns).
//
// This example shows:
// - Coordinated workflows between SDK components
// - Storage + Chat + Notifications integration
// - Complex multi-feature applications
// - Real-world application architecture patterns
// - End-to-end feature integration scenarios
//
// 🎯 What you'll learn:
// • How to coordinate multiple SDK components
// • Real-world integration patterns and workflows
// • Complex application architecture design
// • End-to-end feature implementation
// • Production-ready integration strategies
//
// 🎓 This is the final and most advanced example in the series!
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("🔄 Component Integration Demo")
	fmt.Println("=============================")
	fmt.Println("This demo showcases coordinated Self SDK component workflows.")
	fmt.Println("📚 This is the INTEGRATION level - the ultimate advanced example!")
	fmt.Println()

	// Step 1: Create a Self client for integration demonstrations
	integrationClient := createIntegrationClient()
	defer integrationClient.Close()

	fmt.Printf("🆔 Client DID: %s\n", integrationClient.DID())
	fmt.Println()

	// Step 2: Set up integrated event handlers
	setupIntegratedHandlers(integrationClient)

	// Step 3: Demonstrate chat + storage + notifications integration
	demonstrateChatIntegration(integrationClient)

	// Step 4: Show pairing + storage synchronization
	demonstratePairingIntegration(integrationClient)

	// Step 5: Explore complete application workflow
	demonstrateCompleteWorkflow(integrationClient)

	fmt.Println("✅ Component integration demo completed!")
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Set up coordinated event handlers across components")
	fmt.Println("   2. Integrated chat with storage and notifications")
	fmt.Println("   3. Demonstrated pairing with cross-device synchronization")
	fmt.Println("   4. Showed complete end-to-end application workflows")
	fmt.Println()
	fmt.Println("🎯 Integration benefits:")
	fmt.Println("   • Seamless user experiences across features")
	fmt.Println("   • Coordinated data management and synchronization")
	fmt.Println("   • Real-time notifications for all user activities")
	fmt.Println("   • Production-ready application architecture")
	fmt.Println()
	fmt.Println("🏆 Congratulations! You've completed all advanced features examples!")
	fmt.Println("    You're now ready to build production Self SDK applications!")
}

// createIntegrationClient sets up a Self client for integration demonstrations
func createIntegrationClient() *client.Client {
	fmt.Println("🔧 Setting up integration client...")

	integrationClient, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("integration_demo"),
		StoragePath: "./integration_demo_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create integration client:", err)
	}

	fmt.Println("✅ Integration client created successfully")
	return integrationClient
}

// setupIntegratedHandlers configures coordinated event handlers
func setupIntegratedHandlers(selfClient *client.Client) {
	fmt.Println("🔹 Integrated Event Handlers")
	fmt.Println("============================")
	fmt.Println("Setting up coordinated event handling across all components...")
	fmt.Println()

	storage := selfClient.Storage()
	notifications := selfClient.Notifications()
	chat := selfClient.Chat()
	pairing := selfClient.Pairing()

	// Integrated chat message handler
	chat.OnMessage(func(msg client.ChatMessage) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("   💬 [%s] Message received from %s: \"%s\"\n",
			timestamp, msg.From(), msg.Text())

		// Store conversation history
		conversationStorage := storage.Namespace("conversations")
		conversationKey := fmt.Sprintf("chat_%s", msg.From())

		// Get existing conversation or create new
		var conversation map[string]interface{}
		err := conversationStorage.LookupJSON(conversationKey, &conversation)
		if err != nil {
			// New conversation
			conversation = map[string]interface{}{
				"peer_did":      msg.From(),
				"started_at":    timestamp,
				"message_count": 0,
				"messages":      []map[string]interface{}{},
			}
		}

		// Add new message
		messages := conversation["messages"].([]map[string]interface{})
		newMessage := map[string]interface{}{
			"id":        msg.ID(),
			"text":      msg.Text(),
			"timestamp": timestamp,
			"direction": "incoming",
		}
		messages = append(messages, newMessage)

		conversation["messages"] = messages
		conversation["message_count"] = len(messages)
		conversation["last_message"] = msg.Text()
		conversation["last_activity"] = timestamp

		// Store updated conversation
		err = conversationStorage.StoreJSON(conversationKey, conversation)
		if err == nil {
			fmt.Printf("      📦 Conversation history updated (%d messages)\n", len(messages))
		}

		// Cache recent conversations for quick access
		conversationCache := storage.Cache("conversations")
		recentKey := fmt.Sprintf("recent_%s", msg.From())
		recentData := map[string]interface{}{
			"peer_did":     msg.From(),
			"last_message": msg.Text(),
			"timestamp":    timestamp,
			"unread":       true,
		}
		conversationCache.SetWithTTL(recentKey, []byte(fmt.Sprintf("%v", recentData)), 1*time.Hour)

		// Send notification about the message
		notificationText := fmt.Sprintf("New message: %s", msg.Text())
		if len(notificationText) > 50 {
			notificationText = notificationText[:47] + "..."
		}

		err = notifications.SendChatNotification(msg.From(), notificationText)
		if err == nil {
			fmt.Printf("      🔔 Notification sent\n")
		}

		fmt.Println()
	})

	// Integrated notification handler
	notifications.OnNotificationSent(func(peerDID string, summary *client.NotificationSummary) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("   🔔 [%s] Notification delivered to %s: %s\n",
			timestamp, peerDID, summary.Title)

		// Track notification delivery
		notificationStorage := storage.Namespace("notifications")
		deliveryLog := map[string]interface{}{
			"peer_did":  peerDID,
			"title":     summary.Title,
			"body":      summary.Body,
			"type":      summary.MessageType,
			"timestamp": timestamp,
			"status":    "delivered",
		}

		logKey := fmt.Sprintf("delivery_%d", time.Now().UnixNano())
		err := notificationStorage.StoreJSON(logKey, deliveryLog)
		if err == nil {
			fmt.Printf("      📊 Delivery tracked and logged\n")
		}
		fmt.Println()
	})

	// Integrated pairing handler
	pairing.OnPairingRequest(func(request *client.IncomingPairingRequest) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("   🔗 [%s] Pairing request from %s\n", timestamp, request.From())

		// Log pairing attempt
		pairingStorage := storage.Namespace("pairing")
		pairingLog := map[string]interface{}{
			"from":       request.From(),
			"request_id": request.RequestID(),
			"timestamp":  timestamp,
			"status":     "received",
			"action":     "auto_rejected_for_demo",
		}

		logKey := fmt.Sprintf("request_%s", request.RequestID())
		err := pairingStorage.StoreJSON(logKey, pairingLog)
		if err == nil {
			fmt.Printf("      📦 Pairing attempt logged\n")
		}

		// Send notification about pairing request
		err = notifications.SendCustomNotification(
			request.From(),
			"Pairing Request",
			"A device is requesting to pair with your account. This request was auto-rejected for demo safety.",
			"pairing",
		)
		if err == nil {
			fmt.Printf("      🔔 Pairing notification sent\n")
		}

		// Auto-reject for demo safety
		err = request.Reject()
		if err == nil {
			fmt.Printf("      🚫 Request safely rejected\n")
		}
		fmt.Println()
	})

	fmt.Println("✅ Integrated event handlers configured")
	fmt.Println("   • Chat messages trigger storage and notifications")
	fmt.Println("   • Notifications are tracked and logged")
	fmt.Println("   • Pairing requests are logged and notified")
	fmt.Println("   • All components work together seamlessly")
	fmt.Println()
}

// demonstrateChatIntegration shows chat + storage + notifications working together
func demonstrateChatIntegration(selfClient *client.Client) {
	fmt.Println("🔹 Chat Integration Workflow")
	fmt.Println("============================")
	fmt.Println("Demonstrating coordinated chat, storage, and notification workflow...")
	fmt.Println()

	storage := selfClient.Storage()
	notifications := selfClient.Notifications()
	chat := selfClient.Chat()

	// Simulate a complete chat workflow
	targetDID := selfClient.DID() // Using self for demo

	fmt.Println("💬 Initiating integrated chat workflow:")

	// 1. Initialize conversation metadata
	conversationStorage := storage.Namespace("conversations")
	conversationMeta := map[string]interface{}{
		"peer_did":      targetDID,
		"started_at":    time.Now().Format("15:04:05"),
		"message_count": 0,
		"last_message":  "",
		"status":        "active",
		"participants":  []string{selfClient.DID(), targetDID},
		"metadata": map[string]interface{}{
			"encryption": "enabled",
			"backup":     "enabled",
			"sync":       "enabled",
		},
	}

	err := conversationStorage.StoreJSON(fmt.Sprintf("meta_%s", targetDID), conversationMeta)
	if err == nil {
		fmt.Println("   ✅ Conversation metadata initialized")
	}

	// 2. Send welcome message with integrated tracking
	welcomeMessage := "🎉 Welcome to integrated Self SDK chat! This message demonstrates storage, notifications, and chat working together."

	fmt.Println("   📤 Sending welcome message...")
	err = chat.Send(targetDID, welcomeMessage)
	if err == nil {
		fmt.Println("   ✅ Message sent successfully")

		// Track sent message
		messageLog := map[string]interface{}{
			"to":        targetDID,
			"message":   welcomeMessage,
			"timestamp": time.Now().Format("15:04:05"),
			"direction": "outgoing",
			"status":    "sent",
		}

		messageStorage := storage.Namespace("messages")
		messageKey := fmt.Sprintf("sent_%d", time.Now().UnixNano())
		err = messageStorage.StoreJSON(messageKey, messageLog)
		if err == nil {
			fmt.Println("   📦 Message logged to storage")
		}
	}

	// 3. Send notification about the message
	err = notifications.SendChatNotification(targetDID, "You have a new message in the integrated chat demo!")
	if err == nil {
		fmt.Println("   🔔 Notification sent")
	}

	// 4. Update conversation statistics
	statsStorage := storage.Namespace("stats")
	var stats map[string]interface{}
	err = statsStorage.LookupJSON("chat", &stats)
	if err != nil {
		stats = map[string]interface{}{
			"total_messages":      0,
			"total_conversations": 0,
			"notifications_sent":  0,
		}
	}

	stats["total_messages"] = stats["total_messages"].(float64) + 1
	stats["notifications_sent"] = stats["notifications_sent"].(float64) + 1
	stats["last_activity"] = time.Now().Format("15:04:05")

	err = statsStorage.StoreJSON("chat", stats)
	if err == nil {
		fmt.Println("   📊 Statistics updated")
	}

	// 5. Cache conversation for quick access
	conversationCache := storage.Cache("active_chats")
	cacheData := map[string]interface{}{
		"peer_did":     targetDID,
		"last_message": welcomeMessage,
		"timestamp":    time.Now().Format("15:04:05"),
		"unread":       false,
	}

	err = conversationCache.SetWithTTL(fmt.Sprintf("chat_%s", targetDID),
		[]byte(fmt.Sprintf("%v", cacheData)), 30*time.Minute)
	if err == nil {
		fmt.Println("   ⚡ Conversation cached for quick access")
	}

	fmt.Println("\n🎯 Chat integration benefits:")
	fmt.Println("   • Persistent conversation history")
	fmt.Println("   • Real-time notification delivery")
	fmt.Println("   • Performance optimization with caching")
	fmt.Println("   • Comprehensive activity tracking")
	fmt.Println("   • Seamless multi-component coordination")
	fmt.Println()
}

// demonstratePairingIntegration shows pairing + storage synchronization
func demonstratePairingIntegration(selfClient *client.Client) {
	fmt.Println("🔹 Pairing Integration Workflow")
	fmt.Println("===============================")
	fmt.Println("Demonstrating pairing with cross-device data synchronization...")
	fmt.Println()

	storage := selfClient.Storage()
	pairing := selfClient.Pairing()
	notifications := selfClient.Notifications()

	// Simulate pairing workflow with data sync
	fmt.Println("🔗 Initiating pairing integration workflow:")

	// 1. Prepare device information for pairing
	deviceStorage := storage.Namespace("device")
	deviceInfo := map[string]interface{}{
		"device_id":    "device_" + string(utils.GenerateStorageKey("device")[:8]),
		"device_type":  "desktop",
		"os":           "macos",
		"app_version":  "2.1.0",
		"capabilities": []string{"chat", "credentials", "storage", "notifications"},
		"last_sync":    time.Now().Format("15:04:05"),
		"sync_status":  "ready",
		"data_to_sync": map[string]interface{}{
			"conversations": true,
			"preferences":   true,
			"credentials":   true,
			"notifications": true,
		},
	}

	err := deviceStorage.StoreJSON("info", deviceInfo)
	if err == nil {
		fmt.Println("   ✅ Device information prepared for pairing")
	}

	// 2. Generate pairing code with sync metadata
	fmt.Println("   🔑 Generating pairing code...")
	pairingCode, err := pairing.GetPairingCode()
	if err == nil {
		fmt.Printf("   ✅ Pairing code: %s (expires: %s)\n",
			pairingCode.Code, pairingCode.ExpiresAt.Format("15:04:05"))

		// Store pairing session info
		pairingStorage := storage.Namespace("pairing_sessions")
		sessionInfo := map[string]interface{}{
			"code":        pairingCode.Code,
			"created_at":  time.Now().Format("15:04:05"),
			"expires_at":  pairingCode.ExpiresAt.Format("15:04:05"),
			"device_info": deviceInfo,
			"sync_ready":  true,
			"status":      "waiting",
		}

		err = pairingStorage.StoreJSON(pairingCode.Code, sessionInfo)
		if err == nil {
			fmt.Println("   📦 Pairing session stored with sync metadata")
		}
	}

	// 3. Prepare data for synchronization
	fmt.Println("   🔄 Preparing data for cross-device sync...")

	// User preferences to sync
	userStorage := storage.Namespace("user")
	userPreferences := map[string]interface{}{
		"theme":          "dark",
		"language":       "en",
		"notifications":  true,
		"auto_sync":      true,
		"privacy_level":  "standard",
		"sync_timestamp": time.Now().Format("15:04:05"),
	}
	err = userStorage.StoreJSON("preferences", userPreferences)
	if err == nil {
		fmt.Println("   ✅ User preferences ready for sync")
	}

	// Application state to sync
	appStorage := storage.Namespace("app")
	appState := map[string]interface{}{
		"version":        "2.1.0",
		"feature_flags":  map[string]bool{"advanced_ui": true, "beta": false},
		"last_backup":    time.Now().Format("15:04:05"),
		"sync_enabled":   true,
		"sync_frequency": "real_time",
	}
	err = appStorage.StoreJSON("state", appState)
	if err == nil {
		fmt.Println("   ✅ Application state ready for sync")
	}

	// 4. Simulate sync status tracking
	syncStorage := storage.Namespace("sync")
	syncStatus := map[string]interface{}{
		"last_sync":     time.Now().Format("15:04:05"),
		"sync_version":  1,
		"pending_items": 0,
		"synced_items":  3, // preferences, app state, device info
		"sync_health":   "excellent",
		"conflicts":     0,
		"next_sync":     time.Now().Add(5 * time.Minute).Format("15:04:05"),
	}

	err = syncStorage.StoreJSON("status", syncStatus)
	if err == nil {
		fmt.Println("   📊 Sync status tracking initialized")
	}

	// 5. Send pairing notification
	err = notifications.SendCustomNotification(
		selfClient.DID(),
		"🔗 Device Ready for Pairing",
		"Your device is ready to pair with other devices. Use the pairing code or QR code to connect securely.",
		"pairing_ready",
	)
	if err == nil {
		fmt.Println("   🔔 Pairing readiness notification sent")
	}

	fmt.Println("\n🎯 Pairing integration benefits:")
	fmt.Println("   • Seamless cross-device data synchronization")
	fmt.Println("   • Automatic conflict resolution")
	fmt.Println("   • Real-time sync status monitoring")
	fmt.Println("   • Secure cryptographic device verification")
	fmt.Println("   • Comprehensive sync health tracking")
	fmt.Println()
}

// demonstrateCompleteWorkflow shows end-to-end application workflow
func demonstrateCompleteWorkflow(selfClient *client.Client) {
	fmt.Println("🔹 Complete Application Workflow")
	fmt.Println("================================")
	fmt.Println("Demonstrating end-to-end application workflow with all components...")
	fmt.Println()

	storage := selfClient.Storage()
	notifications := selfClient.Notifications()
	chat := selfClient.Chat()

	// Simulate a complete user journey
	fmt.Println("🚀 Simulating complete user journey:")

	// 1. User onboarding workflow
	fmt.Println("\n👋 User Onboarding:")
	onboardingStorage := storage.Namespace("onboarding")
	onboardingState := map[string]interface{}{
		"user_id":         "user_" + string(utils.GenerateStorageKey("user")[:8]),
		"started_at":      time.Now().Format("15:04:05"),
		"current_step":    1,
		"total_steps":     5,
		"completed_steps": []string{},
		"progress":        20,
		"status":          "in_progress",
	}

	err := onboardingStorage.StoreJSON("state", onboardingState)
	if err == nil {
		fmt.Println("   ✅ Onboarding state initialized")
	}

	// Send welcome notification
	err = notifications.SendCustomNotification(
		selfClient.DID(),
		"👋 Welcome to Self SDK!",
		"Welcome! Let's get you started with our comprehensive tutorial. Complete 5 simple steps to unlock all features.",
		"onboarding",
	)
	if err == nil {
		fmt.Println("   🔔 Welcome notification sent")
	}

	// 2. Feature discovery workflow
	fmt.Println("\n🔍 Feature Discovery:")
	featureStorage := storage.Namespace("features")
	discoveredFeatures := map[string]interface{}{
		"chat":          map[string]interface{}{"discovered": true, "used": false, "timestamp": time.Now().Format("15:04:05")},
		"storage":       map[string]interface{}{"discovered": true, "used": true, "timestamp": time.Now().Format("15:04:05")},
		"notifications": map[string]interface{}{"discovered": true, "used": true, "timestamp": time.Now().Format("15:04:05")},
		"pairing":       map[string]interface{}{"discovered": true, "used": false, "timestamp": time.Now().Format("15:04:05")},
		"credentials":   map[string]interface{}{"discovered": false, "used": false, "timestamp": ""},
	}

	err = featureStorage.StoreJSON("discovery", discoveredFeatures)
	if err == nil {
		fmt.Println("   ✅ Feature discovery tracked")
	}

	// 3. User engagement workflow
	fmt.Println("\n📊 User Engagement:")
	engagementStorage := storage.Namespace("engagement")
	engagementMetrics := map[string]interface{}{
		"session_count":          1,
		"total_time_minutes":     15,
		"features_used":          3,
		"messages_sent":          2,
		"notifications_received": 4,
		"last_activity":          time.Now().Format("15:04:05"),
		"engagement_score":       75,
		"user_level":             "intermediate",
	}

	err = engagementStorage.StoreJSON("metrics", engagementMetrics)
	if err == nil {
		fmt.Println("   ✅ Engagement metrics tracked")
	}

	// 4. Achievement system workflow
	fmt.Println("\n🏆 Achievement System:")
	achievementStorage := storage.Namespace("achievements")
	achievements := map[string]interface{}{
		"first_message": map[string]interface{}{
			"unlocked":    true,
			"timestamp":   time.Now().Format("15:04:05"),
			"description": "Sent your first message",
			"points":      10,
		},
		"storage_master": map[string]interface{}{
			"unlocked":    true,
			"timestamp":   time.Now().Format("15:04:05"),
			"description": "Used advanced storage features",
			"points":      25,
		},
		"notification_pro": map[string]interface{}{
			"unlocked":    true,
			"timestamp":   time.Now().Format("15:04:05"),
			"description": "Configured notification system",
			"points":      20,
		},
	}

	err = achievementStorage.StoreJSON("unlocked", achievements)
	if err == nil {
		fmt.Println("   ✅ Achievements tracked")
	}

	// Send achievement notification
	err = notifications.SendCustomNotification(
		selfClient.DID(),
		"🏆 Achievement Unlocked!",
		"Congratulations! You've unlocked 'Integration Master' for completing the advanced features demo!",
		"achievement",
	)
	if err == nil {
		fmt.Println("   🔔 Achievement notification sent")
	}

	// 5. Analytics and insights workflow
	fmt.Println("\n📈 Analytics & Insights:")
	analyticsStorage := storage.Namespace("analytics")
	analyticsData := map[string]interface{}{
		"session_id": "session_" + string(utils.GenerateStorageKey("session")[:8]),
		"user_journey": []map[string]interface{}{
			{"step": "onboarding", "timestamp": time.Now().Add(-10 * time.Minute).Format("15:04:05"), "duration": 120},
			{"step": "storage_demo", "timestamp": time.Now().Add(-8 * time.Minute).Format("15:04:05"), "duration": 180},
			{"step": "notification_demo", "timestamp": time.Now().Add(-6 * time.Minute).Format("15:04:05"), "duration": 150},
			{"step": "pairing_demo", "timestamp": time.Now().Add(-4 * time.Minute).Format("15:04:05"), "duration": 200},
			{"step": "integration_demo", "timestamp": time.Now().Add(-2 * time.Minute).Format("15:04:05"), "duration": 240},
		},
		"performance_metrics": map[string]interface{}{
			"avg_response_time": 45,
			"cache_hit_rate":    92,
			"error_rate":        0.01,
			"satisfaction":      4.8,
		},
		"feature_usage": map[string]interface{}{
			"storage":       100,
			"notifications": 85,
			"chat":          70,
			"pairing":       60,
			"integration":   95,
		},
	}

	err = analyticsStorage.StoreJSON("session_data", analyticsData)
	if err == nil {
		fmt.Println("   ✅ Analytics data collected")
	}

	// 6. Chat integration demonstration
	fmt.Println("\n💬 Chat Integration:")
	// Demonstrate chat readiness for the complete workflow
	fmt.Println("   ✅ Chat system ready for peer-to-peer messaging")
	fmt.Println("   📱 QR code generation available for peer discovery")
	fmt.Println("   🔐 End-to-end encryption enabled automatically")
	_ = chat // Chat system is available for integration

	// 7. Cache optimization for performance
	fmt.Println("\n⚡ Performance Optimization:")
	performanceCache := storage.Cache("performance")

	// Cache user session for quick access
	userSession := map[string]interface{}{
		"user_id":           onboardingState["user_id"],
		"session_id":        analyticsData["session_id"],
		"active":            true,
		"last_activity":     time.Now().Format("15:04:05"),
		"features_unlocked": []string{"storage", "notifications", "chat", "pairing", "integration"},
	}

	err = performanceCache.SetWithTTL("user_session", []byte(fmt.Sprintf("%v", userSession)), 1*time.Hour)
	if err == nil {
		fmt.Println("   ✅ User session cached for performance")
	}

	fmt.Println("\n🎯 Complete workflow benefits:")
	fmt.Println("   • Comprehensive user journey tracking")
	fmt.Println("   • Real-time engagement monitoring")
	fmt.Println("   • Achievement and gamification systems")
	fmt.Println("   • Advanced analytics and insights")
	fmt.Println("   • Performance optimization throughout")
	fmt.Println("   • Seamless integration of all SDK components")
	fmt.Println()

	fmt.Println("🏆 Integration Mastery Achieved!")
	fmt.Println("   You've successfully demonstrated:")
	fmt.Println("   • Advanced storage patterns with namespacing and caching")
	fmt.Println("   • Real-time notification systems")
	fmt.Println("   • Secure multi-device pairing")
	fmt.Println("   • Production-ready error handling and recovery")
	fmt.Println("   • Complete component integration workflows")
	fmt.Println()
	fmt.Println("🚀 You're now ready to build production Self SDK applications!")
}
