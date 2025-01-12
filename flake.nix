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
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
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
            go
            gopls
            gotools
          ];
        };

        packages = rec {
          bitsnpicas = pkgs.callPackage ./bitsnpicas { };
          bited-build = pkgs.callPackage ./bited-build { };
          bited-img = pkgs.callPackage ./bited-img { };
          bited-scale = pkgs.callPackage ./bited-scale { };
          bited-utils = pkgs.callPackage ./. { };
          default = bited-utils;
        };
      }
    );
}
