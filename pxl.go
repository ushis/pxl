package main

import (
	_ "embed"
	"encoding/xml"
	"errors"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/bits"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var devMode bool
var listenAddr string

func init() {
	flag.BoolVar(&devMode, "dev", false, "development mode")
	flag.StringVar(&listenAddr, "listen", ":9876", "listen address")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serve)
	http.ListenAndServe(listenAddr, nil)
}

const (
	maxCols   = 64
	maxRows   = 64
	maxWidth  = 1024
	maxHeight = 1024
)

func serve(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveIndex(w, r)
		return
	}
	ext := path.Ext(r.URL.Path)
	pth := r.URL.Path[:len(r.URL.Path)-len(ext)]
	pxl, err := decodePxl(pth)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if pxl.Cols() == 0 || pxl.Cols() > maxCols || pxl.Rows() == 0 || pxl.Rows() > maxRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	opts, err := decodeOptions(r.URL.Query())

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch ext {
	case "", ".png":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "image/png")
		encodePng(w, pxl, opts)
	case ".txt":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		encodeTxt(w, pxl)
	case ".svg":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "image/svg+xml")
		encodeSvg(w, pxl, opts)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

//go:embed index.html
var indexHTML []byte

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if devMode {
		http.ServeFile(w, r, "index.html")
	} else {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	}
}

type pxl struct {
	rows []uint64
	cols int
}

func decodePxl(pth string) (*pxl, error) {
	if pth == "" || pth[0] != '/' {
		return nil, errors.New("invalid pxl path")
	}
	prts := strings.Split(pth[1:], "/")
	pxl := &pxl{rows: make([]uint64, len(prts))}

	for i, prt := range prts {
		row, err := strconv.ParseUint(prt, 10, 64)

		if err != nil {
			return nil, err
		}
		pxl.SetRow(i, row)
	}
	return pxl, nil
}

func (p *pxl) SetRow(y int, row uint64) {
	n := bits.Len64(row)

	if n > p.cols {
		p.cols = n
	}
	p.rows[y] = row
}

func (p *pxl) Get(x, y int) bool {
	return (p.rows[y] & (1 << x)) != 0
}

func (p *pxl) Rows() int {
	return len(p.rows)
}

func (p *pxl) Cols() int {
	return p.cols
}

func (p *pxl) MaxScale(w, h int) int {
	scaleX := w / p.Cols()
	scaleY := h / p.Rows()

	if scaleX < scaleY {
		return scaleX
	}
	return scaleY
}

var (
	defaultBg = color.NRGBA{0, 0, 0, 255}
	defaultFg = color.NRGBA{232, 52, 143, 255}
)

type options struct {
	Fg color.NRGBA
	Bg color.NRGBA
}

func decodeOptions(params url.Values) (options, error) {
	opts := options{
		Fg: defaultFg,
		Bg: defaultBg,
	}
	if fg := params.Get("fg"); fg != "" {
		clr, err := decodeColor(fg)

		if err != nil {
			return opts, err
		}
		opts.Fg = clr
	}
	if bg := params.Get("bg"); bg != "" {
		clr, err := decodeColor(bg)

		if err != nil {
			return opts, err
		}
		opts.Bg = clr
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

func encodePng(w io.Writer, pxl *pxl, opts options) error {
	scl := pxl.MaxScale(maxWidth, maxHeight)
	img := image.NewNRGBA(image.Rect(0, 0, pxl.Cols()*scl, pxl.Rows()*scl))

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			clr := opts.Bg

			if pxl.Get(col, row) {
				clr = opts.Fg
			}
			for x := col * scl; x < (col+1)*scl; x++ {
				for y := row * scl; y < (row+1)*scl; y++ {
					img.Set(x, y, clr)
				}
			}
		}
	}
	return png.Encode(w, img)
}

const (
	bc = ' '
	fc = '█'
)

func encodeTxt(w io.Writer, pxl *pxl) error {
	buf := make([]rune, pxl.Cols()+1)
	buf[pxl.Cols()] = '\n'

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			if pxl.Get(col, row) {
				buf[col] = bc
			} else {
				buf[col] = fc
			}
		}
		if _, err := w.Write([]byte(string(buf))); err != nil {
			return err
		}
	}
	return nil
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

func encodeSvg(w io.Writer, pxl *pxl, opts options) error {
	scl := pxl.MaxScale(maxWidth, maxHeight)
	bgr := encodeSvgColor(opts.Bg)
	fgr := encodeSvgColor(opts.Fg)
	svg := svg{
		Xmlns:  "http://www.w3.org/2000/svg",
		Width:  strconv.Itoa(pxl.Cols() * scl),
		Height: strconv.Itoa(pxl.Rows() * scl),
		Bg:     svgRect{X: "0", Y: "0", Width: "100%", Height: "100%", Fill: bgr},
		Fg:     svgGroup{Fill: fgr},
	}
	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			if !pxl.Get(col, row) {
				continue
			}
			svg.Fg.Childs = append(svg.Fg.Childs, svgRect{
				X:      strconv.Itoa(col * scl),
				Y:      strconv.Itoa(row * scl),
				Width:  strconv.Itoa(scl),
				Height: strconv.Itoa(scl),
			})
		}
	}
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(svg)
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
