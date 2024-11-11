package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, error string, code int) {
	r := errResponse{Error: error}
	js, err := json.Marshal(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.Error("Failed to marshal response", "err", err, "response", r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(js); err != nil {
		slog.Error("Failed to write response", "err", err)
	}

	slog.Info("Response", "status", code, "response", string(js))
}
