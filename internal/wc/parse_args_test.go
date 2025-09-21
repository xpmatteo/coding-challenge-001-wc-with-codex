package wc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseArgs(t *testing.T) {
	tests := map[string]struct {
		args      []string
		expectErr bool
		expectCfg Config
	}{
		"no flags": {
			args:      []string{"sample.txt"},
			expectCfg: Config{Files: []string{"sample.txt"}},
		},
		"bytes short flag": {
			args:      []string{"-c", "sample.txt"},
			expectCfg: Config{Files: []string{"sample.txt"}, CountBytes: true, counterOrder: []counterKind{counterBytes}},
		},
		"lines long flag": {
			args:      []string{"--lines", "sample.txt"},
			expectCfg: Config{Files: []string{"sample.txt"}, CountLines: true, counterOrder: []counterKind{counterLines}},
		},
		"words short flag": {
			args:      []string{"-w", "sample.txt"},
			expectCfg: Config{Files: []string{"sample.txt"}, CountWords: true, counterOrder: []counterKind{counterWords}},
		},
		"chars long flag": {
			args:      []string{"--chars", "sample.txt"},
			expectCfg: Config{Files: []string{"sample.txt"}, CountChars: true, counterOrder: []counterKind{counterChars}},
		},
		"flag order preserved": {
			args: []string{"-w", "-l", "sample.txt"},
			expectCfg: Config{
				Files:        []string{"sample.txt"},
				CountLines:   true,
				CountWords:   true,
				counterOrder: []counterKind{counterWords, counterLines},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cfg, err := ParseArgs(tc.args)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectCfg.CountBytes, cfg.CountBytes)
			require.Equal(t, tc.expectCfg.CountLines, cfg.CountLines)
			require.Equal(t, tc.expectCfg.CountWords, cfg.CountWords)
			require.Equal(t, tc.expectCfg.CountChars, cfg.CountChars)
			require.Equal(t, tc.expectCfg.Files, cfg.Files)
			require.Equal(t, tc.expectCfg.counterOrder, cfg.counterOrder)
		})
	}
}
