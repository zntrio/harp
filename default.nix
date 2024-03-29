{ lib, buildGoModule, rev }:
buildGoModule rec {
  pname = "harp";
  version = rev;

  src = lib.cleanSource ./.;

  subPackages = [ "cmd/harp" ];
  vendorHash = "sha256-BpLn0ykQGR2uRoXQwY5tqAo+gm7/Lx8zgh1BHEOFgwQ=";
  ldflags = [
    "-s" "-w" "-buildid="
    "-X zntr.io/harp/v2/build/version.Name=harp"
    "-X zntr.io/harp/v2/build/version.AppName=zntr.io/harp/v2/cmd/harp"
    "-X zntr.io/harp/v2/build/version.Version=nix-${rev}"
    "-X zntr.io/harp/v2/build/version.Commit=${rev}"
    "-X zntr.io/harp/v2/build/version.Branch=main"
    "-X zntr.io/harp/v2/build/version.BuildTags=defaults"
  ];

  meta = with lib; {
    description = "Secret management by contract toolchain";
    homepage = "https://github.com/zntrio/harp";
    license = licenses.asl20;
    platforms = platforms.unix;
  };
}
