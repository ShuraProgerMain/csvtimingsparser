# Decisions & development rules

Short record of why the app is built the way it is, so future work doesn't
re-open settled questions. Newest decisions at the bottom.

## Development rules (follow these)
1. **Keep Windows and the CLI pure Go (`CGO_ENABLED=0`).** `csvparse`, `clip`,
   and the Windows GUI (go-webview2) build with no C compiler. Only the
   macOS/Linux GUI host (`ui_unix.go`, webview_go) uses cgo — keep cgo isolated
   to that file so Windows and the cross-platform CLI never require a compiler.
2. **One source of truth for output.** Both GUI and CLI render rows through
   `csvparse.FormatList` — one `HH:MM:SS - Notes` line per row (no header, no
   color column). Timecodes go through `ShortTimecode`: frames dropped, and the
   hours group dropped when zero (`00:14:21:21` → `14:21`, `01:02:03:04` →
   `01:02:03`). Never format rows ad hoc.
   The GUI copies the text returned by `hostDataFor`, not the on-screen text,
   and the on-screen timecodes are trimmed to match the copy.
3. **Parse defensively.** Real CSVs vary: tolerate column reordering, header
   casing/spacing, BOM, `,`/`;` delimiters, ragged rows. Match columns by name,
   never by fixed index.
4. **Keep the UI minimal.** No controls beyond: a draggable title bar
   (brand + minimize + close), drop zone, color selector, the list, and one
   Copy button. Resist feature creep.
5. **The UI lives in `internal/ui/assets/index.html`** (embedded). Keep it
   self-contained — no network/CDN assets — so the app works offline and stays
   lightweight. Material 3 styling is hand-rolled via CSS tokens.
6. **Update `docs/PROGRESS.md`** whenever a task's state changes.
7. **Tests first for core.** Any parser change ships with a test.

## Decision log
- **Language: Go.** Required by the user; single static binaries, good
  cross-platform story.
- **Clipboard: `atotto/clipboard`** (not `clip.exe`). Uses native Win32 APIs, so
  Cyrillic / Unicode notes copy correctly; pure Go, no cgo.
- **GUI toolkit: WebView2 via `go-webview2` (final).** Chosen after the design
  requirements grew to include a frameless window, a *draggable* custom title
  bar, and true consistent Material Design. It is pure Go on Windows (no C
  compiler), produces a ~3.3 MB binary, renders an HTML/CSS Material UI with
  pixel-level consistency, and handles OS file drops via HTML5. WebView2 ships
  with Windows 11. The HTML front-end is portable to a macOS/Linux
  `webview/webview` host (cgo) later.
  - **macOS/Linux GUI added via `webview/webview_go`.** Same `index.html`, behind
    build tags, using the system WebKit (cgo). The frameless/custom-title-bar
    chrome stays Windows-only (it relies on Win32); mac/linux use the native OS
    title bar, selected at runtime via `CSVCOLOR_FRAMELESS`. No pure-Go,
    cross-platform GUI with OS file-drop exists, so cgo on mac/linux is
    unavoidable; it is confined to `ui_unix.go`.
  - **Fyne was used first, then dropped.** It supports OS file drops and single
    binaries, but: (a) its desktop driver needs cgo + OpenGL + a C compiler;
    (b) its `Window` interface has **no window-move API**, so a borderless
    window can't be dragged by a custom title bar; (c) its widgets resist a
    genuine Material look. All three conflict with the redesign goals.
  - **Gio was evaluated and rejected.** Builds without a C compiler on Windows,
    but v0.10.3's Windows backend implements no `WM_DROPFILES` / `IDropTarget`
    (its `io/transfer` DnD is in-app only), so it cannot receive a file dropped
    from Explorer — fails the central feature.
  - **Wails rejected.** Needs a C compiler to build; no lighter than
    `go-webview2`, which targets the same WebView2 runtime.
- **No build toolchain needed.** With `go-webview2` the app builds at
  `CGO_ENABLED=0`. (A portable w64devkit was installed for the Fyne attempt and
  then removed once the WebView2 pivot made cgo unnecessary.)
- **Portable, no installer.** The deliverable is one self-contained `.exe`; there
  is no setup program and the user need not build anything. Rather than ship an
  installer, the program registers *itself* on the user's PATH on request
  (`csvcolor path add`, user-scope registry, no admin), and can undo it. This
  keeps a tiny tool a tiny tool.

## Issue log
- **"Clipboard copies garbage" (reported, resolved).** The CLI/GUI write via
  `atotto/clipboard`, which uses `CF_UNICODETEXT` (UTF-16) — verified correct by
  dumping the clipboard to a UTF-8 file: full Cyrillic Notes come through intact.
  The `?` seen earlier was a display artifact of piping PowerShell `Get-Clipboard`
  through the MSYS console (OEM codepage), not the actual clipboard contents. The
  GUI copies the exact `FormatTSV` output via the same `clip.Copy`.
