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
                  version = "4d1d970b7dd8b2c09f0dd46cc749faeb0cee1f9d";

                  src = pkgs.fetchFromGitHub {
                    owner = "frapposelli";
                    repo = "wwhrd";
                    rev = "${version}";
                    sha256 = "sha256-OxJFxa833AmjFIOeLI94SxpP/Jlz7d4qglQwQm9TPG4=";
                  };

                  vendorSha256 = null;

                  nativeBuildInputs = [ pkgs.installShellFiles ];
                };

                cyclonedx-gomod = pkgs.buildGoModule rec {
                  pname = "cyclonedx-gomod";
                  version = "8d48a3aabb6623b715555cc53d1ad61bfd9d3e20";

                  src = pkgs.fetchFromGitHub {
                    owner = "CycloneDX";
                    repo = "cyclonedx-gomod";
                    rev = "${version}";
                    sha256 = "sha256-wrBFt06ym3tfSMJW+1KSv6J9dGwGqpqHIgjn50VlU0k=";
                  };

                  vendorSha256 = "sha256-1gEDHn7EcNVY10ahaEirn8KghgxpsS2pLC7HoycE428=";
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
