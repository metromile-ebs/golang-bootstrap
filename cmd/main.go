package main

import (
	"metromile-ebs/streamline-graph-service/internal/dbservice"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	dbservice.ConnectDB()
	router.GET("/graphs", dbservice.GetGraphs)
	router.GET("/graph/:id", dbservice.GetSingleGraph)
	router.PUT("/graph/:id", dbservice.UpdateSingleGraph)
	router.POST("/create", dbservice.CreateRecord)
	router.Run()

}
