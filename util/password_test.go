package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := RandString(6)
	hashedpassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedpassword)
	match := CheckPasswordHash(password, hashedpassword)
	require.Equal(t, match, true)
	wrongPassword := RandString(6)
	match = CheckPasswordHash(wrongPassword, hashedpassword)
	require.NotEqual(t, match, true)
	hashedpassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedpassword1)
	require.NotEqual(t, hashedpassword, hashedpassword1)
}
