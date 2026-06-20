package image_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matbur/image-text/image"
)

func fixtureName(size, bgColor, fgColor, text string) string {
	return size + "-" + bgColor + "-" + fgColor + "-" + text + ".png"
}

func TestDraw(t *testing.T) {
	tests := []struct {
		size    string
		bgColor string
		fgColor string
		text    string
		width   int
		height  int
	}{
		{
			size:    "vga",
			bgColor: "steel_blue",
			fgColor: "yellow",
			text:    "text",
			width:   640,
			height:  480,
		},
		{
			size:    "hd720",
			bgColor: "1c6",
			fgColor: "53f",
			text:    "",
			width:   1280,
			height:  720,
		},
		{
			size:    "700x300",
			bgColor: "278921",
			fgColor: "53f943",
			text:    "qwertyuiop",
			width:   700,
			height:  300,
		},
	}

	for _, tt := range tests {
		name := fixtureName(tt.size, tt.bgColor, tt.fgColor, tt.text)
		t.Run(name, func(t *testing.T) {
			img, err := image.New(tt.size, tt.bgColor, tt.fgColor, tt.text, "")
			require.NoError(t, err)

			var buf bytes.Buffer
			require.NoError(t, img.Draw(&buf))
			assertValidPNG(t, buf.Bytes(), tt.width, tt.height)
			assertMatchesFixture(t, name, buf.Bytes())
		})
	}
}

func TestNewDefaultsEmptyText(t *testing.T) {
	img, err := image.New("vga", "steel_blue", "yellow", "", "")
	require.NoError(t, err)
	assert.Equal(t, "640x480", img.Text())
}

func TestNewInvalidParams(t *testing.T) {
	img, err := image.New("bad", "bad", "bad", "hello", "bad")
	require.Error(t, err)
	require.NotNil(t, img)

	var buf bytes.Buffer
	require.NoError(t, img.Draw(&buf))
	assertValidPNG(t, buf.Bytes(), 640, 480)
}

func TestNewWithFont(t *testing.T) {
	img, err := image.New("320x200", "000", "fff", "hello", "open_sans")
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, img.Draw(&buf))
	assertValidPNG(t, buf.Bytes(), 320, 200)
}

func TestImageAccessors(t *testing.T) {
	img, err := image.New("vga", "steel_blue", "yellow", "hello", "ubuntu_mono")
	require.NoError(t, err)

	assert.Equal(t, "hello", img.Text())
	assert.Equal(t, "640x480", img.Size())
	assert.Equal(t, "4682b4", img.BgColor())
	assert.Equal(t, "ffff00", img.FgColor())
}
