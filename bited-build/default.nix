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
  vendorHash = "sha256-4lG3R9tCOY2tbKIBvGl7AOHowvV35nr72eQ9DZrn858=";

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
