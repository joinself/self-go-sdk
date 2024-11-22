package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "github.com/joinself/self-go-sdk-next/message"

type LogLevel uint32

const (
	LogError LogLevel = C.LOG_ERROR
	LogWarn  LogLevel = C.LOG_WARN
	LogInfo  LogLevel = C.LOG_INFO
	LogDebug LogLevel = C.LOG_DEBUG
	LogTrace LogLevel = C.LOG_TRACE
)

var (
	TargetProduction = &Target{
		Rpc:     defaultRpcProduction,
		Object:  defaultObjectProduction,
		Message: defaultMessageProduction,
	}
	TargetSandbox = &Target{
		Rpc:     defaultRpcSandbox,
		Object:  defaultObjectSandbox,
		Message: defaultMessageSandbox,
	}
)

type Target struct {
	Rpc     string
	Object  string
	Message string
}

var (
	defaultRpcSandbox        = "https://rpc.sandbox.joinself.com/"
	defaultObjectSandbox     = "https://object.sandbox.joinself.com/"
	defaultMessageSandbox    = "wss://message.sandbox.joinself.com/"
	defaultRpcProduction     = "https://rpc.joinself.com/"
	defaultObjectProduction  = "https://object.joinself.com/"
	defaultMessageProduction = "wss://message.joinself.com/"
)

// Config stores config for an account
type Config struct {
	StorageKey  []byte
	StoragePath string
	Environment *Target
	LogLevel    LogLevel
	Callbacks   Callbacks
}

// Callbacks defines callbacks invoked by the account
type Callbacks struct {
	OnConnect         func(account *Account)
	OnDisconnect      func(account *Account, err error)
	OnAcknowledgement func(account *Account, reference *message.Reference)
	OnError           func(account *Account, reference *message.Reference, err error)
	OnMessage         func(account *Account, message *message.Message)
	OnCommit          func(account *Account, commit *message.Commit)
	OnKeyPackage      func(account *Account, keyPackage *message.KeyPackage)
	OnProposal        func(account *Account, proposal *message.Proposal)
	OnWelcome         func(account *Account, welcome *message.Welcome)
}

func (c *Config) defaults() {
	if c.LogLevel == 0 {
		c.LogLevel = LogError
	}

	if c.Environment == nil {
		c.Environment = TargetSandbox
	}
}
