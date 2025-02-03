# bited-build

A command that generates vector (TTF, WOFF2) and bitmap (BDF, PCF, OTB, DFONT)
fonts from bited BDFs. Supports integer scaling and Nerd Font patching.

## Usage

```
bited-build [--nerd]
```

- **nerd**: Whether to compile Nerd Fonts variants.

bited-build reads configuration from `bited-build.toml` in the current working
directory, typically the font project's root. Its structure is documented
[here](bited-build.toml).
