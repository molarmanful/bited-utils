{
  version,

  lib,
  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-pango";
  src = ../.;
  vendorHash = "sha256-oHra99LHFq1bv6p1gP7+sZlbHvaAaOB2yV5ux4mfLL8=";

  modRoot = "bited-pango";
}
