package image

import (
	stdErr "errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"regexp"

	"github.com/pkg/errors"
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
	Canvas     *image.RGBA
}

func New(size, background, foreground, text string) (*Image, error) {
	width, height, err := parseSize(size)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse size")
	}

	bg, err := parseColor(background)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse background color")
	}

	fg, err := parseColor(foreground)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse foreground color")
	}

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	return &Image{
		Width:      width,
		Height:     height,
		Background: bg,
		Foreground: fg,
		Text:       text,
		Canvas:     rgba,
	}, nil
}

func (img *Image) Draw(w io.Writer) error {
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			img.Canvas.Set(x, y, img.Background)
		}
	}

	img.addLabel()

	if err := png.Encode(w, img.Canvas); err != nil {
		return errors.Wrap(err, "failed to encode png")
	}

	return nil
}

func (img *Image) addLabel() {
	const scale = .95

	height := int(scale * float64(minInt(img.Height, 2*img.Width/len(img.Text))))
	face := loadFace(height)

	width := font.MeasureString(face, img.Text).Round()
	x := (img.Width - width) / 2

	y := img.Height/2 + height/4

	point := fixed.Point26_6{
		X: fixed.Int26_6(x << 6),
		Y: fixed.Int26_6(y << 6),
	}

	d := &font.Drawer{
		Dst:  img.Canvas,
		Src:  image.NewUniform(img.Foreground),
		Dot:  point,
		Face: face,
	}
	d.DrawString(img.Text)
}
