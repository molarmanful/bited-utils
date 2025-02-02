{
  version,
  vendorHash,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-utils";
  src = ./.;
}
