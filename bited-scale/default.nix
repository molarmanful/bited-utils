{
  version,
  vendorHash,

  buildGoModule,
  licenses,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-scale";
  src = ../.;
  subPackages = [ "bited-scale" ];

  meta = {
    description = "A command for scaling bited BDFs";
    longDescription = ''
      bited-scale is a command for scaling bited BDFs.
    '';
    homepage = "https://github.com/molarmanful/bited-utils";
    license = licenses.mit;
  };
}
