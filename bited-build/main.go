// bited-build generates vector (TTF, WOFF2) and bitmap (BDF, PCF, OTB, DFONT)
// fonts from bited BDFs. It supports integer scaling and Nerd Font patching.
//
// Usage:
//
//	bited-build [--nerd]
//
// Flags:
//
//	--nerd
//		Whether to compile Nerd Fonts variants.
//
// bited-build reads configuration from [bited-build.toml] in the current
// working directory, typically the font project's root.
//
// [bited-build.toml]: https://github.com/molarmanful/bited-utils/blob/main/bited-build/bited-build.toml
package main

import (
	"flag"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	bitedutils "github.com/molarmanful/bited-utils"
	bitedbuild "github.com/molarmanful/bited-utils/bited-build/lib"
)

func main() {
	nerd := flag.Bool("nerd", false, "whether to compile Nerd Fonts variants")
	flag.Parse()

	k := koanf.New("")
	err := k.Load(file.Provider("bited-build.toml"), toml.Parser())
	bitedutils.Check(err)
	for _, cfg := range k.Slices("fonts") {
		unit, err := bitedbuild.NewUnit(cfg, *nerd)
		bitedutils.Check(err)
		err = unit.Build()
		bitedutils.Check(err)
	}
}
