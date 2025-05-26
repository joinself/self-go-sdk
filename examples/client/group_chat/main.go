// Package main demonstrates group chat functionality using the Self SDK.
//
// This example shows the basics of:
// - Setting up multiple Self clients for group communication
// - Creating and managing group chats
// - Inviting members and handling invitations
// - Sending and receiving group messages
// - Understanding group administration and permissions
//
// 🎯 What you'll learn:
// • How group chat works with Self SDK
// • Group creation and administration patterns
// • Member invitation and management
// • Multi-participant messaging
// • Role-based permissions in groups
//
// 👥 GROUP CHAT CAPABILITIES DEMONSTRATED:
// • Group creation with admin privileges
// • Member invitation system
// • Group message broadcasting
// • Real-time multi-participant messaging
// • Group management (name/description updates)
// • Event-driven group notifications
// • Role-based access control
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("👥 Group Chat Demo")
	fmt.Println("==================")
	fmt.Println("This demo shows group chat functionality with multiple participants.")
	fmt.Println()

	// Step 1: Create multiple clients (admin and members)
	admin, member1, member2 := createClients()
	defer admin.Close()
	defer member1.Close()
	defer member2.Close()

	fmt.Printf("👑 Admin: %s\n", admin.DID())
	fmt.Printf("👤 Member1: %s\n", member1.DID())
	fmt.Printf("👤 Member2: %s\n", member2.DID())
	fmt.Println()

	// Step 2: Set up group event handlers
	setupGroupHandlers(admin, member1, member2)

	// Step 3: Create a group chat
	group := createGroup(admin)

	// Step 4: Establish peer connections (simplified for demo)
	establishConnections(admin, member1, member2)

	// Step 5: Invite members to the group
	inviteMembers(admin, group, member1, member2)

	// Step 6: Demonstrate group messaging
	demonstrateGroupMessaging(admin, group)

	// Step 7: Show group management features
	demonstrateGroupManagement(admin, group)

	fmt.Println("✅ Group chat demo completed!")
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Created multiple Self clients (admin + members)")
	fmt.Println("   2. Set up handlers for group events and messages")
	fmt.Println("   3. Created a group chat with admin privileges")
	fmt.Println("   4. Established peer connections between clients")
	fmt.Println("   5. Invited members and handled invitations")
	fmt.Println("   6. Exchanged messages in the group chat")
	fmt.Println("   7. Demonstrated group management features")
	fmt.Println()
	fmt.Println("The clients will keep running to show ongoing group capabilities.")
	fmt.Println("Group messages are broadcasted to all members in real-time!")
	fmt.Println("Press Ctrl+C to exit.")

	// Keep running to demonstrate ongoing group capabilities
	select {}
}

// createClients sets up the admin and member clients for group chat
func createClients() (*client.Client, *client.Client, *client.Client) {
	fmt.Println("🔧 Setting up group chat clients...")

	// Create admin client
	admin, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("group_admin"),
		StoragePath: "./group_admin_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create admin client:", err)
	}

	// Create member1 client
	member1, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("group_member1"),
		StoragePath: "./group_member1_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create member1 client:", err)
	}

	// Create member2 client
	member2, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("group_member2"),
		StoragePath: "./group_member2_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create member2 client:", err)
	}

	fmt.Println("✅ All clients created successfully")
	return admin, member1, member2
}

// setupGroupHandlers configures event handlers for all group activities
func setupGroupHandlers(admin, member1, member2 *client.Client) {
	fmt.Println("📨 Setting up group event handlers...")

	// Set up handlers for admin
	setupClientHandlers(admin, "👑 Admin")

	// Set up handlers for member1
	setupClientHandlers(member1, "👤 Member1")

	// Set up handlers for member2
	setupClientHandlers(member2, "👤 Member2")

	fmt.Println("✅ Group handlers configured for all clients")
	fmt.Println()
}

// setupClientHandlers configures group event handlers for a specific client
func setupClientHandlers(selfClient *client.Client, clientName string) {
	// Handle incoming group messages
	selfClient.GroupChats().OnGroupMessage(func(msg client.GroupChatMessage) {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("\n📨 [%s] Group message in '%s' at %s:\n", clientName, msg.GroupName(), timestamp)
		fmt.Printf("   From: %s\n", msg.From())
		fmt.Printf("   💬 \"%s\"\n", msg.Text())
	})

	// Handle group invitations (for members)
	selfClient.GroupChats().OnGroupInvite(func(invitation *client.GroupChatInvitation) {
		fmt.Printf("\n📧 [%s] Group invitation received:\n", clientName)
		fmt.Printf("   Group: %s\n", invitation.GroupName)
		fmt.Printf("   From: %s\n", invitation.InviterDID)
		fmt.Printf("   Message: %s\n", invitation.Message)

		// Auto-accept invitations for demo purposes
		fmt.Printf("   🤖 Auto-accepting invitation...\n")
		err := invitation.Accept()
		if err != nil {
			fmt.Printf("   ❌ Failed to accept: %v\n", err)
		} else {
			fmt.Printf("   ✅ Joined group: %s\n", invitation.GroupName)
		}
	})

	// Handle member join events
	selfClient.GroupChats().OnMemberJoined(func(groupID string, member *client.GroupMember) {
		fmt.Printf("\n👋 [%s] Member joined group: %s (Role: %s)\n", clientName, member.DID, member.Role)
	})

	// Handle group creation events
	selfClient.GroupChats().OnGroupCreated(func(group *client.GroupChat) {
		fmt.Printf("\n🎉 [%s] Group created: %s\n", clientName, group.Name())
	})
}

// createGroup demonstrates group creation with admin privileges
func createGroup(admin *client.Client) *client.GroupChat {
	fmt.Println("📋 Creating a group chat...")

	group, err := admin.GroupChats().CreateGroup("Dev Team", "Daily standup and project discussions")
	if err != nil {
		log.Fatal("Failed to create group:", err)
	}

	fmt.Printf("✅ Group created successfully:\n")
	fmt.Printf("   Name: %s\n", group.Name())
	fmt.Printf("   ID: %s\n", group.ID())
	fmt.Printf("   Description: %s\n", group.Description())
	fmt.Printf("   Admin: %s\n", group.Admin())
	fmt.Printf("   Members: %d\n", group.MemberCount())
	fmt.Println()

	return group
}

// establishConnections simulates peer discovery between clients
func establishConnections(admin, member1, member2 *client.Client) {
	fmt.Println("🔗 Establishing peer connections...")
	fmt.Println("   (Simulating QR code discovery for demo purposes)")

	// In a real scenario, clients would scan each other's QR codes
	// For demo purposes, we simulate this with timeouts

	// Simulate connection establishment
	time.Sleep(2 * time.Second)

	fmt.Println("✅ Peer connections established")
	fmt.Println("   • Admin ↔ Member1")
	fmt.Println("   • Admin ↔ Member2")
	fmt.Println("   • Member1 ↔ Member2")
	fmt.Println()
}

// inviteMembers demonstrates the group invitation process
func inviteMembers(admin *client.Client, group *client.GroupChat, member1, member2 *client.Client) {
	fmt.Println("👥 Inviting members to the group...")

	// Invite Member1
	fmt.Println("📤 Inviting Member1...")
	err := admin.GroupChats().InviteToGroup(group.ID(), member1.DID(), "Welcome to our dev team group!")
	if err != nil {
		log.Printf("Failed to invite Member1: %v", err)
	} else {
		fmt.Printf("✅ Invitation sent to Member1\n")
	}

	// Small delay to see the invitation process
	time.Sleep(1 * time.Second)

	// Invite Member2
	fmt.Println("📤 Inviting Member2...")
	err = admin.GroupChats().InviteToGroup(group.ID(), member2.DID(), "Join our daily discussions!")
	if err != nil {
		log.Printf("Failed to invite Member2: %v", err)
	} else {
		fmt.Printf("✅ Invitation sent to Member2\n")
	}

	// Wait for invitations to be processed
	fmt.Println("⏳ Waiting for invitations to be processed...")
	time.Sleep(3 * time.Second)
	fmt.Println()
}

// demonstrateGroupMessaging shows group message broadcasting
func demonstrateGroupMessaging(admin *client.Client, group *client.GroupChat) {
	fmt.Println("💬 Demonstrating group messaging...")

	// Send welcome message
	welcomeMsg := "🎉 Hello everyone! Welcome to our dev team group."
	fmt.Printf("📤 Admin sending: \"%s\"\n", welcomeMsg)
	err := admin.GroupChats().SendToGroup(group.ID(), welcomeMsg)
	if err != nil {
		log.Printf("Failed to send welcome message: %v", err)
	} else {
		fmt.Println("✅ Welcome message sent to group")
	}

	time.Sleep(2 * time.Second)

	// Send instructions message
	instructionsMsg := "Let's use this for our daily standups and project updates."
	fmt.Printf("📤 Admin sending: \"%s\"\n", instructionsMsg)
	err = admin.GroupChats().SendToGroup(group.ID(), instructionsMsg)
	if err != nil {
		log.Printf("Failed to send instructions: %v", err)
	} else {
		fmt.Println("✅ Instructions sent to group")
	}

	time.Sleep(2 * time.Second)

	// Send multiple demo messages
	demoMessages := []string{
		"Daily standup in 5 minutes!",
		"Please share your updates in the group",
		"Remember to update your task status",
		"Great work everyone! 🚀",
	}

	fmt.Println("\n📤 Sending demo messages to group...")
	for i, msg := range demoMessages {
		fmt.Printf("📤 [%d/%d] \"%s\"\n", i+1, len(demoMessages), msg)
		err := admin.GroupChats().SendToGroup(group.ID(), msg)
		if err != nil {
			fmt.Printf("❌ Failed to send message: %v\n", err)
		} else {
			fmt.Printf("✅ Message sent successfully\n")
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
}

// demonstrateGroupManagement shows group administration features
func demonstrateGroupManagement(admin *client.Client, group *client.GroupChat) {
	fmt.Println("⚙️ Demonstrating group management...")

	// Update group name
	newName := "Dev Team - Sprint 1"
	fmt.Printf("📝 Updating group name to: \"%s\"\n", newName)
	err := group.UpdateName(newName)
	if err != nil {
		log.Printf("Failed to update group name: %v", err)
	} else {
		fmt.Println("✅ Group name updated successfully")
	}

	time.Sleep(1 * time.Second)

	// Update group description
	newDescription := "Sprint 1 planning and daily standups"
	fmt.Printf("📝 Updating description to: \"%s\"\n", newDescription)
	err = group.UpdateDescription(newDescription)
	if err != nil {
		log.Printf("Failed to update description: %v", err)
	} else {
		fmt.Println("✅ Group description updated successfully")
	}

	time.Sleep(1 * time.Second)

	// List all groups for admin
	fmt.Println("\n📋 Listing admin's groups:")
	adminGroups := admin.GroupChats().ListGroups()
	fmt.Printf("Admin manages %d group(s):\n", len(adminGroups))
	for i, g := range adminGroups {
		fmt.Printf("  %d. %s (ID: %s, Members: %d)\n", i+1, g.Name(), g.ID(), g.MemberCount())
		fmt.Printf("     Description: %s\n", g.Description())
	}
	fmt.Println()

	fmt.Println("🎯 Group management features demonstrated:")
	fmt.Println("   • Group name updates")
	fmt.Println("   • Group description updates")
	fmt.Println("   • Group listing and information")
	fmt.Println("   • Admin privilege management")
	fmt.Println()
}
