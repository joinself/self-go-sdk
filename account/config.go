package account

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"
import "github.com/joinself/self-go-sdk-next/message"

type LogLevel uint32

const (
	LogError      LogLevel = C.LOG_ERROR
	LogWarn       LogLevel = C.LOG_WARN
	LogInfo       LogLevel = C.LOG_INFO
	LogDebug      LogLevel = C.LOG_DEBUG
	LogTrace      LogLevel = C.LOG_TRACE
	TargetDevelop          = "develop"
	TargetSandbox          = "sandbox"
)

var (
	defaultRpcDevelop     = "http://127.0.0.1:8080/"
	defaultObjectDevelop  = "http://127.0.0.1:8090/"
	defaultMessageDevelop = "ws://127.0.0.1:9000/"
	defaultRpcSandbox     = "https://rpc.next.sandbox.joinself.com/"
	defaultObjectSandbox  = "https://object.next.sandbox.joinself.com/"
	defaultMessageSandbox = "wss://message.next.sandbox.joinself.com/"
)

// Config stores config for an account
type Config struct {
	StorageKey  []byte
	StoragePath string
	Environment string
	LogLevel    LogLevel
	Callbacks   Callbacks
}

// Callbacks defines callbacks invoked by the account
type Callbacks struct {
	OnConnect    func()
	OnDisconnect func(err error)
	OnMessage    func(account *Account, message *message.Message)
	OnCommit     func(account *Account, commit *message.Commit)
	OnKeyPackage func(account *Account, keyPackage *message.KeyPackage)
	OnProposal   func(account *Account, proposal *message.Proposal)
	OnWelcome    func(account *Account, welcome *message.Welcome)
}

func (c *Config) defaults() {
	if c.LogLevel == 0 {
		c.LogLevel = LogError
	}

	if c.Environment != TargetDevelop && c.Environment != TargetSandbox {
		c.Environment = TargetSandbox
	}
}

func (c *Config) rpcURL() string {
	switch c.Environment {
	case TargetDevelop:
		return defaultRpcDevelop
	default:
		return defaultRpcSandbox
	}
}

func (c *Config) objectURL() string {
	switch c.Environment {
	case TargetDevelop:
		return defaultObjectDevelop
	default:
		return defaultObjectSandbox
	}
}

func (c *Config) messageURL() string {
	switch c.Environment {
	case TargetDevelop:
		return defaultMessageDevelop
	default:
		return defaultMessageSandbox
	}
}
