package wc

import "testing"

func TestFormatProducesZeroOutput(t *testing.T) {
	cfg := Config{}
	stats := []Stats{{Name: "sample.txt"}}
	lines, err := Format(cfg, stats)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}
	if len(lines) != 1 {
		t.Fatalf("expected single line, got %d", len(lines))
	}
	if lines[0] != "0 0 0 sample.txt" {
		t.Fatalf("unexpected formatted line: %q", lines[0])
	}
}
