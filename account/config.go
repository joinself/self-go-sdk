package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"

	"github.com/joinself/self-go-sdk/event"
	"github.com/joinself/self-go-sdk/platform"
)

var (
	// TargetSandbox targets the sandbox environment
	TargetSandbox = &Target{
		Rpc:     defaultRpcSandbox,
		Object:  defaultObjectSandbox,
		Message: defaultMessageSandbox,
	}
)

// Target specifies which endpoints the SDK should target
type Target struct {
	Rpc     string
	Object  string
	Message string
}

var (
	defaultRpcSandbox        = "https://rpc-sandbox.joinself.com/"
	defaultObjectSandbox     = "https://object-sandbox.joinself.com/"
	defaultMessageSandbox    = "wss://message-sandbox.joinself.com/"
	defaultRpcProduction     = "https://rpc.joinself.com/"
	defaultObjectProduction  = "https://object.joinself.com/"
	defaultMessageProduction = "wss://message.joinself.com/"
)

// Config stores config for an account
type Config struct {
	SkipReady   bool
	SkipSetup   bool
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
	OnDropped         func(account *Account, dropped *event.Dropped)
	onIntegrity       func(account *Account, requestHash []byte) *platform.Attestation
}

func (c *Config) defaults() {
	if c.LogLevel == 0 {
		c.LogLevel = LogError
	}

	if c.Environment == nil {
		c.Environment = TargetSandbox
	}
}

// DefaultWelcomeAccept automatically accepts any welcome event to join a new group
var DefaultWelcomeAccept = func(account *Account, welcome *event.Welcome) {
	groupAddress, err := account.ConnectionAccept(
		welcome.ToAddress(),
		welcome.Welcome(),
	)

	if err != nil {
		logger()(LogWarn, fmt.Sprintf("failed to accept welcome to encrypted group. error: %s", err))
		return
	}

	logger()(LogInfo, fmt.Sprintf(
		"accepted welcome to encrypted group. group: %s as: %s from: %s",
		groupAddress.String(),
		welcome.ToAddress().String(),
		welcome.FromAddress().String(),
	))
}

// DefaultWelcomeIgnore automatically ignores any welcome event to join a new group
var DefaultWelcomeIgnore = func(account *Account, welcome *event.Welcome) {
	logger()(LogInfo, fmt.Sprintf(
		"ignoring welcome to encrypted group. from: %s",
		welcome.FromAddress().String()),
	)
}

// DefaultKeyPackageAccept automatically accepts any key package and creates a new group with the sender
var DefaultKeyPackageAccept = func(account *Account, keyPackage *event.KeyPackage) {
	groupAddress, err := account.ConnectionEstablish(
		keyPackage.ToAddress(),
		keyPackage.KeyPackage(),
	)

	if err != nil {
		logger()(LogWarn, fmt.Sprintf("failed to create encrypted group from key package. error: %s", err))
		return
	}

	logger()(LogInfo, fmt.Sprintf(
		"created encrypted group from key package. group: %s as: %s from: %s",
		groupAddress.String(),
		keyPackage.ToAddress().String(),
		keyPackage.FromAddress().String(),
	))
}

// DefaultKeyPackageIgnore automatically ignores any key package and will not create a group
var DefaultKeyPackageIgnore = func(account *Account, keyPackage *event.KeyPackage) {
	logger()(LogInfo, fmt.Sprintf(
		"ignoring welcome to encrypted group. from: %s",
		keyPackage.FromAddress().String(),
	))
}

// NOTE mobile specific api, don't export
func setOnIntegrity(callbacks *Callbacks, onIntegrity func(account *Account, requestHash []byte) *platform.Attestation) {
	callbacks.onIntegrity = onIntegrity
}
