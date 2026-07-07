package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	CodeInvalidSize   = "INVALID_SIZE"
	CodeInvalidColor  = "INVALID_COLOR"
	CodeInvalidFont   = "INVALID_FONT"
	CodeEncodeError   = "ENCODE_ERROR"
	CodeBadRequest    = "BAD_REQUEST"
	CodeInternalError = "INTERNAL_ERROR"
	CodeNotFound      = "NOT_FOUND"
)

type apiErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errResponse struct {
	Error apiErr `json:"error"`
}

func writeJSON(w http.ResponseWriter, code string, message string, httpCode int) {
	r := errResponse{Error: apiErr{Code: code, Message: message}}
	js, err := json.Marshal(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.Error("Failed to marshal response", "err", err, "response", r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	if _, err := w.Write(js); err != nil {
		slog.Error("Failed to write response", "err", err)
	}

	slog.Info("Response", "status", httpCode, "response", string(js))
}
