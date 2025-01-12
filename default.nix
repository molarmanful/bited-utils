{
  version,
  callPackage,
  symlinkJoin,
  bited-build ? callPackage ./bited-build { inherit version; },
  bited-img ? callPackage ./bited-img { inherit version; },
  bited-scale ? callPackage ./bited-scale { inherit version; },
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
