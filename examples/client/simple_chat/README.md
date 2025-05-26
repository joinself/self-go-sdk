# Simple Chat Example

A straightforward demonstration of the Self SDK's core chat capabilities.

## What This Example Demonstrates

🎯 **Core Chat Features:**
- Real-time bidirectional messaging
- End-to-end encryption (automatic)
- QR code-based peer discovery
- Message echo functionality
- Simple command handling
- Multi-peer support

## How to Run

1. **Start the first instance:**
   ```bash
   go run main.go
   ```

2. **Start a second instance in another terminal:**
   ```bash
   go run main.go
   ```

3. **Connect the peers:**
   - The first instance will display a QR code
   - Scan the QR code with the second instance (copy/paste the QR text)
   - Wait for the connection to establish

4. **Start chatting:**
   - Send messages from either instance
   - Messages will be echoed back automatically
   - Try the commands: `/help` and `/quit`

## Key Code Sections

- **Client Setup**: Lines 50-60 - Shows minimal Self SDK configuration
- **Chat Handlers**: Lines 100-150 - Demonstrates message processing
- **Peer Discovery**: Lines 180-220 - QR code-based connection
- **Message Sending**: Throughout - Shows how to send chat messages

## What You'll See

```
🚀 Self SDK Simple Chat Example
===============================
📚 This demo showcases the core chat capabilities:
   • Secure peer-to-peer messaging
   • QR code-based connection
   • Real-time message exchange
   • End-to-end encryption

🆔 Your DID: did:self:example123...
   This is your unique decentralized identity

🔍 PEER DISCOVERY PROCESS
=========================
📱 SCAN THIS QR CODE with another Self client to connect:
[QR CODE DISPLAYED HERE]

✅ Chat connection established with: did:self:peer456...
🔐 All messages are automatically end-to-end encrypted

💬 CHAT IS NOW ACTIVE!
======================
📨 This demo will echo back any messages you send
```

## Educational Notes

- **No Central Server**: Peers connect directly to each other
- **Automatic Encryption**: All messages are encrypted end-to-end
- **Decentralized Identity**: Each client has a unique DID
- **Real-time**: Messages are delivered instantly when peers are online
- **Multi-peer**: The same QR code can be used by multiple peers

## Next Steps

After understanding this example, explore:
- `../group_chat/` - Multi-participant chat rooms
- `../advanced_features/` - Additional SDK capabilities
- `../credentials_exchange/` - Identity verification features 
