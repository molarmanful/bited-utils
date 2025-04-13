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
	bitedbbl "github.com/molarmanful/bited-utils/bited-bbl/lib"
)

func main() {
	name := flag.String("name", "", "output font name")
	nerd := flag.Bool("nerd", false, "whether to double Nerd Font widths")
	ceil := flag.Bool("ceil", false, "whether to right-align glyphs")
	flag.Parse()
	err := bitedbbl.Bbl(os.Stdin, os.Stdout, *name, *nerd, *ceil)
	bitedutils.Check(err)
}
