# bited-utils

Pipeline helpers and utilities for building fonts from bited BDFs. Built with Go
and Nix flakes.

- [**bited-build**](bited-build): A command that generates vector (TTF, WOFF2)
  and bitmap (BDF, PCF, OTB, DFONT) fonts from bited BDFs. Supports integer
  scaling and Nerd Font patching.
- [**bited-img**](bited-img): A command that generates image specimens for bited
  BDFs. Features Base16 color support.
- [**bited-clr**](bited-clr): A TUI tool for fine-tuned coloring of TXT/CLR
  pairs.
- [**bited-pango**](bited-pango): A library and command for converting TXT/CLR
  pairs into Go-typed Pango markup nodes.
- [**bited-scale**](bited-scale): A library and command for scaling a bited BDF.
