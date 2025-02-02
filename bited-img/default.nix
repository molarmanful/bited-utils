{
  version,

  bitsnpicas,

  lib,
  buildGoModule,
  makeWrapper,
  bash,
  perl,
  imagemagick,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-img";
  src = ../.;
  vendorHash = "sha256-DzyHsiUguAOZIo15SULtzgg52G9ftAypmdBm0uyk9SE=";
  modRoot = "bited-img";
  nativeBuildInputs = [ makeWrapper ];

  postFixup = ''
    wrapProgram $out/bin/bited-img \
      --set PATH ${
        lib.makeBinPath [
          bitsnpicas
          bash
          perl
          imagemagick
        ]
      }
  '';
}
