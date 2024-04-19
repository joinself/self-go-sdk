package account

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
	fromAddress *PublicKey
	toAddress   *PublicKey
	message     []byte
}

// FromAddress the public key of the sender
func (m *Message) FromAddress() *PublicKey {
	return m.fromAddress
}

// ToAddress the destination address of the recipient
func (m *Message) ToAddress() *PublicKey {
	return m.toAddress
}

// Message the message content
func (m *Message) Message() []byte {
	return m.message
}
