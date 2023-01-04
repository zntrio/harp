{
  description = "harp";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    let
      supportedSystems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];
    in
    utils.lib.eachSystem supportedSystems (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
        rev =
          if (self ? shortRev)
          then self.shortRev
          else "dev";
      in
      rec
      {
        tools.gci = pkgs.buildGoModule rec {
          pname = "gci";
          version = "0a02b8b306dcdb9b8dd57ca5b1d4c161767bb545";

          src = pkgs.fetchFromGitHub {
            owner = "daixiang0";
            repo = "gci";
            rev = "${version}";
            sha256 = "sha256-qWEEcIbTgYmGVnnTW+hM8e8nw5VLWN1TwzdUIZrxF3s=";
          };

          vendorSha256 = "sha256-dlt+i/pEP3RzW4JwndKTU7my2Nn7/2rLFlk8n1sFR60=";

          nativeBuildInputs = [ pkgs.installShellFiles ];
        };

        tools.wwhrd = pkgs.buildGoModule rec {
          pname = "wwhrd";
          version = "13b50a0e5c6316c9126c537d140dd3efba040e41";

          src = pkgs.fetchFromGitHub {
            owner = "frapposelli";
            repo = "wwhrd";
            rev = "${version}";
            sha256 = "sha256-z6xhRKTqsPFT0I76IVZKcCG90kR/0kKNvZbfPNFaJWw=";
          };

          vendorSha256 = null;

          nativeBuildInputs = [ pkgs.installShellFiles ];
        };

        tools.cyclonedx-gomod = pkgs.buildGoModule rec {
          pname = "cyclonedx-gomod";
          version = "0313c5f801f55fcb38d7d18aa5b5f8c68b7667f4";

          src = pkgs.fetchFromGitHub {
            owner = "CycloneDX";
            repo = "cyclonedx-gomod";
            rev = "${version}";
            sha256 = "sha256-BjzZYakILJuK+sU0hkPDSs4Jt/48oYX1t5pN+xzJXTk=";
          };

          vendorSha256 = "sha256-8j39II91QjRqBNHu3jf5p850YAO6PqLcoOKMXsp1wXw=";
          subPackages = [ "cmd/cyclonedx-gomod" ];

          nativeBuildInputs = [ pkgs.installShellFiles ];
        };

        packages.harp = pkgs.callPackage ./default.nix {
          inherit rev;
        };

        defaultPackage = packages.harp;

        apps = {
          harp = utils.lib.mkApp {
            drv = packages.harp;
            exePath = "/bin/harp";
          };

          default = apps.harp;
        };

        devShell = pkgs.mkShell
          {
            buildInputs = [
              pkgs.go_1_19
              pkgs.gopls
              pkgs.gotools
              pkgs.go-tools
              pkgs.gotestsum
              pkgs.gofumpt
              pkgs.golangci-lint
              tools.gci
              tools.cyclonedx-gomod
              tools.wwhrd
              pkgs.mage
              pkgs.mockgen
              pkgs.protobuf
              pkgs.protoc-gen-go
              pkgs.protoc-gen-go-grpc
              pkgs.just
              pkgs.cue
            ];
          };
      });
}
