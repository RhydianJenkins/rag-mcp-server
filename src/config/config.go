package config

type Config struct {
	CollectionName string
	EmbeddingModel string
	OllamaURL      string
	QdrantHost     string
	QdrantPort     int
	QdrantUseTLS   bool
	ServerName     string
	ServerVersion  string
	VectorSize     uint64
}

func DefaultConfig() *Config {
	return &Config{
		CollectionName: "my_collection",
		EmbeddingModel: "nomic-embed-text",
		OllamaURL:      "http://localhost:11434",
		QdrantHost:     "localhost",
		QdrantPort:     6334,
		QdrantUseTLS:   false,
		ServerName:     "rag-mcp-server",
		ServerVersion:  "1.0.0",
		VectorSize:     768,
	}
}
