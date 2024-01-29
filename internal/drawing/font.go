package drawing

import (
	"fmt"
	"github.com/goki/freetype/truetype"
	"github.com/welps/go-frames-scores/assets"
	"go.uber.org/zap"
	"golang.org/x/image/font"
)

// GetFont returns a new font every time because it's not concurrent safe
func GetFont(fontType string, size float64) truetype.IndexableFace {
	embeddedFont, err := assets.Embedded.ReadFile(fmt.Sprintf("%s/%s", assets.FontsPath, fontType))
	if err != nil {
		zap.S().Errorw("Unable to read embedded font", zap.Error(err))
		return nil
	}

	// We use goki/freetype to address this issue: https://github.com/fogleman/gg/issues/153#issuecomment-1849023145
	f, err := truetype.Parse(embeddedFont)
	if err != nil {
		zap.S().Errorw("Unable to parse embedded font", zap.Error(err))
		return nil
	}

	ff := truetype.NewFace(
		f, &truetype.Options{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingNone,
		},
	)
	return ff
}
