package pxl

import (
	"errors"
	"image/color"
	"net/url"
	"strconv"
)

type EncoderOptions struct {
	Fg    color.NRGBA
	Bg    color.NRGBA
	Fps   int
	Scale int
}

var defaultEncoderOptions = EncoderOptions{
	Fg:    color.NRGBA{232, 52, 143, 255},
	Bg:    color.NRGBA{0, 0, 0, 255},
	Fps:   0,
	Scale: 1,
}

func DecodeEncoderOptions(params url.Values) (EncoderOptions, error) {
	opts := EncoderOptions{
		Fg:    defaultEncoderOptions.Fg,
		Bg:    defaultEncoderOptions.Bg,
		Scale: defaultEncoderOptions.Scale,
	}
	if str := params.Get("fg"); str != "" {
		clr, err := DecodeColor(str)

		if err != nil {
			return opts, err
		}
		opts.Fg = clr
	}
	if str := params.Get("bg"); str != "" {
		clr, err := DecodeColor(str)

		if err != nil {
			return opts, err
		}
		opts.Bg = clr
	}
	if str := params.Get("fps"); str != "" {
		fps, err := strconv.ParseUint(str, 10, 8)

		if err != nil {
			return opts, err
		}
		if fps > 100 {
			return opts, errors.New("fps out of range")
		}
		opts.Fps = int(fps)
	}
	if str := params.Get("scale"); str != "" {
		scale, err := strconv.ParseUint(str, 10, 32)

		if err != nil {
			return opts, err
		}
		opts.Scale = int(scale)
	}
	return opts, nil
}

func DecodeColor(s string) (color.NRGBA, error) {
	val, err := strconv.ParseUint(s, 10, 32)

	clr := color.NRGBA{
		R: uint8((val >> 24) & 0xff),
		G: uint8((val >> 16) & 0xff),
		B: uint8((val >> 8) & 0xff),
		A: uint8((val >> 0) & 0xff),
	}
	return clr, err
}
