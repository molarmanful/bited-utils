{
  pname,
  version,
  src,
  nerd ? false,
  release ? false,

  lib,
  stdenvNoCC,
  bited-build,
  zip,
  ...
}:

stdenvNoCC.mkDerivation {
  inherit pname version src;

  buildPhase = ''
    runHook preBuild
    rm -rf out
    ${bited-build}/bin/bited-build ${lib.optionalString nerd "--nerd"}
    runHook postBuild
  '';

  installPhase = ''
    runHook preInstall
    mkdir -p $out/share/fonts
    cp -r out/. $out/share/fonts
    ${lib.optionalString release ''
      pushd $out/share
      ${zip}/bin/zip -r fonts/${pname}_v${version}.zip fonts
      popd
    ''}
    runHook postInstall
  '';
}
