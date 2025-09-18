package wc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeFileCountsBytes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	content := []byte("abc")
	require.NoError(t, os.WriteFile(path, content, 0o644))

	stat, err := AnalyzeFile(path)
	require.NoError(t, err)
	require.Equal(t, path, stat.Name)
	require.Equal(t, len(content), stat.Bytes)
	require.Equal(t, 1, stat.Lines)
}

func TestAnalyzeFileCountsTrailingNewlines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	content := []byte("first\nsecond\n")
	require.NoError(t, os.WriteFile(path, content, 0o644))

	stat, err := AnalyzeFile(path)
	require.NoError(t, err)
	require.Equal(t, 2, stat.Lines)
}
