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

	//t := 3 * time.Second
	//fmt.Fprint(w, "get after now ")
	//time.Sleep(t)
	//fmt.Fprintf(w, "get after %s", t)
	w.Header().Set("Content-Disposition", `inline; filename="fname.png"`)
	img.Draw(w)
}
