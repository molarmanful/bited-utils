{
  version,
  vendorHash,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-pango";
  src = ../.;
  subPackages = [ "bited-pango" ];
}
