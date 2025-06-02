package client

import (
	"sync"

	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
)

// NotificationSummary represents a notification summary
type NotificationSummary struct {
	Title       string
	Body        string
	MessageType string
	FromDID     string
	MessageID   string
}

// Notifications handles push notification functionality
type Notifications struct {
	client *Client

	// Event handlers
	onNotificationSentHandlers []func(peerDID string, summary *NotificationSummary)
	mu                         sync.RWMutex
}

// newNotifications creates a new notifications component
func newNotifications(client *Client) *Notifications {
	return &Notifications{
		client: client,
	}
}

// SendNotification sends a push notification to a peer
func (n *Notifications) SendNotification(peerDID string, summary *NotificationSummary) error {
	if n.client.isClosed() {
		return ErrClientClosed
	}

	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(peerDID)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	// Create a simple chat message to generate a content summary
	// This is a workaround since ContentSummary doesn't have a direct constructor
	chatBuilder := message.NewChat().Message(summary.Body)
	content, err := chatBuilder.Finish()
	if err != nil {
		return err
	}

	// Generate the content summary from the message
	contentSummary, err := content.Summary()
	if err != nil {
		return err
	}

	// Send the notification
	err = n.client.account.NotificationSend(peerAddress, contentSummary)
	if err != nil {
		return err
	}

	// Notify handlers
	n.mu.RLock()
	handlers := make([]func(string, *NotificationSummary), len(n.onNotificationSentHandlers))
	copy(handlers, n.onNotificationSentHandlers)
	n.mu.RUnlock()

	for _, handler := range handlers {
		go handler(peerDID, summary)
	}

	return nil
}

// SendChatNotification sends a notification for a chat message
func (n *Notifications) SendChatNotification(peerDID, messageText string) error {
	summary := &NotificationSummary{
		Title:       "New Message",
		Body:        truncateText(messageText, 100),
		MessageType: "chat",
		FromDID:     n.client.DID(),
	}

	return n.SendNotification(peerDID, summary)
}

// SendGroupChatNotification sends a notification for a group chat message
func (n *Notifications) SendGroupChatNotification(peerDID, groupName, messageText string) error {
	summary := &NotificationSummary{
		Title:       "New Group Message",
		Body:        groupName + ": " + truncateText(messageText, 80),
		MessageType: "group_chat",
		FromDID:     n.client.DID(),
	}

	return n.SendNotification(peerDID, summary)
}

// SendCredentialNotification sends a notification for credential requests
func (n *Notifications) SendCredentialNotification(peerDID, credentialType, action string) error {
	var title, body string

	switch action {
	case "request":
		title = "Credential Request"
		body = "Someone is requesting your " + credentialType + " credential"
	case "response":
		title = "Credential Response"
		body = "You received a " + credentialType + " credential response"
	case "verification":
		title = "Credential Verification"
		body = "Someone wants to verify your " + credentialType + " credential"
	default:
		title = "Credential Update"
		body = "Credential activity: " + action
	}

	summary := &NotificationSummary{
		Title:       title,
		Body:        body,
		MessageType: "credential",
		FromDID:     n.client.DID(),
	}

	return n.SendNotification(peerDID, summary)
}

// SendGroupInviteNotification sends a notification for group invitations
func (n *Notifications) SendGroupInviteNotification(peerDID, groupName, inviterName string) error {
	summary := &NotificationSummary{
		Title:       "Group Invitation",
		Body:        inviterName + " invited you to join " + groupName,
		MessageType: "group_invite",
		FromDID:     n.client.DID(),
	}

	return n.SendNotification(peerDID, summary)
}

// SendCustomNotification sends a custom notification
func (n *Notifications) SendCustomNotification(peerDID, title, body, messageType string) error {
	summary := &NotificationSummary{
		Title:       title,
		Body:        body,
		MessageType: messageType,
		FromDID:     n.client.DID(),
	}

	return n.SendNotification(peerDID, summary)
}

// CreateSummaryFromContent creates a notification summary from message content
func (n *Notifications) CreateSummaryFromContent(content *message.Content) (*NotificationSummary, error) {
	if n.client.isClosed() {
		return nil, ErrClientClosed
	}

	// Use the SDK's built-in summary creation
	contentSummary, err := content.Summary()
	if err != nil {
		return nil, err
	}

	// Extract information from content summary descriptions
	var title, body, messageType string

	descriptions := contentSummary.Descriptions()
	for _, desc := range descriptions {
		switch desc.DescriptionType() {
		case message.ContentSummaryDescriptionTypeChatMessage:
			if chatMsg, ok := desc.ChatMessage(); ok {
				body = chatMsg
				title = "New Message"
				messageType = "chat"
			}
		case message.ContentSummaryDescriptionTypeCredential:
			if credTypes, ok := desc.Credential(); ok && len(credTypes) > 0 {
				title = "Credential Request"
				body = "Credential: " + credTypes[0]
				messageType = "credential"
			}
		case message.ContentSummaryDescriptionTypePresentation:
			if presTypes, ok := desc.Presentation(); ok && len(presTypes) > 0 {
				title = "Presentation Request"
				body = "Presentation: " + presTypes[0]
				messageType = "presentation"
			}
		}
	}

	// Default values if nothing found
	if title == "" {
		title = "Notification"
	}
	if body == "" {
		body = "You have a new notification"
	}
	if messageType == "" {
		messageType = "unknown"
	}

	summary := &NotificationSummary{
		Title:       title,
		Body:        body,
		MessageType: messageType,
		FromDID:     n.client.DID(),
		MessageID:   string(content.ID()),
	}

	return summary, nil
}

// OnNotificationSent registers a handler for when notifications are sent
func (n *Notifications) OnNotificationSent(handler func(peerDID string, summary *NotificationSummary)) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.onNotificationSentHandlers = append(n.onNotificationSentHandlers, handler)
}

// Internal methods for handling events

func (n *Notifications) onConnect() {
	// Connection established - no specific action needed
}

func (n *Notifications) onDisconnect(err error) {
	// Connection lost - no specific action needed
}

func (n *Notifications) onWelcome(from *signing.PublicKey, groupAddress *signing.PublicKey) {
	// New connection established - no specific action needed
}

func (n *Notifications) onKeyPackage(from *signing.PublicKey) {
	// Key package received - no specific action needed
}

func (n *Notifications) onIntroduction(from *signing.PublicKey, tokenCount int) {
	// Introduction received - no specific action needed
}

func (n *Notifications) close() {
	// Clean up any resources if needed
}

// Helper function to truncate text for notifications
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}
