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
	LogError LogLevel = C.LOG_ERROR
	LogWarn  LogLevel = C.LOG_WARN
	LogInfo  LogLevel = C.LOG_INFO
	LogDebug LogLevel = C.LOG_DEBUG
	LogTrace LogLevel = C.LOG_TRACE
)

// Config stores config for an account
type Config struct {
	StorageKey  []byte
	StoragePath string
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
}
