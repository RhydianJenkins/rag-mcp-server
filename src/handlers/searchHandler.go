package handlers

import (
	"fmt"

	"github.com/rhydianjenkins/rag-mcp-server/src"
)

func Search(searchTerm string, ollamaURL string, limit int) error {
	storage, err := src.Connect(ollamaURL)

	if err != nil {
		return err
	}

	searchResult, err := storage.Search(searchTerm, limit)

	if err != nil {
		return err
	}

	fmt.Printf("\nSearch results for: '%s'\n", searchTerm)
	fmt.Printf("Found %d results:\n", len(searchResult))

	for i, result := range searchResult {
		fmt.Printf("\n--- Result %d (Score: %.4f) ---\n", i+1, result.Score)

		if result.Payload != nil {
			if filename, ok := result.Payload["filename"]; ok {
				fmt.Printf("File: %v\n", filename.GetStringValue())
			}
			if chunkIdx, ok := result.Payload["chunk_index"]; ok {
				fmt.Printf("Chunk: %v\n", chunkIdx.GetIntegerValue())
			}
			if content, ok := result.Payload["content"]; ok {
				fmt.Println()
				fmt.Println(content.GetStringValue())
			}
		}

		fmt.Println()
	}

	return nil
}
