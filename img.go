package pxl

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func EncodeJpeg(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	return jpeg.Encode(w, encodeImage(pxl, opts.Scale, opts.Bg, opts.Fg), nil)
}

func EncodePng(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	return png.Encode(w, encodeImage(pxl, opts.Scale, opts.Bg, opts.Fg))
}

func EncodeGif(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	g := &gif.GIF{}

	if opts.Fps > 0 {
		g.Image = []*image.Paletted{
			encodeImage(pxl, opts.Scale, opts.Bg, opts.Fg),
			encodeImage(pxl, opts.Scale, opts.Fg, opts.Bg),
		}
		g.Delay = []int{
			100 / opts.Fps,
			100 / opts.Fps,
		}
	} else {
		g.Image = []*image.Paletted{encodeImage(pxl, opts.Scale, opts.Bg, opts.Fg)}
		g.Delay = []int{0}
	}
	return gif.EncodeAll(w, g)
}

func encodeImage(pxl Pxl, scale int, bg, fg color.Color) *image.Paletted {
	r := image.Rect(0, 0, pxl.Cols()*scale, pxl.Rows()*scale)
	p := color.Palette{bg, fg}
	img := image.NewPaletted(r, p)

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			var c uint8

			if pxl.Get(col, row) {
				c = 1
			} else {
				c = 0
			}
			for x := col * scale; x < (col+1)*scale; x++ {
				for y := row * scale; y < (row+1)*scale; y++ {
					img.SetColorIndex(x, y, c)
				}
			}
		}
	}
	return img
}
