package templates

import "embed"

// Files contains embedded templates.
//
//go:embed index.html
//go:embed index.js
//go:embed styles.css
var Files embed.FS
