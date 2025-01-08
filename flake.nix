{
  description = "Pipeline helpers and utilities for building fonts from bited BDFs";

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

        bited-scale = pkgs.stdenvNoCC.mkDerivation {
          name = "bited-scale";
          src = ./.;

          nativeBuildInputs = with pkgs; [ makeWrapper ];

          installPhase = ''
            runHook preInstall
            mkdir -p $out/share
            cp src/scale.nu $out/share
            makeWrapper ${pkgs.nushell}/bin/nu $out/bin/bited-scale \
              --add-flags "$out/share/scale.nu"
            runHook postInstall
          '';
        };

        bited-build = pkgs.stdenvNoCC.mkDerivation {
          name = "bited-build";
          src = ./.;

          nativeBuildInputs = with pkgs; [ makeWrapper ];

          installPhase = ''
            runHook preInstall
            mkdir -p $out/share/deps $out/bin
            cp src/build.nu $out/share
            cp -r src/deps/build $out/share/deps
            makeWrapper ${pkgs.nushell}/bin/nu $out/bin/bited-build \
              --set PATH ${
                with pkgs;
                lib.makeBinPath [
                  bited-scale
                  git
                  bitsnpicas
                  fontforge
                  xorg.bdftopcf
                  woff2
                  zip
                  nerd-font-patcher
                ]
              } \
              --add-flags "$out/share/build.nu"
            runHook postInstall
          '';
        };

        bited-img = pkgs.stdenvNoCC.mkDerivation {
          name = "bited-img";
          src = ./.;

          nativeBuildInputs = with pkgs; [ makeWrapper ];

          installPhase = ''
            runHook preInstall
            mkdir -p $out/share/deps $out/bin
            cp src/img.nu $out/share
            cp -r src/deps/img $out/share/deps
            makeWrapper ${pkgs.nushell}/bin/nu $out/bin/bited-img \
              --set PATH ${
                with pkgs;
                lib.makeBinPath [
                  bash
                  bitsnpicas
                  imagemagick
                ]
              } \
              --add-flags "$out/share/img.nu"
            runHook postInstall
          '';
        };

        bited-utils = pkgs.symlinkJoin {
          name = "bited-utils";

          paths = [
            bited-scale
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
            yamlfix
            lua-language-server
            stylua
            selene
          ];
        };

        packages = {
          inherit
            bitsnpicas
            bited-scale
            bited-build
            bited-img
            bited-utils
            ;
          default = bited-utils;
        };
      }
    );
}
