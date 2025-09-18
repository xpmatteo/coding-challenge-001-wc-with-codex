package wc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatProducesZeroOutputForDefaultConfig(t *testing.T) {
	cfg := Config{}
	stats := []Stats{{Name: "sample.txt"}}
	lines, err := Format(cfg, stats)
	require.NoError(t, err)
	require.Len(t, lines, 1)
	require.Equal(t, "0 0 0 sample.txt", lines[0])
}

func TestFormatOutputsBytesWhenRequested(t *testing.T) {
	cfg := Config{CountBytes: true}
	stats := []Stats{{Name: "sample.txt", Bytes: 12}}
	lines, err := Format(cfg, stats)
	require.NoError(t, err)
	require.Len(t, lines, 1)
	require.Equal(t, "      12 sample.txt", lines[0])
}

func TestFormatOutputsLinesWhenRequested(t *testing.T) {
	cfg := Config{CountLines: true}
	stats := []Stats{{Name: "sample.txt", Lines: 2}}
	lines, err := Format(cfg, stats)
	require.NoError(t, err)
	require.Len(t, lines, 1)
	require.Equal(t, "       2 sample.txt", lines[0])
}
