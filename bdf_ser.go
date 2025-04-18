package bitedutils

import (
	"fmt"
	"io"
)

func (bdf *BDF) BDF2W(w io.Writer) error {
	bdf.calcAvgWidth()
	bdf.calcBbx()
	if _, err := fmt.Fprintln(w, "STARTFONT 2.1"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "FONT", bdf.XLFD.String()); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "SIZE", bdf.XLFD.PtSize()/10, bdf.XLFD.Res.X, bdf.XLFD.Res.Y); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "FONTBOUNDINGBOX", bdf.bbx.W, bdf.bbx.H, bdf.bbx.X, bdf.bbx.Y); err != nil {
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
	for _, glyph := range bdf.Glyphs {
		if err := glyph.Write(bdf, w); err != nil {
			return err
		}
	}
	return nil
}

func (glyph *Glyph) Write(bdf *BDF, w io.Writer) error {
	if _, err := fmt.Fprintln(w, "STARTCHAR", glyph.Name()); err != nil {
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
	if _, err := fmt.Fprintln(w, "BBX", glyph.w, glyph.h, glyph.X, glyph.Y); err != nil {
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
	if len(glyph.bm) <= 0 {
		return nil
	}
	for i, b := range glyph.Bytes() {
		if _, err := fmt.Fprintf(w, "%02X", b); err != nil {
			return err
		}
		if i.IsRowEnd {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
	}
	return nil
}
