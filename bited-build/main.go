package main

import (
	"flag"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/molarmanful/bited-utils"
	"github.com/molarmanful/bited-utils/bited-build/lib"
)

func main() {
	nerd := flag.Bool("nerd", false, "whether to compile Nerd Fonts variants")
	flag.Parse()

	k := koanf.New("")
	err := k.Load(file.Provider("bited-build.toml"), toml.Parser())
	bitedutils.Check(err)
	for _, name := range k.MapKeys("") {
		cfg, ok := k.Get(name).(map[string]any)
		if !ok {
			panic(name + " is not a map[string]any")
		}
		unit, err := bitedbuild.NewUnit(cfg, name, *nerd)
		bitedutils.Check(err)
		err = unit.Build()
		bitedutils.Check(err)
	}
}
