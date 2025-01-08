{
  description = "A versatile bitmap font with an organic flair";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
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
        pkgs = nixpkgs.legacyPackages.${system};

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

        bited-build = pkgs.writeShellApplication {
          name = "bited-build";

          runtimeInputs = with pkgs; [
            git
            bitsnpicas
            fontforge
            xorg.bdftopcf
            woff2
            nushell
            zip
            nerd-font-patcher
          ];

          text = ''
            nu main.nu "$@"
          '';
        };

        bited-img = pkgs.writeShellApplication {
          name = "bited-img";

          runtimeInputs = with pkgs; [
            bitsnpicas
            imagemagick
            nushell
          ];

          text = ''
            nu img.nu "$@"
          '';
        };

      in
      {

        devShell = pkgs.mkShell {
          packages = with pkgs; [
            nushell
          ];
        };

        packages = {
          inherit
            bitsnpicas
            bited-build
            bited-img
            ;
        };
      }
    );
}
