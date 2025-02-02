{
  version,

  bitsnpicas,

  lib,
  buildGoModule,
  makeWrapper,
  fontforge,
  xorg,
  woff2,
  zip,
  nerd-font-patcher,
  ...
}:

buildGoModule {
  inherit version;
  pname = "bited-build";
  src = ../.;
  vendorHash = "sha256-JPA/3YquIZUV3mIyLj5S/6LiEov6VEaw7aR1NEGEWM8=";

  modRoot = "bited-build";
  nativeBuildInputs = [ makeWrapper ];

  postFixup = ''
    wrapProgram $out/bin/bited-build \
      --set PATH ${
        lib.makeBinPath [
          bitsnpicas
          fontforge
          xorg.bdftopcf
          woff2
          zip
          nerd-font-patcher
        ]
      }
  '';
}
