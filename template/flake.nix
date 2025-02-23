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
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devshell.flakeModule
        inputs.bited-utils.flakeModule
      ];
      systems = import systems;
      perSystem =
        { config, pkgs, ... }:
        {

          bited-utils = {
            name = "FIXME"; # Change this to your font's name
            version = builtins.readFile ./VERSION;
            src = ./.;
            # nerd = true; # Uncomment this line to enable Nerd Font patching
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
              { package = config.bited-utils.bited-clr; }
            ];

            packages = with pkgs; [
              python313Packages.mdformat-gfm
              python313Packages.mdformat-gfm-alerts
            ];
          };
        };
    };
}
