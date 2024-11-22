package image

import (
	"errors"
	"image/color"
	"reflect"
	"testing"
)

func Test_parseColor(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want color.Color
		err  error
	}{
		{
			name: "empty color",
			arg:  "",
			err:  errorMissing,
		}, {
			name: "too few chars",
			arg:  "12",
			err:  errorMalformed,
		}, {
			name: "too many chars",
			arg:  "1234567",
			err:  errorMalformed,
		}, {
			name: "bad char",
			arg:  "12x",
			err:  errorMalformed,
		}, {
			name: "too few chars",
			arg:  "12",
			err:  errorMalformed,
		}, {
			name: "valid 3 char color",
			arg:  "123",
			want: color.RGBA{R: 17, G: 34, B: 51, A: 255},
		}, {
			name: "minimal valid 3 char color",
			arg:  "000",
			want: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		}, {
			name: "max valid 3 char color",
			arg:  "fff",
			want: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}, {
			name: "valid 6 char color",
			arg:  "abcdef",
			want: color.RGBA{R: 171, G: 205, B: 239, A: 255},
		}, {
			name: "minimal valid 6 char color",
			arg:  "000000",
			want: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		}, {
			name: "max valid 6 char color",
			arg:  "ffffff",
			want: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseColor(tt.arg)
			if !errors.Is(err, tt.err) {
				t.Errorf("parseColor() error = %v, err %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
