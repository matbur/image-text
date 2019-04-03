package server

import (
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/sirupsen/logrus"
)

type interceptor func(http.HandlerFunc) http.HandlerFunc

func chain(fns ...interceptor) interceptor {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		for i := len(fns); i > 0; i-- {
			fn = fns[i-1](fn)
		}
		return fn
	}
}

func dumpReq(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dumped, err := httputil.DumpRequest(r, false)
		if err != nil {
			log.WithField("request", r).Infof("failed to dump request: %v", err)
		} else {
			log.WithField("request", string(dumped)).Infof("%s %s", r.Method, r.URL)
		}

		h(w, r)
	}
}

func checkMethod(methods ...string) interceptor {
	msg := "Expected " + strings.Join(methods, ", ")
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !isIn(r.Method, methods) {
				writeJSON(w, msg, http.StatusMethodNotAllowed)
				return
			}

			h(w, r)
		}
	}
}
