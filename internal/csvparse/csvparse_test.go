package csvparse

import (
	"strings"
	"testing"
)

func TestParseSample(t *testing.T) {
	res, err := Parse("../../testdata/sample.csv")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(res.Rows) != 6 {
		t.Fatalf("rows = %d, want 6", len(res.Rows))
	}
	// First-seen color order: Red, Blue, Green.
	want := []string{"Red", "Blue", "Green"}
	if strings.Join(res.Colors, ",") != strings.Join(want, ",") {
		t.Fatalf("colors = %v, want %v", res.Colors, want)
	}
}

func TestRowsForColorCaseInsensitive(t *testing.T) {
	res, err := Parse("../../testdata/sample.csv")
	if err != nil {
		t.Fatal(err)
	}
	rows := res.RowsForColor("red")
	if len(rows) != 3 {
		t.Fatalf("red rows = %d, want 3", len(rows))
	}
	if rows[0].SourceIn != "00:00:12:05" || rows[0].Notes != "Opening shot" {
		t.Fatalf("unexpected first row: %+v", rows[0])
	}
}

func TestFormatList(t *testing.T) {
	rows := []Row{
		{SourceIn: "00:01:03:11", Color: "Red", Notes: "Wide shot"},
		{SourceIn: "01:02:03:04", Color: "Red", Notes: "with hours"},
	}
	got := FormatList(rows)
	want := "01:03 - Wide shot\n01:02:03 - with hours\n"
	if got != want {
		t.Fatalf("FormatList = %q, want %q", got, want)
	}
}

func TestShortTimecode(t *testing.T) {
	cases := map[string]string{
		"00:00:00:00": "00:00", // zero hours -> MM:SS
		"00:14:21:21": "14:21", // zero hours -> MM:SS
		"00:00:26:07": "00:26", // only seconds -> MM:SS
		"01:02:03:04": "01:02:03",
		"10:00:00":    "10:00:00",
		"00:00:00":    "00:00",
		"12:34":       "12:34", // already MM:SS, unchanged
		"weird":       "weird",
	}
	for in, want := range cases {
		if got := ShortTimecode(in); got != want {
			t.Fatalf("ShortTimecode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDelimiterAndHeaderTolerance(t *testing.T) {
	// semicolon-delimited, reordered, odd casing/spacing, BOM prefix.
	in := string(rune(0xFEFF)) + "Notes;COLOR;source  in\n" +
		"hello;Red;00:01\n" +
		"world;blue;00:02\n"
	res, err := parseReader(strings.NewReader(in))
	if err != nil {
		t.Fatalf("parseReader: %v", err)
	}
	if len(res.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(res.Rows))
	}
	if res.Rows[0].SourceIn != "00:01" || res.Rows[0].Color != "Red" || res.Rows[0].Notes != "hello" {
		t.Fatalf("row0 mismatch: %+v", res.Rows[0])
	}
	if strings.Join(res.Colors, ",") != "Red,blue" {
		t.Fatalf("colors = %v", res.Colors)
	}
}

func TestMissingColumn(t *testing.T) {
	_, err := parseReader(strings.NewReader("Source In,Notes\na,b\n"))
	if err == nil || !strings.Contains(err.Error(), "Color") {
		t.Fatalf("expected missing Color error, got %v", err)
	}
}
