package wc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeFile(t *testing.T) {
	tests := map[string]struct {
		content     []byte
		expectBytes int
		expectLines int
	}{
		"single line": {
			content:     []byte("abc"),
			expectBytes: 3,
			expectLines: 1,
		},
		"two lines with trailing newline": {
			content:     []byte("first\nsecond\n"),
			expectBytes: len("first\nsecond\n"),
			expectLines: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "sample.txt")
			require.NoError(t, os.WriteFile(path, tc.content, 0o644))

			stat, err := AnalyzeFile(path)
			require.NoError(t, err)
			require.Equal(t, path, stat.Name)
			require.Equal(t, tc.expectBytes, stat.Bytes)
			require.Equal(t, tc.expectLines, stat.Lines)
		})
	}
}
