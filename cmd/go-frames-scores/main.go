package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"github.com/welps/go-frames-scores/internal/drawing"
	"github.com/welps/go-frames-scores/internal/frame"
	"github.com/welps/go-frames-scores/internal/sports"

	"html/template"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/welps/go-frames-scores/internal/config"
	"github.com/welps/go-frames-scores/internal/constants"
	"github.com/welps/go-frames-scores/templates"
	"go.uber.org/zap"
)

func main() {
	config := config.InitConfig()

	logger := getLogger(config)
	logger.Sugar().Debugw("Config values", zap.Any("config", config))
	// nolint: errcheck
	defer logger.Sync()

	httpClient := resty.NewWithClient(getHTTPClient(config.HTTPClientSettings))
	client, err := sports.NewClient(httpClient, config.SportsAPIConfig.Host, config.SportsAPIConfig.APIKey)
	fatalAndExitOnError(err, "Unable to create sports client")

	service := sports.NewService(client)
	err = service.UpdateMatches(context.Background(), true)
	fatalAndExitOnError(err, "Unable to update matches")
	drawingService := drawing.NewService(service)

	r := getConfiguredRouter(logger)
	r.GET(
		"/healthcheck", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		},
	)

	controller := frame.NewController(config.PublicURL, drawingService)
	r.GET("/", controller.GetRoot)
	r.POST("/", controller.PostRoot)
	r.GET("/generated/:filename", controller.Draw)
	r.GET("/generated/:timestamp/:filename", controller.Draw)

	crontab := cron.New()
	_, err = crontab.AddFunc(
		"@every 1m", func() {
			err = service.UpdateMatches(context.Background(), true)
			if err != nil {
				zap.S().Error("failed to update live scores", zap.Error(err))
			}
		},
	)
	fatalAndExitOnError(err, "Unable to add cron job")
	crontab.Start()

	// Start main server with graceful shutdown
	listenAndServe(
		logger,
		r,
		fmt.Sprintf(":%d", config.Port),
		time.Duration(config.GracefulShutdownMS)*time.Millisecond,
	)
}

func getLogger(config config.Config) *zap.Logger {
	var cfg zap.Config

	if config.Environment == constants.EnvDevelopment {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Unable to start logger: %s", err)
	}

	_ = zap.ReplaceGlobals(logger)

	return logger
}

// listenAndServe acts as http.Server#listenAndServe with additional layer of logging and graceful shutdown
func listenAndServe(logger *zap.Logger, handler http.Handler, address string, gracefulShutdown time.Duration) {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Main server
	srv := &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second, // Large but finite value to prevent Slowloris Attack (G112). Thanks gosec.
	}

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		logger.Sugar().Infof("Starting server on %s", address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fatalAndExitOnError(err, "Server unable to start")
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown
	stop()
	logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdown)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fatalAndExitOnError(err, "Server forced to shutdown")
	}

	logger.Info("Server exiting")
}

func fatalAndExitOnError(err error, message string) {
	if err != nil {
		zap.S().Fatalw(message, zap.Error(err))
	}
}

func getConfiguredRouter(logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(ginzap.RecoveryWithZap(logger, true))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	r.Use(cors.New(corsConfig))

	// Load templates that are *embedded* in binary
	templates := template.Must(template.New("").ParseFS(templates.Embedded, "*.tmpl"))
	r.SetHTMLTemplate(templates)

	return r
}

// getHTTPClient returns a configured HTTP client with sane defaults
func getHTTPClient(settings config.HTTPClientSettings) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = settings.MaxIdleConns
	t.MaxIdleConnsPerHost = settings.MaxIdleConnsPerHost
	t.MaxConnsPerHost = settings.MaxIdleConnsPerHost * 2

	return &http.Client{
		Transport: t,
		Timeout:   time.Duration(settings.RequestTimeoutMS) * time.Millisecond,
	}
}
