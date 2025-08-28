package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T){
// Test: Valid single header
headers := NewHeaders()
data := []byte("Host: localhost:42069\r\n   foo: barrba:rr\r\n\r\n")
n, done, err := headers.Parse(data)
require.NoError(t, err)
require.NotNil(t, headers)
host , ok := headers.Get("HOST")
assert.True(t, ok)
assert.Equal(t, "localhost:42069", host)
foo ,ok := headers.Get("foo")
assert.True(t , ok)
assert.Equal(t, "barrba:rr", foo)
assert.Equal(t, 44, n)
assert.True(t, done)

// Test: Invalid spacing header
headers = NewHeaders()
data = []byte("       Host : localhost:42069       \r\n\r\n")
n, done, err = headers.Parse(data)
require.Error(t, err)
assert.Equal(t, 0, n)
assert.False(t, done)

//Test: Check valid header
headers = NewHeaders()
data = []byte("H©st: localhost:42069\r\n\r\n")
n, done, err = headers.Parse(data)
require.Error(t, err)
assert.Equal(t, 0, n)
assert.False(t, done)

//Test: Valid header with multiple bodies
headers = NewHeaders()
data = []byte("Host: localhost:42069\r\n  Host: satrat\r\n\r\n")
n, done, err = headers.Parse(data)
require.NoError(t, err)
require.NotNil(t, headers)
host , ok = headers.Get("HOST")
assert.True(t, ok)
assert.Equal(t, "localhost:42069,satrat", host)
assert.Equal(t, 41, n)
assert.True(t, done)
}