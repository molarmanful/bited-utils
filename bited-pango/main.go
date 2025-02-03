package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	bitedutils "github.com/molarmanful/bited-utils"
	bitedpango "github.com/molarmanful/bited-utils/bited-pango/lib"
)

func main() {
	base := flag.String("base", "", "base path for txt and clr file")
	clrsStr := flag.String("clrs", "", "newline-separated string of colors 0-15")
	flag.Parse()

	txtF, err := os.Open(*base + ".txt")
	bitedutils.Check(err)
	defer txtF.Close()

	clrF, err := os.Open(*base + ".clr")
	bitedutils.Check(err)
	defer clrF.Close()

	fmt.Println(bitedpango.Pango(txtF, clrF, strings.SplitN(*clrsStr, "\n", 16)).String())
}
