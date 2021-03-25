package main

import (
	_ "embed"
	"flag"
	"net/http"
	"path"

	"github.com/ushis/pxl"
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

//go:embed index.html
var indexHTML []byte

func serve(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if devMode {
			http.ServeFile(w, r, "./pxld/index.html")
		} else {
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.Write(indexHTML)
		}
		return
	}
	p, opts, err := pxl.DecodeURL(r.URL)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if p.Cols() == 0 || p.Cols() > maxCols || p.Rows() == 0 || p.Rows() > maxRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	scaleX := maxWidth / p.Cols()
	scaleY := maxHeight / p.Rows()

	if scaleX < scaleY {
		opts.Scale = scaleX
	} else {
		opts.Scale = scaleY
	}
	switch path.Ext(r.URL.Path) {
	case ".gif":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "image/gif")
		pxl.EncodeGif(w, p, opts)
	case "", ".png":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "image/png")
		pxl.EncodePng(w, p, opts)
	case ".txt":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		pxl.EncodeTxt(w, p, opts)
	case ".svg":
		w.Header().Add("Cache-Control", "public, max-age=86400, immutable")
		w.Header().Add("Content-Type", "image/svg+xml")
		pxl.EncodeSvg(w, p, opts)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
