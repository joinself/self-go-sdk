package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "github.com/joinself/self-go-sdk/event"

var (
	TargetProduction = &Target{
		Rpc:     defaultRpcProduction,
		Object:  defaultObjectProduction,
		Message: defaultMessageProduction,
	}
	TargetProductionSandbox = &Target{
		Rpc:     defaultRpcProductionSandbox,
		Object:  defaultObjectProductionSandbox,
		Message: defaultMessageProductionSandbox,
	}
	TargetStaging = &Target{
		Rpc:     defaultRpcStaging,
		Object:  defaultObjectStaging,
		Message: defaultMessageStaging,
	}
	TargetStagingSandbox = &Target{
		Rpc:     defaultRpcStagingSandbox,
		Object:  defaultObjectStagingSandbox,
		Message: defaultMessageStagingSandbox,
	}
	TargetPreview = &Target{
		Rpc:     defaultRpcPreview,
		Object:  defaultObjectPreview,
		Message: defaultMessagePreview,
	}
	TargetPreviewSandbox = &Target{
		Rpc:     defaultRpcPreviewSandbox,
		Object:  defaultObjectPreviewSandbox,
		Message: defaultMessagePreviewSandbox,
	}
)

type Target struct {
	Rpc     string
	Object  string
	Message string
}

var (
	defaultRpcProduction            = "https://rpc.joinself.com/"
	defaultObjectProduction         = "https://object.joinself.com/"
	defaultMessageProduction        = "wss://message.joinself.com/"
	defaultRpcProductionSandbox     = "https://rpc-sandbox.joinself.com/"
	defaultObjectProductionSandbox  = "https://object-sandbox.joinself.com/"
	defaultMessageProductionSandbox = "wss://message-sandbox.joinself.com/"
	defaultRpcStaging               = "https://rpc.staging.joinself.com/"
	defaultObjectStaging            = "https://object.staging.joinself.com/"
	defaultMessageStaging           = "wss://message.staging.joinself.com/"
	defaultRpcStagingSandbox        = "https://rpc-sandbox.staging.joinself.com/"
	defaultObjectStagingSandbox     = "https://object-sandbox.staging.joinself.com/"
	defaultMessageStagingSandbox    = "wss://message-sandbox.staging.joinself.com/"
	defaultRpcPreview               = "https://rpc.preview.joinself.com/"
	defaultObjectPreview            = "https://object.preview.joinself.com/"
	defaultMessagePreview           = "wss://message.preview.joinself.com/"
	defaultRpcPreviewSandbox        = "https://rpc-sandbox.preview.joinself.com/"
	defaultObjectPreviewSandbox     = "https://object-sandbox.preview.joinself.com/"
	defaultMessagePreviewSandbox    = "wss://message-sandbox.preview.joinself.com/"
)

// Config stores config for an account
type Config struct {
	SkipReady   bool
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
	OnAcknowledgement func(account *Account, reference *event.Reference)
	OnError           func(account *Account, reference *event.Reference, err error)
	OnMessage         func(account *Account, message *event.Message)
	OnCommit          func(account *Account, commit *event.Commit)
	OnKeyPackage      func(account *Account, keyPackage *event.KeyPackage)
	OnProposal        func(account *Account, proposal *event.Proposal)
	OnWelcome         func(account *Account, welcome *event.Welcome)
}

func (c *Config) defaults() {
	if c.LogLevel == 0 {
		c.LogLevel = LogError
	}

	if c.Environment == nil {
		c.Environment = TargetProductionSandbox
	}
}
