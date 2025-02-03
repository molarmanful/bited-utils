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

A `.clr` file is really just lines of spaces and color codes. The color codes
are as follows:

| Code         | Definition                  |
| ------------ | --------------------------- |
| `0-9`, `A-F` | Set color to a Base16 color |
| `.`          | Reset to foreground color   |

All other characters are treated as no-ops and simply pass on the current color.

Codes are interpreted left-to-right, top-to-bottom. The position of a code
matches a position in the TXT file where you wish to color. Say you have the
following TXT:

```
Hello, world!
Testing Testing 123
```

And the following CLR:

```
A    . D    .
5       6       BCE
```

Would produce:

![TXT/CLR output](assets/txtclr_example.png)
