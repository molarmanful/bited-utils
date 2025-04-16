package bitedutils

import (
	"fmt"
	"io"
	"maps"
	"slices"
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
		fmt.Fprintln(w, k, PropString(v))
	}
	for k, v := range bdf.Props.FromOldest() {
		fmt.Fprintln(w, k, PropString(v))
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
