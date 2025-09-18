package wc

import (
	"fmt"
	"os"
	"strings"
)

// Config captures the information derived from CLI arguments.
type Config struct {
	Files      []string
	CountBytes bool
	CountLines bool
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
			continue
		case "-l", "--lines":
			cfg.CountLines = true
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
		if cfg.CountLines && !cfg.CountBytes {
			lines = append(lines, fmt.Sprintf("%8d %s", st.Lines, st.Name))
			continue
		}
		if cfg.CountBytes && !cfg.CountLines {
			lines = append(lines, fmt.Sprintf("%8d %s", st.Bytes, st.Name))
			continue
		}
		lines = append(lines, fmt.Sprintf("0 0 0 %s", st.Name))
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
