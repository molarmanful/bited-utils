{
  version,

  bitsnpicas,
  bited-scale,

  lib,
  buildGoModule,
  makeWrapper,
  fontforge,
  xorg,
  woff2,
  zip,
  nerd-font-patcher,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-build";
  src = ../.;
  vendorHash = "sha256-/yI0zBqOOhN+PrgF8WvHgdU1zwsmxr6gg0PuHDb1Y2Q=";

  modRoot = "bited-build";
  nativeBuildInputs = [ makeWrapper ];

  postFixup = ''
    wrapProgram $out/bin/bited-build \
      --set PATH ${
        lib.makeBinPath [
          bitsnpicas
          bited-scale
          fontforge
          xorg.bdftopcf
          woff2
          zip
          nerd-font-patcher
        ]
      }
  '';
}
