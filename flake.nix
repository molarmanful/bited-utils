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
          vendorHash = "sha256-GQuDse45vv7artjF3aWrp6rRuRsb9odu7Iey5Vxa/V8=";
        } // self.packages.${system};
      in
      {

        packages = rec {
          default = bited-utils;
          bited-utils = pkgs.callPackage ./. o;
          bitsnpicas = pkgs.callPackage ./bitsnpicas.nix { };
          bited-build = pkgs.callPackage ./bited-build o;
          bited-img = pkgs.callPackage ./bited-img o;
          bited-scale = pkgs.callPackage ./bited-scale o;
          bited-pango = pkgs.callPackage ./bited-pango o;
          bited-clr = pkgs.callPackage ./bited-clr o;
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nil
            nixd
            nixfmt-rfc-style
            statix
            deadnix
            taplo
            go
            gopls
            gotools
            golines
            errcheck
            marksman
            mdformat
            python312Packages.mdformat-gfm
            python312Packages.mdformat-frontmatter
            python312Packages.mdformat-footnote
            python312Packages.mdformat-gfm-alerts
          ];
        };

      }
    );
}
