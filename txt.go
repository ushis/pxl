package pxl

import "io"

func EncodeTxt(w io.Writer, pxl Pxl, opts *EncodingOptions) error {
	buf := make([]rune, pxl.Cols()+1)
	buf[pxl.Cols()] = '\n'

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			if pxl.Get(col, row) {
				buf[col] = opts.Fc
			} else {
				buf[col] = opts.Bc
			}
		}
		if _, err := w.Write([]byte(string(buf))); err != nil {
			return err
		}
	}
	return nil
}
