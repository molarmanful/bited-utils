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
        pkgs = nixpkgs.legacyPackages.${system};
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
          bited-build = pkgs.callPackage ./bited-build { inherit version; };
          bited-img = pkgs.callPackage ./bited-img { inherit version; };
          bited-scale = pkgs.callPackage ./bited-scale { inherit version; };
          bited-utils = pkgs.callPackage ./. { inherit version; };
          default = bited-utils;
        };
      }
    );
}
