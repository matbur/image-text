package server

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"sort"
	"time"

	"github.com/a-h/templ"
	"github.com/elliotchance/pie/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/matbur/image-text/i18n"
	"github.com/matbur/image-text/image"
	"github.com/matbur/image-text/resources"
	"github.com/matbur/image-text/templates"
	"github.com/matbur/image-text/version"
	"github.com/matbur/image-text/wasm"
)

type Config struct {
	RateLimitRequests int
	RateLimitWindow   time.Duration
	CacheSize         int
}

func NewServer(cfg Config) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "application/wasm", "application/octet-stream", "text/html", "text/css", "application/javascript", "image/svg+xml"))
	if cfg.RateLimitRequests > 0 {
		r.Use(httprate.LimitByIP(cfg.RateLimitRequests, cfg.RateLimitWindow))
	}

	cache := imageCache(cfg)

	r.Get("/healthz", handleHealthz)
	r.Get("/readyz", handleReadyz)
	r.Get("/", handleMain)
	r.Get("/online", handleOnlinePage)
	r.Post("/online/post", handleOnlinePost)
	r.Get("/offline", handleOfflinePage)
	r.Get("/favicon.ico", handleFavicon)
	r.Get("/docs", handleDocsPage)
	r.Get("/docs.json", handleDocsJSON)
	r.Get("/resources/{filename}", handleStatic)
	r.Get("/resources/fonts/{filename}", handleFontStatic)
	r.Get("/{size}/{bg_color}/{fg_color}", cachedHandleImage(cache))

	return r
}

func imageCache(cfg Config) *lru.Cache[string, []byte] {
	if cfg.CacheSize <= 0 {
		return nil
	}
	c, err := lru.New[string, []byte](cfg.CacheSize)
	if err != nil {
		slog.Warn("Failed to create image cache", "err", err)
		return nil
	}
	return c
}

func cacheKey(size, bgColor, fgColor, text, font, format string) string {
	return fmt.Sprintf("%s/%s/%s?text=%s&font=%s&format=%s", size, bgColor, fgColor, text, font, format)
}

func cachedHandleImage(cache *lru.Cache[string, []byte]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		size := chi.URLParam(r, "size")
		bgColor := chi.URLParam(r, "bg_color")
		fgColor := chi.URLParam(r, "fg_color")
		text := r.URL.Query().Get("text")
		font := r.URL.Query().Get("font")
		format := r.URL.Query().Get("format")

		key := cacheKey(size, bgColor, fgColor, text, font, format)

		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="image.%s"`, image.Extension(format)))
		w.Header().Set("Content-Type", image.ContentType(format))

		if cache != nil {
			if bb, ok := cache.Get(key); ok {
				w.Write(bb)
				slog.With(
					"response", "binary",
					"cache", "hit",
					"duration", time.Since(begin).String(),
					"status", http.StatusOK,
				).Info("Response")
				return
			}
		}

		img, err := image.New(size, bgColor, fgColor, text, font, format)
		if err != nil {
			writeJSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		var buf bytes.Buffer
		if err := img.Draw(&buf); err != nil {
			writeJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bb := buf.Bytes()

		if cache != nil {
			cache.Add(key, bb)
		}

		w.Write(bb)

		slog.With(
			"response", "binary",
			"duration", time.Since(begin).String(),
			"status", http.StatusOK,
		).Info("Response")
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func handleReadyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Expires", time.Now().Add(24*time.Hour).Format(http.TimeFormat))

	fn := r.PathValue("filename")
	if fn == "main.wasm" {
		slog.Info("main.wasm")
		w.Header().Set("Content-Type", "application/wasm")
		w.Write(wasm.Main)
		return
	}
	http.ServeFileFS(w, r, resources.Static, fn)
}

func handleFontStatic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Expires", time.Now().Add(24*time.Hour).Format(http.TimeFormat))

	filename, err := url.PathUnescape(r.PathValue("filename"))
	if err != nil {
		writeJSON(w, "malformed font path", http.StatusBadRequest)
		return
	}

	http.ServeFileFS(w, r, resources.Static, path.Join("fonts", filename))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	locale := i18n.FromRequest(r)
	templ.Handler(templates.IndexPage(templates.IndexPageParams{
		CommitSHA: version.Commit,
		I18n:      locale,
	})).ServeHTTP(w, r)
}

func handleOnlinePage(w http.ResponseWriter, r *http.Request) {
	locale := i18n.FromRequest(r)
	q := r.URL.Query()
	text := q.Get("text")
	bgColor := q.Get("bg_color")
	fgColor := q.Get("fg_color")
	size := q.Get("size")

	font := q.Get("font")
	format := q.Get("format")
	if bgColor == "" || fgColor == "" || size == "" {
		q := url.Values{}
		q.Set("text", text)
		q.Set("bg_color", cmp.Or(bgColor, "steel_blue"))
		q.Set("fg_color", cmp.Or(fgColor, "yellow"))
		q.Set("size", cmp.Or(size, "vga"))
		q.Set("font", cmp.Or(font, "ubuntu_mono"))
		q.Set("format", cmp.Or(format, "png"))

		u := url.URL{Path: "/online", RawQuery: q.Encode()}

		http.Redirect(w, r, u.String(), http.StatusSeeOther)
		return
	}

	imageQuery := url.Values{}
	imageQuery.Set("text", text)
	imageQuery.Set("font", cmp.Or(font, "ubuntu_mono"))
	if format != "" && format != "png" {
		imageQuery.Set("format", format)
	}
	u := &url.URL{Path: "/", RawQuery: imageQuery.Encode()}
	u = u.JoinPath(size, bgColor, fgColor)

	params := templates.OnlinePageParams{
		Text:          text,
		BgColor:       bgColor,
		FgColor:       fgColor,
		Size:          size,
		Font:          cmp.Or(font, "ubuntu_mono"),
		Format:        cmp.Or(format, "png"),
		Image:         u.String(),
		ColorOptions:  pie.Keys(image.KnownColors()),
		SizeOptions:   pie.Keys(image.KnownSizes()),
		FontOptions:   image.KnownFontNames(),
		FormatOptions: image.KnownFormatStrings(),
		I18n:          locale,
	}
	templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
}

func handleOnlinePost(w http.ResponseWriter, r *http.Request) {
	locale := i18n.FromRequest(r)
	var params templates.OnlinePageParams

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read body", "err", err)
		params.I18n = locale
		templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
		return
	}

	if err := r.Body.Close(); err != nil {
		slog.Error("Failed to close body", "err", err)
	}

	if err := json.Unmarshal(bb, &params); err != nil {
		slog.Error("Failed to unmarshal body", "err", err)
		params.I18n = locale
		templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
		return
	}

	params.I18n = locale

	q := url.Values{}
	q.Set("text", params.Text)
	q.Set("bg_color", params.BgColor)
	q.Set("fg_color", params.FgColor)
	q.Set("size", params.Size)
	q.Set("font", cmp.Or(params.Font, "ubuntu_mono"))
	q.Set("format", cmp.Or(params.Format, "png"))

	u := &url.URL{Path: "/online", RawQuery: q.Encode()}
	w.Header().Set("HX-Push-Url", u.String())
	slog.Info("Pushing", "url", u.String())

	imgQuery := url.Values{}
	imgQuery.Set("text", params.Text)
	imgQuery.Set("font", cmp.Or(params.Font, "ubuntu_mono"))
	if params.Format != "" && params.Format != "png" {
		imgQuery.Set("format", params.Format)
	}
	u2 := &url.URL{RawQuery: imgQuery.Encode()}
	u2 = u2.JoinPath(params.Size, params.BgColor, params.FgColor)
	slog.Info("Image url", "url", u2.String())

	params.Image = u2.String()
	params.ColorOptions = pie.Keys(image.KnownColors())
	params.SizeOptions = pie.Keys(image.KnownSizes())
	params.FontOptions = image.KnownFontNames()
	params.FormatOptions = image.KnownFormatStrings()
	if r.Header.Get("HX-Request") != "" {
		templ.Handler(templates.Img(locale, params.Image)).ServeHTTP(w, r)
		return
	}
	templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
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
	Params   map[string]string `json:"params"`
	Examples map[string]string `json:"example"`
	Colors   map[string]string `json:"colors"`
	Sizes    map[string]string `json:"sizes"`
	Fonts    []string          `json:"fonts"`
}{
	Path: "HOST/size/background/foreground?text=rendered+text&font=ubuntu_mono&format=png",
	Examples: map[string]string{
		"with_names": "/hd720/steel_blue/yellow?text=rendered+text&font=ubuntu_mono",
		"with_codes": "/320x200/000/FFFF00?font=open_sans",
		"with_jpeg":  "/hd720/steel_blue/yellow?text=photo&format=jpg",
		"with_webp":  "/hd720/steel_blue/yellow?text=photo&format=webp",
	},
	Colors: image.KnownColorStrings(),
	Sizes:  image.KnownSizeStrings(),
	Fonts:  image.KnownFontNames(),
}

func docsParams(locale i18n.Locale) map[string]string {
	return map[string]string{
		"text":   locale.T("docs.param.text"),
		"font":   locale.T("docs.param.font"),
		"format": locale.T("docs.param.format"),
	}
}

func docsExamples(locale i18n.Locale) map[string]string {
	examples := make(map[string]string, len(docs.Examples))
	for key, value := range docs.Examples {
		examples[locale.T("docs.example."+key)] = value
	}
	return examples
}

func docsEntries(values map[string]string) []templates.DocsEntry {
	names := pie.Keys(values)
	sort.Strings(names)

	entries := make([]templates.DocsEntry, len(names))
	for i, name := range names {
		entries[i] = templates.DocsEntry{
			Name:  name,
			Value: values[name],
		}
	}
	return entries
}

func handleDocsPage(w http.ResponseWriter, r *http.Request) {
	locale := i18n.FromRequest(r)
	params := templates.DocsPageParams{
		Path:         docs.Path,
		Params:       docsParams(locale),
		Examples:     docsExamples(locale),
		ColorEntries: docsEntries(docs.Colors),
		SizeEntries:  docsEntries(docs.Sizes),
		Fonts:        docs.Fonts,
		I18n:         locale,
	}
	templ.Handler(templates.DocsPage(params)).ServeHTTP(w, r)
}

func handleDocsJSON(w http.ResponseWriter, r *http.Request) {
	locale := i18n.FromRequest(r)
	payload := struct {
		Path     string            `json:"path"`
		Params   map[string]string `json:"params"`
		Examples map[string]string `json:"example"`
		Colors   map[string]string `json:"colors"`
		Sizes    map[string]string `json:"sizes"`
		Fonts    []string          `json:"fonts"`
	}{
		Path:     docs.Path,
		Params:   docsParams(locale),
		Examples: docsExamples(locale),
		Colors:   docs.Colors,
		Sizes:    docs.Sizes,
		Fonts:    docs.Fonts,
	}

	js, err := json.Marshal(payload)
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
	locale := i18n.FromRequest(r)
	params := templates.OfflinePageParams{
		Font:          "ubuntu_mono",
		Format:        "png",
		ColorOptions:  pie.Keys(image.KnownColors()),
		SizeOptions:   pie.Keys(image.KnownSizes()),
		FontOptions:   image.KnownFontNames(),
		FormatOptions: image.KnownFormatStrings(),
		FontFiles:     image.KnownFontFilenames(),
		I18n:          locale,
	}
	templ.Handler(templates.OfflinePage(params)).ServeHTTP(w, r)
}
