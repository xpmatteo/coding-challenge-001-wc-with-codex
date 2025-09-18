package wc

import "testing"

func TestParseArgsWithBytesFlag(t *testing.T) {
	cfg, err := ParseArgs([]string{"-c", "sample.txt"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}
	if !cfg.CountBytes {
		t.Fatalf("expected CountBytes to be true")
	}
	if len(cfg.Files) != 1 || cfg.Files[0] != "sample.txt" {
		t.Fatalf("ParseArgs returned unexpected files: %#v", cfg.Files)
	}
}

func TestParseArgsWithoutFlagsLeavesCountersDisabled(t *testing.T) {
	cfg, err := ParseArgs([]string{"sample.txt"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}
	if cfg.CountBytes {
		t.Fatalf("expected CountBytes to be false")
	}
	if len(cfg.Files) != 1 || cfg.Files[0] != "sample.txt" {
		t.Fatalf("ParseArgs returned unexpected files: %#v", cfg.Files)
	}
}
