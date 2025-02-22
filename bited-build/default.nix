{
  version,
  vendorHash,

  bitsnpicas,

  lib,
  buildGoModule,
  makeWrapper,
  fontforge,
  xorg,
  woff2,
  zip,
  nerd-font-patcher,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-build";
  src = ../.;
  subPackages = [ "bited-build" ];
  nativeBuildInputs = [ makeWrapper ];
  postFixup = ''
    wrapProgram $out/bin/bited-build \
      --set PATH ${
        lib.makeBinPath [
          bitsnpicas
          fontforge
          xorg.bdftopcf
          woff2
          zip
          nerd-font-patcher
        ]
      }
  '';

  meta = with lib; {
    description = "A builder for bited BDFs";
    longDescription = ''
      bited-build generates vector (TTF, WOFF2) and bitmap (BDF, PCF, OTB,
      DFONT) fonts from bited BDFs. It supports integer scaling and Nerd Font
      patching.
    '';
    homepage = "https://github.com/molarmanful/bited-utils";
    license = licenses.mit;
  };
}
