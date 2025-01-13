{
  stdenvNoCC,
  fetchFromGitHub,
  makeWrapper,
  temurin-bin,
  temurin-jre-bin,
  ...
}:

stdenvNoCC.mkDerivation rec {
  pname = "bitsnpicas";
  version = "c6804949137229ef5a0c185e1464ef6b9e222601";

  src = fetchFromGitHub {
    owner = "molarmanful"; # TODO: switch back to kreativekorp if pr merges
    repo = pname;
    rev = version;
    sha256 = "sha256-7qwjZj0dTZyTh6GSsDfDY8b/J5yfmRfWzPBpbuS2xhw=";
    sparseCheckout = [ "main/java/BitsNPicas" ];
  };

  nativeBuildInputs = [
    temurin-bin
    makeWrapper
  ];

  preBuild = ''
    cd main/java/BitsNPicas
  '';

  buildFlags = "BitsNPicas.jar";

  installPhase = ''
    runHook preInstall
    mkdir -p $out/share/java $out/bin
    cp BitsNPicas.jar $out/share/java
    makeWrapper ${temurin-jre-bin}/bin/java $out/bin/bitsnpicas \
      --add-flags "-jar $out/share/java/BitsNPicas.jar"
    runHook postInstall
  '';
}
