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
    let
      version = builtins.readFile ./VERSION;
    in

    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        o = {
          inherit version;
          P = self.packages.${system};
        };
      in
      {

        packages = rec {
          bitsnpicas = pkgs.callPackage ./bitsnpicas.nix o;
          bited-build = pkgs.callPackage ./bited-build o;
          bited-img = pkgs.callPackage ./bited-img o;
          bited-scale = pkgs.callPackage ./bited-scale o;
          bited-pangogo = pkgs.callPackage ./bited-pangogo o;
          bited-utils = pkgs.callPackage ./. o;
          default = bited-utils;
        };

        devShells.default = pkgs.mkShell {
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
