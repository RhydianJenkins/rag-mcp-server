package handlers

import (
	"fmt"
	"log"

	"github.com/rhydianjenkins/rag-mcp-server/src"
)

func Search(searchTerm string, ollamaURL string) error {
	storage, err := src.Connect(ollamaURL)

	if err != nil {
		return err
	}

	searchResult, err := storage.Search(searchTerm)

	if err != nil {
		return err
	}

	log.Printf("\nSearch results for: '%s'\n", searchTerm)
	log.Printf("Found %d results:\n", len(searchResult))

	// Display results with metadata
	for i, result := range searchResult {
		log.Printf("\n--- Result %d (Score: %.4f) ---", i+1, result.Score)

		// Extract metadata from payload
		if result.Payload != nil {
			if filename, ok := result.Payload["filename"]; ok {
				log.Printf("File: %v", filename.GetStringValue())
			}
			if chunkIdx, ok := result.Payload["chunk_index"]; ok {
				log.Printf("Chunk: %v", chunkIdx.GetIntegerValue())
			}
			if content, ok := result.Payload["content"]; ok {
				contentStr := content.GetStringValue()
				// Truncate long content for display
				if len(contentStr) > 200 {
					contentStr = contentStr[:200] + "..."
				}
				log.Printf("Content: %s", contentStr)
			}
		}

		fmt.Println() // Empty line for readability
	}

	return nil
}
