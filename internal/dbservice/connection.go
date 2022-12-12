package dbservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	caFilePath = "rds-combined-ca-bundle.pem"

	// Timeout operations after N seconds
	// password        = "4ZxfcI0zah5u"
	// clusterEndpoint = "127.0.0.1:27018"
	// clusterEndpoint = "test-docdb-please-terraform.cdhqe2ld6wv2.us-west-2.docdb.amazonaws.com:27017"
	// connectionStringTemplate = "mongodb://127.0.0.1:27017"
	connectionStringTemplate = "mongodb://%s:%s@%s/?connect=direct&tls=true&tlsCAFile=%s&&tlsInsecure=true"
)

func ConnectDB() *mongo.Client {
	username := os.Getenv("GRAPH_DB_USER")
	password := os.Getenv("GRAPH_DB_PASS")
	clusterEndpoint := os.Getenv("GRAPH_DB_HOST")
	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint, caFilePath)
	// connectionURI := connectionStringTemplate
	fmt.Println("Connection uri is: " + connectionURI + "\nFile path:" + caFilePath)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		// utils.Logger.Error("Could not bind json for graph structure: " + err.Error())
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to DocumentDB!")
	return client
}
