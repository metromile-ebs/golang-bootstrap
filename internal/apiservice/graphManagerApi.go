package apiservice

import (
	"context"
	"fmt"
	"metromile-ebs/streamline-graph-manager/internal/datamodel"
	"metromile-ebs/streamline-graph-manager/internal/dbservice"
	"metromile-ebs/streamline-graph-manager/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type graphManagerService struct {
	MongoClient *mongo.Client
}

var Logger *zap.Logger = utils.FileLogger()

func NewgraphManagerService() graphManagerService {
	g := graphManagerService{}
	return g
}

func (g graphManagerService) Start() {
	g.MongoClient = dbservice.ConnectDB()
	router := gin.Default()

	//controller
	router.GET("/graphs", getAllGraphs(g))
	router.GET("/graphs/:id", GetSingleGraph(g))
	router.PUT("/graphs/:id", UpdateSingleGraph(g))
	router.POST("/graphs", createGraph(g))
	router.Run()
}

func getAllGraphs(graphManagerService graphManagerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var graphs []datamodel.Graph
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		collection := getCollection("graph", graphManagerService.MongoClient)

		results, err := collection.Find(ctx, bson.M{})
		if err != nil {
			err = fmt.Errorf("No graph found. Error is %w\n", err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var graph datamodel.Graph
			if err = results.Decode(&graph); err != nil {
				err = fmt.Errorf("Unable to decode results from db: %w\n", err)
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			graphs = append(graphs, graph)
		}

		c.JSON(http.StatusOK, graphs)
	}
}

func createGraph(graphManagerService graphManagerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var graph map[string]interface{}
		err := c.BindJSON(&graph)
		if err != nil {
			err = fmt.Errorf("Could not bind json for graph structure: %w\n", err)
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		collection := getCollection("graph", graphManagerService.MongoClient)
		insertResult, err := collection.InsertOne(ctx, graph)
		if err != nil {
			err = fmt.Errorf("Failed to insert graph: %w\n", err)
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		c.IndentedJSON(http.StatusCreated, insertResult.InsertedID)
	}
}

func GetSingleGraph(graphManagerService graphManagerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		graphObjectId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			err = fmt.Errorf("could not create object id from idParam: %s. %w\n", idParam, err)
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var graph map[string]interface{}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		collection := getCollection("graph", graphManagerService.MongoClient)

		// fmt.Println("Graph id is:" + graphId)
		err = collection.FindOne(ctx, bson.M{"_id": graphObjectId}).Decode(&graph)
		if err != nil {
			err = fmt.Errorf("Graph ID:"+graphObjectId.String()+".Error is: %w\n", err)
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, graph)
	}
}

func UpdateSingleGraph(graphManagerService graphManagerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		graphObjectId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			err = fmt.Errorf("could not create object id from idParam: %s. %w\n", idParam, err)
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var graph map[string]interface{}
		if err := c.BindJSON(&graph); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		collection := getCollection("graph", graphManagerService.MongoClient)
		result, err := collection.UpdateOne(ctx, bson.M{"_id": graphObjectId}, bson.M{"$set": graph})
		if err != nil {
			err = fmt.Errorf("Could not update the requested graphId %s. Error is %w\n", graphObjectId.String(), err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		var updatedGraph map[string]interface{}
		if result.MatchedCount == 1 {
			err = collection.FindOne(ctx, bson.M{"_id": graphObjectId}).Decode(&updatedGraph)
			if err != nil {
				err = fmt.Errorf("Unable to find updated graph with graph Id: %s. Error is %w\n", graphObjectId.String(), err)
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		c.JSON(http.StatusOK, updatedGraph)
	}
}

func getCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	collection := client.Database("graphservice").Collection(collectionName)
	return collection
}
