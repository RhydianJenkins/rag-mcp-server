package mcp

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (rs *MCPServer) Run(ctx context.Context) error {
	log.Println("Starting RAG MCP server on stdio transport...")

	if err := rs.mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return err
	}

	return nil
}
