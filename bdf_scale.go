package bitedutils

import "fmt"

func (bdf *BDF) Scale(scale int) error {
	if scale < 1 {
		return fmt.Errorf("scale < 1")
	}

	bdf.XLFD.PxSize *= scale

	for _, k := range []string{"FONT_ASCENT", "FONT_DESCENT", "CAP_HEIGHT",
		"X_HEIGHT", "BITED_DWIDTH", "BITED_EDITOR_GRID_SIZE"} {
		if v, ok := bdf.Props.Get(k); ok {
			if n, ok := v.(int); ok {
				bdf.Props.Set(k, n*scale)
			}
		}
	}

	for _, glyph := range bdf.Glyphs {
		glyph.Scale(scale)
	}

	return nil
}

func (glyph *Glyph) Scale(scale int) {
	glyph0 := glyph.Copy()

	glyph.DWidth *= scale
	glyph.X *= scale
	glyph.Y *= scale
	glyph.NewBm(glyph.w*scale, glyph.h*scale)

	for i, b := range glyph0.Bits() {
		if !b {
			continue
		}
		glyph.SetRect(i.I[0]*scale, i.I[1]*scale, scale, scale, true)
	}
}
