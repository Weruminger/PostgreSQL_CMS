package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetenv(t *testing.T) {
	t.Setenv("FOO_BAR", "x")
	require.Equal(t, "x", getenv("FOO_BAR", "d"))
	require.Equal(t, "d", getenv("NOPE", "d"))

	os.Unsetenv("FOO_BAR")
	require.Equal(t, "d", getenv("FOO_BAR", "d"))
}
