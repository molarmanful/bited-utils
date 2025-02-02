module github.com/molarmanful/bited-utils/bited-build

go 1.23.4

require (
	github.com/bitfield/script v0.24.0
	github.com/knadh/koanf/parsers/toml/v2 v2.1.0
	github.com/knadh/koanf/providers/confmap v0.1.0
	github.com/knadh/koanf/providers/file v1.1.2
	github.com/knadh/koanf/providers/structs v0.1.0
	github.com/knadh/koanf/v2 v2.1.2
	github.com/molarmanful/bited-utils v0.0.0-00010101000000-000000000000
	github.com/molarmanful/bited-utils/bited-scale v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.2.0
)

require (
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/itchyny/gojq v0.12.13 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	golang.org/x/sys v0.21.0 // indirect
	mvdan.cc/sh/v3 v3.7.0 // indirect
)

replace github.com/molarmanful/bited-utils => ../

replace github.com/molarmanful/bited-utils/bited-scale => ../bited-scale
