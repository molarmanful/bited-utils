{
  description = "A cool font waiting to be built";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    bited-utils = {
      url = "github:molarmanful/bited-utils";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "flake-utils";
      };
    };
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      bited-utils,
      ...
    }:

    let
      name = "FIXME";
      version = builtins.readFile ./VERSION;
    in

    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        bupkgs = bited-utils.packages.${system};
      in
      rec {

        # Defined here are wrappers around bited-build and bited-img.
        # - bited-build wrappers can be called via `nix build`.
        # - bited-img wrappers around bited-img can be called via `nix run`.
        # You can also define your own wrappers here or customize internals:
        # - bited-build wrapper is defined in `default.nix`.
        # - bited-img wrapper is defined in `img.nix`.
        packages =
          let
            build =
              o:
              pkgs.callPackage ./. (
                {
                  inherit version;
                  inherit (bupkgs) bited-build;
                }
                // o
              );
          in
          {

            # `nix build` OR `nix build .#FIXME`
            # Default build variant.
            ${name} = build { pname = name; };

            # `nix build .#FIXME-nerd`
            # Uncomment the following to enable a build variant with Nerd Font
            # patches.
            #
            # "${name}-nerd" = build {
            #   pname = "${name}-nerd";
            #   nerd = true;
            # };

            # `nix build .#FIXME-release`
            # Release build variant, which includes a versioned ZIP file of
            # build outputs.
            "${name}-release" = build {
              pname = "${name}-release";
              release = true;
              # Uncomment the following to include Nerd Font patches in your
              # release builds.
              #
              # nerd = true;
            };

            # `nix run .#FIXME-img`
            "${name}-img" = pkgs.callPackage ./img.nix {
              inherit (bupkgs) bited-img;
              name = "${name}-img";
            };

            default = packages.${name};
          };

        # Devtools available via `nix develop` or direnv.
        # Add or remove as you wish.
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nil
            nixd
            nixfmt-rfc-style
            statix
            deadnix
            marksman
            markdownlint-cli
            actionlint
            taplo
            bupkgs.bited-clr
          ];
        };

      }
    );
}
