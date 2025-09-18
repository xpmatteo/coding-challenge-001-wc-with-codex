package wc

import "testing"

func TestParseArgsReturnsConfigWithFiles(t *testing.T) {
	cfg, err := ParseArgs([]string{"sample.txt"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}
	if len(cfg.Files) != 1 || cfg.Files[0] != "sample.txt" {
		t.Fatalf("ParseArgs returned unexpected files: %#v", cfg.Files)
	}
}
