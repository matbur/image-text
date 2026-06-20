package image_test

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertValidPNG(t *testing.T, data []byte, width, height int) {
	t.Helper()

	require.NotEmpty(t, data)
	require.Equal(t, []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}, data[:8])

	img, err := png.Decode(bytes.NewReader(data))
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
