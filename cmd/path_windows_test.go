//go:build windows

package cmd

import "testing"

func TestSameDir(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{`C:\tools\csvcolor`, `C:\tools\csvcolor`, true},
		{`C:\tools\csvcolor\`, `C:\tools\csvcolor`, true},  // trailing slash
		{`C:\Tools\CSVColor`, `c:\tools\csvcolor`, true},   // case-insensitive
		{` C:\tools\csvcolor `, `C:\tools\csvcolor`, true}, // surrounding spaces
		{`C:\tools\csvcolor\bin`, `C:\tools\csvcolor`, false},
	}
	for _, c := range cases {
		if got := sameDir(c.a, c.b); got != c.want {
			t.Errorf("sameDir(%q,%q)=%v want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestOnUserPath(t *testing.T) {
	raw := `%USERPROFILE%\bin;C:\tools\csvcolor\bin ;C:\Windows`
	if !onUserPath(`C:\tools\csvcolor\bin`, raw) {
		t.Fatal("expected dir to be detected on PATH")
	}
	if onUserPath(`C:\tools\other`, raw) {
		t.Fatal("unrelated dir should not be detected on PATH")
	}
	if len(pathEntries(raw)) != 3 {
		t.Fatalf("pathEntries = %d, want 3", len(pathEntries(raw)))
	}
	if len(pathEntries("")) != 0 {
		t.Fatal("empty PATH should yield no entries")
	}
}
