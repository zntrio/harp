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
        packages.default = pkgs.callPackage ./default.nix {
          inherit rev;
        };

        devShells = {
          default =
            let
              devtools = {
                gci = pkgs.buildGoModule rec {
                  pname = "gci";
                  version = "8f9a4e94ae2a8db0093d52281bf7ac0c83eed0ce";

                  src = pkgs.fetchFromGitHub {
                    owner = "daixiang0";
                    repo = "gci";
                    rev = "${version}";
                    sha256 = "sha256-qWEEcIbTgYmGVnnTW+hM8e8nw5VLWN1TwzdUIZrxF3s=";
                  };

                  vendorSha256 = "sha256-dlt+i/pEP3RzW4JwndKTU7my2Nn7/2rLFlk8n1sFR60=";

                  nativeBuildInputs = [ pkgs.installShellFiles ];
                };

                wwhrd = pkgs.buildGoModule rec {
                  pname = "wwhrd";
                  version = "a4048a3e900ae413910a8477d7c8c2cf9eb9fc3a";

                  src = pkgs.fetchFromGitHub {
                    owner = "frapposelli";
                    repo = "wwhrd";
                    rev = "${version}";
                    sha256 = "sha256-z6xhRKTqsPFT0I76IVZKcCG90kR/0kKNvZbfPNFaJWw=";
                  };

                  vendorSha256 = null;

                  nativeBuildInputs = [ pkgs.installShellFiles ];
                };

                cyclonedx-gomod = pkgs.buildGoModule rec {
                  pname = "cyclonedx-gomod";
                  version = "2fe0a1da390fbc17df326f35139a5a4d9e1ffe65";

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
              };
            in
            pkgs.mkShell
              {
                buildInputs = [
                  pkgs.go_1_20
                  pkgs.gopls
                  pkgs.gotools
                  pkgs.go-tools
                  pkgs.gotestsum
                  pkgs.gofumpt
                  pkgs.golangci-lint
                  devtools.gci
                  devtools.cyclonedx-gomod
                  devtools.wwhrd
                  pkgs.mage
                  pkgs.buf
                  pkgs.mockgen
                  pkgs.protobuf
                  pkgs.protoc-gen-go
                  pkgs.protoc-gen-go-grpc
                  pkgs.cue
                ];
              };
        };
      }
    );
}
