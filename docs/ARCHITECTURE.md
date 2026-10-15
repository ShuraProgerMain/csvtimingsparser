# Architecture

```
main.go                  entry point: args present -> CLI, else -> GUI
platform_windows.go      AttachConsole so the windowsgui binary prints in a terminal
platform_other.go        no-op AttachConsole for non-Windows
cmd/cli.go               CLI command parsing (parse / colors / path / flags)
cmd/path_windows.go      'path add/remove/status': self-register exe dir on PATH
cmd/path_other.go        'path' guidance for non-Windows (manual shell profile)
internal/csvparse/       CSV reading, column detection, filtering, TSV output
internal/clip/           clipboard write (atotto/clipboard wrapper)
internal/ui/
  ui_host.go             shared Go<->JS bindings (loadCSV/dataFor/editNote/copy)
  ui_assets.go           go:embed of assets/index.html (shared)
  ui_windows.go          frameless WebView2 host (go-webview2, pure Go)
  ui_unix.go             macOS/Linux host (webview/webview_go, cgo; native chrome)
  ui_other.go            stub Run() on other platforms (CLI still works)
  assets/index.html      the Material Design front-end (HTML/CSS/JS, embedded)
testdata/sample.csv      fixture used by tests
```

## GUI = Go host + HTML front-end
The UI is a single webview window. The entire UI is `assets/index.html`,
embedded with `go:embed` and handed to the webview via `SetHtml`. The host
differs per platform but the HTML and the bindings are identical:

- **Windows** (`ui_windows.go`): `go-webview2` (WebView2, built into Win11),
  **pure Go, no cgo**. Frameless + rounded + custom draggable title bar.
- **macOS / Linux** (`ui_unix.go`): `webview/webview_go` (WKWebView / WebKitGTK),
  **cgo**. Keeps the native OS title bar; the page hides its own title bar because
  the host sets `window.CSVCOLOR_FRAMELESS=false` via `Init` before load.

Go exposes a few bindings the page calls like async functions:

| JS call | Go side | Purpose |
|---|---|---|
| `hostLoadCSV(text)` | `csvparse.ParseBytes` | parse dropped/opened file content, return colors |
| `hostDataFor(color)` | `RowsForColorIndexed` + `FormatList` | rows (with row index) + clipboard text |
| `hostEditNote(idx, text)` | `Result.SetNote` | inline Notes edit, persisted in memory |
| `hostCopy(text)` | `clip.Copy` | put text on the clipboard |
| `hostDrag()` | `WM_NCLBUTTONDOWN`/`HTCAPTION` | drag the frameless window |
| `hostMinimize()` / `hostClose()` | `ShowWindow` / `Terminate` | window controls |

The file is read **in the page** (`FileReader.readAsText`) and its text is sent
to Go — so no file path is needed, which also sidesteps WebView2's lack of
`File.path`. Drag-and-drop uses the standard HTML5 `drop` event; the same drop
zone doubles as a click-to-browse `<input type=file>`.

### Inline Notes editing
Each note cell is `contenteditable`. Committing (blur or Enter) calls
`hostEditNote(idx, text)`, which updates `Result.Rows[idx].Notes` in memory, so
the edit survives switching colors within the session (it does **not** write
back to the CSV file). The Copy button builds its text from the live DOM, so it
always reflects the current edits; `Esc` reverts a cell, and pasted text is
flattened to a single line.

### Frameless window (Windows)
`go-webview2` creates a normal `WS_OVERLAPPEDWINDOW`. After creation
`makeFrameless` restyles the HWND to `WS_POPUP` (+ `WS_EX_APPWINDOW` so it stays
in the taskbar) and calls `DwmSetWindowAttribute(..., DWMWCP_ROUND)` for Win11
rounded corners. The style also includes `WS_THICKFRAME` so the frameless window
is **resizable**, with a minimum size set via `SetSize(380, 480, HintMin)`. A
custom title bar in the HTML provides the brand, minimize, and close controls;
the draggable region calls `hostDrag()`.

## Data flow
1. A CSV arrives: GUI drop/browse (as text) or a CLI path argument.
2. `csvparse.Parse` / `ParseBytes` -> `*Result{ Rows, Colors }`.
   - `Colors` is first-seen order; the GUI defaults the selector to `Colors[0]`.
3. A color is chosen (GUI dropdown, or CLI arg).
4. `RowsForColor` filters case-insensitively.
5. `FormatList` renders the shared output: one `<timecode> - Notes` line per row
   (`ShortTimecode` drops frames, and drops zero hours → `MM:SS` else `HH:MM:SS`;
   no header, no color).
6. GUI shows the rows and copies that text on the Copy button; CLI prints it and
   copies it unless `--no-copy`.

## Why the split
- `csvparse` + `clip` are pure Go and unit-tested without a GUI.
- On **Windows** the whole app builds with `CGO_ENABLED=0` (go-webview2 is pure
  Go) — no C compiler, ~3.4 MB binary.
- On **macOS/Linux** only the GUI needs cgo (system WebKit); the CLI still builds
  pure Go. The `host` bindings and `index.html` are shared across all platforms.
