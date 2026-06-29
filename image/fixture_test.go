package image_test

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var formatMagic = map[string][]byte{
	"png":  {0x89, 0x50, 0x4e, 0x47},
	"jpg":  {0xff, 0xd8, 0xff},
	"webp": {0x52, 0x49, 0x46, 0x46},
}

func assertValidImage(t *testing.T, format string, data []byte, width, height int) {
	t.Helper()
	require.NotEmpty(t, data)
	magic, ok := formatMagic[format]
	require.True(t, ok, "unknown format %s", format)
	require.Equal(t, magic, data[:len(magic)])

	img, _, err := image.Decode(bytes.NewReader(data))
	require.NoError(t, err)

	bounds := img.Bounds()
	assert.Equal(t, width, bounds.Dx())
	assert.Equal(t, height, bounds.Dy())
}

func assertMatchesFixture(t *testing.T, name string, got []byte) {
	t.Helper()

	want, err := os.ReadFile("../fixtures/" + name)
	require.NoError(t, err, "missing fixture %s", name)
	assert.Equal(t, want, got)
}
