package image

import (
	"image/color"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

func parseSize(s string) (int, int, error) {
	if s == "" {
		return 0, 0, errors.Wrap(errorMissing, "size is empty")
	}

	if s, ok := sizes[s]; ok {
		return parseSize(s)
	}

	ss := sizePattern.FindStringSubmatch(s)
	if len(ss) != 3 {
		return 0, 0, errors.Wrapf(errorMalformed, "size '%s' is not valid", s)
	}

	width, err := strconv.Atoi(ss[1])
	if err != nil {
		return 0, 0, errors.Wrapf(errorUnexpected, "bad width '%s': %v", ss[1], err)
	}

	height, err := strconv.Atoi(ss[2])
	if err != nil {
		return 0, 0, errors.Wrapf(errorUnexpected, "bad height '%s': %v", ss[2], err)
	}

	return width, height, nil
}

func parseColor(s string) (color.Color, error) {
	if s == "" {
		return nil, errors.Wrap(errorMissing, "color is empty")
	}

	if v, ok := colors[s]; ok {
		return parseColor(v)
	}

	if !colorPattern.MatchString(s) {
		return nil, errors.Wrapf(errorMalformed, "color '%s' is not valid", s)
	}

	var ss []string
	if n := len(s); n == 3 {
		ss = strings.Split(s, "")
		ss[0] += ss[0]
		ss[1] += ss[1]
		ss[2] += ss[2]
	} else if n == 6 {
		ss = []string{s[0:2], s[2:4], s[4:6]}
	} else {
		return nil, errors.Wrapf(errorUnexpected, "bad color '%s'", s)
	}

	red, err := strconv.ParseUint(ss[0], 16, 8)
	if err != nil {
		return nil, errors.Wrapf(errorUnexpected, "bad red '%s': %v", ss[0], err)
	}

	green, err := strconv.ParseUint(ss[1], 16, 8)
	if err != nil {
		return nil, errors.Wrapf(errorUnexpected, "bad green '%s': %v", ss[1], err)
	}

	blue, err := strconv.ParseUint(ss[2], 16, 8)
	if err != nil {
		return nil, errors.Wrapf(errorUnexpected, "bad blue '%s': %v", ss[2], err)
	}

	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: 255,
	}, nil
}

func loadFont(fn string) (*truetype.Font, error) {
	bb, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, errors.Wrapf(errorUnexpected, "failed to read file '%s': %v", fn, err)
	}

	fnt, err := truetype.Parse(bb)
	if err != nil {
		return nil, errors.Wrapf(errorUnexpected, "failed to parse font: %v", err)
	}

	return fnt, nil
}

func loadUbuntuMono() *truetype.Font {
	fc, err := loadFont("res/UbuntuMono-Regular.ttf")
	if err != nil {
		panic(err)
	}

	return fc
}

func loadFace(size int) font.Face {
	opts := &truetype.Options{
		Size: float64(size),
	}

	return truetype.NewFace(ubuntuMono, opts)
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
