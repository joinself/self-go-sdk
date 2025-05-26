package client

import (
	"context"
	"fmt"
	"time"

	"github.com/joinself/self-go-sdk/keypair/signing"
)

// Connection handles direct peer-to-peer connections
type Connection struct {
	client *Client
}

// ConnectionResult represents the result of a connection attempt
type ConnectionResult struct {
	PeerDID   string
	Connected bool
	Error     error
}

// newConnection creates a new connection component
func newConnection(client *Client) *Connection {
	return &Connection{
		client: client,
	}
}

// ConnectToPeer establishes a direct connection to another peer programmatically
// This bypasses the QR code discovery process and directly negotiates a connection
func (c *Connection) ConnectToPeer(peerDID string) (*ConnectionResult, error) {
	return c.ConnectToPeerWithTimeout(peerDID, 30*time.Second)
}

// ConnectToPeerWithTimeout establishes a connection with a custom timeout
func (c *Connection) ConnectToPeerWithTimeout(peerDID string, timeout time.Duration) (*ConnectionResult, error) {
	if c.client.isClosed() {
		return nil, ErrClientClosed
	}

	// Parse the peer DID to get the signing address
	peerAddress := signing.FromAddress(peerDID)
	if peerAddress == nil {
		return &ConnectionResult{
			PeerDID:   peerDID,
			Connected: false,
			Error:     fmt.Errorf("invalid peer DID: %s", peerDID),
		}, nil
	}

	// Get our own address
	ourAddress := signing.FromAddress(c.client.DID())
	if ourAddress == nil {
		return &ConnectionResult{
			PeerDID:   peerDID,
			Connected: false,
			Error:     fmt.Errorf("invalid client DID: %s", c.client.DID()),
		}, nil
	}

	// Set up a channel to track connection establishment
	connectionEstablished := make(chan bool, 1)
	connectionError := make(chan error, 1)

	// Set up temporary handlers to track the connection
	originalOnWelcome := c.client.discovery.onResponseHandlers
	c.client.discovery.OnResponse(func(peer *Peer) {
		if peer.DID() == peerDID {
			connectionEstablished <- true
		}
	})

	// Restore original handlers after connection attempt
	defer func() {
		c.client.discovery.onResponseHandlers = originalOnWelcome
	}()

	// Initiate the connection negotiation
	err := c.client.account.ConnectionNegotiate(
		ourAddress,
		peerAddress,
		time.Now().Add(timeout),
	)
	if err != nil {
		return &ConnectionResult{
			PeerDID:   peerDID,
			Connected: false,
			Error:     fmt.Errorf("failed to negotiate connection: %w", err),
		}, nil
	}

	// Wait for connection establishment or timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-connectionEstablished:
		return &ConnectionResult{
			PeerDID:   peerDID,
			Connected: true,
			Error:     nil,
		}, nil
	case err := <-connectionError:
		return &ConnectionResult{
			PeerDID:   peerDID,
			Connected: false,
			Error:     err,
		}, nil
	case <-ctx.Done():
		return &ConnectionResult{
			PeerDID:   peerDID,
			Connected: false,
			Error:     fmt.Errorf("connection timeout after %v", timeout),
		}, nil
	}
}

// ConnectTwoClients is a utility function to establish a bidirectional connection
// between two clients in the same process (useful for demos and testing)
func ConnectTwoClients(client1, client2 *Client) error {
	return ConnectTwoClientsWithTimeout(client1, client2, 30*time.Second)
}

// ConnectTwoClientsWithTimeout establishes a connection between two clients with a timeout
func ConnectTwoClientsWithTimeout(client1, client2 *Client, timeout time.Duration) error {
	if client1.isClosed() || client2.isClosed() {
		return ErrClientClosed
	}

	// Get addresses for both clients
	address1 := signing.FromAddress(client1.DID())
	address2 := signing.FromAddress(client2.DID())

	if address1 == nil {
		return fmt.Errorf("invalid client1 DID: %s", client1.DID())
	}
	if address2 == nil {
		return fmt.Errorf("invalid client2 DID: %s", client2.DID())
	}

	// Set up connection tracking
	connectionEstablished := make(chan bool, 2)
	connectionCount := 0

	// Set up handlers for both clients
	client1.discovery.OnResponse(func(peer *Peer) {
		if peer.DID() == client2.DID() {
			connectionEstablished <- true
		}
	})

	client2.discovery.OnResponse(func(peer *Peer) {
		if peer.DID() == client1.DID() {
			connectionEstablished <- true
		}
	})

	// Initiate connection from client1 to client2
	err := client1.account.ConnectionNegotiate(
		address1,
		address2,
		time.Now().Add(timeout),
	)
	if err != nil {
		return fmt.Errorf("failed to negotiate connection from client1 to client2: %w", err)
	}

	// Wait for connection establishment
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for connectionCount < 1 { // We only need one direction to establish the connection
		select {
		case <-connectionEstablished:
			connectionCount++
		case <-ctx.Done():
			return fmt.Errorf("connection timeout after %v", timeout)
		}
	}

	return nil
}

// IsConnectedTo checks if this client is connected to a specific peer
func (c *Connection) IsConnectedTo(peerDID string) bool {
	// This would require tracking active connections
	// For now, we'll return false as a placeholder
	// In a full implementation, this would check the connection state
	return false
}

// ListConnectedPeers returns a list of currently connected peer DIDs
func (c *Connection) ListConnectedPeers() []string {
	// This would require tracking active connections
	// For now, we'll return an empty slice as a placeholder
	// In a full implementation, this would return actual connected peers
	return []string{}
}

// Internal methods for handling events

func (c *Connection) onConnect() {
	// Connection to messaging service established
}

func (c *Connection) onDisconnect(err error) {
	// Connection to messaging service lost
}

func (c *Connection) onWelcome(from *signing.PublicKey, groupAddress *signing.PublicKey) {
	// New peer connection established
}

func (c *Connection) onKeyPackage(from *signing.PublicKey) {
	// Key package received from peer
}

func (c *Connection) onIntroduction(from *signing.PublicKey, tokenCount int) {
	// Introduction received from peer
}

func (c *Connection) close() {
	// Cleanup connection resources
}
