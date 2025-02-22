{
  version,
  vendorHash,

  buildGoModule,
  licenses,
  ...
}:

buildGoModule {
  inherit version vendorHash;
  pname = "bited-clr";
  src = ../.;
  subPackages = [ "bited-clr" ];

  meta = {
    description = "A TUI tool for fine-tuned coloring of TXT/CLR pairs";
    longDescription = ''
      bited-clr is a TUI tool for fine-tuned coloring of TXT/CLR pairs.
    '';
    homepage = "https://github.com/molarmanful/bited-utils";
    license = licenses.mit;
  };
}
