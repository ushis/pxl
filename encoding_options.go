package pxl

import (
	"image/color"
)

type EncodingOptions struct {
	Fg    color.NRGBA
	Bg    color.NRGBA
	Fps   int
	Scale int
}

var defaultEncodingOptions = &EncodingOptions{
	Fg:    color.NRGBA{232, 52, 143, 255},
	Bg:    color.NRGBA{0, 0, 0, 255},
	Fps:   0,
	Scale: 1,
}
