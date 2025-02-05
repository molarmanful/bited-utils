# FIXME

A cool font waiting to be built.

## Getting Started

This project's default font name is `FIXME`, which you should replace with your
own font's name. The project's default author is
[ghost](https://github.com/ghost), which you should also replace with your name.

> [!TIP]
> A tool like [rgr](https://github.com/acheronfail/repgrep) can make this
> process easier!

### Updating

```
nix flake update
```

To only update bited-utils:

```
nix flake update bited-utils
```

### Versioning

By default, Nix flake reads versions via the `VERSION` file.
`.github/workflows/pub.yml` bumps the version based on Github releases. These
can all be modified to fit your preferred versioning scheme.

### CI

The default CI caching solution defined is
[FlakeHub Cache](https://flakehub.com/cache); feel free to replace this with
another caching solution of your choice or simply remove it if you don't need
one.
