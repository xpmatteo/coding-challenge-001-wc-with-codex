package wc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeFilesReturnsStatsForEachFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	content := []byte("hello world\n")
	require.NoError(t, os.WriteFile(path, content, 0o644))

	cfg := Config{Files: []string{path}, CountBytes: true}
	stats, err := AnalyzeFiles(cfg)
	require.NoError(t, err)
	require.Len(t, stats, 1)
	require.Equal(t, path, stats[0].Name)
	require.Equal(t, len(content), stats[0].Bytes)
	require.Equal(t, 1, stats[0].Lines)
}
