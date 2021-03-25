package pxl

import (
	"math/bits"
)

type Pxl []uint64

func New(rows int) Pxl {
	return make(Pxl, rows)
}

func (pxl Pxl) Get(x, y int) bool {
	return pxl[y]&(1<<x) != 0
}

func (pxl Pxl) GetRow(y int) uint64 {
	return pxl[y]
}

func (pxl Pxl) Set(x, y int, val bool) {
	if val {
		pxl[y] |= (1 << x)
	} else {
		pxl[y] &= ^(1 << x)
	}
}

func (pxl Pxl) SetRow(y int, row uint64) {
	pxl[y] = row
}

func (pxl Pxl) Rows() int {
	return len(pxl)
}

func (pxl Pxl) Cols() int {
	var max uint64

	for _, row := range pxl {
		if row > max {
			max = row
		}
	}
	return bits.Len64(max)
}
