package main

import (
	"metromile-ebs/streamline-graph-manager/internal/dbservice"
	"metromile-ebs/streamline-graph-manager/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	dbservice.ConnectDB()
	utils.FileLogger()
	router.GET("/graphs", dbservice.GetGraphs)
	router.GET("/graph/:id", dbservice.GetSingleGraph)
	router.PUT("/graph/:id", dbservice.UpdateSingleGraph)
	router.POST("/create", dbservice.CreateRecord)
	router.Run()

}
