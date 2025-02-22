{
  version,
  vendorHash,

  buildGoModule,
  licenses,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-img";
  src = ../.;
  subPackages = [ "bited-img" ];

  meta = {
    description = "An image specimen generator for bited BDF fonts";
    longDescription = ''
      bited-img generates image specimens for bited BDFs. It features Base16
      color support.
    '';
    homepage = "https://github.com/molarmanful/bited-utils";
    license = licenses.mit;
  };
}
