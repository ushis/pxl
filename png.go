package pxl

import (
	"image"
	"image/color"
	"image/png"
	"io"
)

func EncodePng(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	img := image.NewNRGBA(image.Rect(0, 0, pxl.Cols()*opts.Scale, pxl.Rows()*opts.Scale))

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			var clr color.Color

			if pxl.Get(col, row) {
				clr = opts.Fg
			} else {
				clr = opts.Bg
			}
			for x := col * opts.Scale; x < (col+1)*opts.Scale; x++ {
				for y := row * opts.Scale; y < (row+1)*opts.Scale; y++ {
					img.Set(x, y, clr)
				}
			}
		}
	}
	return png.Encode(w, img)
}
