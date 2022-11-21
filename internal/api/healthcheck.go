package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheck struct{}

func (*HealthCheck) Build(api *api) {
	api.Engine.GET("/streamline/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, "OK")
	})
}
