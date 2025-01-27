{
  version,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-clr";
  src = ../.;
  vendorHash = "sha256-Qf4Vidk42Wo2OBQYtI2Jwf/hU9yHTGsPHyXL5/W/LDc=";

  modRoot = "bited-clr";
}
