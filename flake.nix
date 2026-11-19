{
  description = "Pipeline helpers and utilities for building fonts from bited BDFs";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } (
      { flake-parts-lib, withSystem, ... }:

      let
        flakeModule = flake-parts-lib.importApply ./flake-module.nix {
          inherit withSystem;
          inherit (flake-parts-lib) mkPerSystemOption;
          inherit (inputs.nixpkgs) lib;
        };
      in

      {
        imports = [
          inputs.treefmt-nix.flakeModule
          flakeModule
        ];

        flake = {
          inherit flakeModule;
          templates.default = {
            path = ./template;
            description = "bited font project with bited-utils";
          };
        };

        systems = inputs.nixpkgs.lib.systems.flakeExposed;
        perSystem =
          { pkgs, ... }:
          let
            mdformat-plugins =
              ps: with ps; [
                mdformat-gfm
                mdformat-gfm-alerts
              ];
          in
          {
            packages =
              let
                bitsnpicas = pkgs.callPackage ./nix/bitsnpicas.nix { };
                args = {
                  inherit bitsnpicas;
                  version = builtins.readFile ./VERSION;
                  vendorHash = "sha256-Ya4gSRC2blGXMfcMBAhKHOxz4r6Rb0WRNls8g8vXPNw=";
                };
              in
              {
                inherit bitsnpicas;
                bited-build = pkgs.callPackage ./bited-build args;
                bited-img = pkgs.callPackage ./bited-img args;
                bited-scale = pkgs.callPackage ./bited-scale args;
                bited-clr = pkgs.callPackage ./bited-clr args;
                bited-bbl = pkgs.callPackage ./bited-bbl args;
              };

            devShells.default = pkgs.mkShell {
              packages = with pkgs; [
                go
                go-tools
                taplo
                # lsps
                marksman
                gopls
                # formatters
                mdformat
                (python3.withPackages mdformat-plugins)
                golines
                yamlfmt
                actionlint
                # linters
                errcheck
              ];
            };

            treefmt = {
              projectRootFile = "flake.nix";
              programs = {
                nixfmt.enable = true;
                statix.enable = true;
                deadnix.enable = true;
                yamlfmt.enable = true;
                actionlint.enable = true;
                mdformat = {
                  enable = true;
                  plugins = mdformat-plugins;
                  includes = [ "docs/**/*.md" ];
                };
                taplo.enable = true;
                golines.enable = true;
              };
            };
          };
      }
    );
}
