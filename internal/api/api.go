package api

import (
	"github.com/gin-gonic/gin"
)

type api struct {
	Engine *gin.Engine
	Routes []Controller
}

func Init() {
	server := &api{
		Engine: gin.New(),
		Routes: []Controller{
			&Deployments{},
		},
	}

	server.DefaultMiddleware()
	server.BuildRoutes()

	server.Run()
}

func (api *api) Run() {
	api.Engine.Run()
}

func (api *api) DefaultMiddleware() {
	api.Engine.Use(gin.Recovery())
}

func (api *api) BuildRoutes() {
	health := &HealthCheck{}
	health.Build(api)

	api.Engine.Use(gin.Logger())

	for _, controller := range api.Routes {
		controller.Build(api)
	}
}
