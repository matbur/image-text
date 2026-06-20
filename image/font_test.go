package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matbur/image-text/image"
)

func TestKnownFonts(t *testing.T) {
	fonts := image.KnownFonts()
	require.NotEmpty(t, fonts)
	assert.Contains(t, fonts, "ubuntu_mono")
	assert.Contains(t, fonts, "roboto")
	assert.Contains(t, fonts, "open_sans")
}

func TestNewFontFromString(t *testing.T) {
	t.Run("known font", func(t *testing.T) {
		fnt, err := image.NewFontFromString("roboto")
		require.NoError(t, err)
		require.NotNil(t, fnt)
	})

	t.Run("default font", func(t *testing.T) {
		fnt, err := image.NewFontFromString("")
		require.NoError(t, err)
		require.NotNil(t, fnt)
	})

	t.Run("case insensitive", func(t *testing.T) {
		fnt, err := image.NewFontFromString("UBUNTU_MONO")
		require.NoError(t, err)
		require.NotNil(t, fnt)
	})

	t.Run("unknown font", func(t *testing.T) {
		_, err := image.NewFontFromString("not_a_font")
		require.Error(t, err)
	})
}
