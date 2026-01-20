package handlers

import (
	"log"

	"github.com/qdrant/go-client/qdrant"
	"github.com/rhydianjenkins/rag-mcp-server/src"
)

func Index() error {
	storage, err := src.Connect()

	if err != nil {
		log.Fatal("Unable to create storage");
		return err;
	}

	vec1, _ := storage.GetEmbedding("eggs")
	vec2, _ := storage.GetEmbedding("beef")
	vec3, _ := storage.GetEmbedding("cheese")

	points := []*qdrant.PointStruct{
		{
			Id:      qdrant.NewIDNum(1),
			Vectors: qdrant.NewVectors(vec1...),
		},
		{
			Id:      qdrant.NewIDNum(2),
			Vectors: qdrant.NewVectors(vec2...),
		},
		{
			Id:      qdrant.NewIDNum(3),
			Vectors: qdrant.NewVectors(vec3...),
		},
	}

	err = storage.GenerateDb(points);

	if err != nil {
		log.Fatal("Unable to generate db");
		return err;
	}

	return nil;
}
