{
  description = "Raiden - A HTTP(S) Tracer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.stdenv.mkDerivation {
          pname = "raiden";
          version = "1.0.0";

          src = ./.;

          nativeBuildInputs = [ pkgs.go ];

          buildPhase = ''
            go build -o raiden
          '';

          installPhase = ''
            mkdir -p $out/bin
            cp raiden $out/bin/
          '';

          meta = with pkgs.lib; {
            description = "Raiden - A HTTP(S) Tracer";
            license = licenses.mit;
            maintainers = with maintainers; [ damroth ];
          };
        };
      });
}