{
  version,

  buildGoModule,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-pango";
  src = ../.;
  vendorHash = "sha256-n48qF+ZIKhsFC3whpEPHjrftTQWXoQZNaHD8uwJjzWc=";

  modRoot = "bited-pango";
}
