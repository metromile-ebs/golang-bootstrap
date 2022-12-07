package dbservice

import (
	"context"
	"metromile-ebs/streamline-graph-manager/internal/datamodel"
	"metromile-ebs/streamline-graph-manager/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var GraphCollection *mongo.Collection = GetCollection(DB, "graph")

func CreateRecord(c *gin.Context) {
	var graph datamodel.Graph
	err := c.BindJSON(&graph)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		utils.Logger.Error("Could not bind json for graph structure: " + err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	res, err := GraphCollection.InsertOne(ctx, graph)
	if err != nil {
		utils.Logger.Error("Failed to insert graph: " + err.Error())
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}

	// id := res.InsertedID
	utils.Logger.Info("Inserted document", zap.Any("data", res))
	c.IndentedJSON(http.StatusCreated, res.InsertedID)
}

func GetSingleGraph(c *gin.Context) {
	graphId := c.Param("id")
	var graph datamodel.Graph

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := GraphCollection.FindOne(ctx, bson.M{"id": graphId}).Decode(&graph)
	if err != nil {
		utils.Logger.Error("Graph ID:" + graphId + ".Error is:" + err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, graph)
}

func GetGraphs(c *gin.Context) {
	var graphs []datamodel.Graph
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := GraphCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.Logger.Error("No graph found .Error is:" + err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var graph datamodel.Graph
		if err = results.Decode(&graph); err != nil {
			utils.Logger.Error("Unable to decode results from db:" + err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		graphs = append(graphs, graph)
	}

	c.JSON(http.StatusOK, graphs)
}

func UpdateSingleGraph(c *gin.Context) {
	graphId := c.Param("id")
	var graph datamodel.Graph
	if err := c.BindJSON(&graph); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	result, err := GraphCollection.UpdateOne(ctx, bson.M{"id": graphId}, bson.M{"$set": graph})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		utils.Logger.Error("Could not update the requested graphId" + graphId + ".Error is:" + err.Error())
		return
	}

	var updatedGraph datamodel.Graph
	if result.MatchedCount == 1 {
		err := GraphCollection.FindOne(ctx, bson.M{"id": graphId}).Decode(&updatedGraph)
		if err != nil {
			utils.Logger.Error("Unable to find updated graph with graph Id:" + graphId + ".Error is:" + err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, updatedGraph)
}
