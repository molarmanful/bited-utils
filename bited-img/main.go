// bited-img generates image specimens for bited BDFs. It features Base16 color
// support.
//
// Usage:
//
//	bited-img
//
// bited-img reads configuration from [bited-img.toml] in the current
// working directory, typically the font project's root.
//
// [bited-img.toml]: https://github.com/molarmanful/bited-utils/blob/main/bited-img/bited-img.toml
package main

import (
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	bitedutils "github.com/molarmanful/bited-utils"
	bitedimg "github.com/molarmanful/bited-utils/bited-img/lib"
)

func main() {
	k := koanf.New("")
	err := k.Load(file.Provider("bited-img.toml"), toml.Parser())
	bitedutils.Check(err)
	for _, cfg := range k.Slices("fonts") {
		unit, err := bitedimg.NewUnit(cfg)
		bitedutils.Check(err)
		err = unit.Build()
		bitedutils.Check(err)
	}
}
