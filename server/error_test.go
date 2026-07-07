package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	writeJSON(rr, CodeBadRequest, "bad request", http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"error":{"code":"BAD_REQUEST","message":"bad request"}}`, rr.Body.String())
}
