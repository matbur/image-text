package server

import (
	"net/http"

	"github.com/matbur/image-text/img"
)

func Foo() http.HandlerFunc {
	return dumpReq(
		checkMethod("GET")(
			handle))
}

func handle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Disposition", `inline; filename="fname.png"`)
	img.Draw(w)
	// fmt.Fprintf(w, "%+v", r)
}
