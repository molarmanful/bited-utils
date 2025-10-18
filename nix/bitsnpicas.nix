{
  stdenvNoCC,
  fetchFromGitHub,
  makeWrapper,
  temurin-bin,
  temurin-jre-bin,
  ...
}:

stdenvNoCC.mkDerivation (finalAttrs: {
  pname = "bitsnpicas";
  version = "2.1.1";

  src =
    fetchFromGitHub {
      owner = "kreativekorp";
      repo = "bitsnpicas";
      rev = "v${finalAttrs.version}";
      sha256 = "sha256-YC3SZckxB53Kna8M9NPoQRpViWmxEu0YR/OHgxEjPiU=";
      sparseCheckout = [ "main/java/BitsNPicas" ];
    }
    + "/main/java/BitsNPicas";

  nativeBuildInputs = [
    temurin-bin
    makeWrapper
  ];

  buildFlags = "BitsNPicas.jar";

  installPhase = ''
    runHook preInstall
    mkdir -p $out/share/java $out/bin
    cp BitsNPicas.jar $out/share/java
    makeWrapper ${temurin-jre-bin}/bin/java $out/bin/bitsnpicas \
      --add-flags "-jar $out/share/java/BitsNPicas.jar"
    runHook postInstall
  '';
})
