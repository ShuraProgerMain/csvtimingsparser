# csvcolor

A tiny, portable tool that pulls the **Source In**, **Color**, and **Notes**
columns out of a CSV and copies the rows for a chosen color to your clipboard —
from a minimalist drag-and-drop window or from the command line.

On **Windows** it's a **single ~3.4 MB `.exe`** — no installer, no runtime, no
build step; copy it anywhere and run it. The **CLI runs on Windows, macOS, and
Linux**. The **GUI runs on all three too**, but macOS/Linux builds are compiled
from source (they use the system WebKit via cgo) — see
[Building](#building-from-source).

---

## Contents
- [Install (just the one file)](#install-just-the-one-file)
- [Use it from any terminal (PATH)](#use-it-from-any-terminal-path)
- [Using the app (GUI guide)](#using-the-app-gui-guide)
- [Using the command line (CLI guide)](#using-the-command-line-cli-guide)
- [Output format](#output-format)
- [What it accepts in a CSV](#what-it-accepts-in-a-csv)
- [Building from source (optional)](#building-from-source-optional)
- [Project docs](#project-docs)

---

## Install (just the one file)

### Windows
There is nothing to install. `csvcolor` is one self-contained executable.

1. Get **`csvcolor.exe`** (it's in this project's `bin\` folder after a build,
   or whatever folder someone shared it in).
2. Put it wherever you like to keep small tools, e.g. `C:\Tools\csvcolor\` or a
   folder in your user profile. That folder is the "install location" — moving
   the file later is fine.
3. **Run it:**
   - **GUI** — double-click `csvcolor.exe`.
   - **CLI** — open a terminal in that folder and run `csvcolor ...` (or
     `.\csvcolor.exe ...` in PowerShell).

Requirements: Windows 10/11. The GUI uses the **WebView2** runtime, which is
built into Windows 11 (and installed with Edge on Windows 10). Nothing else is
needed.

### macOS / Linux
Build it once from source (the GUI uses the system WebKit) — see
[Building from source](#building-from-source) — then run `./csvcolor`. The CLI
needs no special setup; only the GUI build needs a C compiler / WebKitGTK.

> To use it from **any** folder in the terminal — not just where the file is —
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

## Building from source (optional)

You only need this if you're changing the code. The app is **pure Go** — no C
compiler required.

Requires Go 1.27+.
```powershell
# from the project folder
./build.ps1          # -> bin/csvcolor.exe
./build.ps1 -Run     # build, then launch the GUI
```
Or directly:
```powershell
$env:CGO_ENABLED = "0"
go build -ldflags "-H windowsgui -s -w" -o bin/csvcolor.exe .
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

> Note: these builds are produced on the target OS (or with a matching
> cross-toolchain). You can't cross-compile the cgo GUI from Windows without one.

---

## Project docs
- [docs/PROGRESS.md](docs/PROGRESS.md) — what's done / in progress / todo
- [docs/DECISIONS.md](docs/DECISIONS.md) — stack decisions & dev rules
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) — layout & data flow
- [specs/001-csvcolor/spec.md](specs/001-csvcolor/spec.md) — specification
