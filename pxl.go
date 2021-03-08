package main

import (
	_ "embed"
	"encoding/xml"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/bits"
	"net/http"
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

var (
	bg = color.NRGBA{0, 0, 0, 255}
	fg = color.NRGBA{232, 52, 143, 255}
	bh = "#000000"
	fh = "#e8348f"
	bc = ' '
	fc = '█'
	//go:embed index.html
	indexHTML []byte
)

func serve(w http.ResponseWriter, r *http.Request) {
	pth := r.URL.Path

	if pth == "/" {
		if devMode {
			http.ServeFile(w, r, "index.html")
			return
		}
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
		return
	}
	ext := path.Ext(pth)

	if ext != "" {
		pth = pth[:len(pth)-len(ext)]
	}
	strs := strings.Split(pth[1:], "/")
	rows := len(strs)

	if rows > maxRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	vals := make([]uint64, rows)
	cols := 0

	for i, str := range strs {
		val, err := strconv.ParseUint(str, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		n := bits.Len64(val)

		if n > maxCols {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if n > cols {
			cols = n
		}
		vals[i] = val
	}
	if cols == 0 || rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var scale int
	scaleX := maxWidth / cols
	scaleY := maxHeight / rows

	if scaleX > scaleY {
		scale = scaleY
	} else {
		scale = scaleX
	}
	switch ext {
	case "", ".png":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "image/png")
		encodePng(w, vals, cols, scale)
	case ".txt":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		encodeTxt(w, vals, cols)
	case ".svg":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "image/svg+xml")
		encodeSvg(w, vals, cols, scale)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func encodePng(w io.Writer, vals []uint64, cols int, scale int) error {
	img := image.NewNRGBA(image.Rect(0, 0, cols*scale, len(vals)*scale))

	for row, val := range vals {
		for col := 0; col < cols; col++ {
			var clr color.Color

			if val&(1<<col) == 0 {
				clr = bg
			} else {
				clr = fg
			}
			for x := col * scale; x < (col+1)*scale; x++ {
				for y := row * scale; y < (row+1)*scale; y++ {
					img.Set(x, y, clr)
				}
			}
		}
	}
	return png.Encode(w, img)
}

func encodeTxt(w io.Writer, vals []uint64, cols int) error {
	buf := make([]rune, cols+1)
	buf[cols] = '\n'

	for _, val := range vals {
		for col := 0; col < cols; col++ {
			if val&(1<<col) == 0 {
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
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Rects   []svgRect
}

type svgRect struct {
	XMLName xml.Name `xml:"rect"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	Fill    string   `xml:"fill,attr"`
}

func encodeSvg(w io.Writer, vals []uint64, cols, scale int) error {
	svg := svg{
		Width:  strconv.Itoa(cols * scale),
		Height: strconv.Itoa(len(vals) * scale),
		Xmlns:  "http://www.w3.org/2000/svg",
		Rects:  []svgRect{{X: "0", Y: "0", Width: "100%", Height: "100%", Fill: bh}},
	}
	for row, val := range vals {
		for col := 0; col < cols; col++ {
			if val&(1<<col) == 0 {
				continue
			}
			svg.Rects = append(svg.Rects, svgRect{
				X:      strconv.Itoa(col * scale),
				Y:      strconv.Itoa(row * scale),
				Width:  strconv.Itoa(scale),
				Height: strconv.Itoa(scale),
				Fill:   fh,
			})
		}
	}
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(svg)
}
