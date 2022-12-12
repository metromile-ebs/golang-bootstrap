package dbservice

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Path to the AWS CA file
	caFilePath = "rds-combined-ca-bundle.pem"

	// Timeout operations after N seconds
	connectTimeout           = 5
	queryTimeout             = 30
	username                 = "mongo"
	password                 = "4ZxfcI0zah5u"
	clusterEndpoint          = "127.0.0.1:27018"
	connectionStringTemplate = "mongodb://127.0.0.1:27017"
	// connectionStringTemplate = "mongodb://%s:%s@%s/?connect=direct&tls=true&tlsCAFile=rds-combined-ca-bundle.pem&&tlsInsecure=true"
)

func ConnectDB() *mongo.Client {
	// envUser := os.Environ("DB_USER");
	// os.Getenv()
	// connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)
	connectionURI := connectionStringTemplate
	fmt.Println("Connection uri is: " + connectionURI + "\nFile path:" + caFilePath)
	// utils.Logger.Info("I am logging")

	// tlsConfig, err := getCustomTLSConfig(caFilePath)
	// if err != nil {
	// 	log.Fatalf("Failed getting TLS configuration: %v", err)
	// }

	// client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI).SetTLSConfig(tlsConfig))
	// if err != nil {
	// 	log.Fatalf("Failed to create client: %v", err)
	// }

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		// utils.Logger.Error("Could not bind json for graph structure: " + err.Error())
		// log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
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

// func getCustomTLSConfig(caFile string) (*tls.Config, error) {
// 	tlsConfig := new(tls.Config)
// 	certs, err := ioutil.ReadFile(caFile)

// 	if err != nil {
// 		return tlsConfig, err
// 	}

// 	fmt.Println(certs)

// 	tlsConfig.RootCAs = x509.NewCertPool()
// 	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

// 	if !ok {
// 		return tlsConfig, errors.New("Failed parsing pem file")
// 	}

// 	return tlsConfig, nil
// }

// Disconnect on failure when app is stall.
// The service itself tried to connect it without user interaction.

// func GetCollection(collectionName string) *mongo.Collection {
// 	collection := .Database("graphservice").Collection(collectionName)
// 	return collection
// }
