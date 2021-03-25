package pxl

import (
	"errors"
	"image/color"
	"net/url"
	"path"
	"strconv"
	"strings"
)

func DecodeURL(url *url.URL) (Pxl, *EncodingOptions, error) {
	pxl, err := decodePath(url.Path)

	if err != nil {
		return nil, nil, err
	}
	opts, err := decodeQuery(url.Query())

	if err != nil {
		return nil, nil, err
	}
	return pxl, opts, err
}

func decodePath(pth string) (Pxl, error) {
	pth = pth[:len(pth)-len(path.Ext(pth))]

	if pth == "" || pth[0] != '/' {
		return nil, errors.New("invalid path")
	}
	prts := strings.Split(pth[1:], "/")
	pxl := New(len(prts))

	for y, prt := range prts {
		row, err := strconv.ParseUint(prt, 10, 64)

		if err != nil {
			return nil, err
		}
		pxl.SetRow(y, row)
	}
	return pxl, nil
}

func decodeQuery(params url.Values) (*EncodingOptions, error) {
	opts := &EncodingOptions{
		Fg:    defaultEncodingOptions.Fg,
		Bg:    defaultEncodingOptions.Bg,
		Fps:   defaultEncodingOptions.Fps,
		Scale: defaultEncodingOptions.Scale,
	}
	if str := params.Get("fg"); str != "" {
		clr, err := decodeColor(str)

		if err != nil {
			return opts, err
		}
		opts.Fg = clr
	}
	if str := params.Get("bg"); str != "" {
		clr, err := decodeColor(str)

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

func decodeColor(s string) (color.NRGBA, error) {
	val, err := strconv.ParseUint(s, 10, 32)

	clr := color.NRGBA{
		R: uint8((val >> 24) & 0xff),
		G: uint8((val >> 16) & 0xff),
		B: uint8((val >> 8) & 0xff),
		A: uint8((val >> 0) & 0xff),
	}
	return clr, err
}
