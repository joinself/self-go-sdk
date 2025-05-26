package client

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
)

// DiscoveryQR represents a QR code for discovery
type DiscoveryQR struct {
	client    *Client
	content   *message.Content
	requestID string
	completer chan *Peer
}

// Peer represents a discovered peer
type Peer struct {
	did     string
	address *signing.PublicKey
}

// Discovery handles peer discovery functionality
type Discovery struct {
	client *Client

	// Event handlers
	onResponseHandlers []func(*Peer)
	mu                 sync.RWMutex
}

// newDiscovery creates a new discovery component
func newDiscovery(client *Client) *Discovery {
	return &Discovery{
		client: client,
	}
}

// GenerateQR creates a new discovery QR code
func (d *Discovery) GenerateQR() (*DiscoveryQR, error) {
	return d.GenerateQRWithTimeout(5 * time.Minute)
}

// GenerateQRWithTimeout creates a discovery QR code with custom timeout
func (d *Discovery) GenerateQRWithTimeout(timeout time.Duration) (*DiscoveryQR, error) {
	if d.client.isClosed() {
		return nil, ErrClientClosed
	}

	// Generate key package for out-of-band negotiation
	keyPackage, err := d.client.account.ConnectionNegotiateOutOfBand(
		d.client.inboxAddress,
		time.Now().Add(timeout),
	)
	if err != nil {
		return nil, err
	}

	// Build discovery request
	content, err := message.NewDiscoveryRequest().
		KeyPackage(keyPackage).
		Expires(time.Now().Add(timeout)).
		Finish()
	if err != nil {
		return nil, err
	}

	requestID := hex.EncodeToString(content.ID())
	completer := make(chan *Peer, 1)

	// Store request for response tracking
	d.client.storeRequest(requestID, completer)

	qr := &DiscoveryQR{
		client:    d.client,
		content:   content,
		requestID: requestID,
		completer: completer,
	}

	return qr, nil
}

// OnResponse registers a handler for discovery responses
func (d *Discovery) OnResponse(handler func(*Peer)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.onResponseHandlers = append(d.onResponseHandlers, handler)
}

// Unicode returns the QR code as Unicode text
func (qr *DiscoveryQR) Unicode() (string, error) {
	anonymousMsg := event.NewAnonymousMessage(qr.content)

	// Set environment-specific flags
	if qr.client.config.Environment == Sandbox {
		anonymousMsg.SetFlags(event.MessageFlagTargetSandbox)
	}

	qrCode, err := anonymousMsg.EncodeToQR(event.QREncodingUnicode)
	if err != nil {
		return "", err
	}
	return string(qrCode), nil
}

// SVG returns the QR code as SVG
func (qr *DiscoveryQR) SVG() (string, error) {
	anonymousMsg := event.NewAnonymousMessage(qr.content)

	// Set environment-specific flags
	if qr.client.config.Environment == Sandbox {
		anonymousMsg.SetFlags(event.MessageFlagTargetSandbox)
	}

	qrCode, err := anonymousMsg.EncodeToQR(event.QREncodingSVG)
	if err != nil {
		return "", err
	}
	return string(qrCode), nil
}

// WaitForResponse waits for someone to scan the QR code and respond
func (qr *DiscoveryQR) WaitForResponse(ctx context.Context) (*Peer, error) {
	select {
	case peer := <-qr.completer:
		return peer, nil
	case <-ctx.Done():
		// Clean up the stored request
		qr.client.loadAndDeleteRequest(qr.requestID)
		return nil, ctx.Err()
	}
}

// RequestID returns the unique identifier for this discovery request
func (qr *DiscoveryQR) RequestID() string {
	return qr.requestID
}

// DID returns the peer's decentralized identifier
func (p *Peer) DID() string {
	return p.did
}

// Address returns the peer's signing public key
func (p *Peer) Address() *signing.PublicKey {
	return p.address
}

// Internal methods for handling events

func (d *Discovery) onConnect() {
	// Connection established - no specific action needed
}

func (d *Discovery) onDisconnect(err error) {
	// Connection lost - no specific action needed for discovery
}

func (d *Discovery) onWelcome(from *signing.PublicKey, groupAddress *signing.PublicKey) {
	// New connection established - no specific action needed
}

func (d *Discovery) onKeyPackage(from *signing.PublicKey) {
	// Key package received - no specific action needed
}

func (d *Discovery) onIntroduction(from *signing.PublicKey, tokenCount int) {
	// Introduction received - no specific action needed for discovery
}

func (d *Discovery) onDiscoveryResponse(msg *event.Message) {
	// Decode the discovery response
	discoveryResponse, err := message.DecodeDiscoveryResponse(msg.Content())
	if err != nil {
		return
	}

	requestID := hex.EncodeToString(discoveryResponse.ResponseTo())

	// Find the waiting request
	completerInterface, ok := d.client.loadAndDeleteRequest(requestID)
	if !ok {
		return
	}

	completer, ok := completerInterface.(chan *Peer)
	if !ok {
		return
	}

	// Create peer object
	peer := &Peer{
		did:     msg.FromAddress().String(),
		address: msg.FromAddress(),
	}

	// Send to waiting request
	select {
	case completer <- peer:
	default:
		// Channel full or closed - ignore
	}

	// Notify subscription handlers
	d.mu.RLock()
	handlers := make([]func(*Peer), len(d.onResponseHandlers))
	copy(handlers, d.onResponseHandlers)
	d.mu.RUnlock()

	for _, handler := range handlers {
		go handler(peer) // Run handlers in goroutines to avoid blocking
	}
}

func (d *Discovery) close() {
	// Clean up any pending requests
	// Note: We could iterate through stored requests and close channels,
	// but the current sync.Map doesn't provide an easy way to do this.
	// For now, pending requests will timeout naturally.
}
