package main

import (
	"fmt"

	"github.com/rhydianjenkins/rag-mcp-server/src/handlers"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var ollamaAddress string

	var rootCmd = &cobra.Command{
		Short: "RAG MCP Server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringVar(&ollamaAddress, "ollamaAddress", "http://localhost:11434", "Ollama server address")

	var dataDir string
	var chunkSize int
	var indexCmd = &cobra.Command{
		Use: "index",
		Short: "Index the knowledge base",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			handlers.Index(ollamaAddress, dataDir, chunkSize)
		},
	}
	indexCmd.Flags().StringVar(&dataDir, "dataDir", "", "Directory containing .txt files to index (required)")
	indexCmd.Flags().IntVar(&chunkSize, "chunkSize", 1000, "Maximum chunk size in characters for splitting text")
	indexCmd.MarkFlagRequired("dataDir")
	rootCmd.AddCommand(indexCmd)

	var searchCmd = &cobra.Command{
		Use: "search",
		Short: "Search the knowledge base",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchTerm := args[0]
			handlers.Search(searchTerm, ollamaAddress)
		},
	}
	rootCmd.AddCommand(searchCmd)

	var runCmd = &cobra.Command{
		Use: "run",
		Short: "Run the mcp server",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO")
		},
	}
	rootCmd.AddCommand(runCmd)

	return rootCmd
}

func main() {
	initCmd().Execute()
}
