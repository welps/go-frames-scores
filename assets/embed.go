package assets

import "embed"

const (
	FontsPath     = "fonts"
	FontWorkSans  = "WorkSans-Medium.ttf"
	FontNotoEmoji = "NotoEmoji-Regular.ttf"
)

// Embedded will hold all assets in memory at compile time
//
//go:embed *
var Embedded embed.FS
