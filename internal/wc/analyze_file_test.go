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
		expectWords int
		expectChars int
	}{
		"single line": {
			content:     []byte("abc"),
			expectBytes: 3,
			expectLines: 1,
			expectWords: 1,
			expectChars: 3,
		},
		"two lines with trailing newline": {
			content:     []byte("first\nsecond\n"),
			expectBytes: len("first\nsecond\n"),
			expectLines: 2,
			expectWords: 2,
			expectChars: 13,
		},
		"multiple spaces": {
			content:     []byte("go   gophers\n"),
			expectBytes: len("go   gophers\n"),
			expectLines: 1,
			expectWords: 2,
			expectChars: 13,
		},
		"multi-byte characters": {
			content:     []byte("héllo 🌍\n"),
			expectBytes: len("héllo 🌍\n"),
			expectLines: 1,
			expectWords: 2,
			expectChars: 8,
		},
		"invalid utf8": {
			content:     []byte{0xff, 'h', 'i', '\n'},
			expectBytes: 4,
			expectLines: 1,
			expectWords: 1,
			expectChars: 3,
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
			require.Equal(t, tc.expectWords, stat.Words)
			require.Equal(t, tc.expectChars, stat.Chars)
		})
	}
}
