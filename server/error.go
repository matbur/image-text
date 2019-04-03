package server

import (
	"encoding/json"
	"net/http"

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
