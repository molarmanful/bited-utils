// bited-bbl is a command for proportionalizing bited BDFs.
//
// Usage:
//
//	bited-bbl --name <string> [--nerd]
//
// Flags:
//
//	--name
//		Family name of the output font.
//	--nerd
//		Whether to double Nerd Font widths.
//	--ceil
//		Whether to right-align glyphs during recenter.
//
// bited-bbl accepts input via STDIN and outputs to STDOUT.
package main

import (
	"flag"
	"os"

	bitedutils "github.com/molarmanful/bited-utils"
)

func main() {
	name := flag.String("name", "", "output font name")
	nerd := flag.Bool("nerd", false, "whether to double Nerd Font widths")
	ceil := flag.Bool("ceil", false, "whether to right-align glyphs")
	flag.Parse()

	ceiln := 0
	if *ceil {
		ceiln = 1
	}

	bdf, err := bitedutils.R2BDF(os.Stdin)
	bitedutils.Check(err)
	if *name != "" {
		bdf.XLFD.Family = *name
	}

	for _, glyph := range bdf.Glyphs {
		n := 1
		if glyph.Code >= 0 {
			n = bitedutils.WcWidth(rune(glyph.Code), *nerd)
		}
		glyph.DWidth = max(glyph.DWidth*n, glyph.W())
		glyph.X = max(0, glyph.X+(glyph.DWidth+ceiln)*(n-1)/n/2)
	}

	err = bdf.BDF2W(os.Stdout)
	bitedutils.Check(err)
}
