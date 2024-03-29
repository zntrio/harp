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
                  version = "52e43a48497695a73d350b04a7e03423e3dd4909";

                  src = pkgs.fetchFromGitHub {
                    owner = "daixiang0";
                    repo = "gci";
                    rev = "${version}";
                    sha256 = "sha256-EBklnsZR8VwM89BztligZBDapmgyKuI9Ns8EFFo3ai8=";
                  };

                  vendorHash = "sha256-g7htGfU6C2rzfu8hAn6SGr0ZRwB8ZzSf9CgHYmdupE8=";

                  nativeBuildInputs = [ pkgs.installShellFiles ];
                };

                wwhrd = pkgs.buildGoModule rec {
                  pname = "wwhrd";
                  version = "b3052af659df6c012ddfe5c433c98953005faae8";

                  src = pkgs.fetchFromGitHub {
                    owner = "frapposelli";
                    repo = "wwhrd";
                    rev = "${version}";
                    sha256 = "sha256-jppPciJR2A1uKWhr+CB1Aquh33JXag0gjoBvohdt6+k=";
                  };

                  vendorHash = null;

                  nativeBuildInputs = [ pkgs.installShellFiles ];
                };

                cyclonedx-gomod = pkgs.buildGoModule rec {
                  pname = "cyclonedx-gomod";
                  version = "5e460935632418ad303008501b7dd498eb646b6d";

                  src = pkgs.fetchFromGitHub {
                    owner = "CycloneDX";
                    repo = "cyclonedx-gomod";
                    rev = "${version}";
                    sha256 = "sha256-BTvvVwompM7Ie4BsRcjT4yxLLJc1qH1514A4uzwgGCc=";
                  };

                  vendorHash = "sha256-uJgC44JtBSEwR6iKg8nmVmsFdyc3rtKSYxt8dbCABe0=";
                  subPackages = [ "cmd/cyclonedx-gomod" ];

                  nativeBuildInputs = [ pkgs.installShellFiles ];
                };
              };
            in
            pkgs.mkShell
              {
                buildInputs = [
                  pkgs.go_1_21
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
