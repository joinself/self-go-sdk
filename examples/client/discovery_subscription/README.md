# Discovery Subscription Example

This example demonstrates the powerful subscription-based peer discovery capabilities of the Self SDK. Learn how to create a discovery service that can handle multiple simultaneous peer connections through QR code scanning, with real-time notifications for each discovery event.

## 🚀 Quick Start

| 🎯 Goal | 🏃‍♂️ Command | ⏱️ Time |
|---------|-------------|---------|
| **See discovery subscription in action** | `go run main.go` | 2-3 min |
| **Test with multiple peers** | Run + scan QR codes | 5-10 min |

## 📚 What You'll Learn

### 🎯 Core Concepts

- **Discovery Subscription Pattern** - How to set up handlers that receive real-time notifications
- **Multiple QR Code Generation** - Creating several discovery endpoints simultaneously  
- **Peer Discovery Workflow** - Understanding the complete discovery process
- **Real-time Event Handling** - Processing discovery events as they happen

### 🔄 Discovery Subscription Flow

```
1. Client generates multiple QR codes with different timeouts
2. Each QR code represents a unique discovery endpoint
3. Other Self clients scan the QR codes
4. Original client receives real-time notifications for each discovery
5. Multiple peers can discover simultaneously without interference
```

### 📱 Interactive Demo Features

- **Three QR codes** with different timeout periods (15min, 30min, 60min)
- **Real-time notifications** when peers discover you
- **Detailed discovery information** including peer DID and timestamp
- **Educational explanations** of what's happening at each step

## 🎓 Learning Outcomes

After completing this example, you'll understand:

### Discovery Subscription Concepts
- ✅ How subscription-based discovery works
- ✅ The difference between one-time and subscription discovery
- ✅ Real-time event handling patterns
- ✅ Multiple simultaneous discovery endpoints

### Technical Implementation
- ✅ Setting up discovery event handlers with `OnResponse()`
- ✅ Generating QR codes with `GenerateQRWithTimeout()`
- ✅ Managing multiple discovery sessions
- ✅ Processing peer discovery events

### Real-world Applications
- ✅ Building discovery services for multiple users
- ✅ Creating persistent discovery endpoints
- ✅ Handling concurrent peer connections
- ✅ Implementing discovery-based workflows

## 🚀 Getting Started

### Prerequisites

1. **Go 1.19 or later**
2. **Self SDK dependencies** (automatically handled by go.mod)
3. **Multiple Self clients** for testing (mobile apps, other SDK instances)

### Running the Example

```bash
# Run the discovery subscription demo
go run main.go

# The program will:
# 1. Generate three QR codes with different timeouts
# 2. Display them in the terminal
# 3. Listen for discovery events
# 4. Show real-time notifications when peers connect
```

### 📱 Testing the Subscription

1. **Run the program** - It will generate three QR codes
2. **Use Self mobile apps** or other SDK clients to scan the codes
3. **Try scanning different codes** with different devices
4. **Watch the real-time notifications** appear in the terminal
5. **Test simultaneous discoveries** by having multiple people scan at once

## 🔧 Key SDK Components Covered

### Discovery Management
- `client.Discovery()` - Access discovery functionality
- `OnResponse(func(*client.Peer))` - Set up subscription handler
- `GenerateQRWithTimeout(duration)` - Create discovery QR codes

### Event Handling
- **Subscription Pattern** - Handler called for every discovery
- **Peer Information** - Access to discovered peer details
- **Real-time Processing** - Immediate notification of events

### QR Code Management
- **Multiple QR Codes** - Generate several discovery endpoints
- **Timeout Configuration** - Different validity periods
- **Request ID Tracking** - Unique identifiers for each QR code

## 🎯 Use Cases

### Real-world Applications

1. **Event Registration**
   - Generate QR codes for event check-in
   - Real-time attendee discovery and registration
   - Multiple entry points with different access levels

2. **Networking Events**
   - Business card exchange via QR codes
   - Real-time contact discovery
   - Multiple networking sessions simultaneously

3. **Service Discovery**
   - IoT device discovery and pairing
   - Service endpoint registration
   - Dynamic peer-to-peer connections

4. **Educational Platforms**
   - Student-teacher connections
   - Classroom participation tracking
   - Real-time attendance monitoring

## 🔄 Discovery vs. Other Patterns

### Discovery Subscription vs. One-time Discovery

| Feature | Subscription | One-time |
|---------|-------------|----------|
| **Handler Calls** | Multiple (for each peer) | Single |
| **QR Code Reuse** | Yes, until timeout | No, single use |
| **Concurrent Peers** | Unlimited | One |
| **Use Case** | Services, events | Direct peer connection |

### When to Use Discovery Subscription

- ✅ **Multiple peer connections** expected
- ✅ **Service-like behavior** needed
- ✅ **Real-time notifications** required
- ✅ **Persistent discovery** endpoints
- ✅ **Event-driven architecture**

### When to Use One-time Discovery

- ✅ **Direct peer-to-peer** connection
- ✅ **Single connection** expected
- ✅ **Simple pairing** scenarios
- ✅ **One-off interactions**

## 🛠️ Customization

### Extending the Example

```go
// Custom discovery handler with business logic
client.Discovery().OnResponse(func(peer *client.Peer) {
    // Store peer information
    storePeerInDatabase(peer)
    
    // Send welcome message
    sendWelcomeMessage(peer)
    
    // Initiate credential exchange
    requestCredentials(peer)
    
    // Log discovery event
    logDiscoveryEvent(peer)
})

// Generate QR codes with custom timeouts
timeouts := []time.Duration{
    5 * time.Minute,   // Quick connections
    1 * time.Hour,     // Standard connections  
    24 * time.Hour,    // Long-term availability
}

for _, timeout := range timeouts {
    qr, err := client.Discovery().GenerateQRWithTimeout(timeout)
    // Handle QR code...
}
```

### Integration Patterns

1. **Database Integration**
   - Store discovered peers in database
   - Track discovery events and timestamps
   - Implement peer relationship management

2. **Notification Systems**
   - Send push notifications on discovery
   - Email alerts for new connections
   - Real-time dashboard updates

3. **Workflow Automation**
   - Trigger credential exchange on discovery
   - Initiate chat sessions automatically
   - Start business process workflows

## 🔧 Production Considerations

### Security
- **QR Code Expiration** - Use appropriate timeouts
- **Peer Validation** - Verify discovered peers
- **Rate Limiting** - Prevent discovery spam
- **Access Control** - Implement discovery permissions

### Scalability
- **Handler Performance** - Keep discovery handlers fast
- **Concurrent Connections** - Handle multiple simultaneous discoveries
- **Resource Management** - Clean up expired QR codes
- **Monitoring** - Track discovery metrics

### Error Handling
- **QR Generation Failures** - Graceful degradation
- **Network Issues** - Retry mechanisms
- **Timeout Management** - Handle expired discoveries
- **Peer Validation** - Verify discovery authenticity

## 📚 Next Steps

After mastering discovery subscription:

1. **🔄 Explore Credential Exchange** (`../credentials_exchange/`) - Use discovered peers for credential exchange
2. **💬 Try Simple Chat** (`../simple_chat/`) - Build messaging with discovered peers  
3. **🏗️ Build Discovery Services** - Create your own discovery-based applications
4. **🔗 Integrate with Existing Systems** - Add discovery to current applications

### 🎯 Advanced Discovery Patterns

- **Discovery with Credential Exchange** - Combine discovery and credential sharing
- **Multi-hop Discovery** - Discover peers through intermediaries
- **Discovery Networks** - Build interconnected discovery services
- **Discovery Analytics** - Track and analyze discovery patterns

## 🤝 Contributing

Found ways to improve this example? Have ideas for additional discovery patterns? Contributions are welcome!

## 📖 Additional Resources

- [Self SDK Documentation](https://docs.joinself.com)
- [QR Code Standards](https://www.qrcode.com/en/)
- [Peer-to-Peer Networking](https://en.wikipedia.org/wiki/Peer-to-peer)
- [Event-Driven Architecture](https://en.wikipedia.org/wiki/Event-driven_architecture)
- [Real-time Systems](https://en.wikipedia.org/wiki/Real-time_computing)

---

**Ready to discover? 🔍**

Run `go run main.go` and start exploring the power of subscription-based peer discovery with the Self SDK!

### 🎉 Pro Tips

- **Use multiple devices** to test simultaneous discoveries
- **Try different QR codes** to see timeout behavior
- **Watch the real-time notifications** to understand the subscription pattern
- **Experiment with the handler** to add custom business logic 
