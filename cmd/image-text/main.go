package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/matbur/image-text/server"
	"github.com/matbur/image-text/templates"
)

type config struct {
	Port              string        `envconfig:"PORT" default:"8080"`
	Mode              string        `envconfig:"MODE"`
	RateLimitRequests int           `envconfig:"RATE_LIMIT_REQUESTS" default:"100"`
	RateLimitWindow   time.Duration `envconfig:"RATE_LIMIT_WINDOW" default:"1m"`
	CacheSize         int           `envconfig:"CACHE_SIZE" default:"512"`
	AllowedOrigins    string        `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
}

func init() {
	level, err := parseLogLevel(os.Getenv("LOG_LEVEL"))

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})))

	if err != nil {
		slog.Warn("Invalid LOG_LEVEL, falling back to default", "err", err, "default", level.String())
	}
}

// parseLogLevel parses a log level string (e.g. "debug", "info", "warn",
// "error") into a slog.Level. An empty string defaults to slog.LevelInfo
// without error. On any parse error it also returns slog.LevelInfo.
func parseLogLevel(s string) (slog.Level, error) {
	if s == "" {
		return slog.LevelInfo, nil
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(s)); err != nil {
		return slog.LevelInfo, err
	}

	return level, nil
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		slog.Error("Failed to process envconfig", "err", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	srvCfg := server.Config{
		RateLimitRequests: cfg.RateLimitRequests,
		RateLimitWindow:   cfg.RateLimitWindow,
		CacheSize:         cfg.CacheSize,
		AllowedOrigins:    cfg.AllowedOrigins,
	}

	switch cfg.Mode {
	case "TEST":
		mode2(addr)
	default:
		mode1(addr, srvCfg)
	}
}

func mode1(addr string, srvCfg server.Config) {
	slog.Info("Starting server", "addr", addr)

	if os.Getenv("K_SERVICE") != "" {
		templates.EnvLabel = ""
	} else if _, err := os.Stat("/.dockerenv"); err == nil {
		templates.EnvLabel = "docker"
	} else {
		templates.EnvLabel = "local"
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: server.NewServer(srvCfg),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		slog.Info("Shutting down gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Graceful shutdown failed", "err", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("Server closed")
			return
		}
		slog.Error("Failed to start server", "err", err)
	}
}

// for debug only
func mode2(addr string) {
	go func() {
		mode1(addr, server.Config{})
	}()

	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	u := fmt.Sprintf("http://%s/3000x200/steel_blue/yellow?text=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", addr)
	resp, err := http.Get(u)
	if err != nil {
		slog.Error("Failed to get image", "err", err)
		os.Exit(1)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Failed to close body", "err", err)
		}
	}()

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read body", "err", err)
		os.Exit(1)
	}

	if err := os.WriteFile("image.png", bb, 0644); err != nil {
		slog.Error("Failed to write image", "err", err)
		os.Exit(1)
	}
}
