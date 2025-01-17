{
  version,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-pangogo";
  src = ./.;
  vendorHash = null;
}
