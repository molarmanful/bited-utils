module github.com/molarmanful/bited-utils/bited-img

go 1.23.4

require (
	github.com/bitfield/script v0.24.0
	github.com/knadh/koanf/parsers/toml/v2 v2.1.0
	github.com/knadh/koanf/providers/confmap v0.1.0
	github.com/knadh/koanf/providers/file v1.1.2
	github.com/knadh/koanf/providers/structs v0.1.0
	github.com/knadh/koanf/v2 v2.1.2
	github.com/molarmanful/bited-utils v0.0.0-00010101000000-000000000000
	github.com/molarmanful/bited-utils/bited-pango v0.0.0-00010101000000-000000000000
)

require (
	barista.run v0.0.0-20240418001405-c936f35316af // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/frankban/quicktest v1.14.6 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/itchyny/gojq v0.12.13 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/martinlindhe/unit v0.0.0-20220817221856-f7b595b5f97e // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	mvdan.cc/sh/v3 v3.7.0 // indirect
)

replace github.com/molarmanful/bited-utils => ../

replace github.com/molarmanful/bited-utils/bited-pango => ../bited-pango
