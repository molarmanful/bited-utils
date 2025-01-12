{
  version ? builtins.readFile ./VERSION,
  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-scale";
  src = ./.;
  vendorHash = null;
}
