package client

import (
	"testing"

	"github.com/joinself/self-go-sdk/account"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentMapping(t *testing.T) {
	tests := []struct {
		name               string
		clientEnvironment  Environment
		expectedAccountEnv *account.Target
	}{
		{
			name:               "Sandbox environment",
			clientEnvironment:  Sandbox,
			expectedAccountEnv: account.TargetSandbox,
		},
		{
			name:               "Production environment",
			clientEnvironment:  Production,
			expectedAccountEnv: account.TargetProduction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				StorageKey:  make([]byte, 32),
				StoragePath: "/tmp/test",
				Environment: tt.clientEnvironment,
				LogLevel:    LogWarn,
			}

			accountConfig := config.toAccountConfig()

			assert.Equal(t, tt.expectedAccountEnv, accountConfig.Environment)

			// Also verify the URLs are correct
			if tt.clientEnvironment == Sandbox {
				assert.Equal(t, "https://rpc-sandbox.joinself.com/", accountConfig.Environment.Rpc)
				assert.Equal(t, "https://object-sandbox.joinself.com/", accountConfig.Environment.Object)
				assert.Equal(t, "wss://message-sandbox.joinself.com/", accountConfig.Environment.Message)
			} else if tt.clientEnvironment == Production {
				assert.Equal(t, "https://rpc.joinself.com/", accountConfig.Environment.Rpc)
				assert.Equal(t, "https://object.joinself.com/", accountConfig.Environment.Object)
				assert.Equal(t, "wss://message.joinself.com/", accountConfig.Environment.Message)
			}
		})
	}
}

func TestDefaultEnvironment(t *testing.T) {
	// Test that default environment (zero value) maps to Sandbox
	config := &Config{
		StorageKey:  make([]byte, 32),
		StoragePath: "/tmp/test",
		// Environment not set (zero value)
		LogLevel: LogWarn,
	}

	accountConfig := config.toAccountConfig()

	assert.Equal(t, account.TargetSandbox, accountConfig.Environment)
	assert.Equal(t, "https://rpc-sandbox.joinself.com/", accountConfig.Environment.Rpc)
}

// TestClientCreationWithSandboxEnvironment is commented out because it requires network connectivity
// func TestClientCreationWithSandboxEnvironment(t *testing.T) {
// 	// Test that we can create a client with Sandbox environment
// 	// and it correctly maps to the account configuration
// 	config := Config{
// 		StorageKey:  make([]byte, 32),
// 		StoragePath: t.TempDir() + "/test_storage",
// 		Environment: Sandbox,
// 		LogLevel:    LogWarn,
// 		SkipReady:   true, // Skip ready check for testing
// 		SkipSetup:   true, // Skip setup for testing
// 	}
//
// 	client, err := NewClient(config)
// 	require.NoError(t, err)
// 	require.NotNil(t, client)
// 	defer client.Close()
//
// 	// Verify the client was created successfully
// 	assert.NotNil(t, client.account)
// 	assert.Equal(t, Sandbox, client.config.Environment)
// }
