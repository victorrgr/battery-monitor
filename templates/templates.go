package templates

import "embed"

// Files contains embedded templates.
//
//go:embed report.html
//go:embed index.js
var Files embed.FS
