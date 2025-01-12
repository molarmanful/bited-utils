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
  pname = "bited-utils";

  paths = [
    bited-build
    bited-img
    bited-scale
  ];
}
