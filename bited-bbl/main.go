// bited-bbl is a command for proportionalizing bited BDF DWIDTHs.
//
// Usage:
//
//	bited-bbl
//
// Flags:
//
//	--name
//		Family name of the output font.
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
	flag.Parse()
	err := bitedbbl.Bbl(os.Stdin, os.Stdout, *name)
	bitedutils.Check(err)
}
