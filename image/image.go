package image

import (
	stdErr "errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"regexp"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	errorMissing    = stdErr.New("missing value")
	errorMalformed  = stdErr.New("malformed value")
	errorUnexpected = stdErr.New("unexpected error")

	sizePattern  = regexp.MustCompile(`^([1-9][0-9]{0,5})x([1-9][0-9]{0,5})$`)
	colorPattern = regexp.MustCompile(`^(?:[0-9A-Fa-f]{3}){1,2}$`)

	ubuntuMono = loadUbuntuMono()
)

type Image struct {
	Width      int
	Height     int
	Background color.Color
	Foreground color.Color
	Text       string
	*image.RGBA
}

func New(size, background, foreground, text string) (*Image, error) {
	width, height, err := parseSize(size)
	if err != nil {
		return nil, err
	}

	bg, err := parseColor(background)
	if err != nil {
		return nil, err
	}

	fg, err := parseColor(foreground)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	return &Image{
		Width:      width,
		Height:     height,
		Background: bg,
		Foreground: fg,
		Text:       text,
		RGBA:       rgba,
	}, nil
}

func (img *Image) Draw(w io.Writer) error {
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			img.Set(x, y, img.Background)
		}
	}

	if err := img.addLabel(); err != nil {
		return err
	}

	if err := png.Encode(w, img); err != nil {
		return err
	}

	return nil
}

func (img *Image) addLabel() error {
	size := 32
	face := loadFace(size)

	point := fixed.Point26_6{
		X: fixed.Int26_6(img.Width / 2 << 6),
		Y: fixed.Int26_6((img.Height/2 + size/4) << 6),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(img.Foreground),
		Dot:  point,
		Face: face,
	}
	d.DrawString(img.Text)

	return nil
}
