{
  version,
  vendorHash,

  lib,
  buildGoModule,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-bbl";
  src = ../.;
  subPackages = [ "bited-bbl" ];

  meta = with lib; {
    description = "A command for proportionalizing bited BDF DWIDTHs";
    longDescription = ''
      bited-bbl is a command that adjusts glyphs DWIDTHs bited BDFs based on their respective BBXs.
    '';
    homepage = "https://github.com/molarmanful/bited-utils";
    license = licenses.mit;
  };
}
