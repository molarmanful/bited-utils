{
  version,
  vendorHash,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-scale";
  src = ../.;
  subPackages = [ "bited-scale" ];
}
