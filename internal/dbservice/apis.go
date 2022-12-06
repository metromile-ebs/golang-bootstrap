package dbservice

import (
	"context"
	"metromile-ebs/streamline-graph-service/internal/datamodel"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var GraphCollection *mongo.Collection = GetCollection(DB, "graph")

func CreateRecord(c *gin.Context) {
	var graph datamodel.Graph
	err := c.BindJSON(&graph)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	isRecordInserted := insertRecord(graph)
	if isRecordInserted {
		c.IndentedJSON(http.StatusCreated, "success")
	}
}

func GetSingleGraph(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	graphId := c.Param("id")
	var graph datamodel.Graph
	defer cancel()

	err := GraphCollection.FindOne(ctx, bson.M{"id": graphId}).Decode(&graph)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, graph)
}

func GetGraphs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var graphs []datamodel.Graph
	defer cancel()

	results, err := GraphCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var graph datamodel.Graph
		if err = results.Decode(&graph); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		graphs = append(graphs, graph)
	}

	c.JSON(http.StatusOK, graphs)
}

func UpdateSingleGraph(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	graphId := c.Param("id")
	var graph datamodel.Graph
	defer cancel()
	if err := c.BindJSON(&graph); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := GraphCollection.UpdateOne(ctx, bson.M{"id": graphId}, bson.M{"$set": graph})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var updatedGraph datamodel.Graph
	if result.MatchedCount == 1 {
		err := GraphCollection.FindOne(ctx, bson.M{"id": graphId}).Decode(&updatedGraph)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, updatedGraph)
}
