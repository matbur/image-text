package img

import (
	"image/color"
	"reflect"
	"testing"
	"github.com/pkg/errors"
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
			name: "size cannot be empty",
			arg:  "",
			err:  ErrorEmptySize,
		}, {
			name: "size should contain exactly 2 parts",
			arg:  "234x23x243",
			err:  ErrorMalformedSize,
		}, {
			name: "each part should be an integer",
			arg:  "200x300.243",
			err:  ErrorMalformedSize,
		}, {
			name: "size must contain positive values",
			arg:  "200x0",
			err:  ErrorZeroValue,
		}, {
			name:   "should return proper size",
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
		name    string
		args    string
		want    color.RGBA
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseColor(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseColor() = %v, width %v", got, tt.want)
			}
		})
	}
}

func Test_parseText(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseText(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseText() = %v, width %v", got, tt.want)
			}
		})
	}
}
