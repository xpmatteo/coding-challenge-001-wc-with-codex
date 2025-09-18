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
			expectCfg: Config{Files: []string{"sample.txt"}, CountBytes: true},
		},
		"lines long flag": {
			args:      []string{"--lines", "sample.txt"},
			expectCfg: Config{Files: []string{"sample.txt"}, CountLines: true},
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
			require.Equal(t, tc.expectCfg.Files, cfg.Files)
		})
	}
}
