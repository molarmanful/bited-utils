{
  version,
  vendorHash,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-clr";
  src = ../.;
  subPackages = [ "bited-clr" ];
}
