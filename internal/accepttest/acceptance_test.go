package accepttest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/tools/txtar"
)

type directives struct {
	args  []string
	stdin string
	env   []string
}

func parseDirectives(comment []byte) directives {
	d := directives{}
	lines := strings.Split(string(comment), "\n")
	for _, raw := range lines {
		raw = strings.TrimSpace(raw)
		if raw == "" || !strings.HasPrefix(raw, "#") {
			continue
		}
		trimmed := strings.TrimSpace(strings.TrimPrefix(raw, "#"))
		switch {
		case strings.HasPrefix(trimmed, "args:"):
			fields := strings.Fields(strings.TrimSpace(strings.TrimPrefix(trimmed, "args:")))
			d.args = append(d.args, fields...)
		case strings.HasPrefix(trimmed, "stdin:"):
			d.stdin = strings.TrimSpace(strings.TrimPrefix(trimmed, "stdin:"))
		case strings.HasPrefix(trimmed, "env:"):
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "env:"))
			if val != "" {
				d.env = append(d.env, val)
			}
		}
	}
	return d
}

func TestAcceptanceSuite(t *testing.T) {
	matches, err := filepath.Glob("testdata/*.txtar")
	if err != nil {
		t.Fatalf("glob acceptance files: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("no acceptance fixtures present")
	}

	bin := buildBinary(t)

	for _, path := range matches {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			runCase(t, bin, path)
		})
	}
}

func runCase(t *testing.T, bin, archivePath string) {
	t.Helper()

	data, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("read archive: %v", err)
	}

	ar := txtar.Parse(data)
	dirs := parseDirectives(ar.Comment)

	workdir := t.TempDir()

	expectedStdout := []byte{}
	expectedStderr := []byte{}
	expectedExit := 0

	for _, file := range ar.Files {
		target := filepath.Join(workdir, file.Name)
		switch file.Name {
		case "stdout.txt":
			expectedStdout = append([]byte(nil), file.Data...)
		case "stderr.txt":
			expectedStderr = append([]byte(nil), file.Data...)
		case "exitcode":
			text := strings.TrimSpace(string(file.Data))
			if text != "" {
				val, err := strconv.Atoi(text)
				if err != nil {
					t.Fatalf("parse exit code in %s: %v", archivePath, err)
				}
				expectedExit = val
			}
		default:
			if err := writeFile(workdir, target, file.Data); err != nil {
				t.Fatalf("materialize %s: %v", file.Name, err)
			}
		}
	}

	stdoutExp := applyPlaceholders(expectedStdout, workdir)
	stderrExp := applyPlaceholders(expectedStderr, workdir)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, dirs.args...)
	cmd.Dir = workdir
	cmd.Env = append(os.Environ(), dirs.env...)

	if dirs.stdin != "" {
		stdinPath := filepath.Join(workdir, dirs.stdin)
		stdinFile, err := os.Open(stdinPath)
		if err != nil {
			t.Fatalf("open stdin file %q: %v", dirs.stdin, err)
		}
		defer stdinFile.Close()
		cmd.Stdin = stdinFile
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("execute %s: %v", archivePath, err)
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("command timed out")
	}

	if exitCode != expectedExit {
		t.Errorf("exit code: got %d, want %d", exitCode, expectedExit)
	}

	stdoutGot := stdoutBuf.Bytes()
	stderrGot := stderrBuf.Bytes()

	if !bytes.Equal(stdoutGot, stdoutExp) {
		t.Errorf("stdout mismatch\n--- got ---\n%s\n--- want ---\n%s", stdoutGot, stdoutExp)
	}

	if !bytes.Equal(stderrGot, stderrExp) {
		t.Errorf("stderr mismatch\n--- got ---\n%s\n--- want ---\n%s", stderrGot, stderrExp)
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	output := filepath.Join(tempDir, "wc-accept")

	cmd := exec.Command("go", "build", "-o", output, "./cmd/wc")
	cmd.Env = os.Environ()
	cmd.Dir = "."

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("failed to build ./cmd/wc: %v\n%s", err, out)
	}

	return output
}

func writeFile(root, path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	if err := os.WriteFile(path, data, fs.FileMode(0o644)); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

func applyPlaceholders(input []byte, workdir string) []byte {
	if len(input) == 0 {
		return input
	}
	replaced := strings.ReplaceAll(string(input), "%TMPDIR%", workdir)
	return []byte(replaced)
}
