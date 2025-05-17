{
  description = "Kin flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};
      in {
        formatter = pkgs.alejandra;
        devShells.default = with pkgs;
          mkShell {
            buildInputs = [
              just
              alejandra
              go
              golangci-lint
            ];
          };

        packages.default = with pkgs;
          buildGoModule {
            pname = "kin";
            version = "0.1.0";
            src = ./.;
            # src = fetchFromGitHub {
            #     owner = "VinceDeslo";
            #     repo = "kin";
            #     tag = "v${finalAttrs.version}";
            #     hash = lib.fakeHash;
            # };
            vendorHash = "sha256-aVk0ccByDS4+gs2im4eU6S5daK3OVoRYoBxn3SSgDGw=";
            meta = {
              description = "A pretty k8s cluster access prompt for Teleport";
              homepage = "https://github.com/VinceDeslo/kin";
              license = lib.licenses.mit;
              maintainers = with lib.maintainers; [VinceDeslo];
            };
          };
      }
    );
}
