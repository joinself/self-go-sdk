package keypair

type KeyType byte

const (
	KeyTypeSigning  KeyType = 1
	KeyTypeExchange KeyType = 2
)

type PublicKey interface {
	Type() KeyType
}
