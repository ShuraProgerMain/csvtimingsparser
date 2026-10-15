package ui

import _ "embed"

// indexHTML is the entire Material Design front-end, shared by every platform's
// webview host.
//
//go:embed assets/index.html
var indexHTML string
