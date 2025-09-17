package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMainSmoke(t *testing.T) {
	tmp := t.TempDir()
	sample := filepath.Join(tmp, "sample.txt")
	if err := os.WriteFile(sample, []byte("hi\n"), 0o644); err != nil {
		t.Fatalf("write sample: %v", err)
	}

	cmd := exec.Command("go", "run", ".", sample)
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run . %s failed: %v\n%s", sample, err, out)
	}
	if len(out) == 0 {
		t.Fatalf("expected output, got none")
	}
}
