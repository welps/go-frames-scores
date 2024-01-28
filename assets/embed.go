package assets

import "embed"

// Embedded will hold all assets in memory at compile time
//
//go:embed *
var Embedded embed.FS
