# bited-utils
Pipeline helpers and utilities for building fonts from bited BDFs. Built with
Go and Nix flakes.

- **bited-build**: Given a `bited-build.toml` and bited BDF(s), produces
  vector (TTF, WOFF2) and bitmap (BDF, PCF, OTB, DFONT) fonts. Supports
  optional integer scaling and Nerd Font patching.
- **bited-img**: Given a `bited-img.toml`, bited BDF(s), and TXT/CLR pairs,
  produces image specimens with Base16 color support.
- **bited-clr**: TUI tool for fine-tuned coloring of TXT/CLR pairs.
- **bited-pango**: Library and CLI for converting TXT/CLR pairs into Go-typed
  Pango markup nodes.
- **bited-scale**: Library and CLI for scaling a bited BDF.
