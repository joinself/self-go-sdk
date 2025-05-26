package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
)

func main() {
	// Create three clients to simulate a group chat scenario
	adminClient, err := client.NewClient(client.Config{
		StorageKey:  make([]byte, 32), // In production, use a secure key
		StoragePath: "./admin_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create admin client:", err)
	}
	defer adminClient.Close()

	member1Client, err := client.NewClient(client.Config{
		StorageKey:  make([]byte, 32), // In production, use a secure key
		StoragePath: "./member1_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create member1 client:", err)
	}
	defer member1Client.Close()

	member2Client, err := client.NewClient(client.Config{
		StorageKey:  make([]byte, 32), // In production, use a secure key
		StoragePath: "./member2_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create member2 client:", err)
	}
	defer member2Client.Close()

	fmt.Printf("Admin DID: %s\n", adminClient.DID())
	fmt.Printf("Member1 DID: %s\n", member1Client.DID())
	fmt.Printf("Member2 DID: %s\n", member2Client.DID())

	// Set up group message handlers for all clients
	setupGroupMessageHandlers(adminClient, "Admin")
	setupGroupMessageHandlers(member1Client, "Member1")
	setupGroupMessageHandlers(member2Client, "Member2")

	// Set up group invitation handlers for members
	setupInvitationHandlers(member1Client, "Member1")
	setupInvitationHandlers(member2Client, "Member2")

	// Example 1: Create a group chat
	fmt.Println("\n📋 Creating a group chat...")
	group, err := adminClient.GroupChats().CreateGroup("Dev Team", "Daily standup and project discussions")
	if err != nil {
		log.Fatal("Failed to create group:", err)
	}

	fmt.Printf("✅ Created group: %s (ID: %s)\n", group.Name(), group.ID())
	fmt.Printf("   Description: %s\n", group.Description())
	fmt.Printf("   Admin: %s\n", group.Admin())
	fmt.Printf("   Members: %d\n", group.MemberCount())

	// Example 2: Discovery and connection establishment
	fmt.Println("\n📱 Setting up connections between clients...")

	// For demo purposes, we'll simulate connections by having each client
	// generate QR codes and "scan" each other's codes
	// In a real scenario, users would scan QR codes with their devices

	// Admin generates QR for Member1
	fmt.Println("Admin generating QR for Member1...")
	_, err = adminClient.Discovery().GenerateQRWithTimeout(30 * time.Second)
	if err != nil {
		log.Fatal("Failed to generate QR:", err)
	}

	// Member1 generates QR for Admin (simulating mutual discovery)
	_, err = member1Client.Discovery().GenerateQRWithTimeout(30 * time.Second)
	if err != nil {
		log.Fatal("Failed to generate QR:", err)
	}

	// Wait a moment for connections to establish
	fmt.Println("⏳ Waiting for connections to establish...")
	time.Sleep(3 * time.Second)

	// Example 3: Invite members to the group
	fmt.Println("\n👥 Inviting members to the group...")

	// Invite Member1
	err = adminClient.GroupChats().InviteToGroup(group.ID(), member1Client.DID(), "Welcome to our dev team group!")
	if err != nil {
		log.Printf("Failed to invite Member1: %v", err)
	} else {
		fmt.Printf("✅ Invited Member1 to group: %s\n", group.Name())
	}

	// Invite Member2
	err = adminClient.GroupChats().InviteToGroup(group.ID(), member2Client.DID(), "Join our daily discussions!")
	if err != nil {
		log.Printf("Failed to invite Member2: %v", err)
	} else {
		fmt.Printf("✅ Invited Member2 to group: %s\n", group.Name())
	}

	// Wait for invitations to be processed
	time.Sleep(2 * time.Second)

	// Example 4: Send group messages
	fmt.Println("\n💬 Sending group messages...")

	err = adminClient.GroupChats().SendToGroup(group.ID(), "Hello everyone! Welcome to our dev team group.")
	if err != nil {
		log.Printf("Failed to send admin message: %v", err)
	} else {
		fmt.Println("✅ Admin sent welcome message")
	}

	time.Sleep(1 * time.Second)

	err = adminClient.GroupChats().SendToGroup(group.ID(), "Let's use this for our daily standups and project updates.")
	if err != nil {
		log.Printf("Failed to send admin message: %v", err)
	} else {
		fmt.Println("✅ Admin sent instructions message")
	}

	// Example 5: Group management
	fmt.Println("\n⚙️ Demonstrating group management...")

	// Update group name
	err = group.UpdateName("Dev Team - Sprint 1")
	if err != nil {
		log.Printf("Failed to update group name: %v", err)
	} else {
		fmt.Println("✅ Updated group name")
	}

	time.Sleep(1 * time.Second)

	// Update group description
	err = group.UpdateDescription("Sprint 1 planning and daily standups")
	if err != nil {
		log.Printf("Failed to update group description: %v", err)
	} else {
		fmt.Println("✅ Updated group description")
	}

	// Example 6: List groups
	fmt.Println("\n📋 Listing groups...")
	adminGroups := adminClient.GroupChats().ListGroups()
	fmt.Printf("Admin has %d groups:\n", len(adminGroups))
	for i, g := range adminGroups {
		fmt.Printf("  %d. %s (ID: %s, Members: %d)\n", i+1, g.Name(), g.ID(), g.MemberCount())
	}

	// Example 7: Simulate some group activity
	fmt.Println("\n🎭 Simulating group activity...")

	// Send messages from different perspectives
	messages := []struct {
		client *client.Client
		name   string
		text   string
	}{
		{adminClient, "Admin", "Daily standup in 5 minutes!"},
		{adminClient, "Admin", "Please share your updates in the group"},
		{adminClient, "Admin", "Remember to update your task status"},
	}

	for _, msg := range messages {
		err := msg.client.GroupChats().SendToGroup(group.ID(), msg.text)
		if err != nil {
			log.Printf("Failed to send message from %s: %v", msg.name, err)
		} else {
			fmt.Printf("✅ %s: %s\n", msg.name, msg.text)
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n🎉 Group chat demo completed!")
	fmt.Println("Features demonstrated:")
	fmt.Println("  ✅ Group creation with admin privileges")
	fmt.Println("  ✅ Member invitation system")
	fmt.Println("  ✅ Group message broadcasting")
	fmt.Println("  ✅ Group management (name/description updates)")
	fmt.Println("  ✅ Role-based permissions")
	fmt.Println("  ✅ Event-driven message handling")
	fmt.Println("  ✅ Multi-client group coordination")

	// Keep running to see any delayed messages
	fmt.Println("\n⏳ Waiting for any delayed messages...")
	time.Sleep(5 * time.Second)
}

// Helper function to set up group message handlers
func setupGroupMessageHandlers(selfClient *client.Client, clientName string) {
	selfClient.GroupChats().OnGroupMessage(func(msg client.GroupChatMessage) {
		fmt.Printf("\n[%s] 📨 Group message in '%s':\n", clientName, msg.GroupName())
		fmt.Printf("   From: %s\n", msg.From())
		fmt.Printf("   Message: %s\n", msg.Text())
		fmt.Printf("   Time: %s\n", msg.Timestamp().Format("15:04:05"))
	})

	selfClient.GroupChats().OnGroupCreated(func(group *client.GroupChat) {
		fmt.Printf("\n[%s] 🎉 Group created: %s\n", clientName, group.Name())
	})

	selfClient.GroupChats().OnMemberJoined(func(groupID string, member *client.GroupMember) {
		fmt.Printf("\n[%s] 👋 Member joined group: %s (Role: %s)\n", clientName, member.DID, member.Role)
	})

	selfClient.GroupChats().OnMemberLeft(func(groupID string, memberDID string) {
		fmt.Printf("\n[%s] 👋 Member left group: %s\n", clientName, memberDID)
	})
}

// Helper function to set up invitation handlers
func setupInvitationHandlers(selfClient *client.Client, clientName string) {
	selfClient.GroupChats().OnGroupInvite(func(invitation *client.GroupChatInvitation) {
		fmt.Printf("\n[%s] 📧 Group invitation received:\n", clientName)
		fmt.Printf("   Group: %s\n", invitation.GroupName)
		fmt.Printf("   From: %s\n", invitation.InviterDID)
		fmt.Printf("   Message: %s\n", invitation.Message)
		fmt.Printf("   Expires: %s\n", invitation.ExpiresAt.Format("2006-01-02 15:04:05"))

		// For demo purposes, automatically accept invitations
		fmt.Printf("   🤖 Auto-accepting invitation...\n")
		err := invitation.Accept()
		if err != nil {
			fmt.Printf("   ❌ Failed to accept invitation: %v\n", err)
		} else {
			fmt.Printf("   ✅ Accepted invitation to join: %s\n", invitation.GroupName)
		}
	})
}
