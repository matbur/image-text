package server

import (
	"cmp"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/a-h/templ"
	"github.com/elliotchance/pie/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/matbur/image-text/image"
	"github.com/matbur/image-text/resources"
	"github.com/matbur/image-text/templates"
	"github.com/matbur/image-text/wasm"
)

func NewServer() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handleMain)
	r.Get("/online", handleOnlinePage)
	r.Post("/online/post", handleOnlinePost)
	r.Get("/offline", handleOfflinePage)
	r.Get("/favicon.ico", handleFavicon)
	r.Get("/docs", handleDocs)
	r.Get("/resources/{filename}", handleStatic)
	r.Get("/{size}/{bg_color}/{fg_color}", handleImage)

	return r
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Expires", time.Now().Add(24*time.Hour).Format(http.TimeFormat))

	fn := r.PathValue("filename")
	if fn == "main.wasm" {
		slog.Info("main.wasm")
		w.Write(wasm.Main)
		return
	}
	http.ServeFileFS(w, r, resources.Static, fn)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	templ.Handler(templates.IndexPage(templates.IndexPageParams{})).ServeHTTP(w, r)
}

func handleOnlinePage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	text := q.Get("text")
	bgColor := q.Get("bg_color")
	fgColor := q.Get("fg_color")
	size := q.Get("size")

	if bgColor == "" || fgColor == "" || size == "" {
		q := url.Values{}
		q.Set("text", text)
		q.Set("bg_color", cmp.Or(bgColor, "steel_blue"))
		q.Set("fg_color", cmp.Or(fgColor, "yellow"))
		q.Set("size", cmp.Or(size, "vga"))

		u := url.URL{Path: "/online", RawQuery: q.Encode()}

		http.Redirect(w, r, u.String(), http.StatusSeeOther)
		return
	}

	u := &url.URL{Path: "/", RawQuery: q.Encode()}
	u = u.JoinPath(size, bgColor, fgColor)

	params := templates.OnlinePageParams{
		Text:         text,
		BgColor:      bgColor,
		FgColor:      fgColor,
		Size:         size,
		Image:        u.String(),
		ColorOptions: pie.Keys(image.KnownColors()),
		SizeOptions:  pie.Keys(image.KnownSizes()),
	}
	templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
}

func handleOnlinePost(w http.ResponseWriter, r *http.Request) {
	var params templates.OnlinePageParams

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read body", "err", err)
		templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
		return
	}

	if err := r.Body.Close(); err != nil {
		slog.Error("Failed to close body", "err", err)
	}

	if err := json.Unmarshal(bb, &params); err != nil {
		slog.Error("Failed to unmarshal body", "err", err)
		templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
		return
	}

	q := url.Values{}
	q.Set("text", params.Text)
	q.Set("bg_color", params.BgColor)
	q.Set("fg_color", params.FgColor)
	q.Set("size", params.Size)

	u := &url.URL{Path: "/online", RawQuery: q.Encode()}
	w.Header().Set("HX-Push-Url", u.String())
	slog.Info("Pushing", "url", u.String())

	u2 := &url.URL{RawQuery: url.Values{"text": {params.Text}}.Encode()}
	u2 = u2.JoinPath(params.Size, params.BgColor, params.FgColor)
	slog.Info("Image url", "url", u2.String())

	params.Image = u2.String()
	params.ColorOptions = pie.Keys(image.KnownColors())
	params.SizeOptions = pie.Keys(image.KnownSizes())
	templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	size := chi.URLParam(r, "size")
	bgColor := chi.URLParam(r, "bg_color")
	fgColor := chi.URLParam(r, "fg_color")
	text := r.URL.Query().Get("text")

	w.Header().Set("Content-Disposition", `inline; filename="image.png"`)

	img, err := image.New(size, bgColor, fgColor, text)
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

func handleOfflinePage(w http.ResponseWriter, r *http.Request) {
	params := templates.OfflinePageParams{
		ColorOptions: pie.Keys(image.KnownColors()),
		SizeOptions:  pie.Keys(image.KnownSizes()),
	}
	templ.Handler(templates.OfflinePage(params)).ServeHTTP(w, r)
}
