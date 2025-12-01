package client

import (
	"sync"

	"github.com/joinself/self-go-sdk/account"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
)

// Client provides a high-level interface to the Self SDK
type Client struct {
	account *account.Account
	config  *Config

	// Internal state
	inboxAddress *signing.PublicKey
	closed       bool
	mu           sync.RWMutex

	// Request tracking
	requests sync.Map

	// Sub-components
	discovery     *Discovery
	chat          *Chat
	credentials   *Credentials
	groupChats    *GroupChats
	notifications *Notifications
	storage       *Storage
	pairing       *Pairing
	connection    *Connection
}

// New creates a new Self client
func New(config Config) (*Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	client := &Client{
		config: &config,
	}

	// Convert to account config and set up callbacks
	accountConfig := config.toAccountConfig()
	accountConfig.Callbacks = account.Callbacks{
		OnConnect:    client.onConnect,
		OnDisconnect: client.onDisconnect,
		OnWelcome:    client.onWelcome,
		OnKeyPackage: client.onKeyPackage,
		OnMessage:    client.onMessage,
	}

	// Initialize account
	acc, err := account.New(accountConfig)
	if err != nil {
		return nil, err
	}
	client.account = acc

	// Open inbox
	inboxAddress, err := acc.InboxOpen()
	if err != nil {
		return nil, err
	}
	client.inboxAddress = inboxAddress

	// Initialize sub-components
	client.discovery = newDiscovery(client)
	client.chat = newChat(client)
	client.credentials = newCredentials(client)
	client.groupChats = newGroupChats(client)
	client.notifications = newNotifications(client)
	client.storage = newStorage(client)
	client.pairing = newPairing(client)
	client.connection = newConnection(client)

	return client, nil
}

// DID returns the client's decentralized identifier
func (c *Client) DID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.inboxAddress == nil {
		return ""
	}
	return c.inboxAddress.String()
}

// Discovery returns the discovery component
func (c *Client) Discovery() *Discovery {
	return c.discovery
}

// Chat returns the chat component
func (c *Client) Chat() *Chat {
	return c.chat
}

// Credentials returns the credentials component
func (c *Client) Credentials() *Credentials {
	return c.credentials
}

// GroupChats returns the group chats component
func (c *Client) GroupChats() *GroupChats {
	return c.groupChats
}

// Notifications returns the notifications component
func (c *Client) Notifications() *Notifications {
	return c.notifications
}

// Storage returns the storage component
func (c *Client) Storage() *Storage {
	return c.storage
}

// Pairing returns the pairing component
func (c *Client) Pairing() *Pairing {
	return c.pairing
}

// Connection returns the connection component
func (c *Client) Connection() *Connection {
	return c.connection
}

// Close closes the client and releases resources
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	// Close sub-components
	if c.discovery != nil {
		c.discovery.close()
	}
	if c.chat != nil {
		c.chat.close()
	}
	if c.credentials != nil {
		c.credentials.close()
	}
	if c.groupChats != nil {
		c.groupChats.close()
	}
	if c.notifications != nil {
		c.notifications.close()
	}
	if c.storage != nil {
		c.storage.close()
	}
	if c.pairing != nil {
		c.pairing.close()
	}
	if c.connection != nil {
		c.connection.close()
	}

	// Note: The account doesn't have a close method in the current SDK
	// This might need to be added to the underlying SDK

	return nil
}

// Internal methods for handling account callbacks

func (c *Client) onConnect(account *account.Account) {
	// Connection established - notify sub-components
	if c.discovery != nil {
		c.discovery.onConnect()
	}
	if c.chat != nil {
		c.chat.onConnect()
	}
	if c.credentials != nil {
		c.credentials.onConnect()
	}
	if c.groupChats != nil {
		c.groupChats.onConnect()
	}
	if c.notifications != nil {
		c.notifications.onConnect()
	}
	if c.storage != nil {
		c.storage.onConnect()
	}
	if c.pairing != nil {
		c.pairing.onConnect()
	}
	if c.connection != nil {
		c.connection.onConnect()
	}
}

func (c *Client) onDisconnect(account *account.Account, err error) {
	// Connection lost - notify sub-components
	if c.discovery != nil {
		c.discovery.onDisconnect(err)
	}
	if c.chat != nil {
		c.chat.onDisconnect(err)
	}
	if c.credentials != nil {
		c.credentials.onDisconnect(err)
	}
	if c.groupChats != nil {
		c.groupChats.onDisconnect(err)
	}
	if c.notifications != nil {
		c.notifications.onDisconnect(err)
	}
	if c.storage != nil {
		c.storage.onDisconnect(err)
	}
	if c.pairing != nil {
		c.pairing.onDisconnect(err)
	}
	if c.connection != nil {
		c.connection.onDisconnect(err)
	}
}

func (c *Client) onWelcome(account *account.Account, wlc *event.Welcome) {
	// Accept the connection automatically
	groupAddress, err := account.ConnectionAccept(
		wlc.ToAddress(),
		wlc.Welcome(),
	)
	if err != nil {
		// Log error but don't fail - this is handled internally
		return
	}

	// Notify sub-components of new connection
	if c.discovery != nil {
		c.discovery.onWelcome(wlc.FromAddress(), groupAddress)
	}
	if c.chat != nil {
		c.chat.onWelcome(wlc.FromAddress(), groupAddress)
	}
	if c.credentials != nil {
		c.credentials.onWelcome(wlc.FromAddress(), groupAddress)
	}
	if c.groupChats != nil {
		c.groupChats.onWelcome(wlc.FromAddress(), groupAddress)
	}
	if c.notifications != nil {
		c.notifications.onWelcome(wlc.FromAddress(), groupAddress)
	}
	if c.storage != nil {
		c.storage.onWelcome(wlc.FromAddress(), groupAddress)
	}
	if c.pairing != nil {
		c.pairing.onWelcome(wlc.FromAddress(), groupAddress)
	}
	if c.connection != nil {
		c.connection.onWelcome(wlc.FromAddress(), groupAddress)
	}
}

func (c *Client) onKeyPackage(account *account.Account, kp *event.KeyPackage) {
	// Establish connection automatically
	_, err := account.ConnectionEstablish(kp.ToAddress(), kp.KeyPackage())
	if err != nil {
		// Log error but don't fail - this is handled internally
		return
	}

	// Notify sub-components
	if c.discovery != nil {
		c.discovery.onKeyPackage(kp.FromAddress())
	}
	if c.chat != nil {
		c.chat.onKeyPackage(kp.FromAddress())
	}
	if c.credentials != nil {
		c.credentials.onKeyPackage(kp.FromAddress())
	}
	if c.groupChats != nil {
		c.groupChats.onKeyPackage(kp.FromAddress())
	}
	if c.notifications != nil {
		c.notifications.onKeyPackage(kp.FromAddress())
	}
	if c.storage != nil {
		c.storage.onKeyPackage(kp.FromAddress())
	}
	if c.pairing != nil {
		c.pairing.onKeyPackage(kp.FromAddress())
	}
	if c.connection != nil {
		c.connection.onKeyPackage(kp.FromAddress())
	}
}

func (c *Client) onMessage(account *account.Account, msg *event.Message) {
	// Route messages to appropriate handlers based on content type
	switch event.ContentTypeOf(msg) {
	case message.ContentTypeDiscoveryResponse:
		if c.discovery != nil {
			c.discovery.onDiscoveryResponse(msg)
		}
	case message.ContentTypeChat:
		if c.chat != nil {
			c.chat.onChatMessage(msg)
		}
		if c.groupChats != nil {
			c.groupChats.onChatMessage(msg)
		}
	case message.ContentTypeCredentialPresentationRequest:
		if c.credentials != nil {
			c.credentials.onCredentialPresentationRequest(msg)
		}
	case message.ContentTypeCredentialPresentationResponse:
		if c.credentials != nil {
			c.credentials.onCredentialPresentationResponse(msg)
		}
	case message.ContentTypeCredentialVerificationRequest:
		if c.credentials != nil {
			c.credentials.onCredentialVerificationRequest(msg)
		}
	case message.ContentTypeCredentialVerificationResponse:
		if c.credentials != nil {
			c.credentials.onCredentialVerificationResponse(msg)
		}
	case message.ContentTypeAccountPairingRequest:
		if c.pairing != nil {
			c.pairing.onAccountPairingRequest(msg)
		}
	case message.ContentTypeAccountPairingResponse:
		if c.pairing != nil {
			c.pairing.onAccountPairingResponse(msg)
		}
	case message.ContentTypeIntroduction:
		// Handle introduction messages - these establish tokens for communication
		c.handleIntroduction(msg)
	default:
		// Unknown message type - could be handled by extensions in the future
	}
}

func (c *Client) handleIntroduction(msg *event.Message) {
	introduction, err := message.DecodeIntroduction(msg.Content())
	if err != nil {
		return
	}

	tokens, err := introduction.Tokens()
	if err != nil {
		return
	}

	// Store tokens for future communication
	for _, token := range tokens {
		err = c.account.TokenStore(
			msg.FromAddress(),
			msg.ToAddress(),
			msg.ToAddress(),
			token,
		)
		if err != nil {
			// Log error but continue with other tokens
			continue
		}
	}

	// Notify sub-components of introduction
	if c.discovery != nil {
		c.discovery.onIntroduction(msg.FromAddress(), len(tokens))
	}
	if c.chat != nil {
		c.chat.onIntroduction(msg.FromAddress(), len(tokens))
	}
	if c.credentials != nil {
		c.credentials.onIntroduction(msg.FromAddress(), len(tokens))
	}
	if c.groupChats != nil {
		c.groupChats.onIntroduction(msg.FromAddress(), len(tokens))
	}
	if c.notifications != nil {
		c.notifications.onIntroduction(msg.FromAddress(), len(tokens))
	}
	if c.storage != nil {
		c.storage.onIntroduction(msg.FromAddress(), len(tokens))
	}
	if c.pairing != nil {
		c.pairing.onIntroduction(msg.FromAddress(), len(tokens))
	}
	if c.connection != nil {
		c.connection.onIntroduction(msg.FromAddress(), len(tokens))
	}
}

// Internal helper methods

func (c *Client) isClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

func (c *Client) storeRequest(requestID string, completer interface{}) {
	c.requests.Store(requestID, completer)
}

func (c *Client) loadAndDeleteRequest(requestID string) (interface{}, bool) {
	return c.requests.LoadAndDelete(requestID)
}

func (c *Client) sendMessage(to *signing.PublicKey, content message.Content) error {
	if c.isClosed() {
		return ErrClientClosed
	}
	return c.account.MessageSend(to, &content)
}
