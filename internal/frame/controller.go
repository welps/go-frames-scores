package frame

import (
	"fmt"
	"github.com/welps/go-frames-scores/internal/drawing"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	publicURL      string
	drawingService drawing.Service
}

func NewController(publicURL string, drawingService drawing.Service) *Controller {
	return &Controller{
		publicURL:      publicURL,
		drawingService: drawingService,
	}
}

func (c *Controller) GetRoot(ctx *gin.Context) {
	assetPath := c.drawingService.GetAssetPath(0)

	ctx.HTML(
		http.StatusOK, "index.tmpl", gin.H{
			"image":   fmt.Sprintf("%s/%s", c.publicURL, assetPath),
			"button1": GetFrameButton(1, "🎾 Tennis"),
			"button2": GetFrameButton(2, "🏀 Basketball"),
		},
	)
}

func (c *Controller) PostRoot(ctx *gin.Context) {
	var data Post
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zap.S().Debugw("JSON data", zap.Any("data", data))
	buttonIndex := data.UntrustedData.ButtonIndex

	assetPath := c.drawingService.GetAssetPath(buttonIndex)
	ctx.HTML(
		http.StatusOK, "index.tmpl", gin.H{
			"image": fmt.Sprintf("%s/%s", c.publicURL, assetPath),
		},
	)
}

func (c *Controller) Draw(ctx *gin.Context) {
	filename := ctx.Param("filename")
	if filename == "" {
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	buf, err := c.drawingService.DrawFile(ctx, filename)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusOK, "image/png", buf.Bytes())
}
