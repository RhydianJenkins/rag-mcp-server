package config

import (
	"os"
	"sync"
)

type Config struct {
	CollectionName string
	EmbeddingModel string
	OllamaURL      string
	OllamaAPIKey   string
	QdrantHost     string
	QdrantPort     int
	QdrantAPIKey   string
	ServerName     string
	ServerVersion  string
	VectorSize     uint64
}

var (
	instance *Config
	once     sync.Once
)

func Initialize(cfg *Config) {
	once.Do(func() {
		instance = applyDefaults(cfg)
	})
}

func Get() *Config {
	if instance == nil {
		Initialize(&Config{})
	}
	return instance
}

func applyDefaults(cfg *Config) *Config {
	if cfg.CollectionName == "" {
		cfg.CollectionName = "my_collection"
	}
	if cfg.EmbeddingModel == "" {
		cfg.EmbeddingModel = "nomic-embed-text"
	}
	if cfg.OllamaURL == "" {
		cfg.OllamaURL = "http://localhost:11434"
	}
	if cfg.OllamaAPIKey == "" {
		cfg.OllamaAPIKey = os.Getenv("OLLAMA_API_KEY")
	}
	if cfg.QdrantHost == "" {
		cfg.QdrantHost = "localhost"
	}
	if cfg.QdrantPort == 0 {
		cfg.QdrantPort = 6334
	}
	if cfg.QdrantAPIKey == "" {
		cfg.QdrantAPIKey = os.Getenv("QDRANT_API_KEY")
	}
	if cfg.ServerName == "" {
		cfg.ServerName = "seek"
	}
	if cfg.ServerVersion == "" {
		cfg.ServerVersion = "dev"
	}
	if cfg.VectorSize == 0 {
		cfg.VectorSize = 768
	}
	return cfg
}
