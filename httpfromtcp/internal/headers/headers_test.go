package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFoo:     barbar   \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	str, _ := headers.Get("HOST")
	assert.Equal(t, "localhost:42069", str)
	str, _ = headers.Get("Foo")
	assert.Equal(t, "barbar", str)
	assert.Equal(t, 45, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: localhost:42060\r\n")
	n, done, err = headers.Parse(data)
	// require.NoError(t, err)
	// require.NotNil(t, headers)
	str, _ = headers.Get("HOST")
	assert.Equal(t, "localhost:42069,localhost:42060", str)
	// assert.True(t, done)
}
