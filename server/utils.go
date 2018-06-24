package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/sirupsen/logrus"
)

type errResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, error string, code int) {
	r := errResponse{Error: error}
	js, err := json.Marshal(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.WithField("response", r).Errorf("Failed to marshal response: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(js); err != nil {
		log.Errorf("Failed to write response: %v", err)
	}

	log.WithField("response", string(js)).Infof("Response %d", code)
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

func checkMethod(methods ...string) func(http.HandlerFunc) http.HandlerFunc {
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

func isIn(s string, ss []string) bool {
	for _, i := range ss {
		if s == i {
			return true
		}
	}
	return false
}

func parsePath(s string) (size, bg, fg string, err error) {
	s = strings.Trim(s, "/")
	ss := strings.Split(s, "/")
	if len(ss) != 3 {
		return "", "", "", fmt.Errorf("malformed path '%s'", s)
	}
	return ss[0], ss[1], ss[2], nil
}
