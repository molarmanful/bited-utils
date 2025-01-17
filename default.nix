{
  version,
  P,

  symlinkJoin,
  ...
}:

symlinkJoin {
  inherit version;
  pname = "bited-utils";

  paths = [
    P.bited-build
    P.bited-img
    P.bited-scale
    P.bited-pangogo
  ];
}
