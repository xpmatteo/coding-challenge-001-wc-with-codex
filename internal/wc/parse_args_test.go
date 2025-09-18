package wc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseArgsWithBytesFlag(t *testing.T) {
	cfg, err := ParseArgs([]string{"-c", "sample.txt"})
	require.NoError(t, err)
	require.True(t, cfg.CountBytes)
	require.Equal(t, []string{"sample.txt"}, cfg.Files)
}

func TestParseArgsWithoutFlagsLeavesCountersDisabled(t *testing.T) {
	cfg, err := ParseArgs([]string{"sample.txt"})
	require.NoError(t, err)
	require.False(t, cfg.CountBytes)
	require.Equal(t, []string{"sample.txt"}, cfg.Files)
}
