package image_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/matbur/image-text/image"
	"github.com/stretchr/testify/require"
)

func TestImage(t *testing.T) {
	img, err := image.New("vga", "steel_blue", "yellow", "")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = img.Draw(&buf)
	require.NoError(t, err)
	require.NotEmpty(t, buf.Bytes())

	err = os.WriteFile("image.png", buf.Bytes(), 0644)
	require.NoError(t, err)
}
