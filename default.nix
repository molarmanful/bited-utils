{
  version ? builtins.readFile ./VERSION,
  callPackage,
  symlinkJoin,
  bited-build ? callPackage ./bited-build { },
  bited-img ? callPackage ./bited-img { },
  bited-scale ? callPackage ./bited-scale { },
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
