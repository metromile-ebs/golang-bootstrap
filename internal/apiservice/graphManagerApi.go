package apiservice

import (
	"context"
	"metromile-ebs/streamline-graph-manager/internal/datamodel"
	"metromile-ebs/streamline-graph-manager/internal/dbservice"
	"metromile-ebs/streamline-graph-manager/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
			Logger.Error("No graph found .Error is:" + err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var graph datamodel.Graph
			if err = results.Decode(&graph); err != nil {
				Logger.Error("Unable to decode results from db:" + err.Error())
				c.JSON(http.StatusInternalServerError, err.Error())
			}

			graphs = append(graphs, graph)
		}

		c.JSON(http.StatusOK, graphs)
	}
}

func createGraph(graphManagerService graphManagerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var graph datamodel.Graph
		err := c.BindJSON(&graph)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			Logger.Error("Could not bind json for graph structure: " + err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		collection := getCollection("graph", graphManagerService.MongoClient)
		_, err = collection.InsertOne(ctx, graph)
		if err != nil {
			Logger.Error("Failed to insert graph: " + err.Error())
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		c.IndentedJSON(http.StatusCreated, graph.Id)
	}
}

func GetSingleGraph(graphManagerService graphManagerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		graphId := c.Param("id")
		intGraphId, _ := strconv.Atoi(graphId)
		var graph datamodel.Graph

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		collection := getCollection("graph", graphManagerService.MongoClient)
		// fmt.Println("Graph id is:" + graphId)
		err := collection.FindOne(ctx, bson.M{"id": intGraphId}).Decode(&graph)

		if err != nil {
			Logger.Info("Graph ID:" + graphId + ".Error is:" + err.Error())
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, graph)
	}
}

func UpdateSingleGraph(graphManagerService graphManagerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		graphId, _ := strconv.Atoi(c.Param("id"))
		var graph datamodel.Graph
		if err := c.BindJSON(&graph); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		collection := getCollection("graph", graphManagerService.MongoClient)
		result, err := collection.UpdateOne(ctx, bson.M{"id": graphId}, bson.M{"$set": graph})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			Logger.Info("Could not update the requested graphId" + strconv.Itoa(graphId) + ".Error is:" + err.Error())
			return
		}

		var updatedGraph datamodel.Graph
		if result.MatchedCount == 1 {
			err = collection.FindOne(ctx, bson.M{"id": graphId}).Decode(&updatedGraph)
			if err != nil {
				Logger.Error("Unable to find updated graph with graph Id:" + strconv.Itoa(graphId) + ".Error is:" + err.Error())
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
