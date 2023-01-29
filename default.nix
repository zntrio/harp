{ lib, buildGo119Module, rev }:
buildGo119Module rec {
  pname = "harp";
  version = rev;

  src = lib.cleanSource ./.;

  subPackages = [ "cmd/harp" ];
  vendorSha256 = "sha256-/3FgROnwjMqXJ9hEEJOdVXh24r1hXmiAV+0iump2j3w=";
  ldflags = [ 
    "-s" "-w" "-buildid="
    "-X github.com/zntrio/harp/v2/build/version.Name=harp"
    "-X github.com/zntrio/harp/v2/build/version.AppName=github.com/zntrio/harp/v2/cmd/harp"
		"-X github.com/zntrio/harp/v2/build/version.Version=nix-${rev}"
    "-X github.com/zntrio/harp/v2/build/version.Commit=${rev}"
		"-X github.com/zntrio/harp/v2/build/version.Branch=main"
		"-X github.com/zntrio/harp/v2/build/version.BuildTags=defaults"
  ];

  meta = with lib; {
    description = "Secret management by contract toolchain";
    homepage = "https://github.com/zntrio/harp";
    license = licenses.asl20;
    platforms = platforms.unix;
  };
}