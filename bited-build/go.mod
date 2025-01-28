module github.com/molarmanful/bited-utils/bited-build

go 1.23.4

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/bitfield/script v0.24.0
	github.com/molarmanful/bited-utils v0.0.0-00010101000000-000000000000
	github.com/molarmanful/bited-utils/bited-scale v0.0.0-00010101000000-000000000000
)

require (
	github.com/itchyny/gojq v0.12.13 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	mvdan.cc/sh/v3 v3.7.0 // indirect
)

replace github.com/molarmanful/bited-utils => ../

replace github.com/molarmanful/bited-utils/bited-scale => ../bited-scale
