# bited-utils

Pipeline helpers and utilities for building fonts from bited BDFs. Built with Go
and Nix flakes.

- **bited-build**: Command that generates vector (TTF, WOFF2) and bitmap (BDF,
  PCF, OTB, DFONT) fonts from bited BDFs. Supports optional integer scaling and
  Nerd Font patching.
- **bited-img**: Command that generates image specimens for bited BDFs. Features
  Base16 color support.
- **bited-clr**: TUI tool for fine-tuned coloring of TXT/CLR pairs.
- **bited-pango**: Library and command for converting TXT/CLR pairs into
  Go-typed Pango markup nodes.
- **bited-scale**: Library and command for scaling a bited BDF.
