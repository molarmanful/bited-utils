{
  version,
  symlinkJoin,
  bited-build,
  bited-img,
  bited-scale,
  ...
}:

symlinkJoin {
  inherit version;

  name = "bited-utils";
  paths = [
    bited-build
    bited-img
    bited-scale
  ];
}
