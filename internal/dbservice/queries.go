package dbservice

import (
	"context"
	"log"
	"metromile-ebs/streamline-graph-service/internal/datamodel"
	"time"
)

func insertRecord(graph datamodel.Graph) bool {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	res, err := GraphCollection.InsertOne(ctx, graph)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}

	id := res.InsertedID
	log.Printf("Inserted document ID: %s", id)

	return true
}
