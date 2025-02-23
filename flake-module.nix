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
      options.bited-utils =
        {
          name = lib.mkOption {
            type = lib.types.nullOr lib.types.nonEmptyStr;
            default = null;
            description = "Font family name.";
          };
          src = lib.mkOption {
            type = lib.types.nullOr lib.types.path;
            default = null;
            description = "Location of both `bited-build.toml` and `bited-img.toml`.";
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
        }

        // builtins.listToAttrs (
          builtins.map
            (name: {
              inherit name;
              value = lib.mkOption {
                type = lib.types.package;
                default = withSystem system ({ config, ... }: config.packages.${name});
                description = "The ${name} package to use.";
              };
            })
            [
              "bited-build"
              "bited-img"
              "bited-scale"
              "bited-clr"
            ]
        );

      config = {
        packages = lib.mkIf (cfg.name != null && cfg.src != null) (
          let
            build =
              o:
              cfg.buildTransformer (
                pkgs.callPackage ./nix/build.nix ({ inherit (cfg) src version bited-build; } // o)
              );
            base = build { pname = cfg.name; };
          in
          {
            default = lib.mkDefault base;
            ${cfg.name} = base;
            "${cfg.name}-release" = build {
              inherit (cfg) nerd;
              pname = "${cfg.name}-release";
              release = true;
            };
            "${cfg.name}-img" = cfg.imgTransformer (
              pkgs.callPackage ./nix/img.nix {
                inherit (cfg) bited-img;
                name = "${cfg.name}-img";
              }
            );
          }
          // lib.optionalAttrs cfg.nerd {
            "${cfg.name}-nerd" = build {
              pname = "${cfg.name}-nerd";
              nerd = true;
            };
          }
        );
      };
    }
  );
}
