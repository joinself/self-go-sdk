package account

import "github.com/joinself/self-go-sdk/message"

// Config stores config for an account
type Config struct {
	StorageKey  []byte
	StoragePath string
	Callbacks   Callbacks
}

// Callbacks defines callbacks invoked by the account
type Callbacks struct {
	OnConnect    func()
	OnDisconnect func(err error)
	OnMessage    func(account *Account, message *message.Message)
}
