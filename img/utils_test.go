package img

import (
	"github.com/pkg/errors"
	"image/color"
	"reflect"
	"testing"
)

func Test_parseSize(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		width  int
		height int
		err    error
	}{
		{
			name: "empty size",
			arg:  "",
			err:  errorMissing,
		}, {
			name: "too many parts",
			arg:  "234x23x243",
			err:  errorMalformed,
		}, {
			name: "with float",
			arg:  "200x300.243",
			err:  errorMalformed,
		}, {
			name: "with 0",
			arg:  "200x0",
			err:  errorMalformed,
		}, {
			name: "with too big value",
			arg:  "1x1000000",
			err:  errorMalformed,
		}, {
			name:   "with big value",
			arg:    "1x999999",
			width:  1,
			height: 999999,
		}, {
			name:   "valid size",
			arg:    "200x300",
			width:  200,
			height: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseSize(tt.arg)
			if errors.Cause(err) != tt.err {
				t.Errorf("parseSize() error = %v, err %v", err, tt.err)
				return
			}
			if got != tt.width {
				t.Errorf("parseSize() got = %v, width %v", got, tt.width)
			}
			if got1 != tt.height {
				t.Errorf("parseSize() got1 = %v, height %v", got1, tt.height)
			}
		})
	}
}

func Test_parseColor(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want *color.RGBA
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
			want: &color.RGBA{R: 17, G: 34, B: 51, A: 255},
		}, {
			name: "minimal valid 3 char color",
			arg:  "000",
			want: &color.RGBA{R: 0, G: 0, B: 0, A: 255},
		}, {
			name: "max valid 3 char color",
			arg:  "fff",
			want: &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}, {
			name: "valid 6 char color",
			arg:  "abcdef",
			want: &color.RGBA{R: 171, G: 205, B: 239, A: 255},
		}, {
			name: "minimal valid 6 char color",
			arg:  "000000",
			want: &color.RGBA{R: 0, G: 0, B: 0, A: 255},
		}, {
			name: "max valid 6 char color",
			arg:  "ffffff",
			want: &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseColor(tt.arg)
			if errors.Cause(err) != tt.err {
				t.Errorf("parseColor() error = %v, err %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseText(t *testing.T) {
	var tests []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseText(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseText() = %v, want %v", got, tt.want)
			}
		})
	}
}
