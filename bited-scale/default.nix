{
  version,
  stdenvNoCC,
  makeWrapper,
  nushell,
  ...
}:

stdenvNoCC.mkDerivation {
  inherit version;
  pname = "bited-scale";
  src = ./.;

  nativeBuildInputs = [ makeWrapper ];

  installPhase = ''
    runHook preInstall
    mkdir -p $out/share $out/bin
    cp -r . $out/share
    makeWrapper ${nushell}/bin/nu $out/bin/bited-scale \
      --add-flags "$out/share/scale.nu"
    runHook postInstall
  '';
}
