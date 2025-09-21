package wc

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type counterKind string

const (
	counterLines counterKind = "lines"
	counterWords counterKind = "words"
	counterBytes counterKind = "bytes"
	counterChars counterKind = "chars"
)

// Config captures the information derived from CLI arguments.
type Config struct {
	Files        []string
	CountBytes   bool
	CountLines   bool
	CountWords   bool
	CountChars   bool
	counterOrder []counterKind
}

func (cfg *Config) addCounter(kind counterKind) {
	for _, existing := range cfg.counterOrder {
		if existing == kind {
			return
		}
	}
	cfg.counterOrder = append(cfg.counterOrder, kind)
}

// Stats represents the results for a single analyzed input.
type Stats struct {
	Name  string
	Lines int
	Words int
	Bytes int
	Chars int
}

// ParseArgs converts CLI arguments into a configuration struct.
func ParseArgs(args []string) (Config, error) {
	cfg := Config{}
	for _, arg := range args {
		switch arg {
		case "-c", "--bytes":
			cfg.CountBytes = true
			cfg.addCounter(counterBytes)
			continue
		case "-l", "--lines":
			cfg.CountLines = true
			cfg.addCounter(counterLines)
			continue
		case "-w", "--words":
			cfg.CountWords = true
			cfg.addCounter(counterWords)
			continue
		case "-m", "--chars":
			cfg.CountChars = true
			cfg.addCounter(counterChars)
			continue
		}
		if strings.HasPrefix(arg, "-") && arg != "-" {
			return Config{}, fmt.Errorf("unsupported flag: %s", arg)
		}
		cfg.Files = append(cfg.Files, arg)
	}
	return cfg, nil
}

// AnalyzeFiles collects Stats for each configured file.
func AnalyzeFiles(cfg Config) ([]Stats, error) {
	stats := make([]Stats, 0, len(cfg.Files))
	for _, name := range cfg.Files {
		stat, err := AnalyzeFile(name)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

// AnalyzeFile returns the Stats for a single file.
func AnalyzeFile(name string) (Stats, error) {
	stat := Stats{Name: name}
	data, err := os.ReadFile(name)
	if err != nil {
		return Stats{}, err
	}
	stat.Bytes = len(data)
	stat.Lines = countLines(data)
	stat.Words = countWords(data)
	stat.Chars = countChars(data)
	return stat, nil
}

// AddTotal appends a synthetic total entry when appropriate.
func AddTotal(cfg Config, stats []Stats) ([]Stats, error) {
	cloned := append([]Stats(nil), stats...)
	return cloned, nil
}

// Format renders Stats entries into lines ready for stdout.
func Format(cfg Config, stats []Stats) ([]string, error) {
	lines := make([]string, 0, len(stats))
	for _, st := range stats {
		parts := make([]string, 0, 4)
		if cfg.CountLines {
			parts = append(parts, fmt.Sprintf("%8d", st.Lines))
		}
		if cfg.CountWords {
			parts = append(parts, fmt.Sprintf("%8d", st.Words))
		}
		if cfg.CountChars {
			parts = append(parts, fmt.Sprintf("%8d", st.Chars))
		}
		if cfg.CountBytes {
			parts = append(parts, fmt.Sprintf("%8d", st.Bytes))
		}
		if len(parts) == 0 {
			lines = append(lines, fmt.Sprintf("0 0 0 %s", st.Name))
			continue
		}
		var counts strings.Builder
		for _, part := range parts {
			counts.WriteString(part)
		}
		lines = append(lines, fmt.Sprintf("%s %s", counts.String(), st.Name))
	}
	return lines, nil
}

func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	count := 0
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}
	if data[len(data)-1] != '\n' {
		count++
	}
	return count
}

func countWords(data []byte) int {
	return len(bytes.Fields(data))
}

func countChars(data []byte) int {
	count := 0
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r != utf8.RuneError || size != 1 {
			count++
		}
		data = data[size:]
	}
	return count
}
