{
  pkgs,
  inputs,
}: let
  pre-commit-check = inputs.pre-commit-hooks.lib.${pkgs.system}.run {
    src = ./.;
    hooks = {
      typos = {
        enable = true;
        stages = ["pre-commit"];
      };

      check-toml = {
        enable = true;
        stages = ["pre-commit"];
      };

      nix-check = {
        enable = true;

        name = "nix flake check";
        entry = "nix flake check";

        pass_filenames = false;
        always_run = true;
        stages = ["pre-push"];
      };

      nix-fmt = {
        enable = true;

        name = "nix fmt";
        entry = "nix fmt";

        pass_filenames = false;
        stages = ["pre-commit"];
      };
    };
  };
in
  pkgs.mkShell {
    buildInputs = with pkgs;
      [
        go
      ]
      ++ pre-commit-check.buildInputs;

    shellHook =
      /*
      bash
      */
      ''
        ${pre-commit-check.shellHook}

        if alias 'go' >/dev/null 2>&1; then
        	unalias go;
        fi

        alias fmt="nix fmt"
      '';
  }
