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
	fnt, err := image.NewFontFromString("roboto")
	require.NoError(t, err)
	require.NotNil(t, fnt)

	_, err = image.NewFontFromString("not_a_font")
	require.Error(t, err)
}
