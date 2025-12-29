{
  description = "Pipeline helpers and utilities for building fonts from bited BDFs";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
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
        imports = [ flakeModule ];

        flake = {
          inherit flakeModule;
          templates.default = {
            path = ./template;
            description = "bited font project with bited-utils";
          };
        };

        systems = import inputs.systems;
        perSystem =
          { pkgs, self', ... }:
          {

            packages =
              let
                bitsnpicas = pkgs.callPackage ./nix/bitsnpicas.nix { };
                args = {
                  inherit bitsnpicas;
                  version = builtins.readFile ./VERSION;
                  vendorHash = "sha256-pc7Q6QbnzQPDTMmwLkykwhtJSFW1qnbOyevWAzmr1dk=";
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
                gotools
                taplo
                # lsps
                nil
                nixd
                marksman
                gopls
                # formatters
                nixfmt-rfc-style
                mdformat
                python3Packages.mdformat-gfm
                python3Packages.mdformat-gfm-alerts
                golines
                # linters
                statix
                deadnix
                errcheck
              ];
            };

            formatter = pkgs.writeShellApplication {
              name = "linter";
              runtimeInputs = self'.devShells.default.nativeBuildInputs;
              text = ''
                find . -iname '*.nix' -exec nixfmt {} \; -exec deadnix -e {} \; -exec statix fix {} \;
                find . -iname '*.toml' -exec taplo fmt {} \;
                find . -iname '*.md' -exec mdformat {} \;
                find . -iname '*.go' -exec golines {} \;
                errcheck
              '';
            };
          };
      }
    );
}
