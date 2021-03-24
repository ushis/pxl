package pxl

import (
	"image"
	"image/color"
	"image/gif"
	"io"
)

func EncodeGif(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	g := &gif.GIF{}

	if opts.Fps > 0 {
		g.Image = []*image.Paletted{encodeGif(pxl, opts, false), encodeGif(pxl, opts, true)}
		g.Delay = []int{100.0 / opts.Fps, 100.0 / opts.Fps}
	} else {
		g.Image = []*image.Paletted{encodeGif(pxl, opts, false)}
		g.Delay = []int{0}
	}
	return gif.EncodeAll(w, g)
}

func encodeGif(pxl Pxl, opts *EncodingOptions, invert bool) *image.Paletted {
	rct := image.Rect(0, 0, pxl.Cols()*opts.Scale, pxl.Rows()*opts.Scale)
	plt := color.Palette{opts.Bg, opts.Fg}
	img := image.NewPaletted(rct, plt)

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			fg := pxl.Get(col, row)

			if invert {
				fg = !fg
			}
			clr := opts.Bg

			if fg {
				clr = opts.Fg
			}
			for x := col * opts.Scale; x < (col+1)*opts.Scale; x++ {
				for y := row * opts.Scale; y < (row+1)*opts.Scale; y++ {
					img.Set(x, y, clr)
				}
			}
		}
	}
	return img
}
