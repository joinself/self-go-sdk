package client

import (
	"github.com/joinself/self-go-sdk/account"
)

// Environment represents the target environment
type Environment int

const (
	Sandbox Environment = iota
	Production
)

// LogLevel represents logging verbosity
type LogLevel int

const (
	LogError LogLevel = iota
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

// Config holds configuration for the Self client
type Config struct {
	// StorageKey is the encryption key for local storage (required)
	StorageKey []byte

	// StoragePath is the directory for local storage (required)
	StoragePath string

	// Environment specifies the target environment (default: Sandbox)
	Environment Environment

	// LogLevel specifies logging verbosity (default: LogWarn)
	LogLevel LogLevel

	// SkipReady skips the ready check during initialization
	SkipReady bool

	// SkipSetup skips the setup phase during initialization
	SkipSetup bool
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if len(c.StorageKey) == 0 {
		return ErrStorageKeyRequired
	}
	if c.StoragePath == "" {
		return ErrStoragePathRequired
	}
	return nil
}

// toAccountConfig converts client config to account config
func (c *Config) toAccountConfig() *account.Config {
	cfg := &account.Config{
		StorageKey:  c.StorageKey,
		StoragePath: c.StoragePath,
		SkipReady:   c.SkipReady,
		SkipSetup:   c.SkipSetup,
	}

	// Set environment
	switch c.Environment {
	case Production:
		cfg.Environment = account.TargetProduction
	default:
		cfg.Environment = account.TargetSandbox
	}

	// Set log level
	switch c.LogLevel {
	case LogError:
		cfg.LogLevel = account.LogError
	case LogWarn:
		cfg.LogLevel = account.LogWarn
	case LogInfo:
		cfg.LogLevel = account.LogInfo
	case LogDebug:
		cfg.LogLevel = account.LogDebug
	case LogTrace:
		cfg.LogLevel = account.LogTrace
	default:
		cfg.LogLevel = account.LogWarn
	}

	return cfg
}
