package image

import (
	"fmt"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/matbur/image-text/resources"
)

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
