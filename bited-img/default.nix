{
  version,

  lib,
  stdenvNoCC,
  makeWrapper,
  bitsnpicas,
  nushell,
  bash,
  imagemagick,
  ...
}:

stdenvNoCC.mkDerivation {
  inherit version;
  pname = "bited-img";
  src = ./.;

  nativeBuildInputs = [ makeWrapper ];

  installPhase = ''
    runHook preInstall
    mkdir -p $out/share $out/bin
    cp -r . $out/share
    makeWrapper ${nushell}/bin/nu $out/bin/bited-img \
      --set PATH ${
        lib.makeBinPath [
          bitsnpicas
          bash
          imagemagick
        ]
      } \
      --add-flags "$out/share/img.nu"
    runHook postInstall
  '';
}
