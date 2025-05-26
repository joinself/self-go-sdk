# Self SDK Client Examples

This directory contains comprehensive examples demonstrating various Self SDK client capabilities. Each example is designed to teach specific concepts and can be used as a foundation for building your own applications.

## 🎯 Quick Navigation

| Example | Focus | Complexity | Best For |
|---------|-------|------------|----------|
| **[Credential Issuance](credential_issuance/)** | Creating & signing credentials | 🟢 Beginner-friendly | Learning credential creation |
| **[Credential Exchange](credentials_exchange/)** | Requesting & sharing credentials | 🟢 Beginner-friendly | Learning credential workflows |
| **[Simple Chat](simple_chat/)** | Basic messaging | 🟡 Intermediate | Understanding messaging |
| **[Group Chat](group_chat/)** | Multi-party messaging | 🟠 Advanced | Group communication |
| **[Discovery Subscription](discovery_subscription/)** | Peer discovery | 🟠 Advanced | Connection management |
| **[Advanced Features](advanced_features/)** | Production patterns | 🔴 Expert | Advanced integrations |

## 📚 Advanced Features Subfolders

The **[Advanced Features](advanced_features/)** directory contains specialized examples for production-ready applications:

| Example | Complexity | Description | Time |
|---------|------------|-------------|------|
| **[Storage](advanced_features/storage/)** | 🟠 **5/10** (Advanced) | Data persistence with namespacing & TTL | 15-20 min |
| **[Notifications](advanced_features/notifications/)** | 🟡 **4/10** (Intermediate) | Push notification system | 10-15 min |
| **[Pairing](advanced_features/pairing/)** | 🟠 **5/10** (Advanced) | Multi-device synchronization | 15-20 min |
| **[Production Patterns](advanced_features/production_patterns/)** | 🟠 **6/10** (Advanced) | Real-world storage & session patterns | 20-25 min |
| **[Integration](advanced_features/integration/)** | 🔴 **7/10** (Expert) | Multi-component workflows | 30-45 min |

### 🎯 Advanced Features Learning Path

Complete the advanced examples in this order for the best learning experience:

```
notifications/ → storage/ → pairing/ → production_patterns/ → integration/
  (4/10)        (5/10)     (5/10)      (6/10)                (7/10)
```

Each example builds upon previous concepts and demonstrates increasingly sophisticated patterns.

## 🚀 Getting Started

### For Beginners
Start with these examples to learn the fundamentals:

1. **[Credential Issuance](credential_issuance/)** - Learn how to create and sign verifiable credentials
2. **[Credential Exchange](credentials_exchange/)** - Learn how to request and share credentials

### For Intermediate Users
Once comfortable with basics, explore:

3. **[Simple Chat](simple_chat/)** - Understand secure messaging
4. **[Discovery Subscription](discovery_subscription/)** - Learn peer discovery

### For Advanced Users
For production-ready patterns:

5. **[Group Chat](group_chat/)** - Multi-party communication
6. **[Advanced Features](advanced_features/)** - Complex integrations
   - Start with **[Notifications](advanced_features/notifications/)** (4/10) for user engagement
   - Progress through **[Storage](advanced_features/storage/)** (5/10) for data persistence
   - Master **[Integration](advanced_features/integration/)** (7/10) for complete workflows

## 🎓 Learning Path

### Complete Credential Workflow
For a comprehensive understanding of verifiable credentials:

```
credential_issuance/ → credentials_exchange/
     (Create)              (Share)
```

This combination teaches the complete credential lifecycle from creation to sharing.

### Communication Features
For messaging and discovery capabilities:

```
simple_chat/ → group_chat/ → discovery_subscription/
  (1-to-1)      (Groups)        (Discovery)
```

### Advanced Production Features
For production-ready applications with advanced capabilities:

```
advanced_features/notifications/ → storage/ → pairing/ → production_patterns/ → integration/
        (4/10)                    (5/10)     (5/10)      (6/10)                (7/10)
```

This progression teaches production patterns from basic notifications to complete multi-component integration.

## 🔧 Prerequisites

All examples require:
- Go 1.19 or later
- Self SDK dependencies (handled automatically)

## 📚 Educational Features

Each example directory includes:
- **Progressive complexity** - Examples build from simple to advanced
- **Comprehensive documentation** - Detailed explanations and learning outcomes
- **Standalone operation** - Each example can be run independently
- **Real-world scenarios** - Practical use cases and patterns

## 🤝 Contributing

Found ways to improve these examples? Contributions are welcome! Please follow the existing educational patterns and documentation standards.

## 📖 Additional Resources

- [Self SDK Documentation](https://docs.joinself.com)
- [W3C Verifiable Credentials](https://w3.org/TR/vc-data-model/)
- [Decentralized Identity Foundation](https://identity.foundation)

---

**Happy learning! 🎉**

Choose an example that matches your current knowledge level and goals. Each example is designed to be educational and practical. 
