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
  vendorHash = "sha256-TpBCXEmOrUR4TLNuSYoNndMJNRZrSAAIMx7EvTcaU4s=";

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
