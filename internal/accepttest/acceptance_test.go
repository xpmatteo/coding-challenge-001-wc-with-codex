package accepttest

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v3"

	"wc/internal/cli"
)

type directives struct {
	args    []string
	stdin   string
	env     []string
	skipped bool
}

type stringList []string

func (s *stringList) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.SequenceNode:
		var list []string
		if err := value.Decode(&list); err != nil {
			return err
		}
		*s = list
	case yaml.ScalarNode:
		if value.Tag == "!!null" || value.Value == "" {
			*s = nil
			return nil
		}
		var single string
		if err := value.Decode(&single); err != nil {
			return err
		}
		if single == "" {
			*s = nil
			return nil
		}
		*s = []string{single}
	case yaml.AliasNode:
		if value.Alias == nil {
			*s = nil
			return nil
		}
		return s.UnmarshalYAML(value.Alias)
	default:
		return fmt.Errorf("expected scalar or sequence for string list, got %s", value.ShortTag())
	}
	return nil
}

func parseDirectives(comment []byte) (directives, error) {
	const marker = "#"
	var builder strings.Builder
	lines := strings.Split(strings.ReplaceAll(string(comment), "\r\n", "\n"), "\n")
	for _, raw := range lines {
		trimmed := strings.TrimRight(raw, "\r")
		idx := strings.Index(trimmed, marker)
		if idx == -1 {
			continue
		}
		line := trimmed[idx+len(marker):]
		if len(line) > 0 && line[0] == ' ' {
			line = line[1:]
		}
		if line == "" {
			builder.WriteByte('\n')
			continue
		}
		builder.WriteString(line)
		builder.WriteByte('\n')
	}

	yamlBody := strings.TrimSpace(builder.String())
	if yamlBody == "" {
		return directives{}, nil
	}

	var raw struct {
		Args    stringList `yaml:"args"`
		Stdin   string     `yaml:"stdin"`
		Env     stringList `yaml:"env"`
		Skipped bool       `yaml:"skipped"`
	}

	if err := yaml.Unmarshal([]byte(yamlBody), &raw); err != nil {
		return directives{}, fmt.Errorf("parse directives yaml: %w", err)
	}

	return directives{
		args:    append([]string(nil), raw.Args...),
		stdin:   raw.Stdin,
		env:     append([]string(nil), raw.Env...),
		skipped: raw.Skipped,
	}, nil
}

func TestAcceptanceSuite(t *testing.T) {
	matches, err := filepath.Glob("testdata/*.txtar")
	if err != nil {
		t.Fatalf("glob acceptance files: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("no acceptance fixtures present")
	}

	for _, path := range matches {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			runCase(t, path)
		})
	}
}

func runCase(t *testing.T, archivePath string) {
	t.Helper()

	data, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("read archive: %v", err)
	}

	ar := txtar.Parse(data)
	dirs, err := parseDirectives(ar.Comment)
	if err != nil {
		t.Fatalf("parse directives: %v", err)
	}
	if dirs.skipped {
		t.Skip("skipped by directive")
	}

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

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(workdir); err != nil {
		t.Fatalf("chdir to workdir: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(originalWD); chdirErr != nil {
			t.Fatalf("restore cwd: %v", chdirErr)
		}
	})

	restoreEnv := applyEnv(dirs.env)
	t.Cleanup(func() {
		if err := restoreEnv(); err != nil {
			t.Fatalf("restore env: %v", err)
		}
	})

	var stdin io.Reader
	var stdinFile *os.File
	if dirs.stdin != "" {
		stdinFile, err = os.Open(dirs.stdin)
		if err != nil {
			t.Fatalf("open stdin file %q: %v", dirs.stdin, err)
		}
		t.Cleanup(func() {
			stdinFile.Close() //nolint:errcheck // best-effort close in test cleanup
		})
		stdin = stdinFile
	} else {
		stdin = bytes.NewReader(nil)
	}

	runner := cli.App{}

	var stdoutBuf, stderrBuf bytes.Buffer

	exitCode := runner.Run(dirs.args, stdin, &stdoutBuf, &stderrBuf)

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

func applyEnv(vars []string) func() error {
	if len(vars) == 0 {
		return func() error { return nil }
	}

	originals := make(map[string]*string, len(vars))

	for _, kv := range vars {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		if prev, ok := os.LookupEnv(key); ok {
			copy := prev
			originals[key] = &copy
		} else {
			originals[key] = nil
		}
		_ = os.Setenv(key, val)
	}

	return func() error {
		var restoreErr error
		for key, val := range originals {
			var err error
			if val == nil {
				err = os.Unsetenv(key)
			} else {
				err = os.Setenv(key, *val)
			}
			if err != nil && restoreErr == nil {
				restoreErr = err
			}
		}

		return restoreErr
	}
}
