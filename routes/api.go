package routes

import (
	"github.com/BHAVYAghub/Youtube-API/controller"
	"github.com/BHAVYAghub/Youtube-API/service"
	"github.com/gin-gonic/gin"
)

const (
	GET    = "findAll"
	SEARCH = "find"
)

// TODO: Add constants

func AttachRoutes(r *gin.Engine, ytsrvc service.Service) {
	ctrl := controller.New(ytsrvc)

	api := r.Group("/youtube/")
	{
		api.GET(GET, ctrl.GetData)
		api.GET(SEARCH, ctrl.GetSearchData)
	}
}
