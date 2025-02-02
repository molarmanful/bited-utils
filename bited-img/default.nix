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
  vendorHash = "sha256-MqLXFi9Yc+ds3Mn4pi7/nFfosUSqaAz9a1fIDWqesW0=";
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
