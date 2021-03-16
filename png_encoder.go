package pxl

import (
	"image"
	"image/png"
	"io"
)

type PngEncoder struct {
	w    io.Writer
	opts EncoderOptions
}

func NewPngEncoder(w io.Writer, opts EncoderOptions) PngEncoder {
	return PngEncoder{w, opts}
}

func (enc PngEncoder) Encode(pxl Pxl) error {
	img := image.NewNRGBA(image.Rect(0, 0, pxl.Cols()*enc.opts.Scale, pxl.Rows()*enc.opts.Scale))

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			clr := enc.opts.Bg

			if pxl.Get(col, row) {
				clr = enc.opts.Fg
			}
			for x := col * enc.opts.Scale; x < (col+1)*enc.opts.Scale; x++ {
				for y := row * enc.opts.Scale; y < (row+1)*enc.opts.Scale; y++ {
					img.Set(x, y, clr)
				}
			}
		}
	}
	return png.Encode(enc.w, img)
}
