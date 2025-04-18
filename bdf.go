package bitedutils

import (
	"fmt"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type BDF struct {
	*XLFD
	bbx struct {
		W int
		H int
		X int
		Y int
	}
	Props   *orderedmap.OrderedMap[string, any]
	Glyphs  []*Glyph
	Named   map[string]*Glyph
	Unicode map[rune]*Glyph
}

func (bdf *BDF) calcAvgWidth() {
	sum := 0
	for _, glyph := range bdf.Glyphs {
		sum += glyph.DWidth
	}
	bdf.XLFD.avgW = sum * 10 / len(bdf.Glyphs)
}

func (bdf *BDF) calcBbx() {
	for _, glyph := range bdf.Glyphs {
		bdf.bbx.W = max(bdf.bbx.W, glyph.w)
		bdf.bbx.H = max(bdf.bbx.H, glyph.h)
		bdf.bbx.X = min(bdf.bbx.X, glyph.X)
		bdf.bbx.Y = min(bdf.bbx.Y, glyph.Y)
	}
}

func (bdf *BDF) CleanProps() {
	for _, k := range []string{
		"FOUNDRY", "FAMILY_NAME", "WEIGHT_NAME", "SLANT", "SETWIDTH_NAME",
		"ADD_STYLE_NAME", "PIXEL_SIZE", "POINT_SIZE", "RESOLUTION_X",
		"RESOLUTION_Y", "SPACING", "AVERAGE_WIDTH", "CHARSET_REGISTRY",
		"CHARSET_ENCODING"} {
		bdf.Props.Delete(k)
	}
}

func PropString(v any) string {
	if s, ok := v.(string); ok {
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(s, `"`, `""`))
	}
	return fmt.Sprintf("%v", v)
}
