package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matbur/image-text/server"
)

func TestIntegrationIndexPage(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "text/html")
	body := rr.Body.String()
	assert.Contains(t, body, "Online")
	assert.Contains(t, body, "Offline")
	assert.Contains(t, body, "API docs")
}

func TestIntegrationOfflinePage(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/offline", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Generate image")
}

func TestIntegrationOnlinePageRedirect(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/online", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusSeeOther, rr.Code)
	location := rr.Header().Get("Location")
	assert.Contains(t, location, "/online?")
	assert.Contains(t, location, "bg_color=steel_blue")
	assert.Contains(t, location, "fg_color=yellow")
	assert.Contains(t, location, "size=vga")
}

func TestIntegrationOnlinePageWithParams(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/online?text=hello&bg_color=steel_blue&fg_color=yellow&size=vga&font=ubuntu_mono", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.Contains(t, body, "hello")
	assert.Contains(t, body, "Settings")
}

func TestIntegrationOnlinePost(t *testing.T) {
	body := bytes.NewBufferString(`{"text":"hello","bg_color":"steel_blue","fg_color":"yellow","size":"vga","font":"ubuntu_mono"}`)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/online/post", body)
	req.Header.Set("Content-Type", "application/json")

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("HX-Push-Url"), "/online?")
	assert.Contains(t, rr.Body.String(), "image-result")
}

func TestIntegrationOnlinePostHTMXPartial(t *testing.T) {
	body := bytes.NewBufferString(`{"text":"hello","bg_color":"steel_blue","fg_color":"yellow","size":"vga","font":"ubuntu_mono"}`)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/online/post", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	response := rr.Body.String()
	assert.Contains(t, response, "image-result")
	assert.NotContains(t, response, "Settings")
}

func TestIntegrationImageBadRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vga/steel_blue/yellow?font=not_a_font", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var payload struct {
		Error string `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &payload))
	assert.Contains(t, payload.Error, "font")
}

func TestIntegrationImageHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vga/steel_blue/yellow?text=hello", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `inline; filename="image.png"`, rr.Header().Get("Content-Disposition"))
	assertPNGSignature(t, rr.Body.Bytes())
}

func TestIntegrationImageWithFont(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vga/steel_blue/yellow?text=hello&font=open_sans", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assertPNGSignature(t, rr.Body.Bytes())
}

func TestIntegrationDocsJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/docs.json", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var docs struct {
		Path     string            `json:"path"`
		Params   map[string]string `json:"params"`
		Examples map[string]string `json:"example"`
		Colors   map[string]string `json:"colors"`
		Sizes    map[string]string `json:"sizes"`
		Fonts    []string          `json:"fonts"`
	}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &docs))

	assert.Contains(t, docs.Path, "size/background/foreground")
	assert.Contains(t, docs.Params, "text")
	assert.Contains(t, docs.Params, "font")
	assert.Contains(t, docs.Examples, "with_names")
	assert.Contains(t, docs.Colors, "steel_blue")
	assert.Contains(t, docs.Sizes, "vga")
	assert.Contains(t, docs.Fonts, "ubuntu_mono")
}

func TestIntegrationDocsPage(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.Contains(t, body, "API documentation")
	assert.Contains(t, body, "steel_blue")
	assert.Contains(t, body, "ubuntu_mono")
}

func TestIntegrationFavicon(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Body.Bytes())
}

func TestIntegrationStaticFont(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/resources/fonts/Roboto-VariableFont_wdth,wght.ttf", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, strings.HasSuffix(strings.ToLower(rr.Header().Get("Content-Type")), "font/ttf") ||
		strings.Contains(rr.Header().Get("Content-Type"), "octet-stream"))
	assert.NotEmpty(t, rr.Body.Bytes())
}

func TestIntegrationStaticWASM(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/resources/main.wasm", nil)

	server.NewServer().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Body.Bytes())
}

func assertPNGSignature(t *testing.T, data []byte) {
	t.Helper()
	require.NotEmpty(t, data)
	assert.Equal(t, []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}, data[:8])
}
