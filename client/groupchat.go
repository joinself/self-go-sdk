package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
)

// GroupChatMessage represents a received group chat message
type GroupChatMessage struct {
	from        string
	text        string
	id          string
	refID       string
	groupID     string
	groupName   string
	timestamp   time.Time
	attachments []ChatAttachment
}

// GroupChat represents a group chat session
type GroupChat struct {
	id          string
	name        string
	description string
	members     map[string]*GroupMember
	admin       string
	created     time.Time
	client      *Client
	mu          sync.RWMutex
}

// GroupMember represents a member of a group chat
type GroupMember struct {
	DID      string
	Name     string
	Role     GroupRole
	JoinedAt time.Time
	LastSeen time.Time
	IsOnline bool
}

// GroupRole represents the role of a group member
type GroupRole string

const (
	GroupRoleAdmin     GroupRole = "admin"
	GroupRoleModerator GroupRole = "moderator"
	GroupRoleMember    GroupRole = "member"
)

// GroupChatInvitation represents an invitation to join a group chat
type GroupChatInvitation struct {
	GroupID     string
	GroupName   string
	InviterDID  string
	InviterName string
	Message     string
	ExpiresAt   time.Time
	client      *Client
}

// GroupChats handles group chat functionality
type GroupChats struct {
	client *Client

	// Active group chats
	groups map[string]*GroupChat
	mu     sync.RWMutex

	// Event handlers
	onGroupMessageHandlers []func(GroupChatMessage)
	onGroupInviteHandlers  []func(*GroupChatInvitation)
	onMemberJoinedHandlers []func(groupID string, member *GroupMember)
	onMemberLeftHandlers   []func(groupID string, memberDID string)
	onGroupCreatedHandlers []func(*GroupChat)
	onGroupUpdatedHandlers []func(*GroupChat)
	handlerMu              sync.RWMutex
}

// newGroupChats creates a new group chats component
func newGroupChats(client *Client) *GroupChats {
	return &GroupChats{
		client: client,
		groups: make(map[string]*GroupChat),
	}
}

// CreateGroup creates a new group chat
func (gc *GroupChats) CreateGroup(name, description string) (*GroupChat, error) {
	if gc.client.isClosed() {
		return nil, ErrClientClosed
	}

	groupID := generateGroupID()

	group := &GroupChat{
		id:          groupID,
		name:        name,
		description: description,
		members:     make(map[string]*GroupMember),
		admin:       gc.client.DID(),
		created:     time.Now(),
		client:      gc.client,
	}

	// Add creator as admin
	group.members[gc.client.DID()] = &GroupMember{
		DID:      gc.client.DID(),
		Name:     "Me", // Could be enhanced to get actual name
		Role:     GroupRoleAdmin,
		JoinedAt: time.Now(),
		LastSeen: time.Now(),
		IsOnline: true,
	}

	gc.mu.Lock()
	gc.groups[groupID] = group
	gc.mu.Unlock()

	// Notify handlers
	gc.handlerMu.RLock()
	handlers := make([]func(*GroupChat), len(gc.onGroupCreatedHandlers))
	copy(handlers, gc.onGroupCreatedHandlers)
	gc.handlerMu.RUnlock()

	for _, handler := range handlers {
		go handler(group)
	}

	return group, nil
}

// InviteToGroup invites a peer to join a group chat
func (gc *GroupChats) InviteToGroup(groupID, peerDID, inviteMessage string) error {
	if gc.client.isClosed() {
		return ErrClientClosed
	}

	gc.mu.RLock()
	group, exists := gc.groups[groupID]
	gc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("group not found: %s", groupID)
	}

	// Check if user is admin or moderator
	member, exists := group.members[gc.client.DID()]
	if !exists || (member.Role != GroupRoleAdmin && member.Role != GroupRoleModerator) {
		return fmt.Errorf("insufficient permissions to invite members")
	}

	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(peerDID)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	// Send invitation via chat message with special content
	// TODO: In a real implementation, we'd include structured metadata
	// For now, we'll use a simple text-based approach
	chatBuilder := message.NewChat().Message(fmt.Sprintf("Group Invitation: %s", inviteMessage))

	// TODO: Add custom metadata for group invitations
	// This would require extending the chat message format or using a custom message type

	content, err := chatBuilder.Finish()
	if err != nil {
		return err
	}

	return gc.client.sendMessage(peerAddress, *content)
}

// JoinGroup joins a group chat via invitation
func (gc *GroupChats) JoinGroup(invitation *GroupChatInvitation) error {
	if gc.client.isClosed() {
		return ErrClientClosed
	}

	if time.Now().After(invitation.ExpiresAt) {
		return fmt.Errorf("invitation has expired")
	}

	// Create or update group
	group := &GroupChat{
		id:          invitation.GroupID,
		name:        invitation.GroupName,
		description: "",
		members:     make(map[string]*GroupMember),
		admin:       invitation.InviterDID,
		created:     time.Now(),
		client:      gc.client,
	}

	// Add self as member
	group.members[gc.client.DID()] = &GroupMember{
		DID:      gc.client.DID(),
		Name:     "Me",
		Role:     GroupRoleMember,
		JoinedAt: time.Now(),
		LastSeen: time.Now(),
		IsOnline: true,
	}

	gc.mu.Lock()
	gc.groups[invitation.GroupID] = group
	gc.mu.Unlock()

	// Send join confirmation to group admin
	peerAddress := signing.FromAddress(invitation.InviterDID)
	if peerAddress != nil {
		joinMessage := fmt.Sprintf("I've joined the group: %s", invitation.GroupName)
		chatBuilder := message.NewChat().Message(joinMessage)

		if content, err := chatBuilder.Finish(); err == nil {
			gc.client.sendMessage(peerAddress, *content)
		}
	}

	// Notify handlers
	gc.handlerMu.RLock()
	handlers := make([]func(groupID string, member *GroupMember), len(gc.onMemberJoinedHandlers))
	copy(handlers, gc.onMemberJoinedHandlers)
	gc.handlerMu.RUnlock()

	for _, handler := range handlers {
		go handler(invitation.GroupID, group.members[gc.client.DID()])
	}

	return nil
}

// SendToGroup sends a message to all members of a group
func (gc *GroupChats) SendToGroup(groupID, messageText string) error {
	if gc.client.isClosed() {
		return ErrClientClosed
	}

	gc.mu.RLock()
	group, exists := gc.groups[groupID]
	gc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("group not found: %s", groupID)
	}

	// Check if user is a member
	if _, exists := group.members[gc.client.DID()]; !exists {
		return fmt.Errorf("not a member of group: %s", groupID)
	}

	// Build group message with metadata
	groupMessageText := fmt.Sprintf("[%s] %s", group.name, messageText)

	var errors []error
	successCount := 0

	// Send to all members except self
	for memberDID := range group.members {
		if memberDID == gc.client.DID() {
			continue // Don't send to self
		}

		peerAddress := signing.FromAddress(memberDID)
		if peerAddress == nil {
			errors = append(errors, fmt.Errorf("invalid DID for member: %s", memberDID))
			continue
		}

		chatBuilder := message.NewChat().Message(groupMessageText)

		// TODO: Add group metadata to message
		// This would require extending the message format

		content, err := chatBuilder.Finish()
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to build message for %s: %v", memberDID, err))
			continue
		}

		err = gc.client.sendMessage(peerAddress, *content)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %s: %v", memberDID, err))
		} else {
			successCount++
		}
	}

	if len(errors) > 0 && successCount == 0 {
		return fmt.Errorf("failed to send to any group members: %v", errors)
	}

	return nil
}

// ReplyToGroupMessage replies to a specific message in a group
func (gc *GroupChats) ReplyToGroupMessage(originalMessage GroupChatMessage, replyText string) error {
	if gc.client.isClosed() {
		return ErrClientClosed
	}

	gc.mu.RLock()
	group, exists := gc.groups[originalMessage.groupID]
	gc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("group not found: %s", originalMessage.groupID)
	}

	// Build reply message
	replyMessageText := fmt.Sprintf("[%s] Reply: %s", group.name, replyText)

	// Send reply to all group members
	return gc.SendToGroup(originalMessage.groupID, replyMessageText)
}

// GetGroup returns a group chat by ID
func (gc *GroupChats) GetGroup(groupID string) (*GroupChat, bool) {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	group, exists := gc.groups[groupID]
	return group, exists
}

// ListGroups returns all group chats
func (gc *GroupChats) ListGroups() []*GroupChat {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	groups := make([]*GroupChat, 0, len(gc.groups))
	for _, group := range gc.groups {
		groups = append(groups, group)
	}
	return groups
}

// LeaveGroup leaves a group chat
func (gc *GroupChats) LeaveGroup(groupID string) error {
	if gc.client.isClosed() {
		return ErrClientClosed
	}

	gc.mu.RLock()
	group, exists := gc.groups[groupID]
	gc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("group not found: %s", groupID)
	}

	// Notify other members
	leaveMessage := fmt.Sprintf("I'm leaving the group: %s", group.name)
	gc.SendToGroup(groupID, leaveMessage)

	// Remove from local groups
	gc.mu.Lock()
	delete(gc.groups, groupID)
	gc.mu.Unlock()

	// Notify handlers
	gc.handlerMu.RLock()
	handlers := make([]func(groupID string, memberDID string), len(gc.onMemberLeftHandlers))
	copy(handlers, gc.onMemberLeftHandlers)
	gc.handlerMu.RUnlock()

	for _, handler := range handlers {
		go handler(groupID, gc.client.DID())
	}

	return nil
}

// Event handler registration methods

// OnGroupMessage registers a handler for incoming group messages
func (gc *GroupChats) OnGroupMessage(handler func(GroupChatMessage)) {
	gc.handlerMu.Lock()
	defer gc.handlerMu.Unlock()
	gc.onGroupMessageHandlers = append(gc.onGroupMessageHandlers, handler)
}

// OnGroupInvite registers a handler for group invitations
func (gc *GroupChats) OnGroupInvite(handler func(*GroupChatInvitation)) {
	gc.handlerMu.Lock()
	defer gc.handlerMu.Unlock()
	gc.onGroupInviteHandlers = append(gc.onGroupInviteHandlers, handler)
}

// OnMemberJoined registers a handler for when a member joins a group
func (gc *GroupChats) OnMemberJoined(handler func(groupID string, member *GroupMember)) {
	gc.handlerMu.Lock()
	defer gc.handlerMu.Unlock()
	gc.onMemberJoinedHandlers = append(gc.onMemberJoinedHandlers, handler)
}

// OnMemberLeft registers a handler for when a member leaves a group
func (gc *GroupChats) OnMemberLeft(handler func(groupID string, memberDID string)) {
	gc.handlerMu.Lock()
	defer gc.handlerMu.Unlock()
	gc.onMemberLeftHandlers = append(gc.onMemberLeftHandlers, handler)
}

// OnGroupCreated registers a handler for when a group is created
func (gc *GroupChats) OnGroupCreated(handler func(*GroupChat)) {
	gc.handlerMu.Lock()
	defer gc.handlerMu.Unlock()
	gc.onGroupCreatedHandlers = append(gc.onGroupCreatedHandlers, handler)
}

// OnGroupUpdated registers a handler for when a group is updated
func (gc *GroupChats) OnGroupUpdated(handler func(*GroupChat)) {
	gc.handlerMu.Lock()
	defer gc.handlerMu.Unlock()
	gc.onGroupUpdatedHandlers = append(gc.onGroupUpdatedHandlers, handler)
}

// GroupChat methods

// ID returns the group ID
func (g *GroupChat) ID() string {
	return g.id
}

// Name returns the group name
func (g *GroupChat) Name() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.name
}

// Description returns the group description
func (g *GroupChat) Description() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.description
}

// Members returns all group members
func (g *GroupChat) Members() []*GroupMember {
	g.mu.RLock()
	defer g.mu.RUnlock()

	members := make([]*GroupMember, 0, len(g.members))
	for _, member := range g.members {
		members = append(members, member)
	}
	return members
}

// Admin returns the group admin DID
func (g *GroupChat) Admin() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.admin
}

// Created returns when the group was created
func (g *GroupChat) Created() time.Time {
	return g.created
}

// MemberCount returns the number of members
func (g *GroupChat) MemberCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.members)
}

// UpdateName updates the group name (admin/moderator only)
func (g *GroupChat) UpdateName(newName string) error {
	if g.client.isClosed() {
		return ErrClientClosed
	}

	member, exists := g.members[g.client.DID()]
	if !exists || (member.Role != GroupRoleAdmin && member.Role != GroupRoleModerator) {
		return fmt.Errorf("insufficient permissions to update group name")
	}

	g.mu.Lock()
	oldName := g.name
	g.name = newName
	g.mu.Unlock()

	// Notify group members
	updateMessage := fmt.Sprintf("Group name changed from '%s' to '%s'", oldName, newName)
	g.client.GroupChats().SendToGroup(g.id, updateMessage)

	return nil
}

// UpdateDescription updates the group description (admin/moderator only)
func (g *GroupChat) UpdateDescription(newDescription string) error {
	if g.client.isClosed() {
		return ErrClientClosed
	}

	member, exists := g.members[g.client.DID()]
	if !exists || (member.Role != GroupRoleAdmin && member.Role != GroupRoleModerator) {
		return fmt.Errorf("insufficient permissions to update group description")
	}

	g.mu.Lock()
	g.description = newDescription
	g.mu.Unlock()

	// Notify group members
	updateMessage := fmt.Sprintf("Group description updated: %s", newDescription)
	g.client.GroupChats().SendToGroup(g.id, updateMessage)

	return nil
}

// GroupChatMessage methods

// From returns the sender's DID
func (m GroupChatMessage) From() string {
	return m.from
}

// Text returns the message text
func (m GroupChatMessage) Text() string {
	return m.text
}

// ID returns the message ID
func (m GroupChatMessage) ID() string {
	return m.id
}

// ReferencedID returns the ID of the message this is replying to (if any)
func (m GroupChatMessage) ReferencedID() string {
	return m.refID
}

// GroupID returns the group ID
func (m GroupChatMessage) GroupID() string {
	return m.groupID
}

// GroupName returns the group name
func (m GroupChatMessage) GroupName() string {
	return m.groupName
}

// Timestamp returns when the message was sent
func (m GroupChatMessage) Timestamp() time.Time {
	return m.timestamp
}

// Attachments returns the message attachments
func (m GroupChatMessage) Attachments() []ChatAttachment {
	return m.attachments
}

// GroupChatInvitation methods

// Accept accepts the group invitation
func (inv *GroupChatInvitation) Accept() error {
	return inv.client.GroupChats().JoinGroup(inv)
}

// Decline declines the group invitation
func (inv *GroupChatInvitation) Decline() error {
	// Send decline message to inviter
	peerAddress := signing.FromAddress(inv.InviterDID)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	declineMessage := fmt.Sprintf("I declined the invitation to join: %s", inv.GroupName)
	chatBuilder := message.NewChat().Message(declineMessage)

	content, err := chatBuilder.Finish()
	if err != nil {
		return err
	}

	return inv.client.sendMessage(peerAddress, *content)
}

// Internal methods for handling events

func (gc *GroupChats) onConnect() {
	// Connection established - no specific action needed
}

func (gc *GroupChats) onDisconnect(err error) {
	// Connection lost - mark all members as offline
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	for _, group := range gc.groups {
		group.mu.Lock()
		for _, member := range group.members {
			if member.DID != gc.client.DID() {
				member.IsOnline = false
			}
		}
		group.mu.Unlock()
	}
}

func (gc *GroupChats) onWelcome(from *signing.PublicKey, groupAddress *signing.PublicKey) {
	// New connection established - could be a group member coming online
	fromDID := from.String()

	gc.mu.RLock()
	defer gc.mu.RUnlock()

	for _, group := range gc.groups {
		group.mu.Lock()
		if member, exists := group.members[fromDID]; exists {
			member.IsOnline = true
			member.LastSeen = time.Now()
		}
		group.mu.Unlock()
	}
}

func (gc *GroupChats) onKeyPackage(from *signing.PublicKey) {
	// Key package received - no specific action needed
}

func (gc *GroupChats) onIntroduction(from *signing.PublicKey, tokenCount int) {
	// Introduction received - connection is now ready for group chat
}

func (gc *GroupChats) onChatMessage(msg *event.Message) {
	// Decode the chat message
	chat, err := message.DecodeChat(msg.Content())
	if err != nil {
		return
	}

	messageText := chat.Message()
	fromDID := msg.FromAddress().String()

	// Check if this is a group message by looking for group prefix
	if len(messageText) > 3 && messageText[0] == '[' {
		// Extract group name from message format: [GroupName] Message
		endBracket := -1
		for i := 1; i < len(messageText); i++ {
			if messageText[i] == ']' {
				endBracket = i
				break
			}
		}

		if endBracket > 1 && endBracket < len(messageText)-2 {
			groupName := messageText[1:endBracket]
			actualMessage := messageText[endBracket+2:] // Skip "] "

			// Find the group by name (this is a simplified approach)
			gc.mu.RLock()
			var targetGroup *GroupChat
			for _, group := range gc.groups {
				if group.name == groupName {
					targetGroup = group
					break
				}
			}
			gc.mu.RUnlock()

			if targetGroup != nil {
				// Update member last seen
				targetGroup.mu.Lock()
				if member, exists := targetGroup.members[fromDID]; exists {
					member.LastSeen = time.Now()
					member.IsOnline = true
				}
				targetGroup.mu.Unlock()

				// Create group chat message
				groupMessage := GroupChatMessage{
					from:        fromDID,
					text:        actualMessage,
					id:          string(msg.ID()),
					refID:       string(chat.Referencing()),
					groupID:     targetGroup.id,
					groupName:   targetGroup.name,
					timestamp:   time.Now(),
					attachments: []ChatAttachment{}, // TODO: Handle attachments
				}

				// Notify handlers
				gc.handlerMu.RLock()
				handlers := make([]func(GroupChatMessage), len(gc.onGroupMessageHandlers))
				copy(handlers, gc.onGroupMessageHandlers)
				gc.handlerMu.RUnlock()

				for _, handler := range handlers {
					go handler(groupMessage)
				}
				return
			}
		}
	}

	// Check if this is a group invitation
	if len(messageText) > 17 && messageText[:17] == "Group Invitation:" {
		// This is a simplified invitation detection
		// In a real implementation, we'd use structured message metadata
		invitation := &GroupChatInvitation{
			GroupID:     generateGroupID(), // Simplified - should come from message metadata
			GroupName:   "Invited Group",   // Simplified - should come from message metadata
			InviterDID:  fromDID,
			InviterName: fromDID,          // Could be enhanced with actual names
			Message:     messageText[18:], // Skip "Group Invitation: "
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
			client:      gc.client,
		}

		// Notify handlers
		gc.handlerMu.RLock()
		handlers := make([]func(*GroupChatInvitation), len(gc.onGroupInviteHandlers))
		copy(handlers, gc.onGroupInviteHandlers)
		gc.handlerMu.RUnlock()

		for _, handler := range handlers {
			go handler(invitation)
		}
	}
}

func (gc *GroupChats) close() {
	// Clean up any resources if needed
	gc.mu.Lock()
	defer gc.mu.Unlock()

	// Clear all groups
	gc.groups = make(map[string]*GroupChat)
}

// Helper function to generate group IDs
func generateGroupID() string {
	return fmt.Sprintf("group_%d", time.Now().UnixNano())
}
