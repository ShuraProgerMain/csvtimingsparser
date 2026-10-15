// Package csvparse reads a CSV file and extracts the three columns the app
// cares about (Source In, Color, Notes), tolerating column reordering,
// header casing/spacing differences, and comma- or semicolon-delimited files.
package csvparse

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Row is one extracted record.
type Row struct {
	SourceIn string
	Color    string
	Notes    string
}

// Result holds everything the UI and CLI need after a parse.
type Result struct {
	Rows   []Row    // all rows, in file order
	Colors []string // distinct colors, first-seen order preserved
}

// canonical header names we look for, matched case-insensitively after
// collapsing internal whitespace.
const (
	hdrSourceIn = "source in"
	hdrColor    = "color"
	hdrNotes    = "notes"
)

// normalizeHeader lowercases and collapses whitespace so "Source  In",
// "source in" and "SOURCE IN" all compare equal.
func normalizeHeader(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}

// detectDelimiter sniffs the first non-empty line to pick ',' or ';'.
// Defaults to ',' when ambiguous.
func detectDelimiter(sample string) rune {
	for _, line := range strings.Split(sample, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.Count(line, ";") > strings.Count(line, ",") {
			return ';'
		}
		return ','
	}
	return ','
}

// Parse reads the file at path and returns the extracted result.
func Parse(path string) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()
	return parseReader(f)
}

// ParseBytes parses CSV content held in memory (e.g. a file dropped onto the
// GUI and read by the webview, where no path is available).
func ParseBytes(b []byte) (*Result, error) {
	return parseReader(bytes.NewReader(b))
}

func parseReader(r io.Reader) (*Result, error) {
	// Read a sample for delimiter detection, then re-chain it for the reader.
	br := bufio.NewReader(r)
	peek, _ := br.Peek(4096)
	delim := detectDelimiter(string(peek))

	cr := csv.NewReader(br)
	cr.Comma = delim
	cr.FieldsPerRecord = -1 // tolerate ragged rows
	cr.LazyQuotes = true
	cr.TrimLeadingSpace = true

	header, err := cr.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("csv is empty")
		}
		return nil, fmt.Errorf("read header: %w", err)
	}
	header[0] = strings.TrimPrefix(header[0], string(rune(0xFEFF))) // strip UTF-8 BOM

	idxSource, idxColor, idxNotes := -1, -1, -1
	for i, h := range header {
		switch normalizeHeader(h) {
		case hdrSourceIn:
			idxSource = i
		case hdrColor:
			idxColor = i
		case hdrNotes:
			idxNotes = i
		}
	}

	var missing []string
	if idxSource == -1 {
		missing = append(missing, "Source In")
	}
	if idxColor == -1 {
		missing = append(missing, "Color")
	}
	if idxNotes == -1 {
		missing = append(missing, "Notes")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing column(s): %s (found headers: %s)",
			strings.Join(missing, ", "), strings.Join(header, ", "))
	}

	res := &Result{}
	seen := map[string]bool{}
	at := func(rec []string, i int) string {
		if i >= 0 && i < len(rec) {
			return strings.TrimSpace(rec[i])
		}
		return ""
	}

	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}
		row := Row{
			SourceIn: at(rec, idxSource),
			Color:    at(rec, idxColor),
			Notes:    at(rec, idxNotes),
		}
		res.Rows = append(res.Rows, row)
		if row.Color != "" && !seen[row.Color] {
			seen[row.Color] = true
			res.Colors = append(res.Colors, row.Color)
		}
	}

	if len(res.Colors) == 0 {
		return res, fmt.Errorf("no color values found in Color column")
	}
	return res, nil
}

// RowsForColor returns rows whose Color equals color (case-insensitive).
func (r *Result) RowsForColor(color string) []Row {
	_, rows := r.RowsForColorIndexed(color)
	return rows
}

// RowsForColorIndexed returns the matching rows together with their index in
// Rows, so callers (e.g. the editor) can map a displayed row back to the
// underlying record.
func (r *Result) RowsForColorIndexed(color string) ([]int, []Row) {
	var idxs []int
	var out []Row
	target := strings.ToLower(strings.TrimSpace(color))
	for i, row := range r.Rows {
		if strings.ToLower(row.Color) == target {
			idxs = append(idxs, i)
			out = append(out, row)
		}
	}
	return idxs, out
}

// SetNote updates the Notes of the row at index i. Returns false if i is out of
// range. Used by the GUI's inline editor.
func (r *Result) SetNote(i int, notes string) bool {
	if i < 0 || i >= len(r.Rows) {
		return false
	}
	r.Rows[i].Notes = notes
	return true
}

// ShortTimecode formats a timecode for display/copy:
//   - drops the frames group (HH:MM:SS:FF -> HH:MM:SS);
//   - drops the hours group when it is zero, showing MM:SS instead.
//
// So 00:14:21:21 -> 14:21, 00:00:26:07 -> 00:26, but 01:02:03:04 -> 01:02:03.
// Values that aren't 3 or 4 colon-separated groups are returned unchanged.
func ShortTimecode(s string) string {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) == 4 { // drop frames
		parts = parts[:3]
	}
	if len(parts) == 3 { // HH:MM:SS — collapse to MM:SS when hours are zero
		if n, err := strconv.Atoi(parts[0]); err == nil && n == 0 {
			return parts[1] + ":" + parts[2]
		}
		return parts[0] + ":" + parts[1] + ":" + parts[2]
	}
	return s
}

// FormatList renders rows as the clipboard/CLI output: one line per row,
// "HH:MM:SS - Notes". No header, no color column.
func FormatList(rows []Row) string {
	var b strings.Builder
	for _, r := range rows {
		b.WriteString(ShortTimecode(r.SourceIn))
		b.WriteString(" - ")
		b.WriteString(r.Notes)
		b.WriteByte('\n')
	}
	return b.String()
}

// SortedColors returns the distinct colors sorted alphabetically (used where a
// stable display order is preferred over first-seen order).
func (r *Result) SortedColors() []string {
	out := append([]string(nil), r.Colors...)
	sort.Strings(out)
	return out
}
