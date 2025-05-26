// Package main demonstrates push notification capabilities of the Self SDK.
//
// This is the NOTIFICATIONS level of advanced features examples.
// Prerequisites: Complete ../storage/main.go first to understand storage foundations.
//
// This example shows:
// - Push notification system for real-time user engagement
// - Multiple notification types (chat, credential, custom)
// - Event-driven notification handling
// - Delivery tracking and status management
// - Notification customization and targeting
//
// 🎯 What you'll learn:
// • How to send different types of notifications
// • Event handling for notification delivery
// • Notification customization and targeting
// • Real-time user engagement patterns
// • Integration with other SDK components
//
// 📚 Next steps:
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
	fmt.Println("🔔 Push Notifications Demo")
	fmt.Println("==========================")
	fmt.Println("This demo showcases Self SDK push notification capabilities.")
	fmt.Println("📚 This is the NOTIFICATIONS level - user engagement patterns.")
	fmt.Println()

	// Step 1: Create a Self client for notification demonstrations
	notificationClient := createNotificationClient()
	defer notificationClient.Close()

	fmt.Printf("🆔 Client DID: %s\n", notificationClient.DID())
	fmt.Println()

	// Step 2: Set up notification event handlers
	setupNotificationHandlers(notificationClient)

	// Step 3: Demonstrate different notification types
	demonstrateNotificationTypes(notificationClient)

	// Step 4: Show notification customization
	demonstrateNotificationCustomization(notificationClient)

	// Step 5: Explore notification targeting and delivery
	demonstrateNotificationTargeting(notificationClient)

	fmt.Println("✅ Notification system demo completed!")
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Set up notification event handlers")
	fmt.Println("   2. Sent various types of notifications")
	fmt.Println("   3. Customized notification content and appearance")
	fmt.Println("   4. Demonstrated targeting and delivery tracking")
	fmt.Println()
	fmt.Println("🎯 Notification benefits:")
	fmt.Println("   • Real-time user engagement and alerts")
	fmt.Println("   • Multiple notification types for different use cases")
	fmt.Println("   • Event-driven delivery tracking")
	fmt.Println("   • Seamless integration with other SDK features")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run ../pairing/main.go to learn about account pairing")
	fmt.Println("   • Run ../production_patterns/main.go for real-world patterns")
	fmt.Println("   • Run ../integration/main.go for component integration")
}

// createNotificationClient sets up a Self client for notification demonstrations
func createNotificationClient() *client.Client {
	fmt.Println("🔧 Setting up notification client...")

	notificationClient, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("notification_demo"),
		StoragePath: "./notification_demo_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create notification client:", err)
	}

	fmt.Println("✅ Notification client created successfully")
	return notificationClient
}

// setupNotificationHandlers configures notification event handlers
func setupNotificationHandlers(selfClient *client.Client) {
	fmt.Println("🔹 Setting Up Notification Handlers")
	fmt.Println("===================================")
	fmt.Println("Configuring event-driven notification handling...")
	fmt.Println()

	notifications := selfClient.Notifications()

	// Handler for successful notification delivery
	notifications.OnNotificationSent(func(peerDID string, summary *client.NotificationSummary) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("   📤 [%s] Notification sent to %s\n", timestamp, peerDID)
		fmt.Printf("      📋 Title: %s\n", summary.Title)
		fmt.Printf("      📝 Body: %s\n", summary.Body)
		fmt.Printf("      🏷️  Type: %s\n", summary.MessageType)
		fmt.Printf("      ✅ Delivery confirmed\n")
		fmt.Println()
	})

	fmt.Println("✅ Notification handlers configured")
	fmt.Println("   • Delivery confirmation tracking")
	fmt.Println("   • Event-driven notification handling")
	fmt.Println("   • Real-time notification status updates")
	fmt.Println()
}

// demonstrateNotificationTypes shows different notification patterns
func demonstrateNotificationTypes(selfClient *client.Client) {
	fmt.Println("🔹 Notification Types")
	fmt.Println("=====================")
	fmt.Println("Exploring different notification types and use cases...")
	fmt.Println()

	notifications := selfClient.Notifications()
	targetDID := selfClient.DID() // Using self for demo

	// Chat notification - for messaging alerts
	fmt.Println("💬 Chat Notification:")
	fmt.Println("   Use case: New message alerts, conversation updates")
	err := notifications.SendChatNotification(targetDID, "Hello! You have a new message from Alice.")
	if err != nil {
		log.Printf("Failed to send chat notification: %v", err)
	} else {
		fmt.Println("   ✅ Chat notification sent successfully")
	}
	time.Sleep(1 * time.Second) // Small delay for demo clarity

	// Credential notification - for identity-related alerts
	fmt.Println("\n🆔 Credential Notification:")
	fmt.Println("   Use case: Credential requests, verifications, issuance")
	err = notifications.SendCredentialNotification(targetDID, "identity", "request")
	if err != nil {
		log.Printf("Failed to send credential notification: %v", err)
	} else {
		fmt.Println("   ✅ Credential notification sent successfully")
	}
	time.Sleep(1 * time.Second)

	// Group invite notification - for group management
	fmt.Println("\n👥 Group Invite Notification:")
	fmt.Println("   Use case: Group invitations, membership changes")
	err = notifications.SendGroupInviteNotification(targetDID, "Development Team", "Alice")
	if err != nil {
		log.Printf("Failed to send group invite notification: %v", err)
	} else {
		fmt.Println("   ✅ Group invite notification sent successfully")
	}
	time.Sleep(1 * time.Second)

	// Custom notification - for application-specific alerts
	fmt.Println("\n🔔 Custom Notification:")
	fmt.Println("   Use case: Application-specific alerts, system notifications")
	err = notifications.SendCustomNotification(
		targetDID,
		"System Alert",
		"Your account security settings have been updated. Please review the changes in your security dashboard.",
		"security",
	)
	if err != nil {
		log.Printf("Failed to send custom notification: %v", err)
	} else {
		fmt.Println("   ✅ Custom notification sent successfully")
	}

	fmt.Println("\n📊 Notification Type Summary:")
	fmt.Println("   • Chat: Real-time messaging alerts")
	fmt.Println("   • Credential: Identity and verification alerts")
	fmt.Println("   • Group: Team and collaboration notifications")
	fmt.Println("   • Custom: Application-specific notifications")
	fmt.Println()
}

// demonstrateNotificationCustomization shows notification personalization
func demonstrateNotificationCustomization(selfClient *client.Client) {
	fmt.Println("🔹 Notification Customization")
	fmt.Println("=============================")
	fmt.Println("Personalizing notifications for better user experience...")
	fmt.Println()

	notifications := selfClient.Notifications()
	targetDID := selfClient.DID()

	// Urgent notification with high priority
	fmt.Println("🚨 High Priority Notification:")
	err := notifications.SendCustomNotification(
		targetDID,
		"🚨 URGENT: Security Alert",
		"Suspicious login attempt detected from new device. If this wasn't you, please secure your account immediately.",
		"security_urgent",
	)
	if err != nil {
		log.Printf("Failed to send urgent notification: %v", err)
	} else {
		fmt.Println("   ✅ Urgent security notification sent")
	}
	time.Sleep(1 * time.Second)

	// Informational notification with rich content
	fmt.Println("\n📊 Rich Content Notification:")
	err = notifications.SendCustomNotification(
		targetDID,
		"📊 Weekly Report Available",
		"Your weekly activity report is ready! This week: 47 messages sent, 3 new connections, 2 credentials verified. View full report in your dashboard.",
		"report",
	)
	if err != nil {
		log.Printf("Failed to send report notification: %v", err)
	} else {
		fmt.Println("   ✅ Rich content notification sent")
	}
	time.Sleep(1 * time.Second)

	// Achievement notification with celebration
	fmt.Println("\n🏆 Achievement Notification:")
	err = notifications.SendCustomNotification(
		targetDID,
		"🏆 Achievement Unlocked!",
		"Congratulations! You've successfully completed your first credential exchange. You're now ready to explore advanced Self SDK features.",
		"achievement",
	)
	if err != nil {
		log.Printf("Failed to send achievement notification: %v", err)
	} else {
		fmt.Println("   ✅ Achievement notification sent")
	}
	time.Sleep(1 * time.Second)

	// Reminder notification with action items
	fmt.Println("\n⏰ Reminder Notification:")
	err = notifications.SendCustomNotification(
		targetDID,
		"⏰ Reminder: Credential Expiring Soon",
		"Your professional certification credential expires in 7 days. Renew now to maintain continuous verification status.",
		"reminder",
	)
	if err != nil {
		log.Printf("Failed to send reminder notification: %v", err)
	} else {
		fmt.Println("   ✅ Reminder notification sent")
	}

	fmt.Println("\n🎨 Customization Features:")
	fmt.Println("   • Emoji and visual indicators for quick recognition")
	fmt.Println("   • Priority levels for urgent vs. informational content")
	fmt.Println("   • Rich content with detailed information")
	fmt.Println("   • Category-based organization and filtering")
	fmt.Println("   • Action-oriented messaging for user engagement")
	fmt.Println()
}

// demonstrateNotificationTargeting shows targeting and delivery patterns
func demonstrateNotificationTargeting(selfClient *client.Client) {
	fmt.Println("🔹 Notification Targeting & Delivery")
	fmt.Println("====================================")
	fmt.Println("Advanced targeting and delivery tracking patterns...")
	fmt.Println()

	notifications := selfClient.Notifications()
	targetDID := selfClient.DID()

	// Simulate different user scenarios
	fmt.Println("🎯 Targeted Notification Scenarios:")

	// New user onboarding
	fmt.Println("\n👋 New User Onboarding:")
	err := notifications.SendCustomNotification(
		targetDID,
		"👋 Welcome to Self SDK!",
		"Welcome! Let's get you started with a quick tour. Complete these steps: 1) Verify your identity, 2) Connect with peers, 3) Exchange your first credential.",
		"onboarding",
	)
	if err != nil {
		log.Printf("Failed to send onboarding notification: %v", err)
	} else {
		fmt.Println("   ✅ Onboarding notification sent")
	}
	time.Sleep(1 * time.Second)

	// Re-engagement for inactive users
	fmt.Println("\n🔄 Re-engagement Notification:")
	err = notifications.SendCustomNotification(
		targetDID,
		"🔄 We Miss You!",
		"It's been a while since your last visit. Check out the new features we've added: advanced storage, group chat, and credential templates.",
		"re_engagement",
	)
	if err != nil {
		log.Printf("Failed to send re-engagement notification: %v", err)
	} else {
		fmt.Println("   ✅ Re-engagement notification sent")
	}
	time.Sleep(1 * time.Second)

	// Feature announcement for existing users
	fmt.Println("\n🆕 Feature Announcement:")
	err = notifications.SendCustomNotification(
		targetDID,
		"🆕 New Feature: Advanced Storage",
		"Exciting news! We've launched advanced storage with namespacing and TTL. Organize your data better and improve performance with intelligent caching.",
		"feature_announcement",
	)
	if err != nil {
		log.Printf("Failed to send feature announcement: %v", err)
	} else {
		fmt.Println("   ✅ Feature announcement sent")
	}
	time.Sleep(1 * time.Second)

	// Maintenance notification
	fmt.Println("\n🔧 Maintenance Notification:")
	err = notifications.SendCustomNotification(
		targetDID,
		"🔧 Scheduled Maintenance",
		"Scheduled maintenance tonight 2-4 AM UTC. Services may be briefly unavailable. We're upgrading our infrastructure for better performance.",
		"maintenance",
	)
	if err != nil {
		log.Printf("Failed to send maintenance notification: %v", err)
	} else {
		fmt.Println("   ✅ Maintenance notification sent")
	}

	// Demonstrate notification analytics
	fmt.Println("\n📈 Notification Analytics & Insights:")
	fmt.Println("   • Delivery confirmation rates")
	fmt.Println("   • User engagement metrics")
	fmt.Println("   • Optimal timing for different notification types")
	fmt.Println("   • Category-based performance analysis")
	fmt.Println()

	fmt.Println("🎯 Targeting Best Practices:")
	fmt.Println("   • Segment users based on activity and preferences")
	fmt.Println("   • Time notifications for optimal engagement")
	fmt.Println("   • Personalize content based on user journey stage")
	fmt.Println("   • Use categories for user filtering and preferences")
	fmt.Println("   • Track delivery and engagement for optimization")
	fmt.Println()

	fmt.Println("🔄 Integration Opportunities:")
	fmt.Println("   • Combine with storage for user preference tracking")
	fmt.Println("   • Integrate with chat for conversation notifications")
	fmt.Println("   • Use with pairing for multi-device notification sync")
	fmt.Println("   • Connect to credentials for verification alerts")
	fmt.Println()
}
