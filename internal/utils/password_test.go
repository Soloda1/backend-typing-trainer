package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	const password = "s3cret-pass"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEqual(t, password, hash)

	require.True(t, CheckPasswordHash(password, hash))
	require.False(t, CheckPasswordHash("wrong-password", hash))
}
