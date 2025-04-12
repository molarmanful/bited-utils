{
  description = "Pipeline helpers and utilities for building fonts from bited BDFs";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    systems.url = "systems";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devshell.url = "github:numtide/devshell";
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
          inputs.devshell.flakeModule
          flakeModule
        ];

        flake = {
          inherit flakeModule;
          templates.default = {
            path = ./template;
            description = "bited font project with bited-utils";
          };
        };

        systems = import inputs.systems;
        perSystem =
          {
            pkgs,
            ...
          }:
          {

            packages =
              let
                bitsnpicas = pkgs.callPackage ./nix/bitsnpicas.nix { };
                args = {
                  inherit bitsnpicas;
                  version = builtins.readFile ./VERSION;
                  vendorHash = "sha256-/EcjKt5IBY1tGOFRiL67LovK2y9J+5WaIeCWaNcjrFA=";
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

            devshells.default = {

              commands = with pkgs; [
                {
                  package = nil;
                  category = "lsp";
                }
                {
                  package = nixd;
                  category = "lsp";
                }
                {
                  package = nixfmt-rfc-style;
                  category = "formatter";
                }
                {
                  package = statix;
                  category = "linter";
                }
                {
                  package = deadnix;
                  category = "linter";
                }
                { package = taplo; }
                { package = go; }
                { package = gotools; }
                {
                  package = gopls;
                  category = "lsp";
                }
                {
                  package = golines;
                  category = "formatter";
                }
                {
                  package = errcheck;
                  category = "linter";
                }
                {
                  package = marksman;
                  category = "lsp";
                }
                {
                  package = mdformat;
                  category = "formatter";
                }
              ];

              packages = with pkgs; [
                python313Packages.mdformat-gfm
                python313Packages.mdformat-gfm-alerts
              ];
            };
          };
      }
    );
}
