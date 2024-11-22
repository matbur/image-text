package image

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	errorMissing    = errors.New("missing value")
	errorMalformed  = errors.New("malformed value")
	errorUnexpected = errors.New("unexpected error")

	ubuntuMono = loadUbuntuMono()
)

type Image struct {
	size   Size
	bg     Color
	fg     Color
	text   string
	canvas *image.RGBA
}

func New(size, background, foreground, text string) (*Image, error) {
	s, err := NewSizeFromString(size)
	if err != nil {
		return nil, fmt.Errorf("failed to parse size: %w", err)
	}

	bg, err := NewBackgroundColorFromString(background)
	if err != nil {
		return nil, fmt.Errorf("failed to parse background color: %w", err)
	}

	fg, err := NewForegroundColorFromString(foreground)
	if err != nil {
		return nil, fmt.Errorf("failed to parse foreground color: %w", err)
	}

	if text == "" {
		text = size
	}

	canvas := image.NewRGBA(image.Rect(0, 0, s.Width(), s.Height()))

	return &Image{
		size:   s,
		bg:     bg,
		fg:     fg,
		text:   text,
		canvas: canvas,
	}, nil
}

func (img *Image) Draw(w io.Writer) error {
	for y := 0; y < img.size.Height(); y++ {
		for x := 0; x < img.size.Width(); x++ {
			img.canvas.Set(x, y, img.bg)
		}
	}

	img.addLabel()

	if err := png.Encode(w, img.canvas); err != nil {
		return fmt.Errorf("failed to encode png: %w", err)
	}
	return nil
}

func (img *Image) addLabel() {
	const scale = .95

	height := int(scale * float64(min(img.size.Height(), 2*img.size.Width()/len(img.text))))
	face := loadFace(height)

	width := font.MeasureString(face, img.text).Round()
	x := (img.size.Width() - width) / 2

	y := img.size.Height()/2 + height/4

	point := fixed.Point26_6{
		X: fixed.Int26_6(x << 6),
		Y: fixed.Int26_6(y << 6),
	}

	d := &font.Drawer{
		Dst:  img.canvas,
		Src:  image.NewUniform(img.fg),
		Dot:  point,
		Face: face,
	}
	d.DrawString(img.text)
}
