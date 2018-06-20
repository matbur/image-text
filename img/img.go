package img

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func Draw(w io.Writer) {
	const width, height = 256, 256

	// Create a colored image of the given width and height.
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	addLabel(img, 100, 100, "qwerty")

	if err := png.Encode(w, img); err != nil {
		log.Fatal(err)
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{R: 15, G: 166, B: 194, A: 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x << 6), Y: fixed.Int26_6(y << 6)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
