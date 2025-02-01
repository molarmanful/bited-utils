{
  version,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-scale";
  src = ../.;
  vendorHash = "sha256-4nSj30X4rn7FeK8bXIr1yIL0xIQ1zC6m/ffI3e139Yc=";

  modRoot = "bited-scale";
}
