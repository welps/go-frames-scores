package drawing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/samber/lo"
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
	2: "basketball.png",
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
	case "basketball.png":
		return s.DrawBasketball(ctx)
	default:
		return bytes.Buffer{}, nil
	}
}

func (s *service) DrawRoot() (bytes.Buffer, error) {
	imageContext := gg.NewContext(frameImageX, frameImageY)
	imageContext.SetFontFace(GetFont(assets.FontFiraCode, 72))

	imageContext.SetRGB255(0, 0, 0)
	imageContext.Clear()
	imageContext.SetRGB255(254, 254, 254)
	imageContext.DrawStringAnchored("Live Sports Scores", frameImageX/2, frameImageY/4, 0.5, 0.5)
	imageContext.SetFontFace(GetFont(assets.FontNotoEmoji, 72))

	// Most emojis are busted: https://github.com/fogleman/gg/issues/7
	imageContext.DrawString("⚽⚾⛳⛸️", frameImageX/2.40, frameImageY/2)

	var buf bytes.Buffer
	err := png.Encode(&buf, imageContext.Image())

	return buf, err
}

func (s *service) DrawBasketball(ctx context.Context) (bytes.Buffer, error) {
	matches, err := s.sportsService.GetMatches(ctx, sports.Basketball, true)
	if err != nil {
		return bytes.Buffer{}, err
	}
	buf, err := s.drawSport(ctx, sports.Basketball, matches)
	return buf, err
}

func (s *service) DrawTennis(ctx context.Context) (bytes.Buffer, error) {
	matches, err := s.sportsService.GetMatches(ctx, sports.Tennis, true)
	if err != nil {
		return bytes.Buffer{}, err
	}
	buf, err := s.drawSport(ctx, sports.Tennis, matches)
	return buf, err
}

func (s *service) drawSport(_ context.Context, gameType sports.GameType, matches []sports.Match) (
	bytes.Buffer,
	error,
) {
	imageContext := gg.NewContext(frameImageX, frameImageY)
	imageContext.SetRGB255(0, 0, 0)
	imageContext.Clear()

	// Set title font and color
	titleFont := GetFont(assets.FontFiraCode, 72)
	imageContext.SetFontFace(titleFont)
	imageContext.SetRGB255(254, 254, 254)
	imageContext.DrawStringAnchored(fmt.Sprintf("Live %s Scores", gameType), frameImageX/2, frameImageY/12, 0.5, 0.5)

	if len(matches) == 0 {
		subTitleFont := GetFont(assets.FontFiraCode, 50)
		imageContext.SetFontFace(subTitleFont)
		imageContext.DrawStringAnchored("No live matches found :(", frameImageX/2, frameImageY/3, 0.5, 0.5)

		var buf bytes.Buffer
		err := png.Encode(&buf, imageContext.Image())

		return buf, err

	}
	// Set font for player names and scores
	playerNameFontSize := float64(60)
	playerNameFont := GetFont(assets.FontFiraCode, playerNameFontSize)
	imageContext.SetFontFace(playerNameFont)

	const paddingLeft float64 = 20
	const paddingRight float64 = 20
	var startX, startY float64 = paddingLeft, frameImageY / 6
	const boxHeight float64 = 130
	const boxWidth float64 = float64(frameImageX) / 2 // Two boxes per row
	const paddingBetweenBoxes float64 = 10

	// Determine the maximum score width
	maxScoreWidth := 0.0
	for _, match := range matches {
		homeScoresStr := fmt.Sprintf("%v", match.Score.Home)
		awayScoresStr := fmt.Sprintf("%v", match.Score.Away)
		homeScoreWidth, _ := imageContext.MeasureString(homeScoresStr)
		awayScoreWidth, _ := imageContext.MeasureString(awayScoresStr)
		if homeScoreWidth > maxScoreWidth {
			maxScoreWidth = homeScoreWidth
		}
		if awayScoreWidth > maxScoreWidth {
			maxScoreWidth = awayScoreWidth
		}
	}

	for i, match := range matches {
		if i%2 == 0 && i != 0 { // Move to next row after every 2 matches
			startY += boxHeight + paddingBetweenBoxes
		}

		// Draw rectangle for the current match
		imageContext.DrawRectangle(startX, startY, boxWidth-paddingRight, boxHeight)
		imageContext.SetRGB255(255, 255, 255) // Set color to white for filling
		imageContext.FillPreserve()           // Fill the rectangle and preserve the path for stroking
		imageContext.SetRGB255(0, 0, 0)       // Set color to black for the border
		imageContext.SetLineWidth(2)          // Set the line width for the border
		imageContext.Stroke()                 // Stroke the border

		// Set text color to black for drawing names and scores
		imageContext.SetRGB255(0, 0, 0)

		// Calculate vertical center for the text
		textYHome := startY + boxHeight/4 + playerNameFontSize/3
		textYAway := startY + 3*boxHeight/4 + playerNameFontSize/3

		// Draw names on the left side
		imageContext.DrawString(match.Home.Name, startX, textYHome)
		imageContext.DrawString(match.Away.Name, startX, textYAway)

		// Draw scores on the right side, aligned based on the maximum score width
		homeScoresStr := fmt.Sprintf("%s", reduceScore(match.Score.Home))
		awayScoresStr := fmt.Sprintf("%s", reduceScore(match.Score.Away))

		// Position scores on the right by using the maximum score width
		imageContext.DrawString(homeScoresStr, startX+boxWidth-paddingRight-maxScoreWidth, textYHome)
		imageContext.DrawString(awayScoresStr, startX+boxWidth-paddingRight-maxScoreWidth, textYAway)

		startX += boxWidth // Move to the next column
		if i%2 == 1 {      // At the end of the row, reset startX for the next row
			startX = paddingLeft
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, imageContext.Image())

	return buf, err
}

func reduceScore(score []string) string {
	return lo.Reduce(
		score,
		func(acc string, curr string, _ int) string {
			return fmt.Sprintf("%s %s", acc, curr)
		},
		"",
	)
}
