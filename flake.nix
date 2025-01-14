{
  description = "Pipeline helpers and utilities for building fonts from bited BDFs";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      nixpkgs,
      utils,
      ...
    }:
    let
      o = {
        version = builtins.readFile ./VERSION;
      };
    in

    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system}.extend (

          final: prev: {
            bitsnpicas = final.callPackage ./bitsnpicas.nix o;
            bited-build = final.callPackage ./bited-build o;
            bited-img = final.callPackage ./bited-img o;
            bited-scale = final.callPackage ./bited-scale o;
            bited-utils = final.callPackage ./. o;
          }

        );
      in
      {

        packages = {
          inherit (pkgs)
            bitsnpicas
            bited-build
            bited-img
            bited-scale
            bited-utils
            ;
          default = pkgs.bited-utils;
        };

        devShell = pkgs.mkShell {
          packages = with pkgs; [
            nil
            nixd
            nixfmt-rfc-style
            statix
            deadnix
            nushell
            taplo
            go
            gopls
            gotools
          ];
        };

      }
    );
}
