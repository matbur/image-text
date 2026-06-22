package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matbur/image-text/server"
)

func TestIntegrationStaticFontEncodedComma(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/resources/fonts/Roboto-VariableFont_wdth%2Cwght.ttf", nil)

	server.NewServer(server.Config{}).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	require.Greater(t, len(rr.Body.Bytes()), 1000)
}

func TestIntegrationOfflinePageFontFilesJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/offline", nil)

	server.NewServer(server.Config{}).ServeHTTP(rr, req)

	body := rr.Body.String()
	assert.Contains(t, body, "font-files-data")
	assert.Contains(t, body, "roboto")
	assert.Contains(t, body, "Roboto-VariableFont")
	assert.Contains(t, body, `getElementById("font-files-data")`)
	require.NotContains(t, body, "{ params.FontFilesJSON }")
}

func TestIntegrationStaticFontBitcount(t *testing.T) {
	rr := httptest.NewRecorder()
	filename := "BitcountPropSingle-VariableFont_CRSV,ELSH,ELXP,slnt,wght.ttf"
	req := httptest.NewRequest(http.MethodGet, "/resources/fonts/"+strings.ReplaceAll(filename, ",", "%2C"), nil)

	server.NewServer(server.Config{}).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	require.Greater(t, len(rr.Body.Bytes()), 1000)
}
