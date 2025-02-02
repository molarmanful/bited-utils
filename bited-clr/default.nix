{
  version,

  lib,
  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-clr";
  src = ../.;
  vendorHash = "sha256-YzcJicW1wJcgqDjmcimeZ+rioXxJyQz/K9Kd2E3ajus=";

  modRoot = "bited-clr";
}
