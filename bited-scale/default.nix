{
  version,

  lib,
  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-scale";
  src = ../.;
  vendorHash = "sha256-sbDr0DcZmlOD2OBSRLtASQ1oTSkY8GrJG4X+FayGDJc=";

  modRoot = "bited-scale";
}
