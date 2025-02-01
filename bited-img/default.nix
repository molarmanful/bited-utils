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
  vendorHash = "sha256-zAk4gJethRyW5oGiy6v/Yv+DvrPWk/y4q7VlsIWPwtw=";
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
