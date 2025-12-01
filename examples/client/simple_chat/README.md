# Simple Chat Example

A straightforward demonstration of the Self SDK's core chat capabilities, designed for educational purposes.

## 🚀 Quick Start

```bash
# Start the chat server
go run main.go
```

Then scan the QR code with the Self developer app to connect and start chatting!

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
- ✅ **Mobile app integration** - Connect via Self developer app
- ✅ **Decentralized identity** - No central servers required

### Educational Learning Path
1. **Server Setup** - Initialize Self SDK chat server
2. **Message Handlers** - Configure incoming message processing
3. **Mobile Connection** - Connect securely via QR codes
4. **Chat Demonstration** - Exchange encrypted messages with developer app

## 🏃‍♂️ How to Run

### Step 1: Start the Chat Server
```bash
go run main.go
```

You'll see:
- Your unique DID (decentralized identifier)
- A QR code for peer connection
- Status messages showing the setup process

### Step 2: Connect with Developer App
- Open the Self developer app on your phone
- Use the QR code scanner to scan the displayed QR code
- Wait for the secure connection to establish

### Step 3: Watch the Demo
- The demo automatically sends several test messages to your developer app
- Try sending messages from your developer app to see real-time chat with the server
- Each message gets a smart response based on content

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

📱 SCAN THIS QR CODE with the Self developer app:
   • Open the Self developer app on your phone
   • Use the built-in QR code scanner
   • Point your camera at the QR code below

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
   1. Created a Self chat server for messaging
   2. Set up handlers to receive and process messages
   3. Generated QR code for developer app connection
   4. Exchanged encrypted messages between server and developer app
   5. Demonstrated smart response functionality
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
- **QR Code Discovery**: Secure connection establishment method

### Real-time Features
- **Instant delivery** when developer app is connected
- **Automatic encryption** for all messages
- **Cross-platform support** between Go server and developer app

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
| [`../credentials_exchange/`](../credentials_exchange/) | 6/10 | Identity verification and credential sharing |
| [`../group_chat/`](../group_chat/) | 5/10 | Multi-participant chat rooms |
| [`../file_sharing/`](../file_sharing/) | 7/10 | Secure file transfer between peers |

## 🛠️ Prerequisites

- Go 1.19 or later
- Self-go-sdk dependencies (check the [README.md](/README.md))
- Self developer app installed on your phone

## 💡 Troubleshooting

**Connection Issues:**
- Ensure your developer app is updated to the latest version
- Check that QR code is scanned completely and clearly

**Build Issues:**
- Run `go mod tidy` to ensure dependencies
- Check Go version with `go version`
- Verify you're in the correct directory
- Verify self-go-sdk dependencies are accomplished (check the [README.md](/README.md))
