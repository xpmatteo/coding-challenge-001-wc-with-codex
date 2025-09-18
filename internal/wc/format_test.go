package wc

import "testing"

func TestFormatProducesZeroOutputForDefaultConfig(t *testing.T) {
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

func TestFormatOutputsBytesWhenRequested(t *testing.T) {
	cfg := Config{CountBytes: true}
	stats := []Stats{{Name: "sample.txt", Bytes: 12}}
	lines, err := Format(cfg, stats)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}
	if len(lines) != 1 {
		t.Fatalf("expected single line, got %d", len(lines))
	}
	expected := "      12 sample.txt"
	if lines[0] != expected {
		t.Fatalf("unexpected formatted line: got %q want %q", lines[0], expected)
	}
}
