{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-22.11";
  };
  outputs = { self, nixpkgs }:
    let
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          terraform-provider-cachix = pkgs.buildGoModule {
            pname = "terraform-provider-cachix";
            inherit version;
            src = ./.;
            vendorSha256 = "sha256-bsuWsIOUlgzTOVXqLk2YgI4Nw6yPI1V7D3Bzoew7Y2o=";
            preBuild = ''
              ${pkgs.go-swagger}/bin/swagger generate client -f swagger/swagger-v1.json -T swagger/templates --allow-template-override
            '';
          };
        });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools go-swagger ];
            shellHook = ''
              if [ ! -d "client" ] || [ ! -d "models" ]; then
                ${pkgs.go-swagger}/bin/swagger generate client -f swagger/swagger-v1.json -T swagger/templates --allow-template-override
              fi
            '';
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.terraform-provider-cachix);
    };
}
