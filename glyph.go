package bitedutils

import (
	"encoding/hex"
	"fmt"
	"iter"
	"slices"
)

type Glyph struct {
	name   string
	Code   int
	DWidth int
	w      int
	h      int
	X      int
	Y      int
	bm     []uint32
}

type BmI[T any] struct {
	I        T
	IsRowEnd bool
}

func (glyph *Glyph) Name() string {
	if glyph.Code > 0 {
		return fmt.Sprintf("U+%04X", glyph.Code)
	}
	return glyph.name
}

func (glyph *Glyph) W() int {
	return glyph.w
}

func (glyph *Glyph) H() int {
	return glyph.h
}

func (glyph *Glyph) SWidth(bdf *BDF) int {
	return glyph.DWidth * 72000 / (bdf.XLFD.Res.X * bdf.bbx.W)
}

func (glyph *Glyph) NewBm(w int, h int) {
	glyph.w = w
	glyph.h = h
	glyph.bm = make([]uint32, glyph.h*glyph.Row32())
}

func (glyph *Glyph) Copy() *Glyph {
	glyph1 := *glyph
	glyph1.bm = make([]uint32, len(glyph.bm))
	copy(glyph1.bm, glyph.bm)
	return &glyph1
}

func (glyph *Glyph) Get(i int, j int) (bool, bool) {
	if i < 0 || i >= glyph.h || j < 0 || j >= glyph.w {
		return false, false
	}
	o := i*glyph.Row32() + j/32
	return (glyph.bm[o]>>(31-j%32))&1 != 0, true
}

func (glyph *Glyph) Set(i int, j int, b bool) bool {
	if i < 0 || i >= glyph.h || j < 0 || j >= glyph.w {
		return false
	}
	o := i*glyph.Row32() + j/32
	if b {
		glyph.bm[o] |= 1 << (31 - j%32)
	} else {
		glyph.bm[o] &= ^(1 << (31 - j%32))
	}
	return true
}

func (glyph *Glyph) Flip(i int, j int) {
	if i < 0 || i >= glyph.h || j < 0 || j >= glyph.w {
		return
	}
	o := i*glyph.Row32() + j/32
	glyph.bm[o] ^= 1 << (31 - j%32)
}

func (glyph *Glyph) SetRect(oi int, oj int, w int, h int, b bool) {
	h1 := min(glyph.h, oi+h)
	w1 := min(glyph.w, oj+w)
	for i := oi; i < h1; i++ {
		for j := oj; j < w1; j++ {
			glyph.Set(i, j, b)
		}
	}
}

func (glyph *Glyph) Hex2Row(i int, s string) error {
	if i >= glyph.h {
		return nil
	}

	bs, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	bi := 0
	row32 := glyph.Row32()
	for bc := range slices.Chunk(bs, 4) {
		var n uint32
		for i, b := range bc {
			n |= uint32(b) << ((3 - i) * 8)
		}
		glyph.bm[i*row32+bi] = n
		bi++
	}

	return nil
}

func (glyph *Glyph) Rows() iter.Seq2[int, []uint32] {
	return func(yield func(int, []uint32) bool) {
		if glyph.w <= 0 || glyph.h <= 0 {
			return
		}
		i := 0
		for row := range slices.Chunk(glyph.bm, int(glyph.Row32())) {
			if !yield(i, row) {
				return
			}
			i++
		}
	}
}

func (glyph *Glyph) Uints() iter.Seq2[BmI[[2]int], uint32] {
	return func(yield func(BmI[[2]int], uint32) bool) {
		for i, row := range glyph.Rows() {
			for j, n := range row {
				if !yield(
					BmI[[2]int]{[2]int{i, j}, j == len(row)-1},
					n,
				) {
					return
				}
			}
		}
	}
}

func (glyph *Glyph) Bytes() iter.Seq2[BmI[[2]int], byte] {
	return func(yield func(BmI[[2]int], byte) bool) {
		row8 := glyph.Row8()
		for i, row := range glyph.Rows() {
		row:
			for j, n := range row {
				for o := range 4 {
					bi := j*4 + o
					if bi >= row8 {
						break row
					}
					if !yield(
						BmI[[2]int]{[2]int{i, bi}, bi == row8-1},
						byte(n>>((3-o)*8)),
					) {
						return
					}
				}
			}
		}
	}
}

func (glyph *Glyph) Bits() iter.Seq2[BmI[[2]int], bool] {
	return func(yield func(BmI[[2]int], bool) bool) {
		for i, row := range glyph.Rows() {
		row:
			for j := range row {
				for o := range 32 {
					bi := j*32 + o
					if bi >= glyph.w {
						break row
					}
					b, _ := glyph.Get(i, bi)
					if !yield(BmI[[2]int]{[2]int{i, bi}, bi == glyph.w-1}, b) {
						return
					}
				}
			}
		}
	}
}

func (glyph *Glyph) Row32() int {
	return (glyph.w + 31) / 32
}

func (glyph *Glyph) Row8() int {
	return (glyph.w + 7) / 8
}
