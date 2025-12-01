package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/identity"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/object"
)

// PairingRequest represents an account pairing request
type PairingRequest struct {
	client    *Client
	content   *message.Content
	requestID string
	completer chan *PairingResponse
}

// PairingResponse represents an account pairing response
type PairingResponse struct {
	from      string
	status    message.ResponseStatus
	operation *identity.Operation
	assets    []*object.Object
}

// IncomingPairingRequest represents an incoming pairing request
type IncomingPairingRequest struct {
	from      string
	requestID string
	address   *signing.PublicKey
	roles     identity.Role
	expires   time.Time
	client    *Client
}

// PairingCode represents a pairing code for linking accounts
type PairingCode struct {
	Code      string
	Unpaired  bool
	ExpiresAt time.Time
}

// Pairing handles account pairing and linking functionality
type Pairing struct {
	client *Client

	// Event handlers
	onPairingRequestHandlers  []func(*IncomingPairingRequest)
	onPairingResponseHandlers []func(*PairingResponse)
	mu                        sync.RWMutex
}

// newPairing creates a new pairing component
func newPairing(client *Client) *Pairing {
	return &Pairing{
		client: client,
	}
}

// GetPairingCode returns the SDK pairing code for linking accounts
func (p *Pairing) GetPairingCode() (*PairingCode, error) {
	if p.client.isClosed() {
		return nil, ErrClientClosed
	}

	code, unpaired, err := p.client.account.SDKPairingCode()
	if err != nil {
		return nil, err
	}

	return &PairingCode{
		Code:      code,
		Unpaired:  unpaired,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Default 24 hour expiry
	}, nil
}

// RequestPairing sends a pairing request to another account
func (p *Pairing) RequestPairing(peerDID string, address *signing.PublicKey, roles identity.Role) (*PairingRequest, error) {
	return p.RequestPairingWithTimeout(peerDID, address, roles, 5*time.Minute)
}

// RequestPairingWithTimeout sends a pairing request with a custom timeout
func (p *Pairing) RequestPairingWithTimeout(peerDID string, address *signing.PublicKey, roles identity.Role, timeout time.Duration) (*PairingRequest, error) {
	if p.client.isClosed() {
		return nil, ErrClientClosed
	}

	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(peerDID)
	if peerAddress == nil {
		return nil, ErrInvalidPeerDID
	}

	// Build the pairing request
	expires := time.Now().Add(timeout)
	content, err := message.NewAccountPairingRequest().
		Address(address).
		Roles(roles).
		Expires(expires).
		Finish()
	if err != nil {
		return nil, err
	}

	// Create request tracker
	requestID := hex.EncodeToString(content.ID())
	completer := make(chan *PairingResponse, 1)

	request := &PairingRequest{
		client:    p.client,
		content:   content,
		requestID: requestID,
		completer: completer,
	}

	// Store the request for response matching
	p.client.storeRequest(requestID, completer)

	// Send the request
	err = p.client.sendMessage(peerAddress, *content)
	if err != nil {
		p.client.loadAndDeleteRequest(requestID) // Clean up on error
		return nil, err
	}

	return request, nil
}

// WaitForResponse waits for a pairing response
func (pr *PairingRequest) WaitForResponse(ctx context.Context) (*PairingResponse, error) {
	select {
	case response := <-pr.completer:
		return response, nil
	case <-ctx.Done():
		// Clean up the stored request
		pr.client.loadAndDeleteRequest(pr.requestID)
		return nil, ctx.Err()
	}
}

// RequestID returns the unique request identifier
func (pr *PairingRequest) RequestID() string {
	return pr.requestID
}

// RespondWithOperation responds to a pairing request with an identity operation
func (ipr *IncomingPairingRequest) RespondWithOperation(operation *identity.Operation) error {
	return ipr.RespondWithOperationAndAssets(operation, nil)
}

// RespondWithOperationAndAssets responds to a pairing request with an operation and supporting assets
func (ipr *IncomingPairingRequest) RespondWithOperationAndAssets(operation *identity.Operation, assets []*object.Object) error {
	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(ipr.from)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	// Build the response
	responseBuilder := message.NewAccountPairingResponse().
		ResponseTo([]byte(ipr.requestID)).
		Status(message.ResponseStatusAccepted).
		Operation(operation)

	// Add assets if provided
	for _, asset := range assets {
		responseBuilder = responseBuilder.Asset(asset)
	}

	content, err := responseBuilder.Finish()
	if err != nil {
		return err
	}

	return ipr.client.sendMessage(peerAddress, *content)
}

// Reject rejects the pairing request
func (ipr *IncomingPairingRequest) Reject() error {
	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(ipr.from)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	// Build the rejection response
	content, err := message.NewAccountPairingResponse().
		ResponseTo([]byte(ipr.requestID)).
		Status(message.ResponseStatusForbidden).
		Finish()
	if err != nil {
		return err
	}

	return ipr.client.sendMessage(peerAddress, *content)
}

// From returns the sender's DID
func (ipr *IncomingPairingRequest) From() string {
	return ipr.from
}

// RequestID returns the request ID
func (ipr *IncomingPairingRequest) RequestID() string {
	return ipr.requestID
}

// Address returns the address to be paired
func (ipr *IncomingPairingRequest) Address() *signing.PublicKey {
	return ipr.address
}

// Roles returns the requested roles
func (ipr *IncomingPairingRequest) Roles() identity.Role {
	return ipr.roles
}

// Expires returns when the request expires
func (ipr *IncomingPairingRequest) Expires() time.Time {
	return ipr.expires
}

// PairingResponse methods

// From returns the sender's DID
func (pr *PairingResponse) From() string {
	return pr.from
}

// Status returns the response status
func (pr *PairingResponse) Status() message.ResponseStatus {
	return pr.status
}

// Operation returns the identity operation (if accepted)
func (pr *PairingResponse) Operation() *identity.Operation {
	return pr.operation
}

// Assets returns supporting assets
func (pr *PairingResponse) Assets() []*object.Object {
	return pr.assets
}

// Event handler registration methods

// OnPairingRequest registers a handler for incoming pairing requests
func (p *Pairing) OnPairingRequest(handler func(*IncomingPairingRequest)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.onPairingRequestHandlers = append(p.onPairingRequestHandlers, handler)
}

// OnPairingResponse registers a handler for pairing responses
func (p *Pairing) OnPairingResponse(handler func(*PairingResponse)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.onPairingResponseHandlers = append(p.onPairingResponseHandlers, handler)
}

// Internal methods for handling events

func (p *Pairing) onConnect() {
	// Connection established - no specific action needed
}

func (p *Pairing) onDisconnect(err error) {
	// Connection lost - no specific action needed
}

func (p *Pairing) onWelcome(from *signing.PublicKey, groupAddress *signing.PublicKey) {
	// New connection established - no specific action needed
}

func (p *Pairing) onKeyPackage(from *signing.PublicKey) {
	// Key package received - no specific action needed
}

func (p *Pairing) onIntroduction(from *signing.PublicKey, tokenCount int) {
	// Introduction received - no specific action needed
}

func (p *Pairing) onAccountPairingRequest(msg *event.Message) {
	// Decode the account pairing request
	pairingRequest, err := message.DecodeAccountPairingRequest(msg.Content())
	if err != nil {
		return
	}

	// Create incoming request object
	incomingRequest := &IncomingPairingRequest{
		from:      msg.FromAddress().String(),
		requestID: hex.EncodeToString(msg.ID()),
		address:   pairingRequest.Address(),
		roles:     pairingRequest.Roles(),
		expires:   pairingRequest.Expires(),
		client:    p.client,
	}

	// Notify handlers
	p.mu.RLock()
	handlers := make([]func(*IncomingPairingRequest), len(p.onPairingRequestHandlers))
	copy(handlers, p.onPairingRequestHandlers)
	p.mu.RUnlock()

	for _, handler := range handlers {
		go handler(incomingRequest)
	}
}

func (p *Pairing) onAccountPairingResponse(msg *event.Message) {
	// Decode the account pairing response
	pairingResponse, err := message.DecodeAccountPairingResponse(msg.Content())
	if err != nil {
		return
	}

	requestID := hex.EncodeToString(pairingResponse.ResponseTo())

	// Find the waiting request
	completerInterface, ok := p.client.loadAndDeleteRequest(requestID)
	if !ok {
		return
	}

	completer, ok := completerInterface.(chan *PairingResponse)
	if !ok {
		return
	}

	// Create response object
	response := &PairingResponse{
		from:      msg.FromAddress().String(),
		status:    pairingResponse.Status(),
		operation: pairingResponse.Operation(),
		assets:    pairingResponse.Assets(),
	}

	// Send to waiting request
	select {
	case completer <- response:
	default:
		// Channel full or closed - ignore
	}

	// Notify subscription handlers
	p.mu.RLock()
	handlers := make([]func(*PairingResponse), len(p.onPairingResponseHandlers))
	copy(handlers, p.onPairingResponseHandlers)
	p.mu.RUnlock()

	for _, handler := range handlers {
		go handler(response)
	}
}

func (p *Pairing) close() {
	// Clean up any resources if needed
}

// Helper functions

// GeneratePairingQR generates a QR code containing the pairing code
func (p *Pairing) GeneratePairingQR() (string, error) {
	pairingCode, err := p.GetPairingCode()
	if err != nil {
		return "", err
	}

	// For now, return the code as text
	// In a real implementation, this would generate an actual QR code
	return fmt.Sprintf("SELF_PAIRING:%s", pairingCode.Code), nil
}

// IsPaired checks if the account is currently paired
func (p *Pairing) IsPaired() (bool, error) {
	pairingCode, err := p.GetPairingCode()
	if err != nil {
		return false, err
	}

	return !pairingCode.Unpaired, nil
}
