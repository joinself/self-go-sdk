# Group Chat Example

A comprehensive demonstration of the Self SDK's group chat capabilities, designed for educational purposes.

## 🚀 Quick Start

```bash
# Run the group chat demo
go run main.go
```

The demo creates 3 clients (admin + 2 members) and demonstrates the complete group chat workflow automatically!

## 📊 Complexity Rating

**5/10** (Intermediate) - Perfect for learning group communication patterns

- 🟢 **Clear structure**: Step-by-step progression with multiple clients
- 🟡 **Group concepts**: Admin roles, invitations, member management
- 🟡 **Multi-client coordination**: Handling multiple participants
- 🟠 **Event handling**: Complex group event management

## 🎯 What This Example Demonstrates

### Core Group Chat Features
- ✅ **Group creation** - Admin-controlled group establishment
- ✅ **Member management** - Invitation and acceptance workflow
- ✅ **Group messaging** - Broadcasting messages to all members
- ✅ **Role-based permissions** - Admin vs member privileges
- ✅ **Real-time events** - Live group activity notifications
- ✅ **Group administration** - Name/description management

### Educational Learning Path
1. **Multi-Client Setup** - Create admin and member clients
2. **Event Handler Configuration** - Set up group event processing
3. **Group Creation** - Establish group with admin privileges
4. **Peer Connections** - Connect all participants
5. **Member Invitations** - Invite and accept group membership
6. **Group Messaging** - Broadcast messages to all members
7. **Group Management** - Demonstrate administrative features

## 🏃‍♂️ How to Run

### Single Command Demo
```bash
go run main.go
```

The demo runs automatically and shows:
- Creation of 3 Self clients (1 admin + 2 members)
- Group creation with admin privileges
- Member invitation and acceptance process
- Group message broadcasting
- Administrative group management

### What Happens Automatically
1. **Client Creation**: Admin, Member1, and Member2 clients are created
2. **Handler Setup**: Group event handlers are configured for all clients
3. **Group Creation**: Admin creates "Dev Team" group
4. **Connections**: Peer connections are established (simulated)
5. **Invitations**: Members are invited and auto-accept
6. **Messaging**: Demo messages are sent to the group
7. **Management**: Group name and description are updated

## 📋 What You'll See

```
👥 Group Chat Demo
==================
This demo shows group chat functionality with multiple participants.

🔧 Setting up group chat clients...
✅ All clients created successfully
👑 Admin: did:self:admin123...
👤 Member1: did:self:member1456...
👤 Member2: did:self:member2789...

📨 Setting up group event handlers...
✅ Group handlers configured for all clients

📋 Creating a group chat...
✅ Group created successfully:
   Name: Dev Team
   ID: group_abc123...
   Description: Daily standup and project discussions
   Admin: did:self:admin123...
   Members: 1

🔗 Establishing peer connections...
   (Simulating QR code discovery for demo purposes)
✅ Peer connections established
   • Admin ↔ Member1
   • Admin ↔ Member2
   • Member1 ↔ Member2

👥 Inviting members to the group...
📤 Inviting Member1...
✅ Invitation sent to Member1

📧 [👤 Member1] Group invitation received:
   Group: Dev Team
   From: did:self:admin123...
   Message: Welcome to our dev team group!
   🤖 Auto-accepting invitation...
   ✅ Joined group: Dev Team

💬 Demonstrating group messaging...
📤 Admin sending: "🎉 Hello everyone! Welcome to our dev team group."
✅ Welcome message sent to group

📨 [👤 Member1] Group message in 'Dev Team' at 15:04:05:
   From: did:self:admin123...
   💬 "🎉 Hello everyone! Welcome to our dev team group."

⚙️ Demonstrating group management...
📝 Updating group name to: "Dev Team - Sprint 1"
✅ Group name updated successfully

✅ Group chat demo completed!

🎓 What happened:
   1. Created multiple Self clients (admin + members)
   2. Set up handlers for group events and messages
   3. Created a group chat with admin privileges
   4. Established peer connections between clients
   5. Invited members and handled invitations
   6. Exchanged messages in the group chat
   7. Demonstrated group management features
```

## 🔍 Key Code Sections

| Function | Lines | Purpose |
|----------|-------|---------|
| `main()` | 30-80 | Step-by-step demo orchestration |
| `createClients()` | 85-120 | Multi-client setup (admin + members) |
| `setupGroupHandlers()` | 125-140 | Event handler configuration |
| `setupClientHandlers()` | 145-180 | Individual client event handling |
| `createGroup()` | 185-205 | Group creation with admin privileges |
| `establishConnections()` | 210-225 | Peer discovery simulation |
| `inviteMembers()` | 230-260 | Member invitation workflow |
| `demonstrateGroupMessaging()` | 265-310 | Group message broadcasting |
| `demonstrateGroupManagement()` | 315-350 | Administrative features |

## 🎓 Educational Notes

### Core Concepts
- **Group Admin**: The creator who has management privileges
- **Group Members**: Participants who can send/receive messages
- **Group Invitations**: Secure invitation and acceptance workflow
- **Message Broadcasting**: Messages sent to all group members
- **Event-Driven Architecture**: Real-time group activity notifications

### Group Lifecycle
1. **Creation**: Admin creates group with name and description
2. **Invitation**: Admin invites members with custom messages
3. **Acceptance**: Members receive and accept invitations
4. **Messaging**: All members can send messages to the group
5. **Management**: Admin can update group properties
6. **Events**: All participants receive real-time notifications

### Role-Based Permissions
- **Admin Privileges**:
  - Create and manage groups
  - Invite and remove members
  - Update group name and description
  - Send messages to the group
- **Member Privileges**:
  - Send messages to the group
  - Receive group messages and events
  - Leave the group

### Event Types Handled
- **Group Messages**: Real-time message broadcasting
- **Member Joined**: Notification when someone joins
- **Group Created**: Notification of new group creation
- **Invitations**: Secure invitation delivery and acceptance

## 🔧 Customization Ideas

Try modifying the code to:
- Add more group members (Member3, Member4, etc.)
- Implement member removal functionality
- Add message reactions or threading
- Create multiple groups with different purposes
- Implement group member role changes
- Add group message history and persistence
- Create private vs public group types

## 🚀 Next Steps

After understanding this example, explore:

| Example | Complexity | Description |
|---------|------------|-------------|
| `../simple_chat/` | 4/10 | Basic peer-to-peer messaging |
| `../credentials_exchange/` | 6/10 | Identity verification and credential sharing |
| `../file_sharing/` | 7/10 | Secure file transfer in groups |

## 🛠️ Prerequisites

- Go 1.19 or later
- Self SDK dependencies (handled by go.mod)
- Understanding of basic Self SDK concepts (recommended: try simple_chat first)

## 💡 Troubleshooting

**Group Creation Issues:**
- Ensure admin client is properly initialized
- Check that storage paths are writable
- Verify network connectivity for Self SDK

**Invitation Problems:**
- Confirm peer connections are established
- Check that invitation handlers are set up before sending
- Verify member DIDs are correct

**Message Delivery Issues:**
- Ensure all clients have message handlers configured
- Check that group members have accepted invitations
- Verify clients remain running to receive messages

**Build Issues:**
- Run `go mod tidy` to ensure dependencies
- Check Go version with `go version`
- Verify you're in the correct directory

## 🎯 Key Differences from Simple Chat

| Feature | Simple Chat | Group Chat |
|---------|-------------|------------|
| **Participants** | 2 peers | Multiple members |
| **Discovery** | QR code scanning | Admin invitations |
| **Messaging** | Direct peer-to-peer | Broadcast to group |
| **Roles** | Equal peers | Admin vs members |
| **Management** | None | Group administration |
| **Complexity** | 4/10 | 5/10 |

This example builds upon the simple_chat concepts and adds the complexity of multi-participant coordination and group management! 
