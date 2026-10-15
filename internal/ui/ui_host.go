// Package ui hosts the GUI: a WebView window rendering the Material Design
// HTML/CSS front-end. The Go side exposes CSV parsing, clipboard, and
// window-control bindings; it is a thin shell over the pure-Go core.
//
// This file holds the platform-neutral bindings shared by every webview host
// (Windows: ui_windows.go; macOS/Linux: ui_unix.go).
package ui

import (
	"sync"

	"csvcolor/internal/clip"
	"csvcolor/internal/csvparse"
)

// host carries the parsed CSV between binding calls. It does not hold the
// webview handle, because that type differs per platform.
type host struct {
	mu  sync.Mutex
	res *csvparse.Result
}

type loadResult struct {
	Colors []string `json:"colors"`
	Error  string   `json:"error"`
}

// loadCSV parses file content (read in the page) and returns the colors found.
func (h *host) loadCSV(content string) loadResult {
	res, err := csvparse.ParseBytes([]byte(content))
	if err != nil {
		return loadResult{Error: err.Error()}
	}
	h.mu.Lock()
	h.res = res
	h.mu.Unlock()
	return loadResult{Colors: res.Colors}
}

type rowDTO struct {
	Idx      int    `json:"idx"` // index in Result.Rows, for inline edits
	SourceIn string `json:"sourceIn"`
	Color    string `json:"color"`
	Notes    string `json:"notes"`
}

type dataResult struct {
	Rows []rowDTO `json:"rows"`
	Copy string   `json:"copy"` // exact clipboard text: "<timecode> - Notes" lines
}

// dataFor returns the rows for a color (with original indices for editing) plus
// the exact clipboard text.
func (h *host) dataFor(color string) dataResult {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.res == nil {
		return dataResult{}
	}
	idxs, rows := h.res.RowsForColorIndexed(color)
	out := dataResult{Copy: csvparse.FormatList(rows)}
	for k, r := range rows {
		out.Rows = append(out.Rows, rowDTO{
			Idx:      idxs[k],
			SourceIn: csvparse.ShortTimecode(r.SourceIn), // display matches the copy
			Color:    r.Color,
			Notes:    r.Notes,
		})
	}
	return out
}

// editNote updates a row's Notes in memory (persists for the session, incl.
// across color switches). Returns true on success.
func (h *host) editNote(idx int, notes string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.res == nil {
		return false
	}
	return h.res.SetNote(idx, notes)
}

// copy puts text on the system clipboard. Returns "" on success or the error.
func (h *host) copy(text string) string {
	if err := clip.Copy(text); err != nil {
		return err.Error()
	}
	return ""
}
