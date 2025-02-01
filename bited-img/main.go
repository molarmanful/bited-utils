package main

import (
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/molarmanful/bited-utils"
	"github.com/molarmanful/bited-utils/bited-img/lib"
)

func main() {
	k := koanf.New("")
	err := k.Load(file.Provider("bited-img.toml"), toml.Parser())
	bitedutils.Check(err)
	for _, name := range k.MapKeys("") {
		cfg, ok := k.Get(name).(map[string]any)
		if !ok {
			panic(name + " is not a map[string]any")
		}
		unit, err := bitedimg.NewUnit(cfg, name)
		bitedutils.Check(err)
		err = unit.Build()
		bitedutils.Check(err)
	}
}
