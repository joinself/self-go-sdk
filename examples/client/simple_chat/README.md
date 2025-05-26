# Simple Chat Example

A straightforward demonstration of the Self SDK's core chat capabilities, designed for educational purposes.

## 🚀 Quick Start

```bash
# Terminal 1: Start first chat client
go run main.go

# Terminal 2: Start second chat client (in another terminal)
go run main.go
```

Then scan the QR code from the first terminal with the second terminal to connect and start chatting!

## 📊 Complexity Rating

**4/10** (Simple) - Perfect for beginners learning Self SDK chat capabilities

- 🟢 **Easy to follow**: Clear step-by-step progression
- 🟢 **Simple concepts**: Send, receive, echo messages
- 🟡 **Some async concepts**: Message handlers and peer discovery
- 🟡 **Multiple features**: QR codes, DIDs, smart responses

## 🎯 What This Example Demonstrates

### Core Chat Features
- ✅ **Real-time messaging** - Instant bidirectional communication
- ✅ **End-to-end encryption** - Automatic message encryption
- ✅ **QR code discovery** - Secure peer connection establishment
- ✅ **Smart responses** - Context-aware message handling
- ✅ **Multi-peer support** - Multiple clients can connect
- ✅ **Decentralized identity** - No central servers required

### Educational Learning Path
1. **Client Setup** - Initialize Self SDK for messaging
2. **Message Handlers** - Configure incoming message processing
3. **Peer Discovery** - Connect securely via QR codes
4. **Chat Demonstration** - Send and receive encrypted messages

## 🏃‍♂️ How to Run

### Step 1: Start First Client
```bash
go run main.go
```

You'll see:
- Your unique DID (decentralized identifier)
- A QR code for peer connection
- Status messages showing the setup process

### Step 2: Start Second Client
```bash
# In another terminal window
go run main.go
```

### Step 3: Connect the Peers
- Copy the QR code text from the first terminal
- Paste it when the second terminal prompts for QR scanning
- Wait for the secure connection to establish

### Step 4: Watch the Demo
- The demo automatically sends several test messages
- Each message gets a smart response based on content
- Try sending messages from either terminal to see real-time chat

## 📋 What You'll See

```
💬 Simple Chat Demo
===================
This demo shows basic chat messaging between peers.

🔧 Setting up chat client...
✅ Chat client created successfully
🆔 Your DID: did:self:example123...

📨 Setting up message handlers...
✅ Message handlers configured

🔍 Discovering peer for chat...
🔑 Generating QR code for secure connection...

📱 SCAN THIS QR CODE with another Self client:
   • Run another instance of this program
   • Use the Self mobile app
   • Any Self SDK application

[QR CODE DISPLAYED HERE]

✅ Peer connected: did:self:peer456...
🔐 Secure encrypted channel established

💬 Demonstrating chat messaging...
📤 Sending initial greeting...
✅ Greeting sent: "🎉 Hello! Chat demo started at 15:04:05..."

📤 Sending demo messages...
📤 [1/4] Sending: "Hello there!"
✅ Message sent successfully
📤 [2/4] Sending: "How are you?"
✅ Message sent successfully
...

✅ Basic chat demo completed!

🎓 What happened:
   1. Created a Self client for messaging
   2. Set up handlers to receive and process messages
   3. Used QR code to discover and connect to a peer
   4. Exchanged encrypted messages in real-time
   5. Demonstrated echo functionality
```

## 🔍 Key Code Sections

| Function | Lines | Purpose |
|----------|-------|---------|
| `main()` | 30-60 | Step-by-step demo execution |
| `createClient()` | 65-80 | Self SDK client initialization |
| `setupChatHandlers()` | 85-125 | Message and peer event handling |
| `generateResponse()` | 130-145 | Smart response logic |
| `discoverPeer()` | 150-185 | QR code-based peer discovery |
| `demonstrateChat()` | 190-230 | Automated chat demonstration |

## 🎓 Educational Notes

### Core Concepts
- **Decentralized Identity (DID)**: Each client has a unique identifier
- **End-to-End Encryption**: Messages are automatically encrypted
- **Peer-to-Peer**: No central servers, direct client connections
- **QR Code Discovery**: Secure connection establishment method

### Smart Response System
The demo includes intelligent responses based on message content:
- `"hello"` or `"hi"` → Friendly greeting
- `"how are you"` → Status response
- `"help"` → Available commands
- `"time"` → Current timestamp
- Other messages → Echo with timestamp

### Real-time Features
- **Instant delivery** when peers are online
- **Automatic encryption** for all messages
- **Multi-peer support** using the same QR code
- **Event-driven handlers** for incoming messages

## 🔧 Customization Ideas

Try modifying the code to:
- Add new response patterns in `generateResponse()`
- Change the demo messages in `demonstrateChat()`
- Add message logging or persistence
- Implement custom commands (like `/weather`, `/joke`)
- Add emoji reactions or message formatting

## 🚀 Next Steps

After understanding this example, explore:

| Example | Complexity | Description |
|---------|------------|-------------|
| `../credentials_exchange/` | 6/10 | Identity verification and credential sharing |
| `../group_chat/` | 5/10 | Multi-participant chat rooms |
| `../file_sharing/` | 7/10 | Secure file transfer between peers |

## 🛠️ Prerequisites

- Go 1.19 or later
- Self SDK dependencies (handled by go.mod)
- Two terminal windows for testing

## 💡 Troubleshooting

**Connection Issues:**
- Ensure both terminals are running the same version
- Check that QR code is copied completely
- Wait up to 10 minutes for peer discovery timeout

**Message Not Received:**
- Verify both clients show "Peer connected" status
- Check that message handlers are set up before sending
- Ensure clients remain running to receive messages

**Build Issues:**
- Run `go mod tidy` to ensure dependencies
- Check Go version with `go version`
- Verify you're in the correct directory
