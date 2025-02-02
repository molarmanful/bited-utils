{
  version,
  vendorHash,

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
  inherit version vendorHash;
  pname = "bited-img";
  src = ../.;
  subPackages = [ "bited-img" ];
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
