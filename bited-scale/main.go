// bited-scale is a command for scaling bited BDFs.
//
// Usage:
//
//	bited-scale [-x <int>] [--name <string>]
//
// Flags:
//
//	--x
//		Scaling factor. Defaults to 2.
//	--name
//		Family name of the output font.
//
// bited-scale accepts input via STDIN and outputs to STDOUT.
package main

import (
	"flag"
	"os"

	bitedutils "github.com/molarmanful/bited-utils"
)

func main() {
	scale := flag.Int("x", 2, "scaling factor")
	name := flag.String("name", "", "output font name")
	flag.Parse()
	bdf, err := bitedutils.R2BDF(os.Stdin)
	bitedutils.Check(err)
	if *name != "" {
		bdf.XLFD.Family = *name
	}
	err = bdf.Scale(*scale)
	bitedutils.Check(err)
	err = bdf.BDF2W(os.Stdout)
	bitedutils.Check(err)
}
