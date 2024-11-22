package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matbur/image-text/image"
)

func Test_parseSize(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		width  int
		height int
		err    bool
	}{
		{
			name:   "empty size",
			arg:    "",
			width:  image.DefaultSize().Width(),
			height: image.DefaultSize().Height(),
			err:    true,
		},
		{
			name:   "too many parts",
			arg:    "234x23x243",
			width:  image.DefaultSize().Width(),
			height: image.DefaultSize().Height(),
			err:    true,
		},
		{
			name:   "with float",
			arg:    "200x300.243",
			width:  image.DefaultSize().Width(),
			height: image.DefaultSize().Height(),
			err:    true,
		},
		{
			name:   "with 0",
			arg:    "200x0",
			width:  image.DefaultSize().Width(),
			height: image.DefaultSize().Height(),
			err:    true,
		},
		{
			name:   "with too big value",
			arg:    "1x1000000",
			width:  image.DefaultSize().Width(),
			height: image.DefaultSize().Height(),
			err:    true,
		},
		{
			name:   "with big value",
			arg:    "1x999999",
			width:  1,
			height: 999999,
		},
		{
			name:   "valid size",
			arg:    "200x300",
			width:  200,
			height: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := image.NewSizeFromString(tt.arg)
			assert.Equal(t, tt.err, err != nil)
			assert.Equal(t, tt.width, size.Width())
			assert.Equal(t, tt.height, size.Height())
		})
	}
}
