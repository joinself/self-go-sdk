package account_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/joinself/self-go-sdk/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAccount(t testing.TB) (*account.Account, chan *account.Message) {
	incoming := make(chan *account.Message, 1024)

	cfg := &account.Config{
		StorageKey:  make([]byte, 32),
		StoragePath: t.TempDir() + "self.db",
		Callbacks: account.Callbacks{
			OnConnect: func() {},
			OnDisconnect: func(err error) {
				require.Nil(t, err)
			},
			OnMessage: func(account *account.Account, message *account.Message) {
				fmt.Println(
					"to:", message.ToAddress(),
					"from:", message.FromAddress(),
					"message:", string(message.Message()),
				)

				incoming <- message
			},
		},
	}

	acc, err := account.New(cfg)
	require.Nil(t, err)

	return acc, incoming
}

func wait(t testing.TB, ch chan *account.Message, timeout time.Duration) *account.Message {
	select {
	case <-time.After(timeout):
		require.Nil(t, errors.New("timeout"))
		return nil
	case m := <-ch:
		return m
	}
}

func TestAccount(t *testing.T) {
	alice, aliceInbox := testAccount(t)
	bobby, bobbyInbox := testAccount(t)

	aliceAddress, err := alice.InboxOpen()
	require.Nil(t, err)

	bobbyAddress, err := bobby.InboxOpen()
	require.Nil(t, err)

	fmt.Println("alice:", aliceAddress)
	fmt.Println("bobby:", bobbyAddress)

	err = alice.ConnectionNegotiate(
		aliceAddress,
		bobbyAddress,
	)

	require.Nil(t, err)

	// wait for negotiation to finish
	time.Sleep(time.Millisecond * 100)

	// send a message from alice
	err = alice.MessageSend(
		bobbyAddress,
		[]byte("hello"),
	)

	require.Nil(t, err)

	message := wait(t, bobbyInbox, time.Second)
	assert.Equal(t, aliceAddress.String(), message.FromAddress().String())
	assert.Equal(t, []byte("hello"), message.Message())

	// send a response from bobby
	err = bobby.MessageSend(
		aliceAddress,
		[]byte("hi!"),
	)

	require.Nil(t, err)

	message = wait(t, aliceInbox, time.Second)
	assert.Equal(t, bobbyAddress.String(), message.FromAddress().String())
	assert.Equal(t, []byte("hi!"), message.Message())
}
