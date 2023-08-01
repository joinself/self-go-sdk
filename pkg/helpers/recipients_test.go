package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	m := mockRest{
		body: []byte(`[]`),
	}

	_, err := PrepareRecipients([]string{"alice", "bob"}, nil, &m)
	require.NotNil(t, err)

	m.body = []byte(`["1"]`)

	recipeints, err := PrepareRecipients([]string{"alice", "bob"}, nil, &m)
	require.Nil(t, err)

	assert.Equal(t, []string{"alice:1", "bob:1"}, recipeints)
}
