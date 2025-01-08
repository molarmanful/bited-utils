{
  description = "A versatile bitmap font with an organic flair";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    bitsnpicas-src = {
      url = "github:kreativekorp/bitsnpicas?dir=main/java/BitsNPicas";
      flake = false;
    };
  };

  outputs =
    {
      nixpkgs,
      utils,
      bitsnpicas-src,
      ...
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };

        bitsnpicas = pkgs.stdenvNoCC.mkDerivation {
          name = "bitsnpicas";
          src = bitsnpicas-src;

          nativeBuildInputs = with pkgs; [
            temurin-bin
            makeWrapper
          ];

          preBuild = ''
            cd main/java/BitsNPicas
          '';

          buildFlags = "BitsNPicas.jar";

          installPhase = ''
            runHook preInstall
            mkdir -p $out/share/java $out/bin
            cp BitsNPicas.jar $out/share/java
            makeWrapper ${pkgs.temurin-jre-bin}/bin/java $out/bin/bitsnpicas \
              --add-flags "-jar $out/share/java/BitsNPicas.jar"
            runHook postInstall
          '';
        };

        bited-build = pkgs.stdenvNoCC.mkDerivation {
          name = "bited-build";
          src = ./.;

          buildInputs = with pkgs; [
            git
            bitsnpicas
            fontforge
            xorg.bdftopcf
            woff2
            zip
            nerd-font-patcher
            makeWrapper
          ];

          installPhase = ''
            runHook preInstall
            mkdir -p $out/share $out/bin
            cp -r src/{build.nu,scripts} $out/share
            makeWrapper ${pkgs.nushell}/bin/nu $out/bin/bited-build \
              --add-flags "$out/share/build.nu"
            runHook postInstall
          '';
        };

        bited-img = pkgs.stdenvNoCC.mkDerivation {
          name = "bited-img";
          src = ./.;

          buildInputs = with pkgs; [
            bitsnpicas
            imagemagick
            nushell
            makeWrapper
          ];

          installPhase = ''
            runHook preInstall
            mkdir -p $out/share $out/bin
            cp -r src/img.nu $out/share
            makeWrapper ${pkgs.nushell}/bin/nu $out/bin/bited-img \
              --add-flags "$out/share/img.nu"
            runHook postInstall
          '';
        };

        bited-utils = pkgs.symlinkJoin {
          name = "bited-utils";

          paths = [
            bited-build
            bited-img
          ];
        };

      in
      {

        devShell = pkgs.mkShell {
          packages = with pkgs; [
            nil
            nixd
            nixfmt-rfc-style
            statix
            deadnix
            nushell
          ];
        };

        packages = {
          inherit
            bitsnpicas
            bited-build
            bited-img
            bited-utils
            ;
          default = bited-utils;
        };
      }
    );
}
