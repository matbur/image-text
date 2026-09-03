package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRoutePattern(t *testing.T) {
	t.Run("no chi route context", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/anything", nil)
		assert.Equal(t, "unknown", routePattern(r))
	})

	t.Run("matched route returns pattern, not raw path", func(t *testing.T) {
		var got string

		router := chi.NewRouter()
		router.Get("/{size}/{bg_color}/{fg_color}", func(w http.ResponseWriter, r *http.Request) {
			got = routePattern(r)
		})

		req := httptest.NewRequest(http.MethodGet, "/640x480/steel_blue/yellow", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, "/{size}/{bg_color}/{fg_color}", got)
	})

	t.Run("unmatched route is bucketed as not_found", func(t *testing.T) {
		var got string

		router := chi.NewRouter()
		// A real route must be registered so chi builds its handler chain
		// and sets a RouteContext even for requests that don't match
		// (mirrors NewServer, which always has routes registered).
		router.Get("/known", func(w http.ResponseWriter, r *http.Request) {})
		router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			got = routePattern(r)
		})

		req := httptest.NewRequest(http.MethodGet, "/does/not/exist", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, "not_found", got)
	})
}
