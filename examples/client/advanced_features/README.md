# Advanced Features Example

A comprehensive demonstration of the Self SDK's advanced capabilities for production applications.

## 🚀 Quick Start

```bash
# Run the advanced features demo
go run main.go
```

The demo automatically showcases all advanced Self SDK features in a structured, educational format!

## 📊 Complexity Rating

**7/10** (Advanced) - Perfect for learning production-ready Self SDK patterns

- 🟡 **Advanced concepts**: Storage namespacing, TTL, caching strategies
- 🟠 **Production patterns**: Session management, state persistence
- 🟠 **Multi-component integration**: Storage + notifications + chat
- 🔴 **Complex workflows**: Account pairing, notification systems

## 🎯 What This Example Demonstrates

### Core Advanced Features
- ✅ **Encrypted Storage** - Secure local data persistence with namespacing
- ✅ **Cache Management** - Performance optimization with TTL support
- ✅ **Notification System** - Push notifications for user engagement
- ✅ **Account Pairing** - Multi-device synchronization capabilities
- ✅ **Production Patterns** - Real-world storage and session management
- ✅ **Component Integration** - Coordinated multi-feature workflows

### Educational Learning Path
1. **Advanced Client Setup** - Initialize Self SDK with production configuration
2. **Storage Capabilities** - Explore encrypted storage, namespacing, and TTL
3. **Notification System** - Implement push notification delivery
4. **Account Pairing** - Enable multi-device synchronization
5. **Production Patterns** - Apply real-world storage and caching strategies
6. **Component Integration** - Coordinate multiple SDK features together

## 🏃‍♂️ How to Run

### Single Command Demo
```bash
go run main.go
```

The demo runs automatically and demonstrates:
- Advanced storage operations with encryption and namespacing
- Push notification system with multiple notification types
- Account pairing for multi-device experiences
- Production-ready storage patterns and caching strategies
- Integration between storage, notifications, and chat components

### What Happens Automatically
1. **Client Creation**: Advanced Self client with production configuration
2. **Storage Demo**: Basic, namespaced, temporary, and cached storage
3. **Notifications**: Various notification types and event handling
4. **Account Pairing**: Pairing codes, QR generation, and event handlers
5. **Production Patterns**: Session management and state persistence
6. **Integration**: Coordinated workflow using multiple components

## 📋 What You'll See

```
🚀 Advanced Features Demo
=========================
This demo showcases advanced Self SDK capabilities for production use.

🔧 Setting up advanced Self client...
✅ Advanced client created successfully
🆔 Client DID: did:self:advanced123...

📦 Advanced Storage Capabilities
================================
🔹 Basic Storage Operations
   📝 Storing different data types...
   📖 Retrieving stored data...
   ✅ Retrieved name: Alice Johnson
   ✅ Retrieved profile: Alice Johnson (developer)

🔹 Namespaced Storage
   🗂️ Creating organized storage namespaces...
   ✅ Stored user preferences in namespace
   ✅ Stored application settings in namespace
   ✅ Stored session data in namespace

🔹 Temporary Storage with TTL
   ⏰ Creating temporary storage with expiry...
   ✅ Stored temporary session token (expires in 10 seconds)
   ✅ Stored verification code (expires in 5 minutes)
   ✅ Stored temporary user state (expires in 1 hour)

🔹 Cache Management
   🗄️ Setting up cache management...
   ✅ Cached user list
   ✅ Cached user profile (expires in 30 minutes)
   ✅ Retrieved from cache: 156 users found
   ✅ Cached UI theme setting

✅ Storage capabilities demonstrated

🔔 Notification System
======================
🔹 Setting up notification handlers...
   ✅ Notification handlers configured

🔹 Sending various notification types...
   💬 Sending chat notification...
   🆔 Sending credential notification...
   👥 Sending group invite notification...
   🔔 Sending custom notification...
   ✅ All notification types sent successfully

✅ Notification system demonstrated

🔗 Account Pairing System
=========================
🔹 Getting pairing information...
   🔑 Pairing Code: ABC123DEF456
   📱 Unpaired: true
   ⏰ Expires: 2024-01-15 15:30:00
   📱 QR Code available for mobile pairing
   🔗 QR Data: eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9...
   🔄 Current pairing status: false

🔹 Setting up pairing event handlers...
   ✅ Pairing event handlers configured

✅ Account pairing system demonstrated

🏭 Production Storage Patterns
==============================
🔹 User Session Management
   ✅ User session stored with 24-hour expiry
   ✅ Refresh token stored with 7-day expiry

🔹 Application State Persistence
   ✅ Application configuration persisted
   ✅ User preferences saved

🔹 Performance Optimization
   ✅ User data cached for 1 hour
   ✅ Search results cached for 15 minutes
   ⚡ Cache hit: Retrieved user data (234 bytes)

✅ Production patterns demonstrated

🔄 Component Integration
=======================
🔹 Integrating storage, notifications, and chat...
   ✅ Conversation metadata stored
   ✅ Chat message sent
   ✅ Notification sent
   ✅ Conversation metadata updated
   ✅ Recent conversations cached

🎯 Integration benefits:
   • Persistent conversation history
   • Real-time user notifications
   • Optimized data access with caching
   • Coordinated multi-component workflows

✅ Component integration demonstrated

✅ Advanced features demo completed!

🎓 What happened:
   1. Created Self client with advanced configuration
   2. Explored encrypted storage with namespacing and TTL
   3. Demonstrated push notification system
   4. Showed account pairing for multi-device sync
   5. Implemented production-ready storage patterns
   6. Integrated multiple SDK components

🎯 These features enable:
   • Secure data persistence and caching
   • Real-time user engagement via notifications
   • Seamless multi-device experiences
   • Production-ready application architecture
```

## 🔍 Key Code Sections

| Function | Lines | Purpose |
|----------|-------|---------|
| `main()` | 30-90 | Step-by-step advanced demo orchestration |
| `createAdvancedClient()` | 95-110 | Advanced Self SDK client setup |
| `demonstrateAdvancedStorage()` | 115-135 | Complete storage capabilities overview |
| `demonstrateBasicStorage()` | 140-175 | Fundamental storage operations |
| `demonstrateNamespacedStorage()` | 180-220 | Organized storage with namespaces |
| `demonstrateTemporaryStorage()` | 225-255 | TTL-based temporary storage |
| `demonstrateCacheManagement()` | 260-300 | Performance caching strategies |
| `demonstrateNotificationSystem()` | 305-325 | Push notification capabilities |
| `demonstrateAccountPairing()` | 380-395 | Multi-device synchronization |
| `demonstrateProductionPatterns()` | 450-475 | Real-world storage patterns |
| `demonstrateComponentIntegration()` | 580-650 | Multi-component workflows |

## 🎓 Educational Notes

### Advanced Storage Concepts
- **Namespacing**: Organize data into logical groups (user, app, session)
- **TTL (Time To Live)**: Automatic expiry for temporary data
- **Caching**: Performance optimization with intelligent data retrieval
- **Encryption**: All data is automatically encrypted at rest

### Storage Patterns
- **Session Management**: User sessions with automatic expiry
- **State Persistence**: Application configuration and user preferences
- **Performance Optimization**: Strategic caching for frequently accessed data
- **Data Organization**: Logical separation using namespaces

### Notification System
- **Push Notifications**: Real-time alerts for user engagement
- **Multiple Types**: Chat, credential, group invite, and custom notifications
- **Event Handling**: Callbacks for notification delivery status
- **Integration**: Seamless integration with other SDK components

### Account Pairing
- **Multi-Device Sync**: Synchronize accounts across devices
- **QR Code Pairing**: Easy device connection via QR codes
- **Security**: Secure pairing with cryptographic verification
- **Event Management**: Handle pairing requests and responses

### Production Benefits
- **Scalability**: Efficient data access and caching strategies
- **Performance**: Optimized storage and retrieval patterns
- **Security**: Encrypted storage and secure pairing
- **User Experience**: Real-time notifications and multi-device support

## 🔧 Customization Ideas

Try modifying the code to:
- Implement custom storage namespaces for your application
- Add different notification types for specific use cases
- Create advanced caching strategies for your data patterns
- Implement custom pairing workflows
- Add data synchronization between devices
- Create complex integration workflows
- Implement data backup and recovery patterns

## 🚀 Next Steps

After understanding this example, you're ready for:

| Next Level | Complexity | Description |
|------------|------------|-------------|
| **Production Apps** | 8-9/10 | Build real applications using these patterns |
| **Custom Integration** | 8/10 | Integrate Self SDK into existing applications |
| **Advanced Workflows** | 9/10 | Create complex multi-component workflows |

## 🛠️ Prerequisites

- Go 1.19 or later
- Self SDK dependencies (handled by go.mod)
- Understanding of previous examples (simple_chat, group_chat, credentials_exchange)
- Basic knowledge of storage and caching concepts

## 💡 Troubleshooting

**Storage Issues:**
- Ensure storage paths are writable
- Check that storage keys are properly configured
- Verify namespace usage is consistent

**Notification Problems:**
- Confirm notification handlers are set up before sending
- Check network connectivity for notification delivery
- Verify target DIDs are valid

**Pairing Issues:**
- Ensure pairing handlers are configured
- Check QR code generation and scanning
- Verify cryptographic operations are working

**Performance Issues:**
- Monitor cache hit rates and adjust TTL values
- Optimize namespace usage for your data patterns
- Consider storage cleanup for expired data

**Build Issues:**
- Run `go mod tidy` to ensure dependencies
- Check Go version with `go version`
- Verify you're in the correct directory

## 🎯 Key Differences from Other Examples

| Feature | Simple Chat | Group Chat | Advanced Features |
|---------|-------------|------------|-------------------|
| **Focus** | Basic messaging | Group coordination | Production patterns |
| **Storage** | None | Basic | Advanced with TTL/caching |
| **Notifications** | None | None | Full notification system |
| **Pairing** | QR discovery | Admin invites | Multi-device pairing |
| **Complexity** | 4/10 | 5/10 | **7/10** |
| **Production Ready** | Demo | Demo | **Yes** |

## 🏗️ Architecture Patterns Demonstrated

### Storage Architecture
```
Storage
├── Namespaces (user, app, session)
├── TTL Management (automatic expiry)
├── Cache Layers (performance optimization)
└── Encryption (automatic security)
```

### Integration Architecture
```
Self SDK Components
├── Storage ←→ Chat (conversation history)
├── Storage ←→ Notifications (delivery tracking)
├── Notifications ←→ Chat (message alerts)
└── Pairing ←→ Storage (device synchronization)
```

This example provides the foundation for building production-ready applications with the Self SDK! 🚀 
