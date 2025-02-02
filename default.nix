{
  version,

  bited-build,
  bited-img,
  bited-scale,
  bited-pango,
  bited-clr,

  symlinkJoin,
  ...
}:

symlinkJoin {
  inherit version;
  pname = "bited-utils";
  paths = [
    bited-build
    bited-img
    bited-scale
    bited-pango
    bited-clr
  ];
}
