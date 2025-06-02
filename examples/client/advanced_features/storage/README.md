# Advanced Storage Example

A comprehensive demonstration of the Self SDK's advanced storage capabilities.

## 🚀 Quick Start

```bash
# Run the storage demo
go run main.go
```

The demo automatically showcases all advanced Self SDK storage features in a structured, educational format!

## 📊 Complexity Rating

**5/10** (Intermediate) - Perfect for learning advanced storage patterns

- 🟡 **Storage concepts**: Namespacing, TTL, caching strategies
- 🟡 **Data organization**: Logical separation and management
- 🟡 **Performance optimization**: Caching and efficient retrieval
- 🟢 **Security**: Automatic encryption at rest

## 🎯 What This Example Demonstrates

### Core Storage Features
- ✅ **Basic Storage** - String and JSON data persistence
- ✅ **Namespaced Storage** - Organized data with logical separation
- ✅ **TTL Storage** - Automatic data expiry for temporary information
- ✅ **Cache Management** - Performance optimization with intelligent caching
- ✅ **Encryption** - Automatic security for all stored data

### Educational Learning Path
1. **Basic Operations** - Store and retrieve different data types
2. **Namespacing** - Organize data into logical groups
3. **TTL Management** - Handle temporary data with automatic expiry
4. **Cache Optimization** - Implement performance caching strategies

## 🏃‍♂️ How to Run

### Single Command Demo
```bash
go run main.go
```

The demo runs automatically and demonstrates:
- Basic string and JSON storage operations
- Namespaced storage for organized data management
- TTL-based temporary storage with automatic expiry
- Cache management for performance optimization

### What Happens Automatically
1. **Client Creation**: Storage-focused Self client setup
2. **Basic Storage**: String and JSON data operations
3. **Namespacing**: User, app, and session data organization
4. **TTL Storage**: Temporary data with automatic cleanup
5. **Caching**: Performance optimization demonstrations

## 📋 What You'll See

```
📦 Advanced Storage Demo
========================
This demo showcases advanced Self SDK storage capabilities.

🔧 Setting up storage client...
✅ Storage client created successfully
🆔 Client DID: did:self:storage123...

🔹 Basic Storage Operations
===========================
📝 Storing different data types...
   ✅ Stored string: user name
   ✅ Stored JSON: user profile with nested data

📖 Retrieving stored data...
   ✅ Retrieved name: Alice Johnson
   ✅ Retrieved profile: Alice Johnson (developer)
   🎨 Theme: dark, Language: en

🔹 Namespaced Storage
=====================
👤 User namespace - personal data:
   ✅ Stored user preferences
   ✅ Stored user settings

🔧 Application namespace - app configuration:
   ✅ Stored application configuration

🔑 Session namespace - temporary session data:
   ✅ Stored session data

🔹 Temporary Storage with TTL
=============================
⏰ Creating short-lived data:
   ✅ Stored temporary session token (expires in 10 seconds)
   ✅ Stored verification code (expires in 5 minutes)
   ✅ Stored temporary user state (expires in 1 hour)

🔹 Cache Management
===================
🗄️ Setting up API response cache:
   ✅ Cached user list (no expiry)
   ✅ Cached user profile (expires in 30 minutes)
   ✅ Cached search results (expires in 15 minutes)

⚡ Cache hit: Retrieved user list (341 bytes)
⚡ Cache hit: Retrieved user profile (271 bytes)
```

## 🔍 Key Code Sections

| Function | Purpose |
|----------|---------|
| `main()` | Step-by-step storage demo orchestration |
| `createStorageClient()` | Storage-focused Self SDK client setup |
| `demonstrateBasicStorage()` | Fundamental storage operations |
| `demonstrateNamespacedStorage()` | Organized storage with namespaces |
| `demonstrateTemporaryStorage()` | TTL-based temporary storage |
| `demonstrateCacheManagement()` | Performance caching strategies |

## 🎓 Educational Notes

### Storage Concepts
- **Encryption**: All data is automatically encrypted at rest
- **Namespacing**: Organize data into logical groups (user, app, session)
- **TTL (Time To Live)**: Automatic expiry for temporary data
- **Caching**: Performance optimization with intelligent data retrieval

### Storage Patterns
- **User Data**: Personal preferences and settings
- **Application Data**: Configuration and state information
- **Session Data**: Temporary information with automatic cleanup
- **Cache Data**: Frequently accessed data for performance

### Benefits
- **Security**: Automatic encryption without additional setup
- **Organization**: Clear data separation with namespaces
- **Performance**: Intelligent caching for faster access
- **Cleanup**: Automatic expiry prevents storage bloat

## 🔧 Customization Ideas

Try modifying the code to:
- Create custom namespaces for your application needs
- Experiment with different TTL values for various data types
- Implement custom caching strategies
- Add data validation and error handling
- Create backup and recovery patterns

## 🚀 Next Steps

After understanding this example, continue with:

| Next Example | Complexity | Description |
|-------------|------------|-------------|
| **Notifications** | 4/10 | Push notification system |
| **Pairing** | 5/10 | Multi-device synchronization |
| **Production Patterns** | 6/10 | Real-world storage patterns |
| **Integration** | 7/10 | Multi-component workflows |

## 🛠️ Prerequisites

- Go 1.19 or later
- Self SDK dependencies (handled by go.mod)
- Basic understanding of storage concepts
- Completion of simple_chat and group_chat examples

## 💡 Troubleshooting

**Storage Issues:**
- Ensure storage paths are writable
- Check that storage keys are properly configured
- Verify namespace usage is consistent

**Performance Issues:**
- Monitor cache hit rates and adjust TTL values
- Optimize namespace usage for your data patterns
- Consider storage cleanup for expired data

**Build Issues:**
- Run `go mod tidy` to ensure dependencies
- Check Go version with `go version`
- Verify you're in the correct directory

## 🎯 Key Differences from Basic Examples

| Feature | Simple Chat | Group Chat | **Storage** |
|---------|-------------|------------|-------------|
| **Focus** | Basic messaging | Group coordination | **Data persistence** |
| **Storage** | None | Basic | **Advanced patterns** |
| **Complexity** | 4/10 | 5/10 | **5/10** |
| **Data Organization** | None | None | **Namespacing** |
| **Performance** | Basic | Basic | **Caching optimization** |

This example provides the foundation for advanced data management in Self SDK applications! 📦 
