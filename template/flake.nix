{
  description = "A cool font waiting to be built";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    bited-utils = {
      url = "github:molarmanful/bited-utils";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-parts.follows = "flake-parts";
      };
    };
  };

  outputs =
    inputs@{ systems, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [ inputs.bited-utils.flakeModule ];
      systems = import systems;
      perSystem =
        {
          config,
          pkgs,
          self',
          ...
        }:
        {

          bited-utils = {
            name = "FIXME"; # Change this to your font's name
            version = builtins.readFile ./VERSION;
            src = ./.;
            # Uncomment the following line to enable Nerd Font patching:
            # nerd = true;
          };

          # Devtools available via `nix develop` or direnv.
          # Add or remove as you wish.
          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              config.bited-utils.bited-clr
              taplo
              # lsps
              nil
              nixd
              marksman
              # formatters
              nixfmt-rfc-style
              mdformat
              python3Packages.mdformat-gfm
              python3Packages.mdformat-gfm-alerts
              # linters
              statix
              deadnix
            ];
          };

          formatter = pkgs.writeShellApplication {
            name = "linter";
            runtimeInputs = self'.devShells.default.nativeBuildInputs;
            text = ''
              find . -iname '*.nix' -exec nixfmt {} \; -exec deadnix -e {} \; -exec statix fix {} \;
              find . -iname '*.toml' -exec taplo fmt {} \;
              find . -iname '*.md' -exec mdformat {} \;
            '';
          };
        };
    };
}
