# Feature: csvcolor — color-filtered CSV column extractor

## Summary
Extract the `Source In`, `Color`, and `Notes` columns from a CSV and deliver the
rows for a selected color to the clipboard, via a minimalist drag-and-drop GUI
and a CLI.

## User stories
1. As an editor, I drag `MainNoAds.csv` onto the window and immediately see a
   list of the relevant columns so I can copy them with one click.
2. As an editor, after a scan a color selector appears, offering exactly the
   colors present in the file, defaulting to the first color in the file.
3. As a power user, I run `csvcolor parse color "Red" file.csv` and get the rows
   printed to the console and copied to my clipboard in one step.

## Functional requirements
- **FR1** Read a CSV and locate the three columns by header name
  (`Source In`, `Color`, `Notes`), case- and whitespace-insensitive, regardless
  of column order. Error clearly if any is missing.
- **FR2** Auto-detect `,` vs `;` delimiter; strip a UTF-8 BOM; tolerate ragged
  rows.
- **FR3** Offer the distinct colors in first-seen file order; the GUI selector
  defaults to the first.
- **FR4** Filter rows by the selected color (case-insensitive).
- **FR5** Output/clipboard format is one line per row: `<timecode> - Notes`,
  no header, no color column. The timecode drops the frames group, and drops the
  hours group when it is zero: `HH:MM:SS` when there are hours, otherwise
  `MM:SS`. GUI display and copy, and the CLI, all use the same format.
- **FR9** The GUI window is resizable, with a minimum size (380×480) below which
  it cannot shrink, so the layout always renders correctly.
- **FR10** Notes are editable inline in the list (to fix typos). Edits persist
  for the session (across color switches) and are included in the copied output.
  Edits are not written back to the source CSV file.
- **FR11** The app is a single portable binary needing no installation. It can
  register its own folder on the user's PATH via `csvcolor path add` (user
  scope, no admin) and undo it with `path remove`.
- **FR6** GUI: a single Copy action puts the filtered output on the clipboard.
- **FR7** CLI: `parse color "<c>" <file>` prints and copies; `colors <file>`
  lists colors; `parse <file>` outputs all rows; `--no-copy` suppresses copying.
- **FR8** Clipboard writes handle Unicode (e.g. Cyrillic) correctly.

## Non-functional requirements
- **NFR1** Cross-platform: Windows, macOS, Linux. CLI runs on all three with no
  dependencies. GUI runs on all three (Windows: prebuilt, pure Go; macOS/Linux:
  built from source using the system WebKit via cgo).
- **NFR2** Lightweight: single binary; core logic has no cgo dependency.
- **NFR3** Minimalist UI with rounded corners; no controls beyond drop zone,
  color selector, list view, and one Copy button.

## Out of scope
- Editing the CSV, writing files, multi-file batching, column remapping UI,
  persistent settings.

## Acceptance
- Dropping `testdata/sample.csv` shows colors `Red, Blue, Green` (default `Red`),
  and Red yields 3 rows.
- `csvcolor colors testdata/sample.csv` prints `Red / Blue / Green`.
- `csvcolor parse color "Red" testdata/sample.csv` prints 3 `HH:MM:SS - Notes`
  lines and copies the same text.
