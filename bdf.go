package bitedutils

import (
	"fmt"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type BDF struct {
	*XLFD
	BBX struct {
		W uint64
		H uint64
		X int
		Y int
	}
	Props   *orderedmap.OrderedMap[string, interface{}]
	Glyphs  []*Glyph
	Named   map[string]*Glyph
	Unicode map[rune]*Glyph
}

func (bdf *BDF) CalcAvgWidth() {
	sum := uint64(0)
	for _, glyph := range bdf.Glyphs {
		sum += glyph.DWidth
	}
	bdf.XLFD.AvgW = sum * 10 / uint64(len(bdf.Glyphs))
}

func (bdf *BDF) CalcBBX() {
	for _, glyph := range bdf.Glyphs {
		dw, dh := glyph.Dim()
		bdf.BBX.W = max(bdf.BBX.W, dw)
		bdf.BBX.H = max(bdf.BBX.H, dh)
		bdf.BBX.X = min(bdf.BBX.X, glyph.Off[0])
		bdf.BBX.Y = min(bdf.BBX.Y, glyph.Off[1])
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

func PropString(v interface{}) string {
	if s, ok := v.(string); ok {
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(s, `"`, `""`))
	}
	return fmt.Sprintf("%v", v)
}
