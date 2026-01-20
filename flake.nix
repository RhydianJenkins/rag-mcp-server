{
  description = "RAG MCP Server - A RAG server with MCP support";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    process-compose-flake.url = "github:Platonic-Systems/process-compose-flake";
  };

  outputs = { self, nixpkgs, flake-utils, process-compose-flake }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        # Define the process composition
        packages.default = pkgs.writeShellScriptBin "services" ''
          ${pkgs.process-compose}/bin/process-compose -f ${pkgs.writeText "process-compose.yaml" ''
            version: "0.5"

            log_level: info

            processes:
              qdrant:
                command: ${pkgs.qdrant}/bin/qdrant
                availability:
                  restart: on_failure
                ready_log_line: "Qdrant gRPC listening"

              ollama:
                command: ${pkgs.ollama}/bin/ollama serve
                availability:
                  restart: on_failure
                ready_log_line: "Listening on"

              ollama-init:
                command: |
                  echo "Waiting for Ollama to be ready..."
                  until ${pkgs.curl}/bin/curl -s http://localhost:11434/api/tags > /dev/null 2>&1; do
                    sleep 1
                  done
                  echo "Ollama is ready, checking for nomic-embed-text model..."

                  if ${pkgs.ollama}/bin/ollama list | grep -q "nomic-embed-text"; then
                    echo "Model nomic-embed-text already exists, skipping pull"
                  else
                    echo "Pulling nomic-embed-text model (this may take a few minutes on first run)..."
                    ${pkgs.ollama}/bin/ollama pull nomic-embed-text
                    echo "Model pulled successfully!"
                  fi
                depends_on:
                  ollama:
                    condition: process_started
                availability:
                  restart: "no"
          ''}
        '';

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/services";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            golangci-lint
            qdrant
            ollama
            process-compose
            curl
          ];

          shellHook = ''
            echo "RAG MCP Server development environment"
            echo "Go version: $(go version)"
            echo "Qdrant version: $(qdrant --version)"
            echo "Ollama version: $(ollama --version)"
            echo ""
            echo "Services available:"
            echo "  - Qdrant (vector database)"
            echo "  - Ollama (local embeddings with nomic-embed-text)"
            echo ""
            echo "To start all services (auto-pulls embedding model on first run):"
            echo "  nix run"
            echo ""
            echo "Everything will be set up automatically!"
          '';
        };
      }
    );
}
