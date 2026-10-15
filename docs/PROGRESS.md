# Progress tracker

Living status of the csvcolor app. Update this whenever a task changes state.

Legend: `[x]` done · `[~]` in progress · `[ ]` todo

## Core (pure Go, `CGO_ENABLED=0`)
- [x] CSV parser `internal/csvparse` — detects `Source In` / `Color` / `Notes`
      by header name (case-insensitive, whitespace-tolerant), auto-detects
      `,` vs `;`, strips UTF-8 BOM, tolerates ragged rows & column reordering
- [x] `ParseBytes` for in-memory content (GUI reads dropped file as text)
- [x] First-seen color ordering (`Result.Colors`) — drives the default selection
- [x] `RowsForColor` (case-insensitive filter) + `FormatList`
      (`<timecode> - Notes` lines) + `ShortTimecode` (drops frames; drops zero
      hours → `MM:SS`, else `HH:MM:SS`)
- [x] Clipboard `internal/clip` — `atotto/clipboard` (native Win32 UTF-16, no cgo)
- [x] Unit tests for parser (delimiter/header/BOM/missing-column/TSV) — passing

## CLI
- [x] `parse color "<Color>" <file.csv>` — prints TSV + copies to clipboard
- [x] `colors <file.csv>` — list distinct colors
- [x] `parse <file.csv>` — all rows · `--no-copy` · `-h/--help`
- [x] `path status|add|remove` — self-registers the exe's folder to the user
      PATH (Windows registry, no admin; manual guidance on other OSes)
- [x] Verified against real `MainNoAds.csv` (Cream/Cyan, Cyrillic intact)

## GUI — WebView2 + HTML/CSS Material Design
- [x] Frameless WebView2 window (`WS_POPUP`), Win11 rounded corners (DWM),
      kept in taskbar (`WS_EX_APPWINDOW`)
- [x] Resizable (`WS_THICKFRAME`) with minimum size 380×480 (layout verified
      at the minimum)
- [x] Clipboard/CLI format `HH:MM:SS - Notes` (frames dropped, no header/color)
- [x] Custom title bar: brand + minimize + close; draggable via `hostDrag`
- [x] Material 3 UI (hand-rolled CSS tokens, light + dark), consistent layout
- [x] Drop zone (HTML5 drop) + click-to-browse `<input type=file>`
- [x] Color dropdown with per-color swatches; default = first color in file
- [x] Row list (monospace timecode + notes) + full-width Copy button + snackbar
- [x] Inline Notes editing (contenteditable): fix typos, persists across color
      switches in-session, flows into the copied output; Esc reverts; does not
      write back to the CSV file
- [x] Go↔JS bindings: loadCSV / dataFor / editNote / copy / drag / minimize / close
- [x] Native launch verified (window opens, valid handle, closes clean)
- [x] UI visually verified in the browser pane (empty, loaded, dropdown states)
- [x] macOS / Linux GUI host (`ui_unix.go`, webview/webview_go, cgo) — reuses
      the same `index.html`; native OS title bar (custom frameless bar hidden via
      `CSVCOLOR_FRAMELESS=false`). **Not compiled/tested from the Windows dev box
      — needs verification on a Mac/Linux machine (cgo + WebKitGTK/Xcode).**
- [ ] Optional: remember last window position / size

## Build / tooling
- [x] Windows: pure Go, no C compiler (`CGO_ENABLED=0`); binary ~3.4 MB
- [x] macOS/Linux GUI: cgo + system WebKit (`CGO_ENABLED=1`); CLI stays pure Go
- [x] `build.ps1` (`-H windowsgui -s -w`); `./build.ps1 -Run`
- [x] **CI** (`.github/workflows/ci.yml`): vet+test+build on Windows/macOS/Linux
      every push/PR — this also finally verifies ui_unix.go compiles
- [x] **Release** (`.github/workflows/release.yml`): on a `v*` tag, builds and
      attaches Windows `.zip`, macOS `.app` zips (arm64+Intel), Linux `.tar.gz`
      to a GitHub Release — no local cross-compiling, users never build
- [x] Linux webkit2gtk 4.0→4.1 pkg-config shim (webview_go pins EOL 4.0)
- [x] Windows dual-mode: `-H windowsgui` + `AttachConsole` (clean GUI, CLI prints)
- [x] Removed the w64devkit toolchain (no longer needed after the WebView2 pivot)

## Docs
- [x] PROGRESS.md (this file), DECISIONS.md, ARCHITECTURE.md
- [x] README.md — full user guide: install, PATH setup, GUI + CLI guides,
      output format, CSV acceptance, build
- [x] Spec Kit spec under `specs/001-csvcolor/spec.md`

## Resolved
- [x] **Clipboard "garbage"** — was a PowerShell/MSYS display artifact; actual
      clipboard holds correct UTF-16 (verified). See DECISIONS.md › Issue log.
- [x] **`MainNoAds.csv`** provided and validated.
