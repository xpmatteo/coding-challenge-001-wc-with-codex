package format

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"wc/internal/pipeline/counters"
)

// Column identifies a rendered metric.
type Column int

const (
	ColumnLines Column = iota
	ColumnWords
	ColumnBytes
	ColumnRunes
	ColumnMaxLine
)

// Row represents a single output row.
type Row struct {
	Label  string
	Counts counters.Snapshot
}

// Render prints aligned rows to writer.
func Render(w io.Writer, columns []Column, rows []Row) error {
	if len(rows) == 0 || len(columns) == 0 {
		return nil
	}

	widths := make(map[Column]int, len(columns))
	for _, col := range columns {
		maxLen := 0
		for _, row := range rows {
			width := len(strconv.Itoa(valueFor(col, row.Counts)))
			if width > maxLen {
				maxLen = width
			}
		}
		if maxLen < 7 {
			maxLen = 7
		}
		widths[col] = maxLen
	}

	for _, row := range rows {
		var builder strings.Builder
		for _, col := range columns {
			val := valueFor(col, row.Counts)
			fmt.Fprintf(&builder, "%*d", widths[col], val)
		}
		if row.Label != "" {
			builder.WriteByte(' ')
			builder.WriteString(row.Label)
		}
		builder.WriteByte('\n')
		if _, err := io.WriteString(w, builder.String()); err != nil {
			return err
		}
	}
	return nil
}

func valueFor(col Column, snap counters.Snapshot) int {
	switch col {
	case ColumnLines:
		return snap.Lines
	case ColumnWords:
		return snap.Words
	case ColumnBytes:
		return snap.Bytes
	case ColumnRunes:
		return snap.Runes
	case ColumnMaxLine:
		return snap.MaxLineLen
	default:
		return 0
	}
}
