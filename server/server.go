package server

import (
	"io/ioutil"
	"net/http"

	"github.com/matbur/image-text/image"
)

func Foo() http.HandlerFunc {
	return dumpReq(
		checkMethod("GET")(
			handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Disposition", `inline; filename="fname.png"`)
	// fmt.Fprintf(w, "%+v", r)
}

func HandleFavicon(w http.ResponseWriter, r *http.Request) {
	bb, err := ioutil.ReadFile("res/favicon.png")
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bb)
}
