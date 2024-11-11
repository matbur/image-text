package server

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"
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

func dumpReq(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dumped, err := httputil.DumpRequest(r, false)
		if err != nil {
			slog.Error("Failed to dump request", "request", r, "err", err)
		} else {
			slog.Info("Request", "request", string(dumped))
		}

		next(w, r)
	}
}

func checkMethod(methods ...string) interceptor {
	msg := "Expected " + strings.Join(methods, ", ")
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !isIn(r.Method, methods) {
				writeJSON(w, msg, http.StatusMethodNotAllowed)
				return
			}

			next(w, r)
		}
	}
}
