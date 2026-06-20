package image_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matbur/image-text/image"
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
			img, err := image.New(tt.size, tt.bgColor, tt.fgColor, tt.text, "")
			require.NoError(t, err)

			var buf bytes.Buffer
			err = img.Draw(&buf)
			require.NoError(t, err)
			require.NotEmpty(t, buf.Bytes())

			bb, err := os.ReadFile("../fixtures/" + name)
			require.NoError(t, err)

			if !assert.Equal(t, bb, buf.Bytes()) {
				err := os.WriteFile(name, buf.Bytes(), 0644)
				require.NoError(t, err)
			}
		})
	}
}
