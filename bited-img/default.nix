{
  version,
  vendorHash,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-img";
  src = ../.;
  subPackages = [ "bited-img" ];
}
