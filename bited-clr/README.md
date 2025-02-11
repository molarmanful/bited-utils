# bited-clr

A TUI tool for fine-tuned coloring of TXT/CLR files.

## Usage

```bash
bited-clr --name <string> --stem <string>
```

- **name**: Font name to retrieve colors from in `bited-img.toml`.
- **stem**: TXT/CLR pair to edit (relative to `txt_dir` in `bited-img.toml`).

bited-clr uses `bited-img.toml` in the current working directory, typically the
font project's root.
