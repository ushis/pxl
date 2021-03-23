package pxl

import (
	"image"
	"image/color"
	"image/gif"
	"io"
)

type GifEncoder struct {
	w    io.Writer
	opts EncoderOptions
}

func NewGifEncoder(w io.Writer, opts EncoderOptions) GifEncoder {
	return GifEncoder{w, opts}
}

func (enc GifEncoder) Encode(pxl Pxl) error {
	g := &gif.GIF{}

	if enc.opts.Fps > 0 {
		g.Image = []*image.Paletted{enc.encode(pxl, false), enc.encode(pxl, true)}
		g.Delay = []int{100.0 / enc.opts.Fps, 100.0 / enc.opts.Fps}
	} else {
		g.Image = []*image.Paletted{enc.encode(pxl, false)}
		g.Delay = []int{0}
	}
	return gif.EncodeAll(enc.w, g)
}

func (enc GifEncoder) encode(pxl Pxl, invert bool) *image.Paletted {
	rct := image.Rect(0, 0, pxl.Cols()*enc.opts.Scale, pxl.Rows()*enc.opts.Scale)
	plt := color.Palette{enc.opts.Bg, enc.opts.Fg}
	img := image.NewPaletted(rct, plt)

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			fg := pxl.Get(col, row)

			if invert {
				fg = !fg
			}
			clr := enc.opts.Bg

			if fg {
				clr = enc.opts.Fg
			}
			for x := col * enc.opts.Scale; x < (col+1)*enc.opts.Scale; x++ {
				for y := row * enc.opts.Scale; y < (row+1)*enc.opts.Scale; y++ {
					img.Set(x, y, clr)
				}
			}
		}
	}
	return img
}
