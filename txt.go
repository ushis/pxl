package pxl

import "io"

const (
	bc = ' '
	fc = '█'
	lf = '\n'
)

func EncodeTxt(w io.Writer, pxl Pxl) error {
	buf := make([]rune, pxl.Cols()+1)
	buf[pxl.Cols()] = lf

	for row := 0; row < pxl.Rows(); row++ {
		for col := 0; col < pxl.Cols(); col++ {
			if pxl.Get(col, row) {
				buf[col] = fc
			} else {
				buf[col] = bc
			}
		}
		if _, err := w.Write([]byte(string(buf))); err != nil {
			return err
		}
	}
	return nil
}
