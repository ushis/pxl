package pxl

import (
	"encoding/xml"
	"image/color"
	"io"
	"strconv"
	"strings"
)

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

func EncodeSvg(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	bg := &svgGroup{}
	fg := &svgGroup{}
	bg.Fill, bg.FillOpacity = encodeSvgColor(opts.Bg)
	fg.Fill, fg.FillOpacity = encodeSvgColor(opts.Fg)

	if opts.Fps > 0 {
		bg.Animation = encodeSvgAnimation("fill", []string{bg.Fill, fg.Fill}, opts.Fps)
		fg.Animation = encodeSvgAnimation("fill", []string{fg.Fill, bg.Fill}, opts.Fps)
	}
	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			rect := &svgRect{
				X:      strconv.Itoa(col * opts.Scale),
				Y:      strconv.Itoa(row * opts.Scale),
				Width:  strconv.Itoa(opts.Scale),
				Height: strconv.Itoa(opts.Scale),
			}
			if pxl.Get(col, row) {
				fg.Childs = append(fg.Childs, rect)
			} else {
				bg.Childs = append(bg.Childs, rect)
			}
		}
	}
	svg := &svg{
		Xmlns:  "http://www.w3.org/2000/svg",
		Width:  strconv.Itoa(pxl.Cols() * opts.Scale),
		Height: strconv.Itoa(pxl.Rows() * opts.Scale),
		Childs: []*svgGroup{bg, fg},
	}
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(svg)
}

func encodeSvgAnimation(attrName string, values []string, fps int) *svgAnimate {
	return &svgAnimate{
		AttributeName: attrName,
		CalcMode:      "discrete",
		Dur:           strconv.FormatFloat(1.0/float64(fps), 'f', 4, 64),
		RepeatCount:   "indefinite",
		Values:        strings.Join(values, ";"),
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
