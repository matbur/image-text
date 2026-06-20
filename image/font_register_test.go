package image_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matbur/image-text/image"
)

func TestRegisterFont(t *testing.T) {
	data, err := os.ReadFile("../resources/fonts/ComicNeue-Regular.ttf")
	require.NoError(t, err)

	err = image.RegisterFont("comic_neue_test", data)
	require.NoError(t, err)

	fnt, err := image.NewFontFromString("comic_neue_test")
	require.NoError(t, err)
	require.NotNil(t, fnt)
}

func TestKnownFontFilenames(t *testing.T) {
	filenames := image.KnownFontFilenames()
	require.NotEmpty(t, filenames)
	assert.Equal(t, "UbuntuMono-Regular.ttf", filenames["ubuntu_mono"])
	assert.Contains(t, filenames, "roboto")
}
