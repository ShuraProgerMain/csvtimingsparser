# csvcolor

[![CI](https://github.com/ShuraProgerMain/csvtimingsparser/actions/workflows/ci.yml/badge.svg)](https://github.com/ShuraProgerMain/csvtimingsparser/actions/workflows/ci.yml)

A tiny, portable tool that pulls the **Source In**, **Color**, and **Notes**
columns out of a CSV and copies the rows for a chosen color to your clipboard —
from a minimalist drag-and-drop window or from the command line.

It runs on **Windows, macOS, and Linux**. Prebuilt downloads are published for
every release — **you never build anything.**

---

## Contents
- [Download & run](#download--run)
- [Use it from any terminal (PATH)](#use-it-from-any-terminal-path)
- [Using the app (GUI guide)](#using-the-app-gui-guide)
- [Using the command line (CLI guide)](#using-the-command-line-cli-guide)
- [Output format](#output-format)
- [What it accepts in a CSV](#what-it-accepts-in-a-csv)
- [Building from source (contributors)](#building-from-source-contributors)
- [Project docs](#project-docs)

---

## Download & run

Grab the file for your OS from the **[Releases page][releases]**, unpack it, and
run it. There's no installer and nothing to compile.

[releases]: https://github.com/ShuraProgerMain/csvtimingsparser/releases/latest

### Windows
1. Download **`csvcolor-windows-amd64.zip`** and extract `csvcolor.exe`.
2. Put it wherever you keep small tools (e.g. `C:\Tools\csvcolor\`). Moving it
   later is fine.
3. Run it: **double-click** for the GUI, or run `csvcolor ...` in a terminal
   opened in that folder.

Needs the **WebView2** runtime — built into Windows 11 (and installed with Edge
on Windows 10). A first download may show a SmartScreen prompt (unsigned app):
**More info → Run anyway**.

### macOS (Apple Silicon)
1. Download **`csvcolor-macos-arm64.zip`** and unzip to get **`CSV Color.app`**.
   _(Intel Macs aren't supported by this build yet — an arm64 app can't run on
   Intel. Build from source until an Intel build is published.)_
2. Move it to Applications (optional). First launch: **right-click → Open**, then
   **Open** in the dialog (the app is ad-hoc-signed but not notarized).
3. If macOS says the app **"is damaged"** (it isn't — that's Gatekeeper on an
   unsigned download), clear the quarantine flag once:
   ```bash
   xattr -dr com.apple.quarantine "/Applications/CSV Color.app"   # adjust the path
   ```
   then open it normally.
4. For the CLI, the binary is inside the app at
   `CSV Color.app/Contents/MacOS/csvcolor`.

### Linux
1. Download **`csvcolor-linux-amd64.tar.gz`** and extract: `tar -xzf csvcolor-linux-amd64.tar.gz`.
2. Run `./csvcolor` (GUI) or `./csvcolor parse ...` (CLI).
3. The GUI needs GTK 3 + WebKitGTK and the clipboard needs `xclip`/`xsel`
   (`wl-copy` on Wayland). On most desktops these are already present; if not,
   e.g. Debian/Ubuntu: `sudo apt install libwebkit2gtk-4.1-0 xclip`.

> To run it from **any** folder in the terminal — not just where the file is —
> see the next section.

---

## Use it from any terminal (PATH)

So you can type `csvcolor` in any terminal, add its folder to your PATH. The
program can do this to **itself** — no separate installer, no admin rights:

```bash
csvcolor path add
```

This adds the folder containing `csvcolor.exe` to your **user** PATH. Then
**open a new terminal** and `csvcolor` works from anywhere:

```bash
csvcolor path status     # check whether it's on PATH (prints the folder)
csvcolor path add        # add it to your user PATH
csvcolor path remove     # undo the above
```

Notes:
- It changes only your **user** PATH (HKCU) — no administrator rights, and it
  won't touch the system-wide PATH.
- Already-open terminals keep the old PATH; open a new one after `path add`.
- If you move `csvcolor.exe` later, run `csvcolor path remove` in the old spot
  (or just `path add` from the new spot and remove the stale entry).

<details>
<summary>Prefer to do it by hand?</summary>

Add the folder containing `csvcolor.exe` to PATH yourself:
Windows Settings → *Edit environment variables for your account* → select
**Path** → **New** → paste the folder → OK. Open a new terminal afterwards.
</details>

---

## Using the app (GUI guide)

1. **Open** `csvcolor.exe` (double-click, or run it with no arguments). A small
   frameless window appears. Drag its title bar to move it; resize from the
   edges.
2. **Load a CSV** — drag a `.csv` file onto the window, *or* click the drop zone
   to browse for one.
3. **Pick a color** — the dropdown fills with the colors found in your file and
   starts on the first one. Each color shows a matching swatch. Choosing a color
   filters the list to its rows.
4. **Read the list** — each row shows its timecode and note.
5. **Fix a typo (optional)** — click a note, edit it, then press **Enter** or
   click away (**Esc** cancels). Edits stay for the session and are included
   when you copy. They are **not** written back to your CSV file.
6. **Copy** — click **Copy to clipboard**. The matching rows land on your
   clipboard as `timecode - Notes` lines, ready to paste. A toast confirms how
   many rows were copied.

Window controls live in the custom title bar: **–** minimizes, **×** closes.

---

## Using the command line (CLI guide)

Run `csvcolor <command>`. If it's on your PATH you can run it from anywhere;
otherwise run it from its folder (`.\csvcolor.exe` in PowerShell).

### Copy the rows for one color
```bash
csvcolor parse color "Cyan" "C:\path\to\MainNoAds.csv"
```
Prints the matching rows **and copies them to the clipboard**:
```
14:21 - a note for this cue
19:02 - another note
```

### List the colors in a file
```bash
csvcolor colors "C:\path\to\MainNoAds.csv"
```
```
Cream
Cyan
```

### Print every row (all colors)
```bash
csvcolor parse "C:\path\to\MainNoAds.csv"
```

### Don't touch the clipboard
Add `--no-copy` to any `parse` command to only print:
```bash
csvcolor parse color "Cyan" MainNoAds.csv --no-copy
```

### Help
```bash
csvcolor --help
```

Tips:
- Quote the color and the file path if they contain spaces.
- Color matching is case-insensitive (`cyan` == `Cyan`).
- Exit code is `0` on success, `1` on an error (e.g. color not found — the error
  lists the available colors), `2` on bad usage.

---

## Output format

Both the GUI copy and the CLI produce the same thing: **one line per matching
row**, `timecode - Notes`, with no header and no color column.

The timecode drops the frame group, and drops the hours group when it's zero:

| Source In in the file | Output |
|---|---|
| `00:00:26:07` | `00:26` |
| `00:14:21:21` | `14:21` |
| `01:02:03:04` | `01:02:03` |

So you get `MM:SS` for short cues and `HH:MM:SS` once there are hours.

---

## What it accepts in a CSV

It finds the three columns it needs by **header name**, so your file can have
any number of other columns in any order. It tolerates:
- Column headers in any case/spacing (`Source In`, `source  in`, `SOURCE IN`).
- `,` or `;` as the delimiter (auto-detected).
- A UTF-8 byte-order mark, `CRLF` line endings, quoted fields, and short rows.
- Unicode notes (e.g. Cyrillic) — copied correctly.

If `Source In`, `Color`, or `Notes` is missing, it tells you which, and lists
the headers it did find.

---

## Building from source (contributors)

**End users don't need this** — download a prebuilt binary from
[Releases][releases]. This section is for working on the code. Requires Go 1.27+.

### Windows
Pure Go, no C compiler:
```powershell
./build.ps1          # -> bin/csvcolor.exe
./build.ps1 -Run     # build, then launch the GUI
# or: $env:CGO_ENABLED="0"; go build -ldflags "-H windowsgui -s -w" -o bin/csvcolor.exe .
```

### macOS / Linux
The **CLI** needs nothing special:
```bash
CGO_ENABLED=0 go build -o csvcolor .
```

The **GUI** on macOS/Linux uses the system WebKit through cgo, so it needs a C
compiler and (on Linux) the WebKitGTK dev package:

- **macOS:** install the Xcode command-line tools (`xcode-select --install`),
  then:
  ```bash
  CGO_ENABLED=1 go build -o csvcolor .
  ```
- **Linux:** install a C compiler and WebKitGTK, e.g. on Debian/Ubuntu:
  ```bash
  sudo apt install build-essential libwebkit2gtk-4.1-dev   # or libwebkit2gtk-4.0-dev
  CGO_ENABLED=1 go build -o csvcolor .
  ```
  The clipboard on Linux uses `xclip` or `xsel` (or `wl-copy` on Wayland) — install one.

On macOS/Linux the window keeps the **native OS title bar** (the custom frameless
title bar is Windows-only); everything else — drag-drop, the color picker, inline
note editing, copy — works the same.

### Releasing (maintainers)
Releases are built by GitHub Actions — no local cross-compiling. CI
(`.github/workflows/ci.yml`) builds and tests on Windows, macOS, and Linux on
every push/PR. To cut a release, push a version tag:
```bash
git tag v1.0.0
git push origin v1.0.0
```
`.github/workflows/release.yml` then builds the Windows `.zip`, the macOS `.app`
zip (Apple Silicon), and the Linux `.tar.gz`, and attaches them to a new GitHub
Release. (Intel macOS is omitted for now — its runners queue for a long time.)

---

## Project docs
- [docs/PROGRESS.md](docs/PROGRESS.md) — what's done / in progress / todo
- [docs/DECISIONS.md](docs/DECISIONS.md) — stack decisions & dev rules
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) — layout & data flow
- [specs/001-csvcolor/spec.md](specs/001-csvcolor/spec.md) — specification
