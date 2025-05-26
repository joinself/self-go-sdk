# Self SDK Client Package

The `client` package provides a high-level, easy-to-use interface for the Self SDK. It simplifies common operations like discovery, chat messaging, and connection management while maintaining the flexibility of the underlying SDK.

## Features

- **Simplified Configuration**: Minimal setup with sensible defaults
- **Automatic Connection Management**: Handles discovery, welcome, and key package events automatically
- **Event-Driven Architecture**: Simple callback registration for messages, discovery, and credentials
- **Request-Response Tracking**: Automatic handling of request-response matching
- **Subscription Support**: Listen for discovery responses from multiple QR codes
- **Credential Exchange**: Easy credential presentation and verification requests

## Quick Start

### Basic Client Setup

```go
package main

import (
    "log"
    "github.com/joinself/self-go-sdk/client"
)

func main() {
    // Create a new Self client
    selfClient, err := client.New(client.Config{
        StorageKey:  make([]byte, 32), // Use a secure key in production
        StoragePath: "./my_app_storage",
        Environment: client.Sandbox,   // or client.Production
        LogLevel:    client.LogInfo,
    })
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    defer selfClient.Close()

    // Your DID (Decentralized Identifier)
    fmt.Printf("My DID: %s\n", selfClient.DID())
}
```

### Discovery

#### Generate QR Code and Wait for Response

```go
// Generate a QR code for discovery
qr, err := selfClient.Discovery().GenerateQR()
if err != nil {
    log.Fatal(err)
}

// Display QR code
qrCode, _ := qr.Unicode()
fmt.Println(qrCode)

// Wait for someone to scan
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

peer, err := qr.WaitForResponse(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Connected to: %s\n", peer.DID())
```

#### Subscribe to Discovery Responses

```go
// Set up discovery response handler
selfClient.Discovery().OnResponse(func(peer *client.Peer) {
    fmt.Printf("New peer discovered: %s\n", peer.DID())
    // Handle new connection...
})

// Generate multiple QR codes - all responses will trigger the handler
for i := 0; i < 3; i++ {
    qr, _ := selfClient.Discovery().GenerateQR()
    qrCode, _ := qr.Unicode()
    fmt.Printf("QR Code #%d:\n%s\n", i+1, qrCode)
}
```

### Chat Messaging

#### Send Messages

```go
// Send a simple message
err := selfClient.Chat().Send(peerDID, "Hello, world!")
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

#### Receive Messages

```go
// Set up message handler
selfClient.Chat().OnMessage(func(msg client.ChatMessage) {
    fmt.Printf("Message from %s: %s\n", msg.From(), msg.Text())
    
    // Reply to the message
    selfClient.Chat().Reply(msg, "Thanks for your message!")
})
```

#### Message Information

```go
selfClient.Chat().OnMessage(func(msg client.ChatMessage) {
    fmt.Printf("From: %s\n", msg.From())
    fmt.Printf("Text: %s\n", msg.Text())
    fmt.Printf("Message ID: %s\n", msg.ID())
    
    if msg.ReferencedID() != "" {
        fmt.Printf("Replying to: %s\n", msg.ReferencedID())
    }
    
    fmt.Printf("Attachments: %d\n", len(msg.Attachments()))
})
```

### Credential Exchange

#### Request Credential Presentations

```go
// Define what credentials you want
details := []*client.CredentialDetail{
    {
        CredentialType: credential.CredentialTypeEmail,
        Parameters: []*client.CredentialParameter{
            {
                Operator: message.OperatorNotEquals,
                Field:    "emailAddress",
                Value:    "",
            },
        },
    },
}

// Request presentations with timeout
req, err := selfClient.Credentials().RequestPresentationWithTimeout(
    peerDID, 
    details, 
    30*time.Second,
)
if err != nil {
    log.Printf("Failed to request presentations: %v", err)
    return
}

// Wait for response
ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
defer cancel()

resp, err := req.WaitForResponse(ctx)
if err != nil {
    log.Printf("Request failed: %v", err)
    return
}

fmt.Printf("Received %d presentations\n", len(resp.Presentations()))
```

#### Create Custom Credentials

```go
// Create a custom credential with flexible claims
customCredential, err := selfClient.Credentials().NewCredentialBuilder().
    Type([]string{"VerifiableCredential", "AgreementCredential"}).
    Subject(holderDID).
    Issuer(selfClient.DID()).
    Claim("agreementType", "Service Agreement").
    Claim("version", "1.0").
    Claims(map[string]interface{}{
        "parties": []map[string]interface{}{
            {"type": "issuer", "id": issuerDID},
            {"type": "holder", "id": holderDID},
        },
        "effectiveDate": time.Now().Format("2006-01-02"),
        "permissions": []string{"read", "write", "execute"},
    }).
    ValidFrom(time.Now()).
    SignWith(selfClient.DID(), time.Now()).
    Issue(selfClient)

if err != nil {
    log.Printf("Failed to create credential: %v", err)
    return
}

fmt.Printf("Created credential: %v\n", customCredential.CredentialType())
```

#### Create and Attach Evidence/Assets

```go
// Create an asset (file attachment)
pdfData := []byte("PDF document content")
asset, err := selfClient.Credentials().CreateAsset("agreement.pdf", "application/pdf", pdfData)
if err != nil {
    log.Printf("Failed to create asset: %v", err)
    return
}

// Create evidence for credential requests
evidence := []*client.CredentialEvidence{
    {
        Type:   "terms",
        Object: asset.Object(),
    },
}

// Request verification with evidence
req, err := selfClient.Credentials().RequestVerificationWithEvidence(
    peerDID,
    []string{"VerifiableCredential", "AgreementCredential"},
    evidence,
    nil, // proof presentations
)
```

#### Create Presentations

```go
// Create a presentation from credentials
presentation, err := selfClient.Credentials().CreatePresentation(
    []string{"VerifiablePresentation", "CustomPresentation"},
    []*credential.VerifiableCredential{myCredential},
)
if err != nil {
    log.Printf("Failed to create presentation: %v", err)
    return
}

fmt.Printf("Created presentation: %v\n", presentation.PresentationType())
```

#### Request Credential Verification

```go
// Request verification of a specific credential type
req, err := selfClient.Credentials().RequestVerificationWithTimeout(
    peerDID,
    credential.CredentialTypeEmail,
    30*time.Second,
)
if err != nil {
    log.Printf("Failed to request verification: %v", err)
    return
}

// Wait for response
ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
defer cancel()

resp, err := req.WaitForResponse(ctx)
if err != nil {
    log.Printf("Request failed: %v", err)
    return
}

fmt.Printf("Verification status: %s\n", resp.Status())
```

#### Handle Incoming Credential Requests

```go
// Handle presentation requests
selfClient.Credentials().OnPresentationRequest(func(req *client.IncomingCredentialRequest) {
    fmt.Printf("Presentation request from: %s\n", req.From())
    fmt.Printf("Requesting: %v\n", req.Type())
    
    // Check what's being requested
    for _, detail := range req.Details() {
        fmt.Printf("Credential Type: %v\n", detail.CredentialType)
        for _, param := range detail.Parameters {
            fmt.Printf("  %s %s %s\n", param.Field, param.Operator, param.Value)
        }
    }
    
    // Accept or reject the request
    if shouldAccept(req) {
        presentations := buildPresentations(req)
        err := req.RespondWithPresentations(presentations)
        if err != nil {
            log.Printf("Failed to respond with presentations: %v", err)
        }
    } else {
        err := req.Reject()
        if err != nil {
            log.Printf("Failed to reject request: %v", err)
        }
    }
})

// Handle verification requests
selfClient.Credentials().OnVerificationRequest(func(req *client.IncomingCredentialRequest) {
    fmt.Printf("Verification request from: %s\n", req.From())
    fmt.Printf("Type: %v\n", req.Type())
    
    // Provide credentials for verification
    if hasCredentials(req.Type()) {
        credentials := getCredentials(req.Type())
        err := req.RespondWithCredentials(credentials)
        if err != nil {
            log.Printf("Failed to respond with credentials: %v", err)
        }
    } else {
        err := req.Reject()
        if err != nil {
            log.Printf("Failed to reject request: %v", err)
        }
    }
})
```

#### Handle Credential Responses

```go
// Handle presentation responses
selfClient.Credentials().OnPresentationResponse(func(resp *client.CredentialResponse) {
    fmt.Printf("Presentation response from: %s\n", resp.From())
    fmt.Printf("Status: %s\n", resp.Status())
    
    for i, presentation := range resp.Presentations() {
        fmt.Printf("Presentation %d:\n", i+1)
        fmt.Printf("  Type: %v\n", presentation.PresentationType())
        fmt.Printf("  Credentials: %d\n", len(presentation.Credentials()))
    }
})

// Handle verification responses
selfClient.Credentials().OnVerificationResponse(func(resp *client.CredentialResponse) {
    fmt.Printf("Verification response from: %s\n", resp.From())
    fmt.Printf("Status: %s\n", resp.Status())
    
    for i, cred := range resp.Credentials() {
        fmt.Printf("Credential %d:\n", i+1)
        fmt.Printf("  Type: %v\n", cred.CredentialType())
        fmt.Printf("  Subject: %s\n", cred.CredentialSubject())
    }
})
```

## Configuration

### Environment

```go
client.Config{
    Environment: client.Sandbox,    // For development/testing
    // or
    Environment: client.Production, // For production use
}
```

### Log Levels

```go
client.Config{
    LogLevel: client.LogError,  // Errors only
    LogLevel: client.LogWarn,   // Warnings and errors
    LogLevel: client.LogInfo,   // Info, warnings, and errors
    LogLevel: client.LogDebug,  // Debug and above
    LogLevel: client.LogTrace,  // All logs
}
```

### Storage

```go
client.Config{
    StorageKey:  []byte("your-32-byte-encryption-key"), // Required
    StoragePath: "./app_data",                          // Required
}
```

## Examples

See the `examples/client/` directory for complete working examples:

- `simple_chat/` - Basic chat application
- `discovery_subscription/` - Discovery with subscription handling
- `credentials_exchange/` - Credential presentation and verification requests
- `credential_issuance/` - Enhanced credential features with issuance, evidence, and custom schemas
- `group_chat/` - Multi-party group chat with admin controls and member management
- `advanced_features/` - Notifications, storage abstractions, and account pairing functionality

### Group Chat

#### Create and Manage Groups

```go
// Create a new group chat
group, err := selfClient.GroupChats().CreateGroup("Dev Team", "Daily standup discussions")
if err != nil {
    log.Printf("Failed to create group: %v", err)
    return
}

fmt.Printf("Created group: %s (ID: %s)\n", group.Name(), group.ID())
fmt.Printf("Members: %d\n", group.MemberCount())

// Update group details (admin/moderator only)
err = group.UpdateName("Dev Team - Sprint 1")
if err != nil {
    log.Printf("Failed to update name: %v", err)
}

err = group.UpdateDescription("Sprint 1 planning and daily standups")
if err != nil {
    log.Printf("Failed to update description: %v", err)
}
```

#### Invite Members and Send Messages

```go
// Invite members to the group (admin/moderator only)
err = selfClient.GroupChats().InviteToGroup(group.ID(), peerDID, "Welcome to our team!")
if err != nil {
    log.Printf("Failed to invite member: %v", err)
}

// Send messages to all group members
err = selfClient.GroupChats().SendToGroup(group.ID(), "Hello everyone!")
if err != nil {
    log.Printf("Failed to send group message: %v", err)
}

// Reply to a specific group message
err = selfClient.GroupChats().ReplyToGroupMessage(originalMessage, "Thanks for the update!")
if err != nil {
    log.Printf("Failed to reply: %v", err)
}
```

#### Handle Group Events

```go
// Set up group message handler
selfClient.GroupChats().OnGroupMessage(func(msg client.GroupChatMessage) {
    fmt.Printf("Group message in '%s' from %s: %s\n", 
        msg.GroupName(), msg.From(), msg.Text())
})

// Handle group invitations
selfClient.GroupChats().OnGroupInvite(func(invitation *client.GroupChatInvitation) {
    fmt.Printf("Invited to group: %s by %s\n", 
        invitation.GroupName, invitation.InviterDID)
    
    // Accept or decline the invitation
    if shouldJoin(invitation) {
        err := invitation.Accept()
        if err != nil {
            log.Printf("Failed to accept invitation: %v", err)
        }
    } else {
        err := invitation.Decline()
        if err != nil {
            log.Printf("Failed to decline invitation: %v", err)
        }
    }
})

// Handle member events
selfClient.GroupChats().OnMemberJoined(func(groupID string, member *client.GroupMember) {
    fmt.Printf("Member %s joined group %s with role %s\n", 
        member.DID, groupID, member.Role)
})

selfClient.GroupChats().OnMemberLeft(func(groupID string, memberDID string) {
    fmt.Printf("Member %s left group %s\n", memberDID, groupID)
})
```

#### Group Management

```go
// List all groups
groups := selfClient.GroupChats().ListGroups()
for _, group := range groups {
    fmt.Printf("Group: %s (%d members)\n", group.Name(), group.MemberCount())
    
    // Get group members
    members := group.Members()
    for _, member := range members {
        fmt.Printf("  - %s (%s, online: %v)\n", 
            member.DID, member.Role, member.IsOnline)
    }
}

// Get specific group
group, exists := selfClient.GroupChats().GetGroup(groupID)
if exists {
    fmt.Printf("Found group: %s\n", group.Name())
}

// Leave a group
err = selfClient.GroupChats().LeaveGroup(groupID)
if err != nil {
    log.Printf("Failed to leave group: %v", err)
}
```

## API Reference

### Client

- `New(config Config) (*Client, error)` - Create a new client
- `DID() string` - Get the client's DID
- `Discovery() *Discovery` - Access discovery functionality
- `Chat() *Chat` - Access chat functionality
- `Credentials() *Credentials` - Access credential exchange functionality
- `GroupChats() *GroupChats` - Access group chat functionality
- `Notifications() *Notifications` - Access push notification functionality
- `Storage() *Storage` - Access key-value storage functionality
- `Pairing() *Pairing` - Access account pairing and linking functionality
- `Close() error` - Close the client and cleanup resources

### Discovery

- `GenerateQR() (*DiscoveryQR, error)` - Generate QR code with default timeout
- `GenerateQRWithTimeout(timeout time.Duration) (*DiscoveryQR, error)` - Generate QR code with custom timeout
- `OnResponse(handler func(*Peer))` - Subscribe to discovery responses

### DiscoveryQR

- `Unicode() (string, error)` - Get QR code as Unicode text
- `SVG() (string, error)` - Get QR code as SVG
- `WaitForResponse(ctx context.Context) (*Peer, error)` - Wait for response
- `RequestID() string` - Get unique request identifier

### Chat

- `Send(peerDID string, message string) error` - Send a message
- `Reply(originalMessage ChatMessage, replyText string) error` - Reply to a message
- `OnMessage(handler func(ChatMessage))` - Subscribe to incoming messages

### ChatMessage

- `From() string` - Sender's DID
- `Text() string` - Message text
- `ID() string` - Message ID
- `ReferencedID() string` - ID of referenced message (for replies)
- `Attachments() []ChatAttachment` - Message attachments

### Credentials

- `RequestPresentation(peerDID string, details []*CredentialDetail) (*CredentialRequest, error)` - Request credential presentations
- `RequestPresentationWithTimeout(peerDID string, details []*CredentialDetail, timeout time.Duration) (*CredentialRequest, error)` - Request presentations with timeout
- `RequestPresentationWithEvidence(peerDID string, details []*CredentialDetail, evidence []*CredentialEvidence, proof []*credential.VerifiablePresentation) (*CredentialRequest, error)` - Request presentations with evidence
- `RequestVerification(peerDID string, credentialType []string) (*CredentialRequest, error)` - Request credential verification
- `RequestVerificationWithTimeout(peerDID string, credentialType []string, timeout time.Duration) (*CredentialRequest, error)` - Request verification with timeout
- `RequestVerificationWithEvidence(peerDID string, credentialType []string, evidence []*CredentialEvidence, proof []*credential.VerifiablePresentation) (*CredentialRequest, error)` - Request verification with evidence
- `NewCredentialBuilder() *CredentialBuilder` - Create a new credential builder for custom credentials
- `CreateAsset(name, mimeType string, data []byte) (*CredentialAsset, error)` - Create and upload an asset/file
- `DownloadAsset(asset *CredentialAsset) error` - Download and decrypt an asset
- `CreatePresentation(presentationType []string, credentials []*credential.VerifiableCredential) (*credential.VerifiablePresentation, error)` - Create a verifiable presentation
- `OnPresentationRequest(handler func(*IncomingCredentialRequest))` - Subscribe to presentation requests
- `OnVerificationRequest(handler func(*IncomingCredentialRequest))` - Subscribe to verification requests
- `OnPresentationResponse(handler func(*CredentialResponse))` - Subscribe to presentation responses
- `OnVerificationResponse(handler func(*CredentialResponse))` - Subscribe to verification responses

### CredentialBuilder

- `Type(credentialType []string) *CredentialBuilder` - Set credential type
- `Subject(subjectDID string) *CredentialBuilder` - Set credential subject
- `Issuer(issuerDID string) *CredentialBuilder` - Set credential issuer
- `Claim(key string, value interface{}) *CredentialBuilder` - Add a single claim
- `Claims(claims map[string]interface{}) *CredentialBuilder` - Add multiple claims
- `ValidFrom(validFrom time.Time) *CredentialBuilder` - Set validity start time
- `SignWith(signerDID string, issuedAt time.Time) *CredentialBuilder` - Set signing key and issuance time
- `Issue(client *Client) (*credential.VerifiableCredential, error)` - Create and issue the credential

### CredentialAsset

- `ID() []byte` - Get asset's unique identifier
- `Hash() []byte` - Get hash of unencrypted data
- `Object() *object.Object` - Get underlying object
- `Name` - Asset filename
- `MimeType` - Asset MIME type
- `Data` - Asset data

### CredentialRequest

- `WaitForResponse(ctx context.Context) (*CredentialResponse, error)` - Wait for response
- `RequestID() string` - Get unique request identifier

### CredentialResponse

- `From() string` - Sender's DID
- `Status() message.ResponseStatus` - Response status
- `Presentations() []*credential.VerifiablePresentation` - Credential presentations (for presentation responses)
- `Credentials() []*credential.VerifiableCredential` - Verified credentials (for verification responses)

### IncomingCredentialRequest

- `From() string` - Sender's DID
- `RequestID() string` - Request ID
- `Type() []string` - Requested credential/presentation type
- `Details() []*CredentialDetail` - Credential details (for presentation requests)
- `Evidence() []*CredentialEvidence` - Evidence attached to the request
- `Proof() []*credential.VerifiablePresentation` - Proof presentations attached to the request
- `Expires() time.Time` - Request expiration time
- `IsVerificationRequest() bool` - Check if this is a verification request
- `RespondWithPresentations(presentations []*credential.VerifiablePresentation) error` - Respond with presentations
- `RespondWithCredentials(credentials []*credential.VerifiableCredential) error` - Respond with credentials
- `Reject() error` - Reject the request

### GroupChats

- `CreateGroup(name, description string) (*GroupChat, error)` - Create a new group chat
- `InviteToGroup(groupID, peerDID, message string) error` - Invite a peer to join a group
- `JoinGroup(invitation *GroupChatInvitation) error` - Join a group via invitation
- `SendToGroup(groupID, messageText string) error` - Send a message to all group members
- `ReplyToGroupMessage(originalMessage GroupChatMessage, replyText string) error` - Reply to a group message
- `GetGroup(groupID string) (*GroupChat, bool)` - Get a group by ID
- `ListGroups() []*GroupChat` - List all groups
- `LeaveGroup(groupID string) error` - Leave a group
- `OnGroupMessage(handler func(GroupChatMessage))` - Subscribe to group messages
- `OnGroupInvite(handler func(*GroupChatInvitation))` - Subscribe to group invitations
- `OnMemberJoined(handler func(groupID string, member *GroupMember))` - Subscribe to member join events
- `OnMemberLeft(handler func(groupID string, memberDID string))` - Subscribe to member leave events
- `OnGroupCreated(handler func(*GroupChat))` - Subscribe to group creation events
- `OnGroupUpdated(handler func(*GroupChat))` - Subscribe to group update events

### GroupChat

- `ID() string` - Group ID
- `Name() string` - Group name
- `Description() string` - Group description
- `Members() []*GroupMember` - All group members
- `Admin() string` - Group admin DID
- `Created() time.Time` - When the group was created
- `MemberCount() int` - Number of members
- `UpdateName(newName string) error` - Update group name (admin/moderator only)
- `UpdateDescription(newDescription string) error` - Update group description (admin/moderator only)

### GroupChatMessage

- `From() string` - Sender's DID
- `Text() string` - Message text
- `ID() string` - Message ID
- `ReferencedID() string` - ID of referenced message (for replies)
- `GroupID() string` - Group ID
- `GroupName() string` - Group name
- `Timestamp() time.Time` - When the message was sent
- `Attachments() []ChatAttachment` - Message attachments

### GroupChatInvitation

- `GroupID string` - Group ID
- `GroupName string` - Group name
- `InviterDID string` - Inviter's DID
- `InviterName string` - Inviter's name
- `Message string` - Invitation message
- `ExpiresAt time.Time` - Invitation expiration
- `Accept() error` - Accept the invitation
- `Decline() error` - Decline the invitation

### GroupMember

- `DID string` - Member's DID
- `Name string` - Member's name
- `Role GroupRole` - Member's role (admin, moderator, member)
- `JoinedAt time.Time` - When the member joined
- `LastSeen time.Time` - Last activity time
- `IsOnline bool` - Online status

### Notifications

- `SendNotification(peerDID string, summary *NotificationSummary) error` - Send a push notification
- `SendChatNotification(peerDID, messageText string) error` - Send a chat message notification
- `SendGroupChatNotification(peerDID, groupName, messageText string) error` - Send a group chat notification
- `SendCredentialNotification(peerDID, credentialType, action string) error` - Send a credential-related notification
- `SendGroupInviteNotification(peerDID, groupName, inviterName string) error` - Send a group invitation notification
- `SendCustomNotification(peerDID, title, body, messageType string) error` - Send a custom notification
- `CreateSummaryFromContent(content *message.Content) (*NotificationSummary, error)` - Create notification summary from message content
- `OnNotificationSent(handler func(peerDID string, summary *NotificationSummary))` - Subscribe to notification sent events

### NotificationSummary

- `Title string` - Notification title
- `Body string` - Notification body text
- `MessageType string` - Type of message (chat, credential, etc.)
- `FromDID string` - Sender's DID
- `MessageID string` - Associated message ID

### Storage

- `Store(key string, value []byte) error` - Store a value
- `StoreWithExpiry(key string, value []byte, expires time.Time) error` - Store a value with expiry
- `StoreString(key, value string) error` - Store a string value
- `StoreStringWithExpiry(key, value string, expires time.Time) error` - Store a string with expiry
- `StoreJSON(key string, value interface{}) error` - Store a JSON-serializable value
- `StoreJSONWithExpiry(key string, value interface{}, expires time.Time) error` - Store JSON with expiry
- `Lookup(key string) ([]byte, error)` - Retrieve a value
- `LookupString(key string) (string, error)` - Retrieve a string value
- `LookupJSON(key string, target interface{}) error` - Retrieve and unmarshal JSON
- `Exists(key string) bool` - Check if a key exists
- `Delete(key string) error` - Remove a value
- `StoreTemporary(key string, value []byte, duration time.Duration) error` - Store with relative expiry
- `StoreTemporaryString(key, value string, duration time.Duration) error` - Store string with relative expiry
- `StoreTemporaryJSON(key string, value interface{}, duration time.Duration) error` - Store JSON with relative expiry
- `Namespace(namespace string) *StorageNamespace` - Get namespaced storage
- `Cache(prefix string) *Cache` - Get cache interface

### StorageNamespace

- `Store(key string, value []byte) error` - Store in namespace
- `StoreWithExpiry(key string, value []byte, expires time.Time) error` - Store with expiry in namespace
- `StoreString(key, value string) error` - Store string in namespace
- `StoreJSON(key string, value interface{}) error` - Store JSON in namespace
- `StoreJSONWithExpiry(key string, value interface{}, expires time.Time) error` - Store JSON with expiry in namespace
- `Lookup(key string) ([]byte, error)` - Retrieve from namespace
- `LookupString(key string) (string, error)` - Retrieve string from namespace
- `LookupJSON(key string, target interface{}) error` - Retrieve JSON from namespace
- `Exists(key string) bool` - Check existence in namespace
- `Delete(key string) error` - Remove from namespace
- `StoreTemporary(key string, value []byte, duration time.Duration) error` - Store with relative expiry in namespace

### Cache

- `Set(key string, value []byte) error` - Cache a value (1 hour default TTL)
- `SetWithTTL(key string, value []byte, ttl time.Duration) error` - Cache with custom TTL
- `SetString(key, value string) error` - Cache a string
- `SetJSON(key string, value interface{}) error` - Cache JSON data
- `Get(key string) ([]byte, error)` - Retrieve from cache
- `GetString(key string) (string, error)` - Retrieve string from cache
- `GetJSON(key string, target interface{}) error` - Retrieve JSON from cache
- `Has(key string) bool` - Check if cached
- `Delete(key string) error` - Remove from cache

### Pairing

- `GetPairingCode() (*PairingCode, error)` - Get SDK pairing code for account linking
- `RequestPairing(peerDID string, address *signing.PublicKey, roles identity.Role) (*PairingRequest, error)` - Send pairing request
- `RequestPairingWithTimeout(peerDID string, address *signing.PublicKey, roles identity.Role, timeout time.Duration) (*PairingRequest, error)` - Send pairing request with timeout
- `GeneratePairingQR() (string, error)` - Generate QR code for pairing
- `IsPaired() (bool, error)` - Check if account is paired
- `OnPairingRequest(handler func(*IncomingPairingRequest))` - Subscribe to pairing requests
- `OnPairingResponse(handler func(*PairingResponse))` - Subscribe to pairing responses

### PairingCode

- `Code string` - The pairing code
- `Unpaired bool` - Whether the account is unpaired
- `ExpiresAt time.Time` - When the code expires

### PairingRequest

- `WaitForResponse(ctx context.Context) (*PairingResponse, error)` - Wait for pairing response
- `RequestID() string` - Get request ID

### IncomingPairingRequest

- `From() string` - Sender's DID
- `RequestID() string` - Request ID
- `Address() *signing.PublicKey` - Address to be paired
- `Roles() identity.Role` - Requested roles
- `Expires() time.Time` - Request expiration
- `RespondWithOperation(operation *identity.Operation) error` - Accept with identity operation
- `RespondWithOperationAndAssets(operation *identity.Operation, assets []*object.Object) error` - Accept with operation and assets
- `Reject() error` - Reject the pairing request

### PairingResponse

- `From() string` - Sender's DID
- `Status() message.ResponseStatus` - Response status
- `Operation() *identity.Operation` - Identity operation (if accepted)
- `Assets() []*object.Object` - Supporting assets

## Advanced Features Examples

### Storage Patterns

```go
// Basic storage
storage := selfClient.Storage()
err := storage.StoreString("user:name", "Alice")
name, err := storage.LookupString("user:name")

// JSON storage
userData := map[string]interface{}{
    "name": "Alice",
    "age":  30,
}
err = storage.StoreJSON("user:profile", userData)

var profile map[string]interface{}
err = storage.LookupJSON("user:profile", &profile)

// Namespaced storage
userStorage := storage.Namespace("user:12345")
err = userStorage.StoreString("preferences", "theme=dark")

// Temporary storage
err = storage.StoreTemporaryString("session:token", "abc123", 5*time.Minute)

// Caching
cache := storage.Cache("api")
err = cache.SetString("response:users", `[{"id":1,"name":"Alice"}]`)
if cache.Has("response:users") {
    data, _ := cache.GetString("response:users")
}
```

### Notification Integration

```go
notifications := selfClient.Notifications()

// Register notification handler
notifications.OnNotificationSent(func(peerDID string, summary *client.NotificationSummary) {
    fmt.Printf("Sent %s to %s: %s\n", summary.MessageType, peerDID, summary.Title)
})

// Send different notification types
err := notifications.SendChatNotification(peerDID, "Hello!")
err = notifications.SendCredentialNotification(peerDID, "identity", "request")
err = notifications.SendGroupInviteNotification(peerDID, "Dev Team", "Alice")
err = notifications.SendCustomNotification(peerDID, "Alert", "System update", "system")
```

### Account Pairing

```go
pairing := selfClient.Pairing()

// Get pairing code for account linking
pairingCode, err := pairing.GetPairingCode()
fmt.Printf("Pairing Code: %s\n", pairingCode.Code)

// Generate QR code
qrCode, err := pairing.GeneratePairingQR()
fmt.Println(qrCode)

// Handle pairing requests
pairing.OnPairingRequest(func(request *client.IncomingPairingRequest) {
    fmt.Printf("Pairing request from: %s\n", request.From())
    
    // Accept or reject
    if shouldAccept(request) {
        operation := createIdentityOperation()
        err := request.RespondWithOperation(operation)
    } else {
        err := request.Reject()
    }
})

// Send pairing request
signingKey, _ := signing.NewKey()
request, err := pairing.RequestPairing(targetDID, signingKey.PublicKey(), identity.RoleOwner)
if err == nil {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    response, err := request.WaitForResponse(ctx)
}
```

## Migration from Low-Level SDK

The client package is designed to work alongside the existing low-level SDK. You can:

1. Start with the client package for common operations
2. Access the underlying account via the client for advanced features
3. Gradually migrate existing code to use the simplified API

## Error Handling

The client package defines common errors:

- `ErrStorageKeyRequired` - Storage key not provided
- `ErrStoragePathRequired` - Storage path not provided
- `ErrClientClosed` - Operation on closed client
- `ErrInvalidPeerDID` - Invalid peer DID format
- `ErrDiscoveryTimeout` - Discovery request timed out

## Thread Safety

The client package is designed to be thread-safe. You can safely call methods from multiple goroutines.

## Best Practices

1. **Storage Key Security**: Use a cryptographically secure random key for `StorageKey`
2. **Error Handling**: Always check errors, especially for network operations
3. **Context Usage**: Use contexts with timeouts for discovery operations
4. **Resource Cleanup**: Always call `Close()` when done with the client
5. **Handler Goroutines**: Message and discovery handlers run in separate goroutines - avoid blocking operations 
