{
  version,

  bitsnpicas,
  bited-pango, # FIXME

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
  vendorHash = "sha256-XAWy2JTaKAnPUrORfdXp3dPwKhXdHz/8rHTvzJQ67cA=";

  modRoot = "bited-img";
  nativeBuildInputs = [ makeWrapper ];

  postFixup = ''
    wrapProgram $out/bin/bited-img \
      --set PATH ${
        lib.makeBinPath [
          bitsnpicas
          bited-pango
          bash
          perl
          imagemagick
        ]
      }
  '';
}
