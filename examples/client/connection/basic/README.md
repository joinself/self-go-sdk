# Connection Examples

This example demonstrates the different ways to establish connections between Self SDK clients. It's designed to help you understand the connection lifecycle and choose the right approach for your use case.

## What You'll Learn

- **Programmatic Connections**: Direct peer-to-peer connections without QR codes
- **QR Code Discovery**: Standard connection method for real-world applications  
- **Connection Status Management**: How to check and monitor connection status
- **Troubleshooting**: Common issues and how to resolve them

## Features Demonstrated

### 1. Programmatic Connection (🤖)
Perfect for demos, testing, and same-process scenarios:
- `ConnectTwoClients()` utility function
- `ConnectToPeer()` for specific peer connections
- Custom timeout handling
- Error handling and status reporting

### 2. QR Code Discovery (📱)
The standard method for real-world applications:
- QR code generation with custom timeouts
- Discovery response handling
- Real-world usage patterns
- Security considerations

### 3. Connection Status & Management (📊)
Monitor and manage your connections:
- Check connection status between peers
- List all connected peers
- Connection attempt monitoring
- Best practices for connection management

### 4. Connection Troubleshooting (🔧)
Diagnose and fix common issues:
- Common error patterns and solutions
- Diagnostic checks
- Detailed error analysis
- Debugging tips and best practices

## Running the Example

```bash
cd examples/client/connection/basic
go run main.go
```

The example presents an interactive menu where you can explore different connection scenarios:

```
🔗 Self SDK Connection Examples
===============================
📋 Connection Examples Menu:
1. 🤖 Programmatic Connection (for demos/testing)
2. 📱 QR Code Discovery (real-world scenario)
3. 📊 Connection Status & Management
4. 🔧 Connection Troubleshooting
5. 🚪 Exit

Choose an option (1-5):
```

## Use Cases

### Demo Applications
Use programmatic connections to quickly establish connections without user interaction:
```go
err := client.ConnectTwoClients(client1, client2)
if err != nil {
    log.Printf("Connection failed: %v", err)
}
```

### Production Applications
Use QR code discovery for secure peer-to-peer connections:
```go
qr, err := client.Discovery().GenerateQR()
if err != nil {
    log.Fatal(err)
}

// Display QR code for scanning
qrCode, _ := qr.Unicode()
fmt.Println(qrCode)

// Wait for connection
peer, err := qr.WaitForResponse(ctx)
```

### Testing and Automation
Check connection status and handle failures gracefully:
```go
if client.Connection().IsConnectedTo(peerDID) {
    // Send message or credential
} else {
    // Establish connection first
    result, err := client.Connection().ConnectToPeer(peerDID)
}
```

## Key Concepts

### Connection Types

1. **Programmatic Connections**
   - Direct API calls to establish connections
   - No QR code scanning required
   - Perfect for demos and testing
   - Both clients must be in the same process or network

2. **QR Code Discovery**
   - Industry standard for Self SDK
   - Secure peer-to-peer discovery
   - Works across different devices/networks
   - User-friendly for mobile applications

### Connection Lifecycle

1. **Discovery**: One client generates a QR code or initiates programmatic connection
2. **Negotiation**: Clients exchange cryptographic keys and establish secure channel
3. **Establishment**: Connection is confirmed and ready for message exchange
4. **Maintenance**: Connection persists and can be monitored/managed

### Error Handling

Common connection errors and their solutions:

- **Connection Timeout**: Check network connectivity and increase timeout
- **Keypair Not Found**: Verify DID format and client initialization
- **Sender Address Not Found**: Ensure connection is established before sending messages

## Best Practices

1. **Always Check Connection Status**: Use `IsConnectedTo()` before operations
2. **Handle Timeouts Gracefully**: Set appropriate timeouts and handle failures
3. **Use Appropriate Method**: Programmatic for demos, QR codes for production
4. **Monitor Connection Health**: Implement reconnection logic for critical apps
5. **Enable Logging**: Use debug logging to troubleshoot connection issues

## Next Steps

After understanding connections, explore:
- **Chat Examples**: Send messages between connected clients
- **Credential Exchange**: Share and verify credentials
- **Group Chat**: Multi-party conversations
- **Advanced Features**: Notifications, storage, and pairing

## Troubleshooting

If you encounter issues:

1. **Enable Debug Logging**:
   ```go
   client.Config{
       LogLevel: client.LogDebug,
   }
   ```

2. **Check Network Connectivity**: Ensure both clients can reach the messaging service

3. **Verify Client Initialization**: Make sure both clients are properly created

4. **Try Different Timeouts**: Some networks may require longer connection times

5. **Use QR Code Discovery**: If programmatic connections fail, try QR code method

## Files

- `main.go`: Interactive connection examples with menu-driven interface
- `README.md`: This documentation file

## Dependencies

- Self SDK client package
- Examples utils package (for storage key generation)

This example provides a comprehensive foundation for understanding Self SDK connections and serves as a reference for implementing connection logic in your own applications. 
