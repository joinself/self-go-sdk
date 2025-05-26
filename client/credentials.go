package client

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/keypair/signing"
	"github.com/joinself/self-go-sdk/message"
	"github.com/joinself/self-go-sdk/object"
)

// CredentialRequest represents a credential request
type CredentialRequest struct {
	client    *Client
	content   *message.Content
	requestID string
	completer chan *CredentialResponse
}

// CredentialResponse represents a credential response
type CredentialResponse struct {
	from          string
	status        message.ResponseStatus
	presentations []*credential.VerifiablePresentation
	credentials   []*credential.VerifiableCredential
}

// CredentialDetail represents details for a credential request
type CredentialDetail struct {
	CredentialType []string
	Parameters     []*CredentialParameter
}

// CredentialParameter represents a parameter for credential filtering
type CredentialParameter struct {
	Operator message.ComparisonOperator
	Field    string
	Value    string
}

// CredentialBuilder represents a builder for creating custom credentials
type CredentialBuilder struct {
	credentialType []string
	subject        *credential.Address
	claims         map[string]interface{}
	issuer         *credential.Address
	validFrom      time.Time
	signer         *signing.PublicKey
	issuedAt       time.Time
}

// CredentialEvidence represents evidence attached to credential requests
type CredentialEvidence struct {
	Type   string
	Object *object.Object
}

// CredentialAsset represents an asset/file that can be attached to credentials
type CredentialAsset struct {
	Name     string
	MimeType string
	Data     []byte
	object   *object.Object
}

// Credentials handles credential exchange functionality
type Credentials struct {
	client *Client

	// Event handlers
	onPresentationRequestHandlers  []func(*IncomingCredentialRequest)
	onVerificationRequestHandlers  []func(*IncomingCredentialRequest)
	onPresentationResponseHandlers []func(*CredentialResponse)
	onVerificationResponseHandlers []func(*CredentialResponse)
	mu                             sync.RWMutex
}

// IncomingCredentialRequest represents an incoming credential request
type IncomingCredentialRequest struct {
	from           string
	requestID      string
	reqType        []string
	details        []*CredentialDetail
	evidence       []*CredentialEvidence
	proof          []*credential.VerifiablePresentation
	expires        time.Time
	client         *Client
	isVerification bool
}

// newCredentials creates a new credentials component
func newCredentials(client *Client) *Credentials {
	return &Credentials{
		client: client,
	}
}

// NewCredentialBuilder creates a new credential builder for custom credentials
func (c *Credentials) NewCredentialBuilder() *CredentialBuilder {
	return &CredentialBuilder{
		claims: make(map[string]interface{}),
	}
}

// CreateAsset creates a new credential asset from file data
func (c *Credentials) CreateAsset(name, mimeType string, data []byte) (*CredentialAsset, error) {
	if c.client.isClosed() {
		return nil, ErrClientClosed
	}

	// Create encrypted object
	obj, err := object.New(mimeType, data)
	if err != nil {
		return nil, err
	}

	// Upload to object store
	err = c.client.account.ObjectUpload(obj, false)
	if err != nil {
		return nil, err
	}

	return &CredentialAsset{
		Name:     name,
		MimeType: mimeType,
		Data:     data,
		object:   obj,
	}, nil
}

// DownloadAsset downloads and decrypts an asset
func (c *Credentials) DownloadAsset(asset *CredentialAsset) error {
	if c.client.isClosed() {
		return ErrClientClosed
	}

	return c.client.account.ObjectDownload(asset.object)
}

// RequestPresentationWithEvidence requests credential presentations with evidence attachments
func (c *Credentials) RequestPresentationWithEvidence(peerDID string, details []*CredentialDetail, evidence []*CredentialEvidence, proof []*credential.VerifiablePresentation) (*CredentialRequest, error) {
	return c.RequestPresentationWithEvidenceAndTimeout(peerDID, details, evidence, proof, 5*time.Minute)
}

// RequestPresentationWithEvidenceAndTimeout requests credential presentations with evidence and custom timeout
func (c *Credentials) RequestPresentationWithEvidenceAndTimeout(peerDID string, details []*CredentialDetail, evidence []*CredentialEvidence, proof []*credential.VerifiablePresentation, timeout time.Duration) (*CredentialRequest, error) {
	if c.client.isClosed() {
		return nil, ErrClientClosed
	}

	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(peerDID)
	if peerAddress == nil {
		return nil, ErrInvalidPeerDID
	}

	// Build the credential presentation request
	builder := message.NewCredentialPresentationRequest().
		Type([]string{"VerifiablePresentation", "CustomPresentation"}).
		Expires(time.Now().Add(timeout))

	// Add details for each credential type
	for _, detail := range details {
		params := make([]*message.CredentialPresentationDetailParameter, len(detail.Parameters))
		for i, param := range detail.Parameters {
			params[i] = message.NewCredentialPresentationDetailParameter(
				param.Operator,
				param.Field,
				param.Value,
			)
		}
		builder.Details(detail.CredentialType, params)
	}

	// Add proof presentations
	for _, p := range proof {
		builder.Proof(p)
	}

	content, err := builder.Finish()
	if err != nil {
		return nil, err
	}

	requestID := hex.EncodeToString(content.ID())
	completer := make(chan *CredentialResponse, 1)

	// Store request for response tracking
	c.client.storeRequest(requestID, completer)

	// Send the request
	err = c.client.sendMessage(peerAddress, *content)
	if err != nil {
		c.client.loadAndDeleteRequest(requestID)
		return nil, err
	}

	req := &CredentialRequest{
		client:    c.client,
		content:   content,
		requestID: requestID,
		completer: completer,
	}

	return req, nil
}

// RequestVerificationWithEvidence requests credential verification with evidence attachments
func (c *Credentials) RequestVerificationWithEvidence(peerDID string, credentialType []string, evidence []*CredentialEvidence, proof []*credential.VerifiablePresentation) (*CredentialRequest, error) {
	return c.RequestVerificationWithEvidenceAndTimeout(peerDID, credentialType, evidence, proof, 5*time.Minute)
}

// RequestVerificationWithEvidenceAndTimeout requests credential verification with evidence and custom timeout
func (c *Credentials) RequestVerificationWithEvidenceAndTimeout(peerDID string, credentialType []string, evidence []*CredentialEvidence, proof []*credential.VerifiablePresentation, timeout time.Duration) (*CredentialRequest, error) {
	if c.client.isClosed() {
		return nil, ErrClientClosed
	}

	// Parse the peer DID to get the signing key
	peerAddress := signing.FromAddress(peerDID)
	if peerAddress == nil {
		return nil, ErrInvalidPeerDID
	}

	// Build the credential verification request
	builder := message.NewCredentialVerificationRequest().
		Type(credentialType).
		Expires(time.Now().Add(timeout))

	// Add evidence
	for _, ev := range evidence {
		builder.Evidence(ev.Type, ev.Object)
	}

	// Add proof presentations
	for _, p := range proof {
		builder.Proof(p)
	}

	content, err := builder.Finish()
	if err != nil {
		return nil, err
	}

	requestID := hex.EncodeToString(content.ID())
	completer := make(chan *CredentialResponse, 1)

	// Store request for response tracking
	c.client.storeRequest(requestID, completer)

	// Send the request
	err = c.client.sendMessage(peerAddress, *content)
	if err != nil {
		c.client.loadAndDeleteRequest(requestID)
		return nil, err
	}

	req := &CredentialRequest{
		client:    c.client,
		content:   content,
		requestID: requestID,
		completer: completer,
	}

	return req, nil
}

// RequestPresentation requests credential presentations from a peer
func (c *Credentials) RequestPresentation(peerDID string, details []*CredentialDetail) (*CredentialRequest, error) {
	return c.RequestPresentationWithTimeout(peerDID, details, 5*time.Minute)
}

// RequestPresentationWithTimeout requests credential presentations with custom timeout
func (c *Credentials) RequestPresentationWithTimeout(peerDID string, details []*CredentialDetail, timeout time.Duration) (*CredentialRequest, error) {
	return c.RequestPresentationWithEvidenceAndTimeout(peerDID, details, nil, nil, timeout)
}

// RequestVerification requests credential verification from a peer
func (c *Credentials) RequestVerification(peerDID string, credentialType []string) (*CredentialRequest, error) {
	return c.RequestVerificationWithTimeout(peerDID, credentialType, 5*time.Minute)
}

// RequestVerificationWithTimeout requests credential verification with custom timeout
func (c *Credentials) RequestVerificationWithTimeout(peerDID string, credentialType []string, timeout time.Duration) (*CredentialRequest, error) {
	return c.RequestVerificationWithEvidenceAndTimeout(peerDID, credentialType, nil, nil, timeout)
}

// OnPresentationRequest registers a handler for incoming credential presentation requests
func (c *Credentials) OnPresentationRequest(handler func(*IncomingCredentialRequest)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onPresentationRequestHandlers = append(c.onPresentationRequestHandlers, handler)
}

// OnVerificationRequest registers a handler for incoming credential verification requests
func (c *Credentials) OnVerificationRequest(handler func(*IncomingCredentialRequest)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onVerificationRequestHandlers = append(c.onVerificationRequestHandlers, handler)
}

// OnPresentationResponse registers a handler for credential presentation responses
func (c *Credentials) OnPresentationResponse(handler func(*CredentialResponse)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onPresentationResponseHandlers = append(c.onPresentationResponseHandlers, handler)
}

// OnVerificationResponse registers a handler for credential verification responses
func (c *Credentials) OnVerificationResponse(handler func(*CredentialResponse)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onVerificationResponseHandlers = append(c.onVerificationResponseHandlers, handler)
}

// CredentialBuilder methods

// Type sets the credential type
func (b *CredentialBuilder) Type(credentialType []string) *CredentialBuilder {
	b.credentialType = credentialType
	return b
}

// Subject sets the credential subject
func (b *CredentialBuilder) Subject(subjectDID string) *CredentialBuilder {
	subjectAddress := signing.FromAddress(subjectDID)
	if subjectAddress != nil {
		b.subject = credential.AddressKey(subjectAddress)
	}
	return b
}

// Claim adds a claim to the credential
func (b *CredentialBuilder) Claim(key string, value interface{}) *CredentialBuilder {
	b.claims[key] = value
	return b
}

// Claims sets multiple claims at once
func (b *CredentialBuilder) Claims(claims map[string]interface{}) *CredentialBuilder {
	for k, v := range claims {
		b.claims[k] = v
	}
	return b
}

// Issuer sets the credential issuer
func (b *CredentialBuilder) Issuer(issuerDID string) *CredentialBuilder {
	issuerAddress := signing.FromAddress(issuerDID)
	if issuerAddress != nil {
		b.issuer = credential.AddressKey(issuerAddress)
	}
	return b
}

// ValidFrom sets when the credential becomes valid
func (b *CredentialBuilder) ValidFrom(validFrom time.Time) *CredentialBuilder {
	b.validFrom = validFrom
	return b
}

// SignWith sets the signing key and issuance time
func (b *CredentialBuilder) SignWith(signerDID string, issuedAt time.Time) *CredentialBuilder {
	signerAddress := signing.FromAddress(signerDID)
	if signerAddress != nil {
		b.signer = signerAddress
		b.issuedAt = issuedAt
	}
	return b
}

// Issue creates and issues the credential
func (b *CredentialBuilder) Issue(client *Client) (*credential.VerifiableCredential, error) {
	if client.isClosed() {
		return nil, ErrClientClosed
	}

	// Build the unsigned credential
	credBuilder := credential.NewCredential().
		CredentialType(b.credentialType).
		CredentialSubject(b.subject).
		CredentialSubjectClaims(b.claims).
		Issuer(b.issuer).
		ValidFrom(b.validFrom).
		SignWith(b.signer, b.issuedAt)

	unsignedCredential, err := credBuilder.Finish()
	if err != nil {
		return nil, err
	}

	// Issue the credential
	return client.account.CredentialIssue(unsignedCredential)
}

// WaitForResponse waits for a response to the credential request
func (req *CredentialRequest) WaitForResponse(ctx context.Context) (*CredentialResponse, error) {
	select {
	case response := <-req.completer:
		return response, nil
	case <-ctx.Done():
		// Clean up the stored request
		req.client.loadAndDeleteRequest(req.requestID)
		return nil, ctx.Err()
	}
}

// RequestID returns the unique identifier for this credential request
func (req *CredentialRequest) RequestID() string {
	return req.requestID
}

// From returns the sender's DID
func (resp *CredentialResponse) From() string {
	return resp.from
}

// Status returns the response status
func (resp *CredentialResponse) Status() message.ResponseStatus {
	return resp.status
}

// Presentations returns the credential presentations (for presentation responses)
func (resp *CredentialResponse) Presentations() []*credential.VerifiablePresentation {
	return resp.presentations
}

// Credentials returns the verified credentials (for verification responses)
func (resp *CredentialResponse) Credentials() []*credential.VerifiableCredential {
	return resp.credentials
}

// From returns the sender's DID
func (req *IncomingCredentialRequest) From() string {
	return req.from
}

// RequestID returns the request ID
func (req *IncomingCredentialRequest) RequestID() string {
	return req.requestID
}

// Type returns the requested credential/presentation type
func (req *IncomingCredentialRequest) Type() []string {
	return req.reqType
}

// Details returns the credential details (for presentation requests)
func (req *IncomingCredentialRequest) Details() []*CredentialDetail {
	return req.details
}

// Evidence returns the evidence attached to the request
func (req *IncomingCredentialRequest) Evidence() []*CredentialEvidence {
	return req.evidence
}

// Proof returns the proof presentations attached to the request
func (req *IncomingCredentialRequest) Proof() []*credential.VerifiablePresentation {
	return req.proof
}

// Expires returns when the request expires
func (req *IncomingCredentialRequest) Expires() time.Time {
	return req.expires
}

// IsVerificationRequest returns true if this is a verification request
func (req *IncomingCredentialRequest) IsVerificationRequest() bool {
	return req.isVerification
}

// RespondWithPresentations responds to a presentation request with presentations
func (req *IncomingCredentialRequest) RespondWithPresentations(presentations []*credential.VerifiablePresentation) error {
	if req.isVerification {
		return ErrInvalidResponse
	}

	// Parse the sender DID to get the signing key
	peerAddress := signing.FromAddress(req.from)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	// Build the response
	builder := message.NewCredentialPresentationResponse().
		ResponseTo([]byte(req.requestID)).
		Status(message.ResponseStatusAccepted)

	for _, presentation := range presentations {
		builder.VerifiablePresentation(presentation)
	}

	content, err := builder.Finish()
	if err != nil {
		return err
	}

	return req.client.sendMessage(peerAddress, *content)
}

// RespondWithCredentials responds to a verification request with credentials
func (req *IncomingCredentialRequest) RespondWithCredentials(credentials []*credential.VerifiableCredential) error {
	if !req.isVerification {
		return ErrInvalidResponse
	}

	// Parse the sender DID to get the signing key
	peerAddress := signing.FromAddress(req.from)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	// Build the response
	builder := message.NewCredentialVerificationResponse().
		ResponseTo([]byte(req.requestID)).
		Status(message.ResponseStatusAccepted)

	for _, credential := range credentials {
		builder.VerifiableCredential(credential)
	}

	content, err := builder.Finish()
	if err != nil {
		return err
	}

	return req.client.sendMessage(peerAddress, *content)
}

// Reject rejects the credential request
func (req *IncomingCredentialRequest) Reject() error {
	// Parse the sender DID to get the signing key
	peerAddress := signing.FromAddress(req.from)
	if peerAddress == nil {
		return ErrInvalidPeerDID
	}

	var content *message.Content
	var err error

	if req.isVerification {
		content, err = message.NewCredentialVerificationResponse().
			ResponseTo([]byte(req.requestID)).
			Status(message.ResponseStatusForbidden).
			Finish()
	} else {
		content, err = message.NewCredentialPresentationResponse().
			ResponseTo([]byte(req.requestID)).
			Status(message.ResponseStatusForbidden).
			Finish()
	}

	if err != nil {
		return err
	}

	return req.client.sendMessage(peerAddress, *content)
}

// CredentialAsset methods

// ID returns the asset's unique identifier
func (a *CredentialAsset) ID() []byte {
	return a.object.Id()
}

// Hash returns the hash of the unencrypted data
func (a *CredentialAsset) Hash() []byte {
	return a.object.Hash()
}

// Object returns the underlying object
func (a *CredentialAsset) Object() *object.Object {
	return a.object
}

// Internal methods for handling events

func (c *Credentials) onConnect() {
	// Connection established - no specific action needed
}

func (c *Credentials) onDisconnect(err error) {
	// Connection lost - no specific action needed for credentials
}

func (c *Credentials) onWelcome(from *signing.PublicKey, groupAddress *signing.PublicKey) {
	// New connection established - no specific action needed
}

func (c *Credentials) onKeyPackage(from *signing.PublicKey) {
	// Key package received - no specific action needed
}

func (c *Credentials) onIntroduction(from *signing.PublicKey, tokenCount int) {
	// Introduction received - connection is now ready for credential exchange
}

func (c *Credentials) onCredentialPresentationRequest(msg *event.Message) {
	// Decode the credential presentation request
	presentationRequest, err := message.DecodeCredentialPresentationRequest(msg.Content())
	if err != nil {
		return
	}

	// Convert details
	details := make([]*CredentialDetail, len(presentationRequest.Details()))
	for i, detail := range presentationRequest.Details() {
		params := make([]*CredentialParameter, len(detail.Parameters()))
		for j, param := range detail.Parameters() {
			params[j] = &CredentialParameter{
				Operator: param.Operator(),
				Field:    param.ClaimField(),
				Value:    param.ClaimValue(),
			}
		}
		details[i] = &CredentialDetail{
			CredentialType: detail.CredentialType(),
			Parameters:     params,
		}
	}

	// Create incoming request object
	incomingRequest := &IncomingCredentialRequest{
		from:           msg.FromAddress().String(),
		requestID:      hex.EncodeToString(msg.ID()),
		reqType:        presentationRequest.Type(),
		details:        details,
		evidence:       []*CredentialEvidence{}, // Presentation requests don't have evidence
		proof:          presentationRequest.Proof(),
		expires:        presentationRequest.Expires(),
		client:         c.client,
		isVerification: false,
	}

	// Notify handlers
	c.mu.RLock()
	handlers := make([]func(*IncomingCredentialRequest), len(c.onPresentationRequestHandlers))
	copy(handlers, c.onPresentationRequestHandlers)
	c.mu.RUnlock()

	for _, handler := range handlers {
		go handler(incomingRequest) // Run handlers in goroutines to avoid blocking
	}
}

func (c *Credentials) onCredentialVerificationRequest(msg *event.Message) {
	// Decode the credential verification request
	verificationRequest, err := message.DecodeCredentialVerificationRequest(msg.Content())
	if err != nil {
		return
	}

	// Convert evidence
	evidence := make([]*CredentialEvidence, len(verificationRequest.Evidence()))
	for i, ev := range verificationRequest.Evidence() {
		evidence[i] = &CredentialEvidence{
			Type:   ev.EvidenceType(),
			Object: ev.Object(),
		}
	}

	// Create incoming request object
	incomingRequest := &IncomingCredentialRequest{
		from:           msg.FromAddress().String(),
		requestID:      hex.EncodeToString(msg.ID()),
		reqType:        verificationRequest.Type(),
		details:        []*CredentialDetail{}, // Verification requests don't have details
		evidence:       evidence,
		proof:          verificationRequest.Proof(),
		expires:        verificationRequest.Expires(),
		client:         c.client,
		isVerification: true,
	}

	// Notify handlers
	c.mu.RLock()
	handlers := make([]func(*IncomingCredentialRequest), len(c.onVerificationRequestHandlers))
	copy(handlers, c.onVerificationRequestHandlers)
	c.mu.RUnlock()

	for _, handler := range handlers {
		go handler(incomingRequest) // Run handlers in goroutines to avoid blocking
	}
}

func (c *Credentials) onCredentialPresentationResponse(msg *event.Message) {
	// Decode the credential presentation response
	presentationResponse, err := message.DecodeCredentialPresentationResponse(msg.Content())
	if err != nil {
		return
	}

	requestID := hex.EncodeToString(presentationResponse.ResponseTo())

	// Find the waiting request
	completerInterface, ok := c.client.loadAndDeleteRequest(requestID)
	if !ok {
		return
	}

	completer, ok := completerInterface.(chan *CredentialResponse)
	if !ok {
		return
	}

	// Create response object
	response := &CredentialResponse{
		from:          msg.FromAddress().String(),
		status:        presentationResponse.Status(),
		presentations: presentationResponse.Presentations(),
		credentials:   []*credential.VerifiableCredential{},
	}

	// Send to waiting request
	select {
	case completer <- response:
	default:
		// Channel full or closed - ignore
	}

	// Notify subscription handlers
	c.mu.RLock()
	handlers := make([]func(*CredentialResponse), len(c.onPresentationResponseHandlers))
	copy(handlers, c.onPresentationResponseHandlers)
	c.mu.RUnlock()

	for _, handler := range handlers {
		go handler(response) // Run handlers in goroutines to avoid blocking
	}
}

func (c *Credentials) onCredentialVerificationResponse(msg *event.Message) {
	// Decode the credential verification response
	verificationResponse, err := message.DecodeCredentialVerificationResponse(msg.Content())
	if err != nil {
		return
	}

	requestID := hex.EncodeToString(verificationResponse.ResponseTo())

	// Find the waiting request
	completerInterface, ok := c.client.loadAndDeleteRequest(requestID)
	if !ok {
		return
	}

	completer, ok := completerInterface.(chan *CredentialResponse)
	if !ok {
		return
	}

	// Create response object
	response := &CredentialResponse{
		from:          msg.FromAddress().String(),
		status:        verificationResponse.Status(),
		presentations: []*credential.VerifiablePresentation{},
		credentials:   verificationResponse.Credentials(),
	}

	// Send to waiting request
	select {
	case completer <- response:
	default:
		// Channel full or closed - ignore
	}

	// Notify subscription handlers
	c.mu.RLock()
	handlers := make([]func(*CredentialResponse), len(c.onVerificationResponseHandlers))
	copy(handlers, c.onVerificationResponseHandlers)
	c.mu.RUnlock()

	for _, handler := range handlers {
		go handler(response) // Run handlers in goroutines to avoid blocking
	}
}

func (c *Credentials) close() {
	// Clean up any pending requests
	// Note: We could iterate through stored requests and close channels,
	// but the current sync.Map doesn't provide an easy way to do this.
	// For now, pending requests will timeout naturally.
}

// CreatePresentation creates a verifiable presentation from credentials
func (c *Credentials) CreatePresentation(presentationType []string, credentials []*credential.VerifiableCredential) (*credential.VerifiablePresentation, error) {
	if c.client.isClosed() {
		return nil, ErrClientClosed
	}

	// Build the unsigned presentation
	builder := credential.NewPresentation().
		PresentationType(presentationType).
		Holder(credential.AddressKey(c.client.inboxAddress))

	// Add all credentials
	for _, cred := range credentials {
		builder.CredentialAdd(cred)
	}

	unsignedPresentation, err := builder.Finish()
	if err != nil {
		return nil, err
	}

	// Issue the presentation
	return c.client.account.PresentationIssue(unsignedPresentation)
}
