{
  withSystem,
  mkPerSystemOption,
  lib,
  ...
}:
_: {
  options.perSystem = mkPerSystemOption (
    {
      config,
      pkgs,
      system,
      ...
    }:

    let
      cfg = config.bited-utils;
    in

    {
      options.bited-utils = {
        name = lib.mkOption {
          type = lib.types.nullOr lib.types.nonEmptyStr;
          default = null;
          description = "Font family name.";
        };
        version = lib.mkOption {
          type = lib.types.nonEmptyStr;
          default = "v0.0.0-0";
          description = "Font package version.";
        };
        nerd = lib.mkEnableOption "Nerd Fonts patching";
        buildTransformer = lib.mkOption {
          type = lib.types.functionTo lib.types.package;
          default = build: build;
          description = "A function to transform the bited-build wrapper derivation.";
        };
        imgTransformer = lib.mkOption {
          type = lib.types.functionTo lib.types.package;
          default = img: img;
          description = "A function to transform the bited-img wrapper derivation.";
        };
      };

      config = {
        packages =
          let
            build =
              o:
              cfg.buildTransformer (
                pkgs.callPackage ./nix/build.nix (
                  {
                    inherit (cfg) version;
                    inherit (config.packages) bited-build;
                  }
                  // o
                )
              );
          in

          (builtins.listToAttrs (
            builtins.map
              (name: {
                inherit name;
                value = withSystem system ({ config, ... }: config.packages.${name});
              })
              [
                "bited-build"
                "bited-img"
                "bited-scale"
                "bited-clr"
              ]
          ))

          // lib.mkIf (cfg.name != null) {
            ${cfg.name} = build { pname = cfg.name; };
            "${cfg.name}-nerd" = build {
              pname = "${cfg.name}-nerd";
              nerd = true;
            };
            "${cfg.name}-release" = build {
              inherit (cfg) nerd;
              pname = "${cfg.name}-release";
              release = true;
            };
            "${cfg.name}-img" = cfg.imgTransformer (
              pkgs.callPackage ./nix/img.nix {
                inherit (config.packages) bited-img;
                name = "${cfg.name}-img";
              }
            );
          };
      };
    }
  );
}
