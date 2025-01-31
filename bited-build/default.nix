{
  version,

  bitsnpicas,

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
  vendorHash = "sha256-H5YDpshOZlSfntBxQIMwo00eDWRwrJ4xbtjxjDWxjHM=";

  modRoot = "bited-build";
  nativeBuildInputs = [ makeWrapper ];

  postFixup = ''
    wrapProgram $out/bin/bited-build \
      --set PATH ${
        lib.makeBinPath [
          bitsnpicas
          fontforge
          xorg.bdftopcf
          woff2
          zip
          nerd-font-patcher
        ]
      }
  '';
}
