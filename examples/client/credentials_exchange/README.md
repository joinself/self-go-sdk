# Credential Exchange Examples

This directory contains a comprehensive tutorial series for learning credential exchange using the Self SDK. The examples are designed with educational progression, starting from basic concepts and building up to advanced real-world scenarios.

## 🚀 Quick Start

| 🎯 Goal | 📁 Example | 🏃‍♂️ Command |
|---------|------------|-------------|
| **Just want to see it work?** | Root | `go run main.go` |
| **Learn the basics** | Basic | `cd basic && go run main.go` |
| **Handle multiple credentials** | Multi | `cd multi && go run main.go` |
| **Advanced filtering** | Advanced | `cd advanced && go run main.go` |
| **Real peer connections** | Discovery | `cd discovery && go run main.go` |

## 📚 Educational Progression

### 🎯 Learning Path

Complete the examples in this order for the best learning experience:

| Example | Complexity | Description | Time |
|---------|------------|-------------|------|
| **`basic/main.go`** | 🟢 **9/10** (Very Simple) | Foundation concepts | 5-10 min |
| **`multi/main.go`** | 🟡 **7/10** (Intermediate) | Multiple credential types | 10-15 min |
| **`advanced/main.go`** | 🟠 **5/10** (Advanced) | Complex parameters & verification | 15-20 min |
| **`discovery/main.go`** | 🟠 **6/10** (Expert) | QR code discovery integration | 20-25 min |

### 🔄 Current main.go

The current `main.go` contains the simplified basic exchange example (**9/10 simplicity**). For the full educational progression, use the individual files above.

### 📊 Complexity Guide

- 🟢 **8-10/10** (Very Simple) - Perfect for beginners, minimal concepts
- 🟡 **6-7/10** (Intermediate) - Some complexity, builds on basics  
- 🟠 **4-5/10** (Advanced) - Complex concepts, production patterns
- 🔴 **1-3/10** (Expert) - Very complex, requires deep understanding

## 📖 Example Descriptions

### 1️⃣ Basic Exchange (`basic/main.go`)
**Start here if you're new to credential exchange**

- **What you'll learn:**
  - How credential exchange works between two parties
  - Basic request/response patterns
  - Simple credential creation and sharing
  
- **Key concepts:**
  - Client setup (issuer and holder)
  - Simple credential creation
  - Handler configuration
  - Basic exchange workflow

- **Complexity:** 🟢 **9/10** (Very Simple)
- **Time to complete:** 5-10 minutes

### 2️⃣ Multi-Credential Exchange (`multi/main.go`)
**Prerequisites: Complete basic/main.go first**

- **What you'll learn:**
  - How to handle multiple credential types
  - Creating credentials with different claim structures
  - Multi-credential request patterns
  - Processing complex responses

- **Key concepts:**
  - Email, profile, and education credentials
  - Multi-credential requests
  - Complex response processing
  - Batch credential handling

- **Complexity:** 🟡 **7/10** (Intermediate)
- **Time to complete:** 10-15 minutes

### 3️⃣ Advanced Exchange (`advanced/main.go`)
**Prerequisites: Complete basic/main.go and multi/main.go**

- **What you'll learn:**
  - Complex credential filtering with operators
  - Difference between presentation and verification requests
  - Advanced response processing patterns
  - Production-ready error handling

- **Key concepts:**
  - Complex parameter filtering (>, <, >=, <=, ==, !=)
  - Verification vs presentation requests
  - Nested claim filtering
  - Advanced response processing

- **Complexity:** 🟠 **5/10** (Advanced)
- **Time to complete:** 15-20 minutes

### 4️⃣ Discovery Exchange (`discovery/main.go`)
**Prerequisites: Complete all previous examples**

- **What you'll learn:**
  - QR code-based peer discovery
  - Real-time credential exchange with live peers
  - Production discovery workflows
  - Integration of discovery with credential exchange

- **Key concepts:**
  - QR code generation and scanning
  - Peer-to-peer connections
  - Live credential exchange
  - Discovery-based workflows

- **Complexity:** 🟠 **6/10** (Expert)
- **Time to complete:** 20-25 minutes

## 🚀 Getting Started

### Prerequisites

1. Go 1.19 or later
2. Self SDK dependencies (automatically handled by go.mod)

### Running the Examples

Each example is a standalone Go program. Run them in order:

```bash
# 1. Basic Exchange (start here)
cd basic && go run main.go

# 2. Multi-Credential Exchange
cd ../multi && go run main.go

# 3. Advanced Exchange
cd ../advanced && go run main.go

# 4. Discovery Exchange
cd ../discovery && go run main.go

# Or run the simplified version in the root directory
cd .. && go run main.go
```

### 🔧 Build Requirements

Each subdirectory is a standalone Go module with its own `go.mod` file. The examples use a local replace directive to reference the Self SDK, so they should build and run without any additional setup.

### What Each Example Does

All examples create two clients (issuer and holder) and demonstrate different aspects of credential exchange:

- **Issuer**: Creates and signs credentials, requests credentials from holders
- **Holder**: Receives credentials, responds to credential requests
- **Exchange**: The interactive process of requesting and sharing credentials

## 🎓 Learning Outcomes

After completing all examples, you'll understand:

### Core Concepts
- ✅ How credential exchange works between parties
- ✅ The difference between issuers, holders, and verifiers
- ✅ Request/response patterns in credential exchange
- ✅ Handler configuration for interactive workflows

### Technical Skills
- ✅ Self SDK client setup and configuration
- ✅ Credential creation using the builder pattern
- ✅ Complex parameter filtering and querying
- ✅ Multi-credential request handling
- ✅ QR code-based peer discovery
- ✅ Real-time credential exchange

### Production Readiness
- ✅ Error handling and timeout management
- ✅ Security considerations for credential exchange
- ✅ Scalable exchange patterns
- ✅ Integration with discovery mechanisms

## 🔧 Key SDK Components Covered

### Client Management
- `client.NewClient()` - Client initialization
- `client.Config` - Configuration options
- Storage and environment setup

### Credential Operations
- `NewCredentialBuilder()` - Credential creation
- `Type()`, `Subject()`, `Issuer()` - Basic credential properties
- `Claim()`, `Claims()` - Adding credential data
- `SignWith()`, `Issue()` - Credential finalization

### Exchange Operations
- `RequestPresentationWithTimeout()` - Request credentials
- `RequestVerificationWithTimeout()` - Verify credentials
- `OnPresentationRequest()` - Handle incoming requests
- `OnPresentationResponse()` - Handle responses

### Discovery Operations
- `Discovery().GenerateQR()` - Create discovery QR codes
- `WaitForResponse()` - Wait for peer connections
- Peer-to-peer credential exchange

## 🛠️ Customization

Each example can be customized for your specific use case:

### Credential Types
- Modify credential types to match your domain
- Add custom claims relevant to your application
- Implement domain-specific validation logic

### Exchange Patterns
- Customize request/response handlers
- Implement business logic for credential validation
- Add audit logging and compliance tracking

### Discovery Integration
- Integrate with existing identity systems
- Customize QR code generation and display
- Implement custom peer discovery mechanisms

## 📚 Next Steps

After completing these examples:

1. **🎓 Explore credential issuance** in `../credential_issuance/` - Learn how to create the credentials you're exchanging
2. **📖 Review the Self SDK documentation** for advanced features
3. **🏗️ Build your own credential exchange application**
4. **🔗 Integrate with existing identity management systems**
5. **🎯 Combine issuance and exchange** - Create end-to-end credential workflows

### 🔄 Credential Lifecycle

Understanding the complete credential lifecycle:
- **Issuance** (`../credential_issuance/`) - Creating and signing credentials
- **Exchange** (this tutorial) - Requesting and sharing credentials  
- **Verification** - Validating credential authenticity and claims
- **Revocation** - Managing credential lifecycle and updates

## 🤝 Contributing

If you find ways to improve these examples or have suggestions for additional educational content, please contribute back to the Self SDK project.

## 📖 Additional Resources

- [Self SDK Documentation](https://docs.joinself.com)
- [W3C Verifiable Credentials](https://w3.org/TR/vc-data-model/)
- [W3C Verifiable Presentations](https://w3.org/TR/vc-data-model/#presentations)
- [DIDComm Messaging](https://identity.foundation/didcomm-messaging/)
- [Decentralized Identity Foundation](https://identity.foundation)

---

**Happy learning! 🎉**

Start with `basic_exchange.go` and work your way through the progression. Each example builds upon the previous ones, providing a comprehensive understanding of credential exchange with the Self SDK. 
