package templates

import "embed"

// Files contains embedded templates.
//
//go:embed report.html
//go:embed index.js
//go:embed style.css
var Files embed.FS
