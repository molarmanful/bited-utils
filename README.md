# bited-utils

Pipeline helpers and utilities for building fonts from bited BDFs. Built with Go
and Nix flakes.

- [**bited-build**](bited-build): A command that generates vector (TTF, WOFF2)
  and bitmap (BDF, PCF, OTB, DFONT) fonts from bited BDFs. Supports integer
  scaling and Nerd Font patching.
- [**bited-img**](bited-img): A command that generates image specimens for bited
  BDFs. Features Base16 color support.
- [**bited-clr**](bited-clr): A TUI tool for fine-tuned coloring of TXT/CLR
  files.
- [**bited-pango**](bited-pango): A library and command for converting TXT/CLR
  files into Go-typed Pango markup nodes.
- [**bited-scale**](bited-scale): A library and command for scaling a bited BDF.

## TXT/CLR files

To generate colorful images, bited-img accepts pairs of `.txt` and `.clr` files.
For example, if you have a file `test.txt` that you wish to color, you would
include `test.clr` containing your desired color codes in the same directory.
