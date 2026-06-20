package server

import (
	"cmp"
	"encoding/json"
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

	"github.com/matbur/image-text/image"
	"github.com/matbur/image-text/resources"
	"github.com/matbur/image-text/templates"
	"github.com/matbur/image-text/wasm"
)

func NewServer() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "application/wasm", "application/octet-stream", "text/html", "text/css", "application/javascript", "image/svg+xml"))

	r.Get("/", handleMain)
	r.Get("/online", handleOnlinePage)
	r.Post("/online/post", handleOnlinePost)
	r.Get("/offline", handleOfflinePage)
	r.Get("/favicon.ico", handleFavicon)
	r.Get("/docs", handleDocsPage)
	r.Get("/docs.json", handleDocsJSON)
	r.Get("/resources/{filename}", handleStatic)
	r.Get("/resources/fonts/{filename}", handleFontStatic)
	r.Get("/{size}/{bg_color}/{fg_color}", handleImage)

	return r
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
	templ.Handler(templates.IndexPage(templates.IndexPageParams{})).ServeHTTP(w, r)
}

func handleOnlinePage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	text := q.Get("text")
	bgColor := q.Get("bg_color")
	fgColor := q.Get("fg_color")
	size := q.Get("size")

	font := q.Get("font")
	if bgColor == "" || fgColor == "" || size == "" {
		q := url.Values{}
		q.Set("text", text)
		q.Set("bg_color", cmp.Or(bgColor, "steel_blue"))
		q.Set("fg_color", cmp.Or(fgColor, "yellow"))
		q.Set("size", cmp.Or(size, "vga"))
		q.Set("font", cmp.Or(font, "ubuntu_mono"))

		u := url.URL{Path: "/online", RawQuery: q.Encode()}

		http.Redirect(w, r, u.String(), http.StatusSeeOther)
		return
	}

	imageQuery := url.Values{}
	imageQuery.Set("text", text)
	imageQuery.Set("font", cmp.Or(font, "ubuntu_mono"))
	u := &url.URL{Path: "/", RawQuery: imageQuery.Encode()}
	u = u.JoinPath(size, bgColor, fgColor)

	params := templates.OnlinePageParams{
		Text:         text,
		BgColor:      bgColor,
		FgColor:      fgColor,
		Size:         size,
		Font:         cmp.Or(font, "ubuntu_mono"),
		Image:        u.String(),
		ColorOptions: pie.Keys(image.KnownColors()),
		SizeOptions:  pie.Keys(image.KnownSizes()),
		FontOptions:  pie.Keys(image.KnownFonts()),
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
	q.Set("font", cmp.Or(params.Font, "ubuntu_mono"))

	u := &url.URL{Path: "/online", RawQuery: q.Encode()}
	w.Header().Set("HX-Push-Url", u.String())
	slog.Info("Pushing", "url", u.String())

	u2 := &url.URL{RawQuery: url.Values{
		"text": {params.Text},
		"font": {cmp.Or(params.Font, "ubuntu_mono")},
	}.Encode()}
	u2 = u2.JoinPath(params.Size, params.BgColor, params.FgColor)
	slog.Info("Image url", "url", u2.String())

	params.Image = u2.String()
	params.ColorOptions = pie.Keys(image.KnownColors())
	params.SizeOptions = pie.Keys(image.KnownSizes())
	params.FontOptions = pie.Keys(image.KnownFonts())
	if r.Header.Get("HX-Request") != "" {
		templ.Handler(templates.Img(params.Image)).ServeHTTP(w, r)
		return
	}
	templ.Handler(templates.OnlinePage(params)).ServeHTTP(w, r)
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	size := chi.URLParam(r, "size")
	bgColor := chi.URLParam(r, "bg_color")
	fgColor := chi.URLParam(r, "fg_color")
	text := r.URL.Query().Get("text")
	font := r.URL.Query().Get("font")

	w.Header().Set("Content-Disposition", `inline; filename="image.png"`)

	img, err := image.New(size, bgColor, fgColor, text, font)
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
	Params   map[string]string `json:"params"`
	Examples map[string]string `json:"example"`
	Colors   map[string]string `json:"colors"`
	Sizes    map[string]string `json:"sizes"`
	Fonts    []string          `json:"fonts"`
}{
	Path: "HOST/size/background/foreground?text=rendered+text&font=ubuntu_mono",
	Params: map[string]string{
		"text": "rendered text (optional, defaults to size)",
		"font": "font name (optional, defaults to ubuntu_mono)",
	},
	Examples: map[string]string{
		"with_names": "/hd720/steel_blue/yellow?text=rendered+text&font=ubuntu_mono",
		"with_codes": "/320x200/000/FFFF00?font=open_sans",
	},
	Colors: image.KnownColorStrings(),
	Sizes:  image.KnownSizeStrings(),
	Fonts:  pie.Keys(image.KnownFonts()),
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
	params := templates.DocsPageParams{
		Path:         docs.Path,
		Params:       docs.Params,
		Examples:     docs.Examples,
		ColorEntries: docsEntries(docs.Colors),
		SizeEntries:  docsEntries(docs.Sizes),
		Fonts:        docs.Fonts,
	}
	sort.Strings(params.Fonts)
	templ.Handler(templates.DocsPage(params)).ServeHTTP(w, r)
}

func handleDocsJSON(w http.ResponseWriter, r *http.Request) {
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
		Font:         "ubuntu_mono",
		ColorOptions: pie.Keys(image.KnownColors()),
		SizeOptions:  pie.Keys(image.KnownSizes()),
		FontOptions:  pie.Keys(image.KnownFonts()),
		FontFiles:    image.KnownFontFilenames(),
	}
	templ.Handler(templates.OfflinePage(params)).ServeHTTP(w, r)
}
