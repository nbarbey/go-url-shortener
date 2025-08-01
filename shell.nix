let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
  hardeningDisable = [ "fortify" ];
  nativeBuildInputs = with pkgs; [
    pkgs.binutils
    pkgs.go
    pkgs.delve
    pkgs.gcc
    pkgs.k3s
  ];
}
