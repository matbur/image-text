package image

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/skrashevich/go-webp"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
)

type Format string

const (
	FormatPNG  Format = "png"
	FormatJPG  Format = "jpg"
	FormatWebP Format = "webp"
)

func KnownFormats() []Format {
	return []Format{FormatPNG, FormatJPG, FormatWebP}
}

func KnownFormatStrings() []string {
	return []string{string(FormatPNG), string(FormatJPG), string(FormatWebP)}
}

var (
	errorMissing    = errors.New("missing value")
	errorMalformed  = errors.New("malformed value")
	errorUnexpected = errors.New("unexpected error")
)

type Image struct {
	size   Size
	bg     Color
	fg     Color
	font   *truetype.Font
	text   string
	format Format
	canvas *image.RGBA
}

func New(size, background, foreground, text, fontName, format string) (*Image, error) {
	var err error

	s, err := NewSizeFromString(size)
	if err != nil {
		err = multierror.Append(err, fmt.Errorf("failed to create size: %w", err))
	}

	bg, err := NewBackgroundColorFromString(background)
	if err != nil {
		err = multierror.Append(err, fmt.Errorf("failed to create background color: %w", err))
	}

	fg, err := NewForegroundColorFromString(foreground)
	if err != nil {
		err = multierror.Append(err, fmt.Errorf("failed to create foreground color: %w", err))
	}

	fnt, err := NewFontFromString(fontName)
	if err != nil {
		err = multierror.Append(err, fmt.Errorf("failed to create font: %w", err))
	}

	if text == "" {
		text = s.String()
	}

	f := parseFormat(format)

	canvas := image.NewRGBA(image.Rect(0, 0, s.Width(), s.Height()))

	return &Image{
		size:   s,
		bg:     bg,
		fg:     fg,
		font:   fnt,
		text:   text,
		format: f,
		canvas: canvas,
	}, err
}

func parseFormat(s string) Format {
	switch s {
	case "jpg", "jpeg":
		return FormatJPG
	case "webp":
		return FormatWebP
	default:
		return FormatPNG
	}
}

func ContentType(f string) string {
	switch parseFormat(f) {
	case FormatJPG:
		return "image/jpeg"
	case FormatWebP:
		return "image/webp"
	default:
		return "image/png"
	}
}

func Extension(f string) string {
	switch parseFormat(f) {
	case FormatJPG:
		return "jpg"
	case FormatWebP:
		return "webp"
	default:
		return "png"
	}
}

func (img *Image) Draw(w io.Writer) error {
	img.drawBackground()
	img.drawLabel()

	switch img.format {
	case FormatJPG:
		if err := jpeg.Encode(w, img.canvas, &jpeg.Options{Quality: 85}); err != nil {
			return fmt.Errorf("failed to encode jpeg: %w", err)
		}
	case FormatWebP:
		if err := webp.Encode(w, img.canvas, &webp.Options{Lossy: true, Quality: 85}); err != nil {
			return fmt.Errorf("failed to encode webp: %w", err)
		}
	default:
		if err := png.Encode(w, img.canvas); err != nil {
			return fmt.Errorf("failed to encode png: %w", err)
		}
	}
	return nil
}

func (img *Image) drawBackground() {
	for y := 0; y < img.size.Height(); y++ {
		for x := 0; x < img.size.Width(); x++ {
			img.canvas.Set(x, y, img.bg)
		}
	}
}

func (img *Image) drawLabel() {
	const scale = .95

	height := int(scale * float64(min(img.size.Height(), 2*img.size.Width()/max(len(img.text), 1))))
	face := loadFace(img.font, height)

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

func (img *Image) Text() string    { return img.text }
func (img *Image) Size() string    { return img.size.String() }
func (img *Image) BgColor() string { return img.bg.String() }
func (img *Image) FgColor() string { return img.fg.String() }
func (img *Image) Format() Format  { return img.format }
func (img *Image) ContentType() string {
	switch img.format {
	case FormatJPG:
		return "image/jpeg"
	case FormatWebP:
		return "image/webp"
	default:
		return "image/png"
	}
}
func (img *Image) Extension() string {
	switch img.format {
	case FormatJPG:
		return "jpg"
	case FormatWebP:
		return "webp"
	default:
		return "png"
	}
}
