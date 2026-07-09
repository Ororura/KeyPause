package tray

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

func IconPNG() []byte {
	img := image.NewNRGBA(image.Rect(0, 0, 18, 18))
	black := color.NRGBA{A: 255}

	fillRect(img, 2, 4, 16, 13, black)
	fillRect(img, 4, 6, 14, 11, color.NRGBA{})
	fillRect(img, 4, 13, 14, 15, black)
	fillRect(img, 7, 15, 11, 16, black)

	fillRect(img, 5, 7, 7, 9, black)
	fillRect(img, 8, 7, 10, 9, black)
	fillRect(img, 11, 7, 13, 9, black)
	fillRect(img, 5, 10, 13, 11, black)

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func fillRect(img *image.NRGBA, x0, y0, x1, y1 int, c color.NRGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.SetNRGBA(x, y, c)
		}
	}
}
