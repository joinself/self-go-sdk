package account

type Config struct {
	StorageKey  []byte
	StoragePath string
	Callbacks   Callbacks
}

type Callbacks struct {
	OnConnect    func()
	OnDisconnect func(err error)
	OnMessage    func(account *Account, message *Message)
}

type Message struct {
	fromAddress *PublicKey
	toAddress   *PublicKey
	message     []byte
}

func (m *Message) FromAddress() *PublicKey {
	return m.fromAddress
}

func (m *Message) ToAddress() *PublicKey {
	return m.toAddress
}

func (m *Message) Message() []byte {
	return m.message
}
