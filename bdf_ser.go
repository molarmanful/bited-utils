package bitedutils

import (
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/makiuchi-d/gozxing"
)

func (bdf *BDF) BDF2W(w io.Writer) error {
	bdf.CalcAvgWidth()
	bdf.CalcBBX()
	if _, err := fmt.Fprintln(w, "STARTFONT 2.1"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "FONT", bdf.XLFD.String()); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "SIZE", bdf.XLFD.PtSize()/10, bdf.XLFD.Res[0], bdf.XLFD.Res[0]); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "FONTBOUNDINGBOX", bdf.BBX.W, bdf.BBX.H, bdf.BBX.X, bdf.BBX.Y); err != nil {
		return err
	}
	if err := bdf.WriteProps(w); err != nil {
		return err
	}
	if err := bdf.WriteChars(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "ENDFONT"); err != nil {
		return err
	}
	return nil
}

func (bdf *BDF) WriteProps(w io.Writer) error {
	bdf.CleanProps()
	if _, err := fmt.Fprintln(w, "STARTPROPERTIES", bdf.Props.Len()+14); err != nil {
		return err
	}
	for k, v := range bdf.XLFD.Props().FromOldest() {
		if _, err := fmt.Fprintln(w, k, PropString(v)); err != nil {
			return err
		}
	}
	for k, v := range bdf.Props.FromOldest() {
		if _, err := fmt.Fprintln(w, k, PropString(v)); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "ENDPROPERTIES"); err != nil {
		return err
	}
	return nil
}

func (bdf *BDF) WriteChars(w io.Writer) error {
	if _, err := fmt.Fprintln(w, "CHARS", len(bdf.Glyphs)); err != nil {
		return err
	}
	for _, k := range slices.Sorted(maps.Keys(bdf.Named)) {
		if err := bdf.Named[k].Write(bdf, w); err != nil {
			return err
		}
	}
	for _, k := range slices.Sorted(maps.Keys(bdf.Unicode)) {
		if err := bdf.Unicode[k].Write(bdf, w); err != nil {
			return err
		}
	}
	return nil
}

func (glyph *Glyph) Write(bdf *BDF, w io.Writer) error {
	if _, err := fmt.Fprintln(w, "STARTCHAR", glyph.GetName()); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "ENCODING", glyph.Code); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "SWIDTH", glyph.SWidth(bdf), 0); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "DWIDTH", glyph.DWidth, 0); err != nil {
		return err
	}
	dw, dh := glyph.Dim()
	if _, err := fmt.Fprintln(w, "BBX", dw, dh, glyph.Off[0], glyph.Off[1]); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "BITMAP"); err != nil {
		return err
	}
	if err := glyph.WriteBm(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "ENDCHAR"); err != nil {
		return err
	}
	return nil
}

func (glyph *Glyph) WriteBm(w io.Writer) error {
	if glyph.Bm == nil {
		return nil
	}
	dw := glyph.Bm.GetWidth()
	dh := glyph.Bm.GetHeight()
	row := gozxing.NewBitArray(dw)
	dwbs := row.GetSizeInBytes()

	for i := range int(dh) {
		glyph.Bm.GetRow(i, row)
	row:
		for i, n := range row.GetBitArray() {
			for o := range 4 {
				if i*4+o >= dwbs {
					break row
				}
				if _, err := fmt.Fprintf(w, "%02X", n>>((3-o)*8)); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	return nil
}
