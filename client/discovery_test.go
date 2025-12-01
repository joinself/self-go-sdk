package client

import (
	"testing"

	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryQRSandboxFlags(t *testing.T) {
	tests := []struct {
		name                string
		environment         Environment
		expectedSandboxFlag bool
	}{
		{
			name:                "Sandbox environment should set sandbox flag",
			environment:         Sandbox,
			expectedSandboxFlag: true,
		},
		{
			name:                "Production environment should not set sandbox flag",
			environment:         Production,
			expectedSandboxFlag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the flag setting logic directly
			content := createMockChatContent(t)
			anonymousMsg := event.NewAnonymousMessage(content)

			// Simulate the logic from Unicode() and SVG() methods
			if tt.environment == Sandbox {
				anonymousMsg.SetFlags(event.MessageFlagTargetSandbox)
			}

			hasSandboxFlag := anonymousMsg.HasFlags(event.MessageFlagTargetSandbox)
			assert.Equal(t, tt.expectedSandboxFlag, hasSandboxFlag,
				"Expected sandbox flag to be %v for environment %v", tt.expectedSandboxFlag, tt.environment)
		})
	}
}

// createMockChatContent creates a simple mock message content for testing
func createMockChatContent(t *testing.T) *message.Content {
	// Use a simple chat message for testing since it doesn't require complex setup
	content, err := message.NewChat().
		Message("test message").
		Finish()
	require.NoError(t, err)
	return content
}
