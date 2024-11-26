package server_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matbur/image-text/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImage(t *testing.T) {
	tests := []struct {
		size, bgColor, fgColor, text string
	}{
		{
			size:    "vga",
			bgColor: "steel_blue",
			fgColor: "yellow",
			text:    "text",
		},
		{
			size:    "hd720",
			bgColor: "1c6",
			fgColor: "53f",
			text:    "",
		},
		{
			size:    "700x300",
			bgColor: "278921",
			fgColor: "53f943",
			text:    "qwertyuiop",
		},
	}

	for _, tt := range tests {
		name := tt.size + "-" + tt.bgColor + "-" + tt.fgColor + "-" + tt.text + ".png"
		t.Run(name, func(t *testing.T) {

			u := fmt.Sprintf("/%s/%s/%s?text=%s", tt.size, tt.bgColor, tt.fgColor, tt.text)
			r := httptest.NewRequest("GET", u, nil)

			rr := httptest.NewRecorder()
			server.NewServer().ServeHTTP(rr, r)

			var buf bytes.Buffer
			_, err := buf.ReadFrom(rr.Result().Body)
			require.NoError(t, err)

			bb, err := os.ReadFile("../fixtures/" + name)
			require.NoError(t, err)

			if !assert.Equal(t, bb, buf.Bytes()) {
				err := os.WriteFile(name, buf.Bytes(), 0644)
				require.NoError(t, err)
			}
		})
	}
}
