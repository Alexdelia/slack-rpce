_: {
  projectRootFile = ".git/config";

  programs = {
    gofmt.enable = true;

    alejandra.enable = true;
    deadnix.enable = true;
    statix.enable = true;

    prettier = {
      enable = true;
      settings = {
        useTabs = true;
      };
    };
  };

  settings.global.excludes = [
    ".gitingore"

    "*.lock"

    ".env*"

    "*.toml"
  ];
}
