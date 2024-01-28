package templates

import "embed"

// EmbeddedTemplates will hold all templates in memory at compile time
//
//go:embed *
var EmbeddedTemplates embed.FS
