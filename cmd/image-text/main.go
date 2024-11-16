package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/matbur/image-text/server"
)

type config struct {
	Port string `envconfig:"PORT" default:"8080"`
	Mode string `envconfig:"MODE"`
}

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		slog.Error("Failed to process envconfig", "err", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	switch cfg.Mode {
	case "TEST":
		mode2(addr)
	default:
		mode1(addr)
	}
}

func mode1(addr string) {
	slog.Info("Starting server", "addr", addr)

	http.HandleFunc("/favicon.ico", server.HandleFavicon)
	http.HandleFunc("/healthz", server.HandleHealthz)
	http.HandleFunc("/", server.HandleMain())
	if err := http.ListenAndServe(addr, nil); err != nil {
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
		mode1(addr)
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
