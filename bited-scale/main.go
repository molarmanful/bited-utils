package main

import (
	"flag"
	"os"

	bitedutils "github.com/molarmanful/bited-utils"
	bitedscale "github.com/molarmanful/bited-utils/bited-scale/lib"
)

func main() {
	scale := flag.Int("n", 2, "scaling factor")
	name := flag.String("name", "", "scaled font name")
	flag.Parse()
	err := bitedscale.Scale(os.Stdin, os.Stdout, *scale, *name)
	bitedutils.Check(err)
}
