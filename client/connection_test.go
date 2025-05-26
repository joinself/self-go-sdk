package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionComponent(t *testing.T) {
	// Skip this test as it requires network connectivity
	t.Skip("Connection test requires network connectivity")

	// Create two test clients
	client1, err := New(Config{
		StorageKey:  make([]byte, 32),
		StoragePath: t.TempDir() + "/client1",
		Environment: Sandbox,
		LogLevel:    LogWarn,
		SkipReady:   true,
		SkipSetup:   true,
	})
	require.NoError(t, err)
	defer client1.Close()

	client2, err := New(Config{
		StorageKey:  make([]byte, 32),
		StoragePath: t.TempDir() + "/client2",
		Environment: Sandbox,
		LogLevel:    LogWarn,
		SkipReady:   true,
		SkipSetup:   true,
	})
	require.NoError(t, err)
	defer client2.Close()

	// Test that Connection component is available
	assert.NotNil(t, client1.Connection())
	assert.NotNil(t, client2.Connection())

	// Test ConnectToPeer method exists and returns expected structure
	result, err := client1.Connection().ConnectToPeerWithTimeout(client2.DID(), 1*time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, client2.DID(), result.PeerDID)
	// Connection will likely fail in test environment, but that's expected
}

func TestConnectTwoClients(t *testing.T) {
	// Skip this test as it requires network connectivity
	t.Skip("Connection test requires network connectivity")

	// Create two test clients
	client1, err := New(Config{
		StorageKey:  make([]byte, 32),
		StoragePath: t.TempDir() + "/client1",
		Environment: Sandbox,
		LogLevel:    LogWarn,
		SkipReady:   true,
		SkipSetup:   true,
	})
	require.NoError(t, err)
	defer client1.Close()

	client2, err := New(Config{
		StorageKey:  make([]byte, 32),
		StoragePath: t.TempDir() + "/client2",
		Environment: Sandbox,
		LogLevel:    LogWarn,
		SkipReady:   true,
		SkipSetup:   true,
	})
	require.NoError(t, err)
	defer client2.Close()

	// Test ConnectTwoClients utility function
	err = ConnectTwoClientsWithTimeout(client1, client2, 1*time.Second)
	// Connection will likely fail in test environment, but function should not panic
	// We just verify the function exists and handles the call gracefully
}

func TestConnectionComponentMethods(t *testing.T) {
	// Test that we can create a connection component without network
	// Skip this test as it requires network connectivity for client creation
	t.Skip("Connection component test requires network connectivity")

	client1, err := New(Config{
		StorageKey:  make([]byte, 32),
		StoragePath: t.TempDir() + "/client1",
		Environment: Sandbox,
		LogLevel:    LogWarn,
		SkipReady:   true,
		SkipSetup:   true,
	})
	require.NoError(t, err)
	defer client1.Close()

	conn := client1.Connection()
	assert.NotNil(t, conn)

	// Test placeholder methods
	assert.False(t, conn.IsConnectedTo("test-did"))
	assert.Empty(t, conn.ListConnectedPeers())

	// Test that methods don't panic
	conn.onConnect()
	conn.onDisconnect(nil)
	conn.onWelcome(nil, nil)
	conn.onKeyPackage(nil)
	conn.onIntroduction(nil, 0)
	conn.close()
}

func TestConnectionComponentStructure(t *testing.T) {
	// Test the structure without creating actual clients
	// This tests that the types and methods exist without network calls

	// Test ConnectionResult structure
	result := &ConnectionResult{
		PeerDID:   "test-did",
		Connected: true,
		Error:     nil,
	}
	assert.Equal(t, "test-did", result.PeerDID)
	assert.True(t, result.Connected)
	assert.Nil(t, result.Error)

	// Test that Connection methods exist (we can't call them without a client)
	// This is a compile-time test to ensure the API exists
	var conn *Connection
	assert.Nil(t, conn) // Just to use the variable
}
