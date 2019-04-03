package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/matbur/image-text/image"
)

func HandleMain() http.HandlerFunc {
	return chain(
		dumpReq,
		checkMethod(http.MethodGet),
	)(handle)
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

	log.
		WithField("response", "binary").
		WithField("duration", time.Since(begin).String()).
		Infof("Response %d", http.StatusOK)
}

func HandleFavicon(w http.ResponseWriter, r *http.Request) {
	bb, err := ioutil.ReadFile("res/favicon.png")
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bb)
}

var docs = struct {
	Path     string            `json:"path"`
	Examples map[string]string `json:"example"`
	Colors   map[string]string `json:"colors"`
	Sizes    map[string]string `json:"sizes"`
}{
	Path: "HOST/size/background/foreground?text=rendered+text",
	Examples: map[string]string{
		"with_names": "http://localhost:8021/hd720/steel_blue/yellow?text=rendered+text",
		"with_codes": "http://localhost:8021/320x200/000/FFFF00",
	},
	Colors: image.Colors,
	Sizes:  image.Sizes,
}

func handleDocs(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(docs)
	if err != nil {
		msg := "Internal Server Error"
		log.WithError(err).Error(msg)
		writeJSON(w, msg, http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(js)
}
