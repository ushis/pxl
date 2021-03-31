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
	return jpeg.Encode(w, encodeNRGBA(pxl, opts.Scale, opts.Bg, opts.Fg), nil)
}

func EncodePng(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	return png.Encode(w, encodeNRGBA(pxl, opts.Scale, opts.Bg, opts.Fg))
}

func encodeNRGBA(pxl Pxl, scale int, bg, fg color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, pxl.Cols()*scale, pxl.Rows()*scale))
	encodeImg(img, pxl, scale, bg, fg)
	return img
}

func EncodeGif(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	g := &gif.GIF{}

	if opts.Fps > 0 {
		g.Image = []*image.Paletted{
			encodePaletted(pxl, opts.Scale, opts.Bg, opts.Fg),
			encodePaletted(pxl, opts.Scale, opts.Fg, opts.Bg),
		}
		g.Delay = []int{
			100 / opts.Fps,
			100 / opts.Fps,
		}
	} else {
		g.Image = []*image.Paletted{encodePaletted(pxl, opts.Scale, opts.Bg, opts.Fg)}
		g.Delay = []int{0}
	}
	return gif.EncodeAll(w, g)
}

func encodePaletted(pxl Pxl, scale int, bg, fg color.Color) *image.Paletted {
	plt := color.Palette{bg, fg}
	img := image.NewPaletted(image.Rect(0, 0, pxl.Cols()*scale, pxl.Rows()*scale), plt)
	encodeImg(img, pxl, scale, bg, fg)
	return img
}

type img interface {
	image.Image
	Set(x, y int, c color.Color)
}

func encodeImg(img img, pxl Pxl, scale int, bg, fg color.Color) {
	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			var c color.Color

			if pxl.Get(col, row) {
				c = fg
			} else {
				c = bg
			}
			for x := col * scale; x < (col+1)*scale; x++ {
				for y := row * scale; y < (row+1)*scale; y++ {
					img.Set(x, y, c)
				}
			}
		}
	}
}
