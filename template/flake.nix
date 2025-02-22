{
  description = "A cool font waiting to be built";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    systems.url = "systems";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devshell.url = "github:numtide/devshell";
    bited-utils.url = "github:molarmanful/bited-utils";
  };

  outputs =
    inputs@{ systems, flake-parts, ... }:

    let
      name = "FIXME";
      version = builtins.readFile ./VERSION;
    in

    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devshell.flakeModule
        inputs.bited-utils.flakeModule
      ];
      systems = import systems;
      perSystem =
        { config, pkgs, ... }:
        {
          inherit name version;
          # nerd = true;

          packages = {
            default = config.packages.${name};
          };

          # Devtools available via `nix develop` or direnv.
          # Add or remove as you wish.
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
              {
                package = marksman;
                category = "lsp";
              }
              {
                package = mdformat;
                category = "formatter";
              }
              {
                package = config.packages.bited-clr;
              }
            ];

            packages = with pkgs; [
              python313Packages.mdformat-gfm
              python313Packages.mdformat-gfm-alerts
            ];
          };
        };
    };
}
