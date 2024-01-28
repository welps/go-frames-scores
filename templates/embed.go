package templates

import "embed"

// Embedded will hold all templates in memory at compile time
//
//go:embed *
var Embedded embed.FS
