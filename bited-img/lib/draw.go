package bitedimg

import (
	"bufio"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
	"github.com/fogleman/gg"
)

// drawTCs draws a TXT/CLR pair to PNG.
func (unit *Unit) drawTCs(stem string) error {
	txt, err := script.File(filepath.Join(unit.TxtDir, stem+".txt")).String()
	if err != nil {
		return err
	}
	txt = strings.TrimSuffix(txt, "\n")
	lines := strings.Split(txt, "\n")

	dummy := gg.NewContext(1, 1)
	dummy.SetFontFace(unit.BDF)
	fsz := int(dummy.FontHeight())
	w, h := dummy.MeasureMultilineString(txt, 1)
	pad := fsz
	W := int(w) + pad*2
	H := int(h) + pad*2

	ctx := gg.NewContext(int(W), int(H))
	ctx.SetHexColor(unit.Clrs.Bg)
	ctx.Clear()
	ctx.SetFontFace(unit.BDF)

	X := pad
	Y := pad + unit.Ascent
	clrScan := bufio.NewScanner(script.File(filepath.Join(unit.TxtDir, stem+".clr")))
	ctx.SetHexColor(unit.Clrs.Fg)
	for _, line := range lines {
		var clrLine []rune
		if clrScan.Scan() {
			clrLine = []rune(clrScan.Text())
		}
		for i, c := range []rune(line) {
			cS := string(c)
			if len(clrLine) > i {
				if clr, ok := unit.ClrsMap[clrLine[i]]; ok {
					ctx.SetHexColor(clr)
				}
			}
			ctx.DrawString(cS, float64(X), float64(Y))
			w, _ := ctx.MeasureString(cS)
			X += int(w)
		}
		Y += fsz
		X = pad
	}

	ctx.SavePNG(filepath.Join(unit.OutDir, stem+".png"))
	return nil
}
