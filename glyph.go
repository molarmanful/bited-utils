package bitedutils

import (
	"encoding/hex"
	"fmt"
	"slices"

	"github.com/makiuchi-d/gozxing"
)

type Glyph struct {
	Name   string
	Code   int
	DWidth uint64
	Off    [2]int
	Bm     *gozxing.BitMatrix
}

func (glyph *Glyph) Dim() (uint64, uint64) {
	if glyph.Bm == nil {
		return 0, 0
	}
	return uint64(glyph.Bm.GetWidth()), uint64(glyph.Bm.GetHeight())
}

func (glyph *Glyph) GetName() string {
	if glyph.Code > 0 {
		return fmt.Sprintf("U+%04X", glyph.Code)
	}
	return glyph.Name
}

func (glyph *Glyph) SWidth(bdf *BDF) uint64 {
	return glyph.DWidth * 72000 / (bdf.XLFD.Res[0] * bdf.BBX.W)
}

func (glyph *Glyph) SetRow(i int, s string) error {
	if glyph.Bm == nil {
		return nil
	}
	w := glyph.Bm.GetWidth()
	h := glyph.Bm.GetHeight()
	if i >= h {
		return nil
	}

	bs, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	row := gozxing.NewBitArray(w)
	bi := 0
	for bc := range slices.Chunk(bs, 4) {
		var n uint32
		for i, b := range bc {
			n |= uint32(b) << ((3 - i) * 8)
		}
		row.SetBulk(bi*32, n)
	}

	glyph.Bm.SetRow(i, row)
	return nil
}
