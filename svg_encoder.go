package pxl

import (
	"encoding/xml"
	"image/color"
	"io"
	"strconv"
)

type SvgEncoder struct {
	w    io.Writer
	opts EncoderOptions
}

func NewSvgEncoder(w io.Writer, opts EncoderOptions) SvgEncoder {
	return SvgEncoder{w, opts}
}

type svg struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	Bg      svgRect
	Fg      svgGroup
}

type svgRect struct {
	XMLName xml.Name `xml:"rect"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	Fill    string   `xml:"fill,attr,omitempty"`
}

type svgGroup struct {
	XMLName xml.Name `xml:"g"`
	Fill    string   `xml:"fill,attr,omitempty"`
	Childs  []svgRect
}

func (enc SvgEncoder) Encode(pxl Pxl) error {
	bgr := encodeSvgColor(enc.opts.Bg)
	fgr := encodeSvgColor(enc.opts.Fg)
	svg := svg{
		Xmlns:  "http://www.w3.org/2000/svg",
		Width:  strconv.Itoa(pxl.Cols() * enc.opts.Scale),
		Height: strconv.Itoa(pxl.Rows() * enc.opts.Scale),
		Bg:     svgRect{X: "0", Y: "0", Width: "100%", Height: "100%", Fill: bgr},
		Fg:     svgGroup{Fill: fgr},
	}
	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			if !pxl.Get(col, row) {
				continue
			}
			svg.Fg.Childs = append(svg.Fg.Childs, svgRect{
				X:      strconv.Itoa(col * enc.opts.Scale),
				Y:      strconv.Itoa(row * enc.opts.Scale),
				Width:  strconv.Itoa(enc.opts.Scale),
				Height: strconv.Itoa(enc.opts.Scale),
			})
		}
	}
	if _, err := enc.w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	return xml.NewEncoder(enc.w).Encode(svg)
}

const hextable = "0123456789abcdef"

func encodeSvgColor(clr color.NRGBA) string {
	rgba := []byte{
		'#',
		hextable[clr.R>>4],
		hextable[clr.R&0x0f],
		hextable[clr.G>>4],
		hextable[clr.G&0x0f],
		hextable[clr.B>>4],
		hextable[clr.B&0x0f],
		hextable[clr.A>>4],
		hextable[clr.A&0x0f],
	}
	return string(rgba)
}
