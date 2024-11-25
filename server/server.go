package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/matbur/image-text/image"
	"github.com/matbur/image-text/resources"
	"github.com/matbur/image-text/templates"
)

func NewServer() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handleMain)
	r.Post("/post", handlePost)
	r.Get("/favicon.ico", handleFavicon)
	r.Get("/docs", handleDocs)
	r.Get("/*", handle)

	return r
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	params := templates.IndexPageParams{
		Text:    r.URL.Query().Get("text"),
		BgColor: r.URL.Query().Get("bg_color"),
		FgColor: r.URL.Query().Get("fg_color"),
		Size:    r.URL.Query().Get("size"),
	}

	img, err := image.New(params.Size, params.BgColor, params.FgColor, params.Text)
	if err != nil {
		slog.Error("Failed to create image", "err", err)
		templ.Handler(templates.IndexPage(params)).ServeHTTP(w, r)
		return
	}

	buf := bytes.Buffer{}
	if err := img.Draw(&buf); err != nil {
		slog.Error("Failed to draw image", "err", err)
		templ.Handler(templates.IndexPage(params)).ServeHTTP(w, r)
		return
	}

	params.Image = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	templ.Handler(templates.IndexPage(params)).ServeHTTP(w, r)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var params templates.IndexPageParams

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read body", "err", err)
		templ.Handler(templates.IndexPage(params)).ServeHTTP(w, r)
		return
	}

	if err := r.Body.Close(); err != nil {
		slog.Error("Failed to close body", "err", err)
	}

	if err := json.Unmarshal(bb, &params); err != nil {
		slog.Error("Failed to unmarshal body", "err", err)
		templ.Handler(templates.IndexPage(params)).ServeHTTP(w, r)
		return
	}

	img, err := image.New(params.Size, params.BgColor, params.FgColor, params.Text)
	if err != nil {
		slog.Error("Failed to create image", "err", err)
		templ.Handler(templates.IndexPage(params)).ServeHTTP(w, r)
		return
	}

	buf := bytes.Buffer{}
	if err := img.Draw(&buf); err != nil {
		slog.Error("Failed to draw image", "err", err)
		templ.Handler(templates.IndexPage(params)).ServeHTTP(w, r)
		return
	}

	params.Image = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	q := url.Values{}
	q.Set("text", params.Text)
	q.Set("bg_color", params.BgColor)
	q.Set("fg_color", params.FgColor)
	q.Set("size", params.Size)

	u := url.URL{
		Path:     "/",
		RawQuery: q.Encode(),
	}

	w.Header().Set("HX-Push-Url", u.String())
	templ.Handler(templates.IndexPage(params)).ServeHTTP(w, r)
}

func handle(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	if r.URL.Path == "/" {
		handleDocs(w, r)
		return
	}

	w.Header().Set("Content-Disposition", `inline; filename="image.png"`)

	size, bg, fg, err := parsePath(r.URL.Path)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	text := r.URL.Query().Get("text")
	if text == "" {
		text = size
	}

	img, err := image.New(size, bg, fg, text)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := img.Draw(w); err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.With(
		"response", "binary",
		"duration", time.Since(begin).String(),
		"status", http.StatusOK,
	).Info("Response")
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	bb, err := resources.Static.ReadFile("favicon.png")
	if err != nil {
		slog.Error("Failed to load favicon", "err", err)
		writeJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(bb); err != nil {
		slog.Error("Failed to write favicon", "err", err)
		return
	}
}

var docs = struct {
	Path     string            `json:"path"`
	Examples map[string]string `json:"example"`
	Colors   map[string]string `json:"colors"`
	Sizes    map[string]string `json:"sizes"`
}{
	Path: "HOST/size/background/foreground?text=rendered+text",
	Examples: map[string]string{
		"with_names": "/hd720/steel_blue/yellow?text=rendered+text",
		"with_codes": "/320x200/000/FFFF00",
	},
}

func handleDocs(w http.ResponseWriter, r *http.Request) {
	sizes := map[string]string{}
	for k, v := range image.KnownSizes() {
		sizes[k] = v.String()
	}
	docs.Sizes = sizes

	colors := map[string]string{}
	for k, v := range image.KnownColors() {
		colors[k] = v.String()
	}
	docs.Colors = colors

	js, err := json.Marshal(docs)
	if err != nil {
		slog.Error("Failed to marshal docs", "err", err)
		writeJSON(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(js); err != nil {
		slog.Error("Failed to write docs", "err", err)
	}
}
