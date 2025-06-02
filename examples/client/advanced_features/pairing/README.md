# Advanced Pairing Example

A comprehensive demonstration of the Self SDK's account pairing system for multi-device synchronization.

## 🚀 Quick Start

```bash
# Run the pairing demo
go run main.go
```

The demo automatically showcases all Self SDK pairing features in a structured, educational format!

## 📊 Complexity Rating

**5/10** (Intermediate) - Perfect for learning multi-device synchronization

- 🟡 **Pairing concepts**: QR codes, device verification, account synchronization
- 🟡 **Security**: Cryptographic device verification and secure pairing
- 🟡 **Multi-device**: Cross-device state management and synchronization
- 🟢 **Event handling**: Pairing requests, responses, and status management

## 🎯 What This Example Demonstrates

### Core Pairing Features
- ✅ **QR Code Pairing** - Easy device connection via QR codes
- ✅ **Pairing Codes** - Alternative pairing method with codes
- ✅ **Device Verification** - Secure cryptographic device authentication
- ✅ **Event Handling** - Pairing requests, responses, and status updates
- ✅ **Multi-Device Sync** - Account synchronization across devices

### Educational Learning Path
1. **Pairing Setup** - Configure pairing handlers and generate codes
2. **QR Code Generation** - Create QR codes for device pairing
3. **Pairing Management** - Handle pairing requests and responses
4. **Device Synchronization** - Manage multi-device account state

## 🏃‍♂️ How to Run

### Single Command Demo
```bash
go run main.go
```

The demo runs automatically and demonstrates:
- Setting up pairing event handlers
- Generating QR codes and pairing codes
- Managing pairing requests and responses
- Demonstrating device verification and synchronization

### What Happens Automatically
1. **Client Creation**: Pairing-focused Self client setup
2. **Handler Setup**: Configure pairing event handlers
3. **Code Generation**: Generate QR codes and pairing codes
4. **Event Management**: Handle pairing requests and responses
5. **Status Monitoring**: Track pairing status and device connections

## 📋 What You'll See

```
🔗 Advanced Pairing Demo
========================
This demo showcases Self SDK pairing capabilities.

🔧 Setting up pairing client...
✅ Pairing client created successfully
🆔 Client DID: did:self:pairing123...

🔹 Setting up Pairing Handlers
==============================
🔗 Configuring pairing event handlers...
   ✅ Pairing request handler configured
   ✅ Pairing response handler configured
   ✅ Pairing status handler configured

🔹 Pairing Code Generation
===========================
🔑 Generating pairing information...
   ✅ Pairing Code: ABC123DEF456
   📱 Unpaired: true
   ⏰ Expires: 2024-01-15 15:30:00
   🔄 Current pairing status: false

🔹 QR Code Generation
=====================
📱 Generating QR code for device pairing...
   ✅ QR code generated successfully
   🔗 QR Data: eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9...
   📱 Scan with another device to pair

🔹 Pairing Event Management
===========================
🔗 Monitoring pairing events...
   📨 Waiting for pairing requests...
   ✅ Event handlers active and monitoring
   🔄 Pairing status: ready for connections

🔹 Device Verification
======================
🔐 Demonstrating secure device verification...
   ✅ Cryptographic verification enabled
   🔑 Secure key exchange protocols active
   🛡️ Device authentication ready
```

## 🔍 Key Code Sections

| Function | Purpose |
|----------|---------|
| `main()` | Step-by-step pairing demo orchestration |
| `createPairingClient()` | Pairing-focused Self SDK client setup |
| `setupPairingHandlers()` | Configure pairing event handlers |
| `demonstratePairingCodes()` | Generate and manage pairing codes |
| `demonstrateQRGeneration()` | Create QR codes for device pairing |
| `demonstratePairingEvents()` | Handle pairing requests and responses |
| `demonstrateDeviceVerification()` | Secure device authentication |

## 🎓 Educational Notes

### Pairing Concepts
- **QR Code Pairing**: Visual pairing method using QR codes
- **Pairing Codes**: Alternative text-based pairing method
- **Device Verification**: Cryptographic authentication of devices
- **Multi-Device Sync**: Account state synchronization across devices

### Security Features
- **Cryptographic Verification**: Secure device authentication
- **Key Exchange**: Safe sharing of encryption keys
- **Identity Verification**: Confirm device ownership
- **Secure Channels**: Encrypted communication during pairing

### Benefits
- **Seamless Experience**: Easy device connection and setup
- **Security**: Cryptographic verification prevents unauthorized access
- **Synchronization**: Consistent account state across all devices
- **Flexibility**: Multiple pairing methods (QR, codes, manual)

## 🔧 Customization Ideas

Try modifying the code to:
- Implement custom pairing workflows for your application
- Add device management and removal capabilities
- Create pairing approval workflows
- Implement device naming and identification
- Add pairing analytics and monitoring

## 🚀 Next Steps

After understanding this example, continue with:

| Next Example | Complexity | Description |
|-------------|------------|-------------|
| **Production Patterns** | 6/10 | Real-world storage patterns |
| **Integration** | 7/10 | Multi-component workflows |

## 🛠️ Prerequisites

- Go 1.19 or later
- Self SDK dependencies (handled by go.mod)
- Basic understanding of cryptographic concepts
- Completion of simple_chat and group_chat examples
- Understanding of storage and notifications examples (recommended)

## 💡 Troubleshooting

**Pairing Issues:**
- Ensure pairing handlers are configured before generating codes
- Check QR code generation and scanning functionality
- Verify cryptographic operations are working correctly

**Connection Issues:**
- Confirm network connectivity between devices
- Check that pairing codes haven't expired
- Verify device compatibility and SDK versions

**Security Issues:**
- Ensure cryptographic libraries are properly installed
- Check device authentication and verification
- Verify secure key exchange is functioning

**Build Issues:**
- Run `go mod tidy` to ensure dependencies
- Check Go version with `go version`
- Verify you're in the correct directory

## 🎯 Key Differences from Other Examples

| Feature | Simple Chat | Group Chat | Storage | Notifications | **Pairing** |
|---------|-------------|------------|---------|---------------|-------------|
| **Focus** | Basic messaging | Group coordination | Data persistence | User engagement | **Multi-device sync** |
| **Security** | Basic | Group security | Encryption | Delivery security | **Device verification** |
| **Complexity** | 4/10 | 5/10 | 5/10 | 4/10 | **5/10** |
| **Device Management** | Single | Single | Single | Single | **Multi-device** |
| **Synchronization** | None | Group state | Local only | None | **Cross-device** |

## 🔗 Pairing Architecture

### Pairing Flow
```
Device A → Generate QR/Code → Device B Scans → Verification → Paired
    ↓              ↓              ↓              ↓         ↓
Event Handler ← Status Update ← Key Exchange ← Auth ← Sync State
```

### Pairing Methods
- **QR Code**: Visual scanning for easy pairing
- **Pairing Code**: Text-based code entry
- **Manual**: Direct device identification
- **Invitation**: Peer-to-peer pairing requests

### Security Layers
- **Cryptographic Verification**: Device identity confirmation
- **Key Exchange**: Secure sharing of encryption keys
- **Authentication**: Multi-factor device verification
- **Authorization**: Permission-based access control

This example provides the foundation for multi-device experiences in Self SDK applications! 🔗 
