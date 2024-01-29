package drawing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/welps/go-frames-scores/assets"
	"github.com/welps/go-frames-scores/internal/sports"
	"image/png"
)

const (
	frameImageX        = 2062
	frameImageY        = 1080
	generatedDirectory = "generated"
)

var assetMapping = map[int]string{
	0: "root.png",
	1: "tennis.png",
}

type Service interface {
	GetAssetPath(buttonIndex int) string
	DrawFile(ctx context.Context, filename string) (bytes.Buffer, error)
}

func NewService(sportsService sports.Service) Service {
	return &service{
		sportsService: sportsService,
	}
}

type service struct {
	sportsService sports.Service
}

func (s *service) GetAssetPath(buttonIndex int) string {
	return fmt.Sprintf("%s/%s", generatedDirectory, assetMapping[buttonIndex])
}

func (s *service) DrawFile(ctx context.Context, filename string) (bytes.Buffer, error) {
	switch filename {
	case "root.png":
		return s.DrawRoot()
	case "tennis.png":
		return s.DrawTennis(ctx)
	default:
		return bytes.Buffer{}, nil
	}
}

func (s *service) DrawRoot() (bytes.Buffer, error) {
	imageContext := gg.NewContext(frameImageX, frameImageY)
	imageContext.SetFontFace(GetFont(assets.FontWorkSans, 72))

	imageContext.SetRGB255(0, 0, 0)
	imageContext.Clear()
	imageContext.SetRGB255(254, 254, 254)
	imageContext.DrawStringAnchored("Sports Scores", frameImageX/2, frameImageY/8, 0.5, 0.5)
	imageContext.SetFontFace(GetFont(assets.FontNotoEmoji, 72))

	// Most emojis are busted: https://github.com/fogleman/gg/issues/7
	imageContext.DrawString("⚽⚾⛳⛸️", frameImageX/2.40, frameImageY/4)

	var buf bytes.Buffer
	err := png.Encode(&buf, imageContext.Image())

	return buf, err
}

func (s *service) DrawTennis(ctx context.Context) (bytes.Buffer, error) {
	matches, err := s.sportsService.GetMatches(ctx, sports.Tennis)
	if err != nil {
		return bytes.Buffer{}, err
	}

	imageContext := gg.NewContext(frameImageX, frameImageY)
	imageContext.SetFontFace(GetFont(assets.FontWorkSans, 72))

	imageContext.SetRGB255(0, 0, 0)
	imageContext.Clear()
	imageContext.SetRGB255(254, 254, 254)
	imageContext.DrawStringAnchored("Tennis Scores", frameImageX/2, frameImageY/8, 0.5, 0.5)

	imageContext.SetFontFace(GetFont(assets.FontWorkSans, 18))

	var startX, startY float64 = 10, frameImageY / 4
	const boxHeight float64 = 50
	const boxWidth float64 = float64(frameImageX) / 3 // Divide the width of the image by 3 to fit three entries in a row
	const padding float64 = 10
	const scoreSpacing float64 = 5

	for i, match := range matches {
		if i%3 == 0 && i != 0 { // Move to next row after every 3 matches
			startX = 10
			startY += boxHeight + padding
		}

		// Draw rectangle for the current match
		imageContext.DrawRectangle(startX, startY, boxWidth-padding, boxHeight)
		imageContext.Fill()
		imageContext.SetRGB255(0, 0, 0)

		// Draw names and scores for Home player
		homeNameWidth, _ := imageContext.MeasureString(match.Home.Name)
		imageContext.DrawString(match.Home.Name, startX, startY+20)
		scoreX := startX + homeNameWidth + scoreSpacing
		for _, score := range match.Score.Home {
			imageContext.DrawString(fmt.Sprintf("%v", score), scoreX, startY+20)
			scoreWidth, _ := imageContext.MeasureString(fmt.Sprintf("%v", score))
			scoreX += scoreWidth + scoreSpacing
		}

		// Draw names and scores for Away player
		awayNameWidth, _ := imageContext.MeasureString(match.Away.Name)
		imageContext.DrawString(match.Away.Name, startX, startY+38)
		scoreX = startX + awayNameWidth + scoreSpacing
		for _, score := range match.Score.Away {
			imageContext.DrawString(fmt.Sprintf("%v", score), scoreX, startY+38)
			scoreWidth, _ := imageContext.MeasureString(fmt.Sprintf("%v", score))
			scoreX += scoreWidth + scoreSpacing
		}

		startX += boxWidth // Move to next column
		imageContext.SetRGB255(254, 254, 254)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, imageContext.Image())

	return buf, err
}
