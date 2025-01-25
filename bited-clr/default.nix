{
  version,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-clr";
  src = ./.;
  vendorHash = "sha256-U7cUBVgnF6xEpYJfo2Qt9bLnSSVnaLNLvE+pQb308qY=";
}
