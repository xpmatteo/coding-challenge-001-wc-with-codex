package wc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeFiles(t *testing.T) {
	dir := t.TempDir()

	tests := map[string]struct {
		contents      map[string][]byte
		cfg           Config
		expectedStats []Stats
	}{
		"single file bytes": {
			contents: map[string][]byte{
				"sample.txt": []byte("hello world\n"),
			},
			cfg:           Config{CountBytes: true},
			expectedStats: []Stats{{Name: "sample.txt", Bytes: len("hello world\n"), Lines: 1}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			paths := make([]string, 0, len(tc.contents))
			for filename, content := range tc.contents {
				path := filepath.Join(dir, filename)
				require.NoError(t, os.WriteFile(path, content, 0o644))
				paths = append(paths, path)
			}

			cfg := tc.cfg
			cfg.Files = append(cfg.Files, paths...)

			stats, err := AnalyzeFiles(cfg)
			require.NoError(t, err)
			require.Len(t, stats, len(tc.expectedStats))
			for i, expected := range tc.expectedStats {
				require.Equal(t, expected.Bytes, stats[i].Bytes)
				require.Equal(t, expected.Lines, stats[i].Lines)
				require.Equal(t, expected.Name, filepath.Base(stats[i].Name))
			}
		})
	}
}
