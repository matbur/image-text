package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matbur/image-text/image"
)

func Test_parseColor(t *testing.T) {
	tests := []struct {
		name       string
		arg        string
		r, g, b, a uint8
		color      image.Color
		err        bool
	}{
		{
			name:  "empty color",
			arg:   "",
			color: image.DefaultBackgroundColor(),
			err:   true,
		},
		{
			name:  "too few chars",
			arg:   "12",
			color: image.DefaultBackgroundColor(),
			err:   true,
		},
		{
			name:  "too many chars",
			arg:   "1234567",
			color: image.DefaultBackgroundColor(),
			err:   true,
		},
		{
			name:  "bad char",
			arg:   "12x",
			color: image.DefaultBackgroundColor(),
			err:   true,
		},
		{
			name:  "valid 3 char color",
			arg:   "123",
			color: image.NewColor(17, 34, 51),
		},
		{
			name:  "minimal valid 3 char color",
			arg:   "000",
			color: image.NewColor(0, 0, 0),
		},
		{
			name:  "max valid 3 char color",
			arg:   "fff",
			color: image.NewColor(255, 255, 255),
		},
		{
			name:  "valid 6 char color",
			arg:   "abcdef",
			color: image.NewColor(171, 205, 239),
		},
		{
			name:  "minimal valid 6 char color",
			arg:   "000000",
			color: image.NewColor(0, 0, 0),
		},
		{
			name:  "max valid 6 char color",
			arg:   "ffffff",
			color: image.NewColor(255, 255, 255),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := image.NewColorFromString(tt.arg)
			assert.Equal(t, tt.err, err != nil)
			assert.Equal(t, tt.color, color)
		})
	}
}

func TestNamedColor(t *testing.T) {
	color, err := image.NewColorFromString("steel_blue")
	require.NoError(t, err)
	assert.Equal(t, image.NewColor(70, 130, 180), color)
}

func TestKnownColorStrings(t *testing.T) {
	names := image.KnownColorStrings()
	assert.Contains(t, names, "steel_blue")
	assert.Equal(t, "4682b4", names["steel_blue"])
	assert.Len(t, names, len(image.KnownColors()))
}

func TestNewForegroundColorFromStringInvalid(t *testing.T) {
	color, err := image.NewForegroundColorFromString("not-a-color")
	require.Error(t, err)
	assert.Equal(t, image.DefaultForegroundColor(), color)
}
