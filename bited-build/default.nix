{
  version,

  bitsnpicas,
  bited-scale,

  lib,
  stdenvNoCC,
  makeWrapper,
  nushell,
  fontforge,
  xorg,
  woff2,
  zip,
  nerd-font-patcher,
  ...
}:

stdenvNoCC.mkDerivation {
  inherit version;
  pname = "bited-build";
  src = ./.;

  nativeBuildInputs = [ makeWrapper ];

  installPhase = ''
    runHook preInstall
    mkdir -p $out/share $out/bin
    cp -r . $out/share
    makeWrapper ${nushell}/bin/nu $out/bin/bited-build \
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
      } \
      --add-flags "$out/share/build.nu"
    runHook postInstall
  '';
}
