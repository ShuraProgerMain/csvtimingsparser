# Build csvcolor (GUI + CLI) on Windows. Pure Go — no C compiler required.
# Usage:  ./build.ps1          -> builds bin/csvcolor.exe
#         ./build.ps1 -Run     -> build then launch the GUI
param([switch]$Run)

$ErrorActionPreference = "Stop"
$env:CGO_ENABLED = "0"

New-Item -ItemType Directory -Force -Path bin | Out-Null
Write-Host "Building bin/csvcolor.exe ..."
# -H windowsgui: no console window on double-click; the CLI reattaches to the
#   parent terminal at runtime (see platform_windows.go).
# -s -w: strip debug info to shrink the binary.
go build -ldflags "-H windowsgui -s -w" -o bin/csvcolor.exe .
Write-Host "Done: bin/csvcolor.exe"

if ($Run) { & ".\bin\csvcolor.exe" }
