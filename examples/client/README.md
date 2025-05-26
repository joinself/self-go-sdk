# Self SDK Client Examples

Welcome to the Self SDK client examples! This directory contains a comprehensive learning journey from basic concepts to production-ready applications. Whether you're new to Self SDK or building advanced features, we've got you covered.

## 🚀 Quick Start

**Just want to see something work right now?**

```bash
# Try the simplest example first
cd simple_chat && go run main.go
```

**Want to jump to a specific feature?**

| I want to... | Go here | Complexity |
|--------------|---------|------------|
| 💬 **Send messages** | [`simple_chat/`](simple_chat/) | 🟢 Beginner |
| 🆔 **Create credentials** | [`credential_issuance/`](credential_issuance/) | 🟢 Beginner |
| 🔄 **Share credentials** | [`credentials_exchange/`](credentials_exchange/) | 🟢 Beginner |
| 🔍 **Find peers** | [`discovery_subscription/`](discovery_subscription/) | 🟢 Beginner |
| 👥 **Build group chat** | [`group_chat/`](group_chat/) | 🟡 Intermediate |
| 💾 **Store data securely** | [`advanced_features/storage/`](advanced_features/storage/) | 🟠 Advanced |
| 🔔 **Send notifications** | [`advanced_features/notifications/`](advanced_features/notifications/) | 🟡 Intermediate |
| 🔗 **Sync across devices** | [`advanced_features/pairing/`](advanced_features/pairing/) | 🟠 Advanced |
| 🏭 **Build production apps** | [`advanced_features/integration/`](advanced_features/integration/) | 🔴 Expert |

## 🎓 Learning Paths

### 🌱 New to Self SDK? Start Here!

**Path 1: Messaging Basics** (30-45 minutes)
```
1. simple_chat/          (🟢 4/10) → Learn secure messaging
2. discovery_subscription/ (🟢 4/10) → Find and connect to peers
3. group_chat/           (🟡 5/10) → Multi-party communication  
```

**Path 2: Credential Fundamentals** (30-45 minutes)
```
1. credential_issuance/   (🟢 Beginner) → Create digital credentials
2. credentials_exchange/  (🟢 Beginner) → Share and verify credentials
```

### 🚀 Ready for Advanced Features?

**Path 3: Production Applications** (60-90 minutes)
```
1. advanced_features/notifications/     (🟡 4/10) → User engagement
2. advanced_features/storage/          (🟠 5/10) → Data persistence
3. advanced_features/pairing/          (🟠 5/10) → Multi-device sync
4. advanced_features/production_patterns/ (🟠 6/10) → Real-world patterns
5. advanced_features/integration/      (🔴 7/10) → Complete workflows
```

### 🎯 Goal-Oriented Learning

**I want to build a chat app:**
```
simple_chat/ → group_chat/ → advanced_features/storage/ → advanced_features/notifications/
```

**I want to work with credentials:**
```
credential_issuance/ → credentials_exchange/ → advanced_features/storage/
```

**I want production-ready patterns:**
```
advanced_features/notifications/ → storage/ → production_patterns/ → integration/
```

## 📁 All Examples Overview

### 🟢 Beginner Examples (Perfect for getting started)

| Example | What it teaches | Time | Key concepts |
|---------|----------------|------|--------------|
| **[Simple Chat](simple_chat/)** | Basic secure messaging | 15 min | P2P messaging, QR discovery, encryption |
| **[Credential Issuance](credential_issuance/)** | Creating digital credentials | 20 min | Credential creation, signing, claims |
| **[Credential Exchange](credentials_exchange/)** | Sharing credentials | 20 min | Credential requests, verification, sharing |

### 🟡 Intermediate Examples (Building on the basics)

| Example | What it teaches | Time | Key concepts |
|---------|----------------|------|--------------|
| **[Discovery Subscription](discovery_subscription/)** | Finding and connecting to peers | 20 min | Peer discovery, QR codes, event handling |
| **[Group Chat](group_chat/)** | Multi-party messaging | 25 min | Group management, roles, invitations |
| **[Notifications](advanced_features/notifications/)** | Push notifications | 15 min | User engagement, event handling, alerts |

### 🟠 Advanced Examples (Production-ready patterns)

| Example | What it teaches | Time | Key concepts |
|---------|----------------|------|--------------|
| **[Storage](advanced_features/storage/)** | Data persistence | 20 min | Encryption, namespacing, TTL, caching |
| **[Pairing](advanced_features/pairing/)** | Multi-device sync | 20 min | Device verification, cross-device state |
| **[Production Patterns](advanced_features/production_patterns/)** | Real-world patterns | 25 min | Sessions, state management, optimization |

### 🔴 Expert Examples (Complex integrations)

| Example | What it teaches | Time | Key concepts |
|---------|----------------|------|--------------|
| **[Integration](advanced_features/integration/)** | Multi-component workflows | 45 min | Component coordination, complete applications |

## 🔍 Find Examples by Feature

<details>
<summary><strong>💬 Messaging & Communication</strong></summary>

- **[Simple Chat](simple_chat/)** - 1-to-1 secure messaging
- **[Group Chat](group_chat/)** - Multi-party group communication
- **[Discovery Subscription](discovery_subscription/)** - Finding and connecting to peers
- **[Notifications](advanced_features/notifications/)** - Push notifications for engagement

</details>

<details>
<summary><strong>🆔 Credentials & Identity</strong></summary>

- **[Credential Issuance](credential_issuance/)** - Creating and signing credentials
- **[Credential Exchange](credentials_exchange/)** - Requesting and sharing credentials

</details>

<details>
<summary><strong>💾 Data & Storage</strong></summary>

- **[Storage](advanced_features/storage/)** - Encrypted data persistence with TTL
- **[Production Patterns](advanced_features/production_patterns/)** - Session and state management

</details>

<details>
<summary><strong>🔗 Device & Connectivity</strong></summary>

- **[Pairing](advanced_features/pairing/)** - Multi-device synchronization
- **[Discovery Subscription](discovery_subscription/)** - Peer discovery and connection

</details>

<details>
<summary><strong>🏭 Production & Integration</strong></summary>

- **[Production Patterns](advanced_features/production_patterns/)** - Real-world application patterns
- **[Integration](advanced_features/integration/)** - Complete multi-component workflows

</details>

## ⚡ Quick Commands

```bash
# Run any example with one command
cd <example_name> && go run main.go

# Examples:
cd simple_chat && go run main.go                    # Start with messaging
cd credential_issuance && go run main.go            # Learn credentials
cd advanced_features/storage && go run main.go      # Try advanced storage
cd advanced_features/integration && go run main.go  # See full integration
```

## 🎯 Choose Your Adventure

### 👋 "I'm completely new to Self SDK"
**Start here:** [`simple_chat/`](simple_chat/) → [`credential_issuance/`](credential_issuance/)

### 💼 "I want to build a real application"
**Start here:** [`advanced_features/notifications/`](advanced_features/notifications/) → [`advanced_features/storage/`](advanced_features/storage/) → [`advanced_features/integration/`](advanced_features/integration/)

### 🔍 "I need a specific feature"
**Use the feature finder above** ☝️ or check the [Quick Start table](#-quick-start)

### 🎓 "I want to learn everything systematically"
**Follow the complete learning path:** All beginner examples → All intermediate → All advanced → Expert

## 🛠️ Prerequisites

- **Go 1.19+** - All examples require Go 1.19 or later
- **Self SDK** - Dependencies handled automatically by `go.mod`
- **5-10 minutes** - Most examples run in under 10 minutes

## 💡 Tips for Success

- ✅ **Start simple** - Begin with `simple_chat/` even if you're experienced
- ✅ **Follow the progression** - Each example builds on previous concepts
- ✅ **Read the READMEs** - Each example has detailed documentation
- ✅ **Experiment** - Modify the code to understand how it works
- ✅ **Check complexity ratings** - Don't skip ahead too quickly

## 🤝 Need Help?

- 📖 **Each example has a detailed README** with troubleshooting
- 🔧 **Build issues?** Run `go mod tidy` in the example directory
- 🐛 **Something not working?** Check the prerequisites and error messages
- 💬 **Questions?** Check the [Self SDK Documentation](https://docs.joinself.com)

## 📚 Additional Resources

- [Self SDK Documentation](https://docs.joinself.com)
- [W3C Verifiable Credentials](https://w3.org/TR/vc-data-model/)
- [Decentralized Identity Foundation](https://identity.foundation)

---

**Ready to start?** Pick an example above and run `cd <example> && go run main.go` 🚀
