package image

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/matbur/image-text/resources"
)

func parseSize(s string) (int, int, error) {
	if s == "" {
		return 0, 0, fmt.Errorf("size is empty: %w", errorMissing)
	}

	if s, ok := Sizes[s]; ok {
		return parseSize(s)
	}

	ss := sizePattern.FindStringSubmatch(s)
	if len(ss) != 3 {
		return 0, 0, fmt.Errorf("size '%s' is not valid: %w", s, errorMalformed)
	}

	width, err := strconv.Atoi(ss[1])
	if err != nil {
		return 0, 0, fmt.Errorf("bad width '%s': %v: %w", ss[1], err, errorUnexpected)
	}

	height, err := strconv.Atoi(ss[2])
	if err != nil {
		return 0, 0, fmt.Errorf("bad height '%s': %v: %w", ss[2], err, errorUnexpected)
	}

	return width, height, nil
}

func parseColor(s string) (color.Color, error) {
	if s == "" {
		return nil, fmt.Errorf("color is empty: %w", errorMissing)
	}

	if v, ok := Colors[s]; ok {
		return parseColor(v)
	}

	if !colorPattern.MatchString(s) {
		return nil, fmt.Errorf("color '%s' is not valid: %w", s, errorMalformed)
	}

	var ss []string
	switch len(s) {
	case 3:
		ss = strings.Split(s, "")
		ss[0] += ss[0]
		ss[1] += ss[1]
		ss[2] += ss[2]
	case 6:
		ss = []string{s[0:2], s[2:4], s[4:6]}
	default:
		return nil, fmt.Errorf("bad color '%s': %w", s, errorUnexpected)
	}

	red, err := strconv.ParseUint(ss[0], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("bad red '%s': %v: %w", ss[0], err, errorUnexpected)
	}

	green, err := strconv.ParseUint(ss[1], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("bad green '%s': %v: %w", ss[1], err, errorUnexpected)
	}

	blue, err := strconv.ParseUint(ss[2], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("bad blue '%s': %v: %w", ss[2], err, errorUnexpected)
	}

	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: 255,
	}, nil
}

func loadFont(fn string) (*truetype.Font, error) {
	bb, err := resources.Static.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %v: %w", fn, err, errorUnexpected)
	}

	fnt, err := truetype.Parse(bb)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font '%s': %v: %w", fn, err, errorUnexpected)
	}

	return fnt, nil
}

func loadUbuntuMono() *truetype.Font {
	fc, err := loadFont("UbuntuMono-Regular.ttf")
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
