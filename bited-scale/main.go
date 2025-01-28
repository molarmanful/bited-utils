package main

import (
	"flag"
	"os"

	"github.com/molarmanful/bited-utils/bited-scale/lib"
)

var scale = flag.Int("n", 2, "scaling factor")
var name = flag.String("name", "", "scaled font name")

func main() {
	flag.Parse()
	if err := bitedscale.Scale(os.Stdin, os.Stdout, *scale, *name); err != nil {
		panic(err)
	}
}
