package wc

import "fmt"

// Config captures the information derived from CLI arguments.
type Config struct {
	Files []string
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
	files := append([]string(nil), args...)
	return Config{Files: files}, nil
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
	return Stats{Name: name}, nil
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
		lines = append(lines, fmt.Sprintf("%d %d %d %s", st.Lines, st.Words, st.Bytes, st.Name))
	}
	return lines, nil
}
