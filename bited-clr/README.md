# bited-clr

A TUI tool for fine-tuned coloring of TXT/CLR files.

![screenshot of bited-clr](screen.png)

## Usage

```bash
bited-clr --name <string> --stem <string>
```

- **name**: Font name to retrieve colors from in `bited-img.toml`.
- **stem**: TXT/CLR pair to edit (relative to `txt_dir` in `bited-img.toml`).

bited-clr uses `bited-img.toml` in the current working directory, typically the
font project's root.

### Keybindings

| Keys               | Action           |
| ------------------ | ---------------- |
| `hjkl`, arrow keys | move cursor      |
| `.`                | foreground color |
| `0-9`, `a-f`       | Base16 colors    |
