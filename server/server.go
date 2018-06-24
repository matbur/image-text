package server

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/matbur/image-text/image"
)

func HandleMain() http.HandlerFunc {
	return dumpReq(checkMethod(http.MethodGet)(handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
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
	log.WithField("response", "binary").Infof("Response %d", http.StatusOK)
}

func HandleFavicon(w http.ResponseWriter, r *http.Request) {
	bb, err := ioutil.ReadFile("res/favicon.png")
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bb)
}
