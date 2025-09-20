package wc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	tests := map[string]struct {
		cfg    Config
		stats  []Stats
		expect []string
	}{
		"default skeleton": {
			cfg:    Config{},
			stats:  []Stats{{Name: "sample.txt"}},
			expect: []string{"0 0 0 sample.txt"},
		},
		"bytes only": {
			cfg:    Config{CountBytes: true},
			stats:  []Stats{{Name: "sample.txt", Bytes: 12}},
			expect: []string{"      12 sample.txt"},
		},
		"lines only": {
			cfg:    Config{CountLines: true},
			stats:  []Stats{{Name: "sample.txt", Lines: 2}},
			expect: []string{"       2 sample.txt"},
		},
		"words only": {
			cfg:    Config{CountWords: true},
			stats:  []Stats{{Name: "sample.txt", Words: 3}},
			expect: []string{"       3 sample.txt"},
		},
		"lines and words": {
			cfg: Config{CountLines: true, CountWords: true},
			stats: []Stats{{
				Name:  "sample.txt",
				Lines: 2,
				Words: 5,
			}},
			expect: []string{"       2       5 sample.txt"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			lines, err := Format(tc.cfg, tc.stats)
			require.NoError(t, err)
			require.Equal(t, tc.expect, lines)
		})
	}
}
