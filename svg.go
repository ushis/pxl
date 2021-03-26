package pxl

import (
	"encoding/xml"
	"fmt"
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
	Childs  []*svgPath
}

type svgPath struct {
	XMLName     xml.Name `xml:"path"`
	D           string   `xml:"d,attr"`
	Fill        string   `xml:"fill,attr"`
	FillOpacity string   `xml:"fill-opacity,attr"`
	Animation   *svgAnimate
}

type svgAnimate struct {
	XMLName       xml.Name `xml:"animate"`
	AttributeName string   `xml:"attributeName,attr"`
	CalcMode      string   `xml:"calcMode,attr"`
	Dur           string   `xml:"dur,attr"`
	RepeatCount   string   `xml:"repeatCount,attr"`
	Values        string   `xml:"values,attr"`
}

func EncodeSvg(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	bg := &svgPath{}
	fg := &svgPath{}
	bg.D, fg.D = encodeSvgPathData(pxl, opts.Scale)
	bg.Fill, bg.FillOpacity = encodeSvgColor(opts.Bg)
	fg.Fill, fg.FillOpacity = encodeSvgColor(opts.Fg)

	if opts.Fps > 0 {
		bg.Animation = encodeSvgAnimation("fill", []string{bg.Fill, fg.Fill}, opts.Fps)
		fg.Animation = encodeSvgAnimation("fill", []string{fg.Fill, bg.Fill}, opts.Fps)
	}
	svg := &svg{
		Xmlns:  "http://www.w3.org/2000/svg",
		Width:  strconv.Itoa(pxl.Cols() * opts.Scale),
		Height: strconv.Itoa(pxl.Rows() * opts.Scale),
		Childs: []*svgPath{bg, fg},
	}
	return xml.NewEncoder(w).Encode(svg)
}

func encodeSvgPathData(pxl Pxl, scale int) (string, string) {
	bg := &strings.Builder{}
	fg := &strings.Builder{}

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			d := fmt.Sprintf("M%[1]d,%[2]dh%[3]dv%[3]dh-%[3]dz", col*scale, row*scale, scale)

			if pxl.Get(col, row) {
				fg.WriteString(d)
			} else {
				bg.WriteString(d)
			}
		}
	}
	return bg.String(), fg.String()
}

func encodeSvgAnimation(attrName string, values []string, fps int) *svgAnimate {
	return &svgAnimate{
		AttributeName: attrName,
		CalcMode:      "discrete",
		Dur:           strconv.FormatFloat(float64(len(values))/float64(fps), 'f', 4, 64),
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
