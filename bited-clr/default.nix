{
  version,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-clr";
  src = ../.;
  vendorHash = "sha256-rJwGPn912n83hJnEogB8mRVfVFknNQOPRYAP4mapebQ=";

  modRoot = "bited-clr";
}
