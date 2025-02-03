# bited-img

A command that generates image specimens for bited BDFs. Features Base16 color
support.

## Usage

```
bited-img
```

bited-img reads configuration from `bited-img.toml` in the current working
directory, typically the font project's root. Its structure is documented
[here](bited-img.toml).

## TXT/CLR files

To generate colorful images, bited-img accepts pairs of `.txt` and `.clr` files.
For example, if you have a file `test.txt` that you wish to color, you would
include `test.clr` containing your desired color codes in the same directory.
