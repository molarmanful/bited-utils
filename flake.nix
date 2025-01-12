{
  description = "Pipeline helpers and utilities for building fonts from bited BDFs";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
      ...
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = self.shortRev or self.dirtyShortRev;
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
          bitsnpicas = pkgs.callPackage ./bitsnpicas { inherit version; };
          bited-build = pkgs.callPackage ./bited-build { inherit version bitsnpicas bited-scale; };
          bited-img = pkgs.callPackage ./bited-img { inherit version bitsnpicas; };
          bited-scale = pkgs.callPackage ./bited-scale { inherit version; };
          bited-utils = pkgs.callPackage ./. {
            inherit
              version
              bited-build
              bited-img
              bited-scale
              ;
          };
          default = bited-utils;
        };
      }
    );
}
