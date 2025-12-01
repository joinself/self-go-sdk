package client

import (
	"sync"

	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
)

// ChatMessage represents a received chat message
type ChatMessage struct {
	from        string
	text        string
	id          string
	refID       string
	attachments []ChatAttachment
}

// ChatAttachment represents a file attachment in a chat message
type ChatAttachment struct {
	name string
	data []byte
	mime string
}

// Chat handles chat messaging functionality
type Chat struct {
	client *Client

	// Event handlers
	onMessageHandlers []func(ChatMessage)
	mu                sync.RWMutex
}

// newChat creates a new chat component
func newChat(client *Client) *Chat {
	return &Chat{
		client: client,
	}
}

// OnMessage registers a handler for incoming chat messages
func (c *Chat) OnMessage(handler func(ChatMessage)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onMessageHandlers = append(c.onMessageHandlers, handler)
}

// Send sends a chat message to a peer
func (c *Chat) Send(peerDID string, messageText string) error {
	return c.SendWithAttachments(peerDID, messageText, nil)
}

// SendWithAttachments sends a chat message with file attachments
func (c *Chat) SendWithAttachments(peerDID string, messageText string, attachments []ChatAttachment) error {
	if c.client.isClosed() {
		return ErrClientClosed
	}

	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(peerDID)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	// Build the chat message
	chatBuilder := message.NewChat().Message(messageText)

	// Add attachments if any
	for _, attachment := range attachments {
		// TODO: Implement attachment handling
		// This would involve uploading the attachment to the object store
		// and including the reference in the message
		_ = attachment // Suppress unused variable warning for now
	}

	content, err := chatBuilder.Finish()
	if err != nil {
		return err
	}

	return c.client.sendMessage(peerAddress, *content)
}

// Reply sends a reply to a specific message
func (c *Chat) Reply(originalMessage ChatMessage, replyText string) error {
	if c.client.isClosed() {
		return ErrClientClosed
	}

	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(originalMessage.from)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	// Build the reply message with reference to original
	content, err := message.NewChat().
		Message(replyText).
		Reference([]byte(originalMessage.id)).
		Finish()
	if err != nil {
		return err
	}

	return c.client.sendMessage(peerAddress, *content)
}

// From returns the sender's DID
func (m ChatMessage) From() string {
	return m.from
}

// Text returns the message text
func (m ChatMessage) Text() string {
	return m.text
}

// ID returns the message ID
func (m ChatMessage) ID() string {
	return m.id
}

// ReferencedID returns the ID of the message this is replying to (if any)
func (m ChatMessage) ReferencedID() string {
	return m.refID
}

// Attachments returns the message attachments
func (m ChatMessage) Attachments() []ChatAttachment {
	return m.attachments
}

// Name returns the attachment filename
func (a ChatAttachment) Name() string {
	return a.name
}

// Data returns the attachment data
func (a ChatAttachment) Data() []byte {
	return a.data
}

// MimeType returns the attachment MIME type
func (a ChatAttachment) MimeType() string {
	return a.mime
}

// Internal methods for handling events

func (c *Chat) onConnect() {
	// Connection established - no specific action needed
}

func (c *Chat) onDisconnect(err error) {
	// Connection lost - no specific action needed for chat
}

func (c *Chat) onWelcome(from *signing.PublicKey, groupAddress *signing.PublicKey) {
	// New connection established - no specific action needed
}

func (c *Chat) onKeyPackage(from *signing.PublicKey) {
	// Key package received - no specific action needed
}

func (c *Chat) onIntroduction(from *signing.PublicKey, tokenCount int) {
	// Introduction received - connection is now ready for chat
}

func (c *Chat) onChatMessage(msg *event.Message) {
	// Decode the chat message
	chat, err := message.DecodeChat(msg.Content())
	if err != nil {
		return
	}

	// Create ChatMessage object
	chatMessage := ChatMessage{
		from:  msg.FromAddress().String(),
		text:  chat.Message(),
		id:    string(msg.ID()),
		refID: string(chat.Referencing()),
		// TODO: Handle attachments
		attachments: []ChatAttachment{},
	}

	// Notify handlers
	c.mu.RLock()
	handlers := make([]func(ChatMessage), len(c.onMessageHandlers))
	copy(handlers, c.onMessageHandlers)
	c.mu.RUnlock()

	for _, handler := range handlers {
		go handler(chatMessage) // Run handlers in goroutines to avoid blocking
	}
}

func (c *Chat) close() {
	// Clean up any resources if needed
}
