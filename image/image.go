package image

import (
	"image"
	"image/color"
	"image/png"
	"io"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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
	// Create a colored image of the given width and height.

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			img.Set(x, y, img.Background)
		}
	}

	img.addLabel()

	if err := png.Encode(w, img); err != nil {
		return err
	}

	return nil
}

func (img *Image) addLabel() {
	point := fixed.Point26_6{
		X: fixed.Int26_6(img.Width / 2 << 6),
		Y: fixed.Int26_6(img.Height / 2 << 6),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(img.Foreground),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(img.Text)
}
