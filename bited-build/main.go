package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/molarmanful/bited-utils"
	"github.com/molarmanful/bited-utils/bited-build/lib"
)

func main() {
	cfgF := flag.String("cfg", "bited-build.toml", "TOML config path")
	nerd := flag.Bool("nerd", false, "whether to compile Nerd Fonts variants")
	flag.Parse()

	var pre map[string]toml.Primitive
	md, err := toml.DecodeFile(*cfgF, &pre)
	bitedutils.Check(err)

	for name, pv := range pre {
		err := bitedbuild.FullUnit(md, pv, name, *nerd)
		bitedutils.Check(err)
	}
}
