// bited-scale is a command for scaling bited BDFs.
//
// Usage:
//
//	bited-scale [-x <int>] [--name <string>]
//
// Flags:
//
//	--x
//		Scaling factor.
//	--name
//		Family name of the scaled font.
//
// bited-scale accepts input via STDIN and outputs to STDOUT.
package main

import (
	"flag"
	"os"

	bitedutils "github.com/molarmanful/bited-utils"
	bitedscale "github.com/molarmanful/bited-utils/bited-scale/lib"
)

func main() {
	scale := flag.Int("x", 2, "scaling factor")
	name := flag.String("name", "", "scaled font name")
	flag.Parse()
	err := bitedscale.Scale(os.Stdin, os.Stdout, *scale, *name)
	bitedutils.Check(err)
}
