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
	Childs  []*svgGroup
}

type svgGroup struct {
	XMLName     xml.Name `xml:"g"`
	Fill        string   `xml:"fill,attr"`
	FillOpacity string   `xml:"fill-opacity,attr"`
	Animation   *svgAnimate
	Childs      []*svgRect
}

type svgAnimate struct {
	XMLName       xml.Name `xml:"animate"`
	AttributeName string   `xml:"attributeName,attr"`
	CalcMode      string   `xml:"calcMode,attr"`
	Dur           string   `xml:"dur,attr"`
	RepeatCount   string   `xml:"repeatCount,attr"`
	Values        string   `xml:"values,attr"`
}

type svgRect struct {
	XMLName xml.Name `xml:"rect"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

func (enc SvgEncoder) Encode(pxl Pxl) error {
	bg, bgOpacity := encodeSvgColor(enc.opts.Bg)
	fg, fgOpacity := encodeSvgColor(enc.opts.Fg)

	svg := svg{
		Xmlns:  "http://www.w3.org/2000/svg",
		Width:  strconv.Itoa(pxl.Cols() * enc.opts.Scale),
		Height: strconv.Itoa(pxl.Rows() * enc.opts.Scale),
		Childs: []*svgGroup{
			{Fill: bg, FillOpacity: bgOpacity},
			{Fill: fg, FillOpacity: fgOpacity},
		},
	}
	if enc.opts.Fps > 0 {
		svg.Childs[0].Animation = encodeSvgAnimation(bg, fg, enc.opts.Fps)
		svg.Childs[1].Animation = encodeSvgAnimation(fg, bg, enc.opts.Fps)
	}
	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			rect := &svgRect{
				X:      strconv.Itoa(col * enc.opts.Scale),
				Y:      strconv.Itoa(row * enc.opts.Scale),
				Width:  strconv.Itoa(enc.opts.Scale),
				Height: strconv.Itoa(enc.opts.Scale),
			}
			if pxl.Get(col, row) {
				svg.Childs[1].Childs = append(svg.Childs[1].Childs, rect)
			} else {
				svg.Childs[0].Childs = append(svg.Childs[0].Childs, rect)
			}
		}
	}
	if _, err := enc.w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	return xml.NewEncoder(enc.w).Encode(svg)
}

func encodeSvgAnimation(fillA, fillB string, fps int) *svgAnimate {
	return &svgAnimate{
		AttributeName: "fill",
		CalcMode:      "discrete",
		Dur:           strconv.FormatFloat(1.0/float64(fps), 'f', 4, 64),
		RepeatCount:   "indefinite",
		Values:        fillA + ";" + fillB,
	}
}

const hextable = "0123456789abcdef"

func encodeSvgColor(clr color.NRGBA) (string, string) {
	rgb := []byte{
		'#',
		hextable[clr.R>>4],
		hextable[clr.R&0x0f],
		hextable[clr.G>>4],
		hextable[clr.G&0x0f],
		hextable[clr.B>>4],
		hextable[clr.B&0x0f],
	}
	return string(rgb), strconv.FormatFloat(float64(clr.A)/0xff, 'f', 4, 64)
}
