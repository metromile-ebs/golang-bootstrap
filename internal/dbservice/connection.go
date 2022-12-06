package dbservice

import (
	"context"
	"fmt"
	"log"
	"time"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

const (
	// Path to the AWS CA file
	caFilePath = "rds-combined-ca-bundle.pem"

	// Timeout operations after N seconds
	connectTimeout = 5
	queryTimeout   = 30
	username       = "mongo"
	password       = "4ZxfcI0zah5u"
	// clusterEndpoint = "test-docdb-please-terraform.cluster-cdhqe2ld6wv2.us-west-2.docdb.amazonaws.com:27017"
	clusterEndpoint = "test-docdb-please-terraform.cdhqe2ld6wv2.us-west-2.docdb.amazonaws.com:27017"

	// Which instances to read from
	readPreference = "primary"

	connectionStringTemplate = "mongodb://127.0.0.1:27017"

	// connectionStringTemplate = "mongodb://%s:%s@%s/?tls=true&ssl_ca_certs=rds-combined-ca-bundle.pem&replicaSet=rs0"
	// connectionStringTemplate = "mongodb://" + username + ":" + password + "@" + clusterEndpoint + "/?ssl=false&ssl_ca_certs=rds-combined-ca-bundle.pem&replicaSet=rs0&readPreference=" + readPreference + "&retryWrites=false"
)

func ConnectDB() *mongo.Client {
	// connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)
	connectionURI := connectionStringTemplate
	fmt.Println("Connection uri is: " + connectionURI)

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
		log.Fatalf("Failed to create client: %v", err)
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

var DB *mongo.Client = ConnectDB()

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("graphservice").Collection(collectionName)
	return collection
}
