{
  pkgs,
  lib,
}:
 pkgs.buildGoModule rec  {
    pname = "go-url-shortener";
    version = "1.0.0";

    src = lib.fileset.toSource {
      root = ./.;
      fileset = ./.;
    };

    meta = {
      description = "Simple url-shortener service, written in Go";
      homepage = "https://github.com/nbarbey/go-url-shortener";
    };

    vendorHash = null;
}
