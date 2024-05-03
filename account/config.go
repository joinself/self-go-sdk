package account

import "github.com/joinself/self-go-sdk/keypair/signing"

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
	OnMessage    func(account *Account, message *Message)
}

// Message message recevied from another sender
type Message struct {
	fromAddress *signing.PublicKey
	toAddress   *signing.PublicKey
	message     []byte
}

// FromAddress the public key of the sender
func (m *Message) FromAddress() *signing.PublicKey {
	return m.fromAddress
}

// ToAddress the destination address of the recipient
func (m *Message) ToAddress() *signing.PublicKey {
	return m.toAddress
}

// Message the message content
func (m *Message) Message() []byte {
	return m.message
}
