package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdrant/go-client/qdrant"
	"github.com/rhydianjenkins/rag-mcp-server/src"
)

func chunkText(text string, maxChunkSize int) []string {
	paragraphs := strings.Split(text, "\n\n")

	var chunks []string
	currentChunk := ""

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		if len(currentChunk) > 0 && len(currentChunk)+len(para) > maxChunkSize {
			chunks = append(chunks, currentChunk)
			currentChunk = para
		} else {
			if len(currentChunk) > 0 {
				currentChunk += "\n\n" + para
			} else {
				currentChunk = para
			}
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

func readTextFiles(dataDir string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".txt") {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v", path, err)
				return nil
			}

			relPath, _ := filepath.Rel(dataDir, path)
			files[relPath] = string(content)
		}

		return nil
	})

	return files, err
}

func Index(ollamaURL string, dataDir string, chunkSize int) error {
	storage, err := src.Connect(ollamaURL)

	if err != nil {
		log.Fatal("Unable to create storage")
		return err
	}

	files, err := readTextFiles(dataDir)
	if err != nil {
		log.Fatalf("Unable to read files from directory %s: %v", dataDir, err)
		return err
	}

	if len(files) == 0 {
		log.Printf("Warning: No .txt files found in %s", dataDir)
	}

	log.Printf("Found %d files to index (chunk size: %d chars)", len(files), chunkSize)

	var points []*qdrant.PointStruct
	pointID := uint64(1)

	for filename, content := range files {
		log.Printf("Processing file: %s (%d bytes)", filename, len(content))

		chunks := chunkText(content, chunkSize)

		log.Printf("  Split into %d chunks", len(chunks))

		for chunkIdx, chunk := range chunks {
			embedding, err := storage.GetEmbedding(chunk)
			if err != nil {
				log.Printf("Error generating embedding for %s chunk %d: %v", filename, chunkIdx, err)
				continue
			}

			payload := map[string]interface{}{
				"filename":    filename,
				"chunk_index": chunkIdx,
				"content":     chunk,
			}

			point := &qdrant.PointStruct{
				Id:      qdrant.NewIDNum(pointID),
				Vectors: qdrant.NewVectors(embedding...),
				Payload: qdrant.NewValueMap(payload),
			}

			points = append(points, point)
			pointID++
		}
	}

	if len(points) == 0 {
		return fmt.Errorf("no points to index")
	}

	log.Printf("Indexing %d total chunks into database", len(points))

	err = storage.GenerateDb(points)

	if err != nil {
		log.Fatal("Unable to generate db")
		return err
	}

	log.Printf("Successfully indexed %d chunks from %d files", len(points), len(files))

	return nil
}
