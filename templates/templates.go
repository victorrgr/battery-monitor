package templates

import "embed"

// Files contains embedded templates.
//
//go:embed report.gohtml
var Files embed.FS
