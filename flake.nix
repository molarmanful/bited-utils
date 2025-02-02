{
  description = "Pipeline helpers and utilities for building fonts from bited BDFs";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    let
      version = builtins.readFile ./VERSION;
    in

    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        o = {
          inherit version;
        } // self.packages.${system};
      in
      {

        packages = {
          bitsnpicas = pkgs.callPackage ./bitsnpicas.nix o;
          bited-build = pkgs.callPackage ./bited-build o;
          bited-img = pkgs.callPackage ./bited-img o;
          bited-scale = pkgs.callPackage ./bited-scale o;
          bited-pango = pkgs.callPackage ./bited-pango o;
          bited-clr = pkgs.callPackage ./bited-clr o;
          default = pkgs.callPackage ./. o;
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
            errcheck
          ];
        };

      }
    );
}
