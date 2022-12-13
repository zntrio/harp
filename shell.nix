{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    go_1_19
    mage
    protobuf
    go-task
  ];
  shellHook =
  ''
    export PATH=$(pwd)/tools/bin:$PATH
  '';
}
