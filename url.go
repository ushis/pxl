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
	pxl, err := decodeURLPath(url.Path)

	if err != nil {
		return nil, nil, err
	}
	opts, err := decodeURLQuery(url.Query())

	if err != nil {
		return nil, nil, err
	}
	return pxl, opts, err
}

func decodeURLPath(pth string) (Pxl, error) {
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

func decodeURLQuery(params url.Values) (*EncodingOptions, error) {
	opts := &EncodingOptions{
		Fg:    defaultEncodingOptions.Fg,
		Bg:    defaultEncodingOptions.Bg,
		Fc:    defaultEncodingOptions.Fc,
		Bc:    defaultEncodingOptions.Bc,
		Fps:   defaultEncodingOptions.Fps,
		Scale: defaultEncodingOptions.Scale,
	}
	if str := params.Get("fg"); str != "" {
		clr, err := decodeURLColor(str)

		if err != nil {
			return opts, err
		}
		opts.Fg = clr
	}
	if str := params.Get("bg"); str != "" {
		clr, err := decodeURLColor(str)

		if err != nil {
			return opts, err
		}
		opts.Bg = clr
	}
	if str := params.Get("fc"); str != "" {
		chars := []rune(str)

		if len(chars) != 1 {
			return opts, errors.New("fc too long")
		}
		opts.Fc = chars[0]
	}
	if str := params.Get("bc"); str != "" {
		chars := []rune(str)

		if len(chars) != 1 {
			return opts, errors.New("bc too long")
		}
		opts.Bc = chars[0]
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

func decodeURLColor(s string) (color.NRGBA, error) {
	val, err := strconv.ParseUint(s, 10, 32)

	clr := color.NRGBA{
		R: uint8((val >> 24) & 0xff),
		G: uint8((val >> 16) & 0xff),
		B: uint8((val >> 8) & 0xff),
		A: uint8((val >> 0) & 0xff),
	}
	return clr, err
}

func EncodeURL(pxl Pxl, opts *EncodingOptions) *url.URL {
	return &url.URL{
		Path:     encodeURLPath(pxl),
		RawQuery: encodeURLQuery(opts).Encode(),
	}
}

func encodeURLPath(pxl Pxl) string {
	b := &strings.Builder{}

	for _, row := range pxl {
		b.WriteByte('/')
		b.WriteString(strconv.FormatUint(row, 10))
	}
	return b.String()
}

func encodeURLQuery(opts *EncodingOptions) url.Values {
	return url.Values{
		"bg":    {encodeURLColor(opts.Bg)},
		"fg":    {encodeURLColor(opts.Fg)},
		"bc":    {string(opts.Bc)},
		"fc":    {string(opts.Fc)},
		"fps":   {strconv.Itoa(opts.Fps)},
		"scale": {strconv.Itoa(opts.Scale)},
	}
}

func encodeURLColor(clr color.NRGBA) string {
	var num uint64
	num |= uint64(clr.R) << 24
	num |= uint64(clr.G) << 16
	num |= uint64(clr.B) << 8
	num |= uint64(clr.A) << 0
	return strconv.FormatUint(num, 10)
}
