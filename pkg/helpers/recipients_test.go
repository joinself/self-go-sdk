package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockRest struct {
	body []byte
	err  error
}

func (m *mockRest) Get(path string) ([]byte, error) {
	return m.body, m.err
}

func TestPrepareRecipients(t *testing.T) {
	var m mockRest

	_, err := PrepareRecipients([]string{"alice", "bob"}, nil, &m)
	require.NotNil(t, err)

	recipients, err := PrepareRecipients([]string{"alice", "bob"}, nil, &m)
	require.Nil(t, err)

	fmt.Println(recipients)

}
