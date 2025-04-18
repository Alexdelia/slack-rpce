{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    treefmt-nix.url = "github:numtide/treefmt-nix";

    pre-commit-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
        };
        treefmtEval = inputs.treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
      in {
        # nix build
        # packages.default = pkgs.callPackage ./default.nix {inherit pkgs;};

        # nix flake check
        checks = {
          formatting = treefmtEval.config.build.check self;
        };

        # nix develop
        devShells.default = pkgs.callPackage ./shell.nix {inherit pkgs inputs;};

        # nix fmt
        formatter = treefmtEval.config.build.wrapper;
      }
    );
}
