package storage

import (
	"errors"
	"testing"

	selfcrypto "github.com/joinself/self-crypto-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageNew(t *testing.T) {
	pki := newTestPKI(t)

	_, err := New(t.TempDir(), "key", pki)
	require.Nil(t, err)
}

func TestStorageAccount(t *testing.T) {
	pki := newTestPKI(t)

	s, err := New(t.TempDir(), "key", pki)
	require.Nil(t, err)

	account, err := selfcrypto.NewAccount("alice")
	require.Nil(t, err)

	err = s.AccountCreate("alice", account)
	require.Nil(t, err)

	// test a failed action which should not update state
	err = s.AccountExecute("alice", func(account *selfcrypto.Account) error {
		account.GenerateOneTimeKeys(100)
		return errors.New("intended failure")
	})
	require.NotNil(t, err)

	// check that it's not been updated
	err = s.AccountExecute("alice", func(account *selfcrypto.Account) error {
		oneTimeKeys, err := account.OneTimeKeys()
		require.Nil(t, err)
		assert.Len(t, oneTimeKeys.Curve25519, 0)
		return nil
	})
	require.Nil(t, err)

	// update the accounts state
	err = s.AccountExecute("alice", func(account *selfcrypto.Account) error {
		return account.GenerateOneTimeKeys(100)
	})
	require.Nil(t, err)

	// check that it's been updated
	err = s.AccountExecute("alice", func(account *selfcrypto.Account) error {
		oneTimeKeys, err := account.OneTimeKeys()
		require.Nil(t, err)
		assert.Len(t, oneTimeKeys.Curve25519, 100)
		return nil
	})
	require.Nil(t, err)
}

func TestStorageEncryptAndDecrypt(t *testing.T) {
	pki := newTestPKI(t)

	s1, err := New(t.TempDir(), "key", pki)
	require.Nil(t, err)

	err = s1.AccountCreate("alice:1", registerUser(t, pki, "alice:1"))
	require.Nil(t, err)

	s2, err := New(t.TempDir(), "key", pki)
	require.Nil(t, err)

	err = s2.AccountCreate("bob:1", registerUser(t, pki, "bob:1"))
	require.Nil(t, err)

	s3, err := New(t.TempDir(), "key", pki)
	require.Nil(t, err)

	err = s3.AccountCreate("carol:1", registerUser(t, pki, "carol:1"))
	require.Nil(t, err)

	ciphertext, err := s1.Encrypt("alice:1", []any{"bob:1"}, []byte("hello"))
	require.Nil(t, err)

	plaintext, err := s2.Decrypt("alice:1", "bob:1", 1, ciphertext)
	require.Nil(t, err)
	assert.Equal(t, []byte("hello"), plaintext)

	offset, err := s2.AccountOffset("bob:1")
	require.Nil(t, err)
	assert.Equal(t, int64(1), offset)
}
